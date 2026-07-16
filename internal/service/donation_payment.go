package service

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/azmiagr/garudahacks-hackathon/entity"
	"github.com/azmiagr/garudahacks-hackathon/internal/repository"
	"github.com/azmiagr/garudahacks-hackathon/model"
	"github.com/azmiagr/garudahacks-hackathon/pkg/config"
	"github.com/azmiagr/garudahacks-hackathon/pkg/database/mariadb"
	apperrors "github.com/azmiagr/garudahacks-hackathon/pkg/errors"
	"github.com/google/uuid"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"gorm.io/gorm"
)

const (
	paymentMethodQRIS           = "qris"
	paymentMethodVirtualAccount = "virtual_account"

	donationStatusPending  = "pending"
	donationStatusApproved = "approved"
	donationStatusRejected = "rejected"

	paymentStatusPending    = "pending"
	paymentStatusSettlement = "settlement"
	paymentStatusCapture    = "capture"
	paymentStatusExpire     = "expire"
	paymentStatusCancel     = "cancel"
	paymentStatusDeny       = "deny"
	paymentStatusFailure    = "failure"

	walletTransactionTypeDeposit = "deposit"
	minDonationAmount            = int64(10000)
)

type IDonationPaymentService interface {
	CreateDonationPayment(user *entity.User, req model.CreateDonationPaymentRequest) (*model.CreateDonationPaymentResponse, error)
	HandleMidtransNotification(req model.MidtransNotificationRequest) (*model.DonationLockStatusResponse, error)
}

type DonationPaymentService struct {
	db                           *gorm.DB
	requestRepository            repository.IRequestRepository
	itemRepository               repository.IItemRepository
	walletRepository             repository.IWalletRepository
	walletTransactionRepository  repository.IWalletTransactionRepository
	donationRepository           repository.IDonationRepository
	paymentTransactionRepository repository.IPaymentTransactionRepository
	orderRepository              repository.IOrderRepository
	orderItemRepository          repository.IOrderItemRepository
	custodyLogRepository         repository.ICustodyLogRepository
	pointService                 IPointService
	midtransConfig               *config.MidtransConfig
	coreClient                   coreapi.Client
}

func NewDonationPaymentService(
	requestRepository repository.IRequestRepository,
	itemRepository repository.IItemRepository,
	walletRepository repository.IWalletRepository,
	walletTransactionRepository repository.IWalletTransactionRepository,
	donationRepository repository.IDonationRepository,
	paymentTransactionRepository repository.IPaymentTransactionRepository,
	orderRepository repository.IOrderRepository,
	orderItemRepository repository.IOrderItemRepository,
	custodyLogRepository repository.ICustodyLogRepository,
	pointService IPointService,
	midtransConfig *config.MidtransConfig,
) IDonationPaymentService {
	return &DonationPaymentService{
		db:                           mariadb.Connection,
		requestRepository:            requestRepository,
		itemRepository:               itemRepository,
		walletRepository:             walletRepository,
		walletTransactionRepository:  walletTransactionRepository,
		donationRepository:           donationRepository,
		paymentTransactionRepository: paymentTransactionRepository,
		orderRepository:              orderRepository,
		orderItemRepository:          orderItemRepository,
		custodyLogRepository:         custodyLogRepository,
		pointService:                 pointService,
		midtransConfig:               midtransConfig,
		coreClient:                   midtransConfig.NewCoreAPIClient(),
	}
}

func (s *DonationPaymentService) CreateDonationPayment(user *entity.User, req model.CreateDonationPaymentRequest) (*model.CreateDonationPaymentResponse, error) {
	if user == nil {
		return nil, apperrors.Unauthorized("user is required")
	}
	if req.Amount < minDonationAmount {
		return nil, apperrors.BadRequest("minimum donation amount is Rp10.000")
	}

	paymentMethod := strings.ToLower(strings.TrimSpace(req.PaymentMethod))
	bank := strings.ToLower(strings.TrimSpace(req.Bank))
	if paymentMethod != paymentMethodQRIS && paymentMethod != paymentMethodVirtualAccount {
		return nil, apperrors.BadRequest("payment method is not supported")
	}
	if paymentMethod == paymentMethodVirtualAccount && !isSupportedVABank(bank) {
		return nil, apperrors.BadRequest("virtual account bank is not supported")
	}

	tx := s.db.Begin()
	defer tx.Rollback()

	request, err := s.requestRepository.GetRequest(tx, model.GetRequestParam{RequestID: req.RequestID})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NotFound("request not found")
		}
		return nil, err
	}
	if request.RequestStatus != "approved" {
		return nil, apperrors.BadRequest("request is not available for donation")
	}

	wallet, err := s.walletRepository.GetWallet(tx, model.GetWalletParam{UserID: user.UserID})
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		wallet = &entity.Wallets{
			WalletID:        uuid.New(),
			UserID:          user.UserID,
			Balance:         0,
			LockedBalance:   0,
			ReservedBalance: 0,
		}
		if err := s.walletRepository.CreateWallet(tx, wallet); err != nil {
			return nil, err
		}
	}

	orderID := buildMidtransOrderID()
	donationID := uuid.New()
	walletTransactionID := uuid.New()
	paymentTransactionID := uuid.New()
	amount := float64(req.Amount)

	walletTransaction := &entity.WalletTransactions{
		WalletTransactionID: walletTransactionID,
		WalletID:            wallet.WalletID,
		Amount:              amount,
		BalanceBefore:       wallet.Balance,
		BalanceAfter:        wallet.Balance,
		TransactionType:     walletTransactionTypeDeposit,
		TransactionStatus:   donationStatusPending,
	}
	err = s.walletTransactionRepository.CreateWalletTransaction(tx, walletTransaction)
	if err != nil {
		return nil, err
	}

	donation := &entity.Donations{
		DonationID:          donationID,
		RequestID:           req.RequestID,
		DonatedBy:           user.UserID,
		WalletTransactionID: walletTransactionID,
		DonationAmount:      amount,
		DonationStatus:      donationStatusPending,
	}
	err = s.donationRepository.CreateDonation(tx, donation)
	if err != nil {
		return nil, err
	}

	payment := &entity.PaymentTransactions{
		PaymentTransactionID: paymentTransactionID,
		OrderID:              orderID,
		UserID:               user.UserID,
		RequestID:            req.RequestID,
		DonationID:           donationID,
		WalletTransactionID:  walletTransactionID,
		Amount:               amount,
		PaymentMethod:        paymentMethod,
		PaymentChannel:       bank,
		TransactionStatus:    paymentStatusPending,
	}
	err = s.paymentTransactionRepository.CreatePaymentTransaction(tx, payment)
	if err != nil {
		return nil, err
	}

	chargeReq := buildMidtransChargeRequest(orderID, req.Amount, paymentMethod, bank, user)
	chargeResp, midtransErr := s.coreClient.ChargeTransaction(chargeReq)
	if midtransErr != nil {
		return nil, apperrors.InternalServer(midtransErr.Message)
	}

	applyChargeResponseToPayment(payment, chargeResp)
	if err := s.paymentTransactionRepository.UpdatePaymentTransaction(tx, payment); err != nil {
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return buildCreateDonationPaymentResponse(payment), nil
}

func (s *DonationPaymentService) HandleMidtransNotification(req model.MidtransNotificationRequest) (*model.DonationLockStatusResponse, error) {
	if strings.TrimSpace(req.OrderID) == "" {
		return nil, apperrors.BadRequest("order_id is required")
	}
	if !s.isValidMidtransSignature(req) {
		return nil, apperrors.Forbidden("invalid midtrans signature")
	}

	tx := s.db.Begin()
	defer tx.Rollback()

	payment, err := s.paymentTransactionRepository.GetPaymentTransaction(tx, model.GetPaymentTransactionParam{OrderID: req.OrderID})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NotFound("payment transaction not found")
		}
		return nil, err
	}

	rawPayload, _ := json.Marshal(req)
	payment.RawNotificationPayload = stringPtr(string(rawPayload))
	payment.TransactionStatus = strings.ToLower(strings.TrimSpace(req.TransactionStatus))
	payment.FraudStatus = strings.ToLower(strings.TrimSpace(req.FraudStatus))
	payment.MidtransStatusCode = req.StatusCode
	payment.PaymentMethod = req.PaymentType
	if req.Bank != "" {
		payment.PaymentChannel = req.Bank
	}
	if req.PermataVANumber != "" {
		payment.PermataVANumber = req.PermataVANumber
	}
	if len(req.VANumbers) > 0 {
		payment.VABank = req.VANumbers[0].Bank
		payment.VANumber = req.VANumbers[0].VANumber
	}

	if isFailedMidtransNotification(req.TransactionStatus) {
		return s.rejectPayment(tx, payment)
	}

	if !isSuccessfulMidtransNotification(req) {
		if err := s.paymentTransactionRepository.UpdatePaymentTransaction(tx, payment); err != nil {
			return nil, err
		}
		if err := tx.Commit().Error; err != nil {
			return nil, err
		}
		return &model.DonationLockStatusResponse{
			PaymentOrderID: payment.OrderID,
			DonationID:     payment.DonationID,
			Amount:         int64(math.Round(payment.Amount)),
			FundStatus:     "PENDING",
			Processed:      false,
		}, nil
	}

	if payment.PaidAt != nil {
		if err := s.paymentTransactionRepository.UpdatePaymentTransaction(tx, payment); err != nil {
			return nil, err
		}
		if err := tx.Commit().Error; err != nil {
			return nil, err
		}
		return &model.DonationLockStatusResponse{
			PaymentOrderID: payment.OrderID,
			DonationID:     payment.DonationID,
			Amount:         int64(math.Round(payment.Amount)),
			FundStatus:     "LOCKED",
			Processed:      false,
		}, nil
	}

	result, err := s.lockDonation(tx, payment)
	if err != nil {
		return nil, err
	}

	err = tx.Commit().Error
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *DonationPaymentService) lockDonation(tx *gorm.DB, payment *entity.PaymentTransactions) (*model.DonationLockStatusResponse, error) {
	now := time.Now().UTC()
	payment.PaidAt = &now

	donation, err := s.donationRepository.GetDonation(tx, model.GetDonationParam{DonationID: payment.DonationID})
	if err != nil {
		return nil, err
	}

	walletTransaction, err := s.walletTransactionRepository.GetWalletTransaction(tx, model.GetWalletTransactionParam{
		WalletTransactionID: payment.WalletTransactionID,
	})
	if err != nil {
		return nil, err
	}

	lockContext, err := s.requestRepository.GetDonationLockContext(tx, payment.RequestID)
	if err != nil {
		return nil, err
	}

	items, err := s.itemRepository.GetItemsByRequestID(tx, model.GetItemParam{RequestID: payment.RequestID})
	if err != nil {
		return nil, err
	}

	lockedOrderID := uuid.New()
	allocation := buildLockedAllocation(lockedOrderID, items, payment.Amount)
	if len(allocation.OrderItems) == 0 {
		return nil, apperrors.BadRequest("donation amount cannot be allocated to available items")
	}

	order := &entity.Orders{
		OrderID:         lockedOrderID,
		RequestID:       payment.RequestID,
		StoreID:         uuid.Nil,
		CourierID:       uuid.Nil,
		OrderCode:       buildDonationTransactionCode(payment.DonationID),
		OrderStatus:     "pending",
		TotalAmount:     allocation.TotalAmount,
		BroadcastRadius: 0,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	if err := s.orderRepository.CreateOrder(tx, order); err != nil {
		return nil, err
	}
	if err := s.orderItemRepository.CreateOrderItems(tx, allocation.OrderItems); err != nil {
		return nil, err
	}

	donation.DonationStatus = donationStatusApproved
	walletTransaction.TransactionStatus = donationStatusApproved

	err = s.donationRepository.UpdateDonation(tx, donation)
	if err != nil {
		return nil, err
	}
	err = s.walletTransactionRepository.UpdateWalletTransaction(tx, walletTransaction)
	if err != nil {
		return nil, err
	}
	err = s.requestRepository.IncrementFundedAmount(tx, payment.RequestID, allocation.TotalAmount)
	if err != nil {
		return nil, err
	}

	custodyLog, err := s.createInitialCustodyLog(tx, payment.UserID, lockContext.AdminUserID, lockedOrderID, lockContext.Latitude, lockContext.Longitude, now)
	if err != nil {
		return nil, err
	}

	err = s.paymentTransactionRepository.UpdatePaymentTransaction(tx, payment)
	if err != nil {
		return nil, err
	}
	if s.pointService != nil {
		if err := s.pointService.AwardDonationPaymentPoints(tx, payment); err != nil {
			return nil, err
		}
	}

	return &model.DonationLockStatusResponse{
		PaymentOrderID:  payment.OrderID,
		DonationID:      payment.DonationID,
		LockedOrderID:   lockedOrderID,
		TransactionCode: order.OrderCode,
		PostName:        lockContext.PostName,
		Amount:          int64(math.Round(allocation.TotalAmount)),
		FundStatus:      "LOCKED",
		AllocationText:  allocation.AllocationText,
		LedgerHash:      custodyLog.CurrentHash,
		ShortLedgerHash: shortenLedgerHash(custodyLog.CurrentHash),
		Processed:       true,
	}, nil
}

func (s *DonationPaymentService) rejectPayment(tx *gorm.DB, payment *entity.PaymentTransactions) (*model.DonationLockStatusResponse, error) {
	donation, err := s.donationRepository.GetDonation(tx, model.GetDonationParam{DonationID: payment.DonationID})
	if err != nil {
		return nil, err
	}

	walletTransaction, err := s.walletTransactionRepository.GetWalletTransaction(tx, model.GetWalletTransactionParam{
		WalletTransactionID: payment.WalletTransactionID,
	})
	if err != nil {
		return nil, err
	}

	donation.DonationStatus = donationStatusRejected
	walletTransaction.TransactionStatus = donationStatusRejected

	err = s.donationRepository.UpdateDonation(tx, donation)
	if err != nil {
		return nil, err
	}
	err = s.walletTransactionRepository.UpdateWalletTransaction(tx, walletTransaction)
	if err != nil {
		return nil, err
	}
	err = s.paymentTransactionRepository.UpdatePaymentTransaction(tx, payment)
	if err != nil {
		return nil, err
	}
	err = tx.Commit().Error
	if err != nil {
		return nil, err
	}

	return &model.DonationLockStatusResponse{
		PaymentOrderID: payment.OrderID,
		DonationID:     payment.DonationID,
		Amount:         int64(math.Round(payment.Amount)),
		FundStatus:     "REJECTED",
		Processed:      true,
	}, nil
}

type lockedAllocation struct {
	OrderItems     []entity.OrderItems
	AllocationText string
	TotalAmount    float64
}

func buildLockedAllocation(orderID uuid.UUID, items []entity.Items, amount float64) lockedAllocation {
	remaining := amount
	orderItems := make([]entity.OrderItems, 0)
	labels := make([]string, 0)

	for _, item := range items {
		available := item.QuantityNeeded - item.QuantityFulfilled
		if available <= 0 || item.Price <= 0 || remaining < item.Price {
			continue
		}

		qty := int(math.Floor(remaining / item.Price))
		if qty > available {
			qty = available
		}
		if qty <= 0 {
			continue
		}

		subtotal := float64(qty) * item.Price
		remaining -= subtotal

		orderItems = append(orderItems, entity.OrderItems{
			OrderItemID: uuid.New(),
			OrderID:     orderID,
			ItemID:      item.ItemID,
			Quantity:    qty,
			Unit:        qty,
			UnitPrice:   item.Price,
			Subtotal:    subtotal,
		})
		labels = append(labels, fmt.Sprintf("%s %d", item.Name, qty))
	}

	return lockedAllocation{
		OrderItems:     orderItems,
		AllocationText: strings.Join(labels, " · "),
		TotalAmount:    amount - remaining,
	}
}

func (s *DonationPaymentService) createInitialCustodyLog(tx *gorm.DB, fromActorID uuid.UUID, toActorID uuid.UUID, orderID uuid.UUID, latitude float64, longitude float64, now time.Time) (*entity.CustodyLogs, error) {
	latestLog, err := s.custodyLogRepository.GetLatestCustodyLog(tx)

	sequence := 1
	prevHash := "GENESIS"
	if err == nil && latestLog != nil {
		sequence = latestLog.Sequence + 1
		prevHash = latestLog.CurrentHash
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	log := &entity.CustodyLogs{
		LogsID:      uuid.NewString(),
		OrderID:     orderID.String(),
		FromActorID: fromActorID,
		ToActorID:   toActorID,
		Sequence:    sequence,
		Latitude:    latitude,
		Longitude:   longitude,
		PrevHash:    prevHash,
		CreatedAt:   now,
	}
	log.CurrentHash = buildCustodyHash(*log)

	err = s.custodyLogRepository.CreateCustodyLog(tx, log)
	if err != nil {
		return nil, err
	}

	return log, nil
}

func buildCustodyHash(log entity.CustodyLogs) string {
	payload := fmt.Sprintf(
		"%s|%s|%s|%s|%d|%.8f|%.8f|%s|%s",
		log.LogsID,
		log.OrderID,
		log.FromActorID.String(),
		log.ToActorID.String(),
		log.Sequence,
		log.Latitude,
		log.Longitude,
		log.PrevHash,
		log.CreatedAt.UTC().Format(time.RFC3339Nano),
	)

	sum := sha256.Sum256([]byte(payload))
	return "0x" + hex.EncodeToString(sum[:])
}

func buildMidtransOrderID() string {
	return fmt.Sprintf("ARK-DON-%s", strings.ReplaceAll(uuid.NewString(), "-", ""))
}

func buildDonationTransactionCode(donationID uuid.UUID) string {
	clean := strings.ReplaceAll(strings.ToUpper(donationID.String()), "-", "")
	return "DN-" + clean[:4] + "-" + clean[4:9]
}

func shortenLedgerHash(hash string) string {
	if len(hash) <= 12 {
		return hash
	}
	return hash[:6] + "..." + hash[len(hash)-4:]
}

func buildMidtransChargeRequest(orderID string, amount int64, paymentMethod string, bank string, user *entity.User) *coreapi.ChargeReq {
	req := &coreapi.ChargeReq{
		PaymentType: coreapi.PaymentTypeQris,
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderID,
			GrossAmt: amount,
		},
		CustomerDetails: &midtrans.CustomerDetails{
			FName: user.Name,
			Email: user.Email,
		},
	}

	if paymentMethod == paymentMethodVirtualAccount {
		req.PaymentType = coreapi.PaymentTypeBankTransfer
		req.BankTransfer = &coreapi.BankTransferDetails{
			Bank: midtrans.Bank(bank),
		}
		return req
	}

	req.Qris = &coreapi.QrisDetails{Acquirer: "gopay"}
	return req
}

func applyChargeResponseToPayment(payment *entity.PaymentTransactions, resp *coreapi.ChargeResponse) {
	raw, _ := json.Marshal(resp)
	payment.RawChargeResponse = stringPtr(string(raw))
	payment.TransactionStatus = strings.ToLower(strings.TrimSpace(resp.TransactionStatus))
	payment.FraudStatus = strings.ToLower(strings.TrimSpace(resp.FraudStatus))
	payment.MidtransStatusCode = resp.StatusCode
	payment.QRString = resp.QRString
	payment.PermataVANumber = resp.PermataVaNumber

	if resp.PaymentType != "" {
		payment.PaymentMethod = resp.PaymentType
	}
	if len(resp.VaNumbers) > 0 {
		payment.VABank = resp.VaNumbers[0].Bank
		payment.VANumber = resp.VaNumbers[0].VANumber
	}
	for _, action := range resp.Actions {
		if action.Name == "generate-qr-code" || action.Name == "generate-qr-code-v2" {
			payment.QRURL = action.URL
			break
		}
	}
	if resp.ExpiryTime != "" {
		if expiredAt, err := time.ParseInLocation("2006-01-02 15:04:05", resp.ExpiryTime, time.Local); err == nil {
			utc := expiredAt.UTC()
			payment.ExpiredAt = &utc
		}
	}
}

func stringPtr(value string) *string {
	return &value
}

func buildCreateDonationPaymentResponse(payment *entity.PaymentTransactions) *model.CreateDonationPaymentResponse {
	return &model.CreateDonationPaymentResponse{
		OrderID:              payment.OrderID,
		DonationID:           payment.DonationID,
		PaymentTransactionID: payment.PaymentTransactionID,
		RequestID:            payment.RequestID,
		Amount:               int64(math.Round(payment.Amount)),
		PaymentMethod:        payment.PaymentMethod,
		PaymentChannel:       payment.PaymentChannel,
		TransactionStatus:    payment.TransactionStatus,
		QRString:             payment.QRString,
		QRURL:                payment.QRURL,
		VANumber:             payment.VANumber,
		VABank:               payment.VABank,
		PermataVANumber:      payment.PermataVANumber,
		ExpiredAt:            payment.ExpiredAt,
	}
}

func (s *DonationPaymentService) isValidMidtransSignature(req model.MidtransNotificationRequest) bool {
	input := req.OrderID + req.StatusCode + req.GrossAmount + s.midtransConfig.ServerKey
	sum := sha512.Sum512([]byte(input))
	return strings.EqualFold(hex.EncodeToString(sum[:]), req.SignatureKey)
}

func isSuccessfulMidtransNotification(req model.MidtransNotificationRequest) bool {
	status := strings.ToLower(strings.TrimSpace(req.TransactionStatus))
	fraud := strings.ToLower(strings.TrimSpace(req.FraudStatus))

	if req.StatusCode != "200" {
		return false
	}
	if status != paymentStatusSettlement && status != paymentStatusCapture {
		return false
	}
	return fraud == "" || fraud == "accept"
}

func isFailedMidtransNotification(status string) bool {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case paymentStatusExpire, paymentStatusCancel, paymentStatusDeny, paymentStatusFailure:
		return true
	default:
		return false
	}
}

func isSupportedVABank(bank string) bool {
	switch strings.ToLower(strings.TrimSpace(bank)) {
	case "bca", "bni", "bri", "mandiri", "permata":
		return true
	default:
		return false
	}
}
