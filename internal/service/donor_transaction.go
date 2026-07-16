package service

import (
	"errors"
	"math"
	"strings"
	"time"

	"github.com/azmiagr/garudahacks-hackathon/entity"
	"github.com/azmiagr/garudahacks-hackathon/internal/repository"
	"github.com/azmiagr/garudahacks-hackathon/model"
	"github.com/azmiagr/garudahacks-hackathon/pkg/database/mariadb"
	apperrors "github.com/azmiagr/garudahacks-hackathon/pkg/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IDonorTransactionService interface {
	GetTransactions(user *entity.User, param model.DonorDonationTransactionParam) (*model.DonorDonationTransactionListResponse, error)
	GetTransactionDetail(user *entity.User, donationID uuid.UUID) (*model.DonorDonationTransactionDetailResponse, error)
}

type DonorTransactionService struct {
	db                 *gorm.DB
	donationRepository repository.IDonationRepository
}

func NewDonorTransactionService(donationRepository repository.IDonationRepository) IDonorTransactionService {
	return &DonorTransactionService{
		db:                 mariadb.Connection,
		donationRepository: donationRepository,
	}
}

func (s *DonorTransactionService) GetTransactions(user *entity.User, param model.DonorDonationTransactionParam) (*model.DonorDonationTransactionListResponse, error) {
	if user == nil {
		return nil, apperrors.Unauthorized("user is required")
	}

	repoParam := model.DonorDonationTransactionListParam{
		UserID: user.UserID,
		Status: normalizeDonorTransactionStatus(param.Status),
		Limit:  normalizeDonorTransactionLimit(param.Limit),
		Offset: normalizeDonorTransactionOffset(param.Offset),
	}

	rows, err := s.donationRepository.GetDonorDonationTransactions(s.db, repoParam)
	if err != nil {
		return nil, err
	}

	total, err := s.donationRepository.CountDonorDonationTransactions(s.db, repoParam)
	if err != nil {
		return nil, err
	}

	items := make([]model.DonorDonationTransactionListItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, buildDonorDonationTransactionListItem(row))
	}

	return &model.DonorDonationTransactionListResponse{
		Items:  items,
		Total:  total,
		Limit:  repoParam.Limit,
		Offset: repoParam.Offset,
	}, nil
}

func (s *DonorTransactionService) GetTransactionDetail(user *entity.User, donationID uuid.UUID) (*model.DonorDonationTransactionDetailResponse, error) {
	if user == nil {
		return nil, apperrors.Unauthorized("user is required")
	}
	if donationID == uuid.Nil {
		return nil, apperrors.BadRequest("donation_id is required")
	}

	param := model.DonorDonationTransactionDetailParam{
		UserID:     user.UserID,
		DonationID: donationID,
	}

	row, err := s.donationRepository.GetDonorDonationTransactionDetail(s.db, param)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NotFound("donation transaction not found")
		}
		return nil, err
	}

	itemRows, err := s.donationRepository.GetDonorDonationTransactionItems(s.db, param)
	if err != nil {
		return nil, err
	}

	custodyRows, err := s.donationRepository.GetDonorDonationTransactionCustodyLogs(s.db, param)
	if err != nil {
		return nil, err
	}

	base := buildDonorDonationTransactionListItem(row.DonorDonationTransactionListRow)
	percentage := calculateDonorFundingPercentage(row.FundedAmount, row.FundingTarget)

	return &model.DonorDonationTransactionDetailResponse{
		DonorDonationTransactionListItem: base,
		PostAddress:                      row.PostAddress,
		Latitude:                         row.Latitude,
		Longitude:                        row.Longitude,
		FundingTarget:                    row.FundingTarget,
		FundedAmount:                     row.FundedAmount,
		FundingPercentage:                percentage,
		FundingText:                      buildDonorFundingText(row.FundedAmount, row.FundingTarget, percentage),
		DonorCount:                       row.DonorCount,
		TotalItemCount:                   row.TotalItemCount,
		Items:                            buildDonorDonationTransactionItems(itemRows),
		CustodyLogs:                      buildDonorDonationTransactionCustodyLogs(custodyRows),
	}, nil
}

func buildDonorDonationTransactionListItem(row model.DonorDonationTransactionListRow) model.DonorDonationTransactionListItem {
	status := resolveDonorTransactionStatus(row)
	occurredAt := resolveDonorTransactionOccurredAt(row)

	return model.DonorDonationTransactionListItem{
		DonationID:           row.DonationID,
		PaymentTransactionID: row.PaymentTransactionID,
		PaymentOrderID:       row.PaymentOrderID,
		RequestID:            row.RequestID,
		LockedOrderID:        row.LockedOrderID,
		TransactionCode:      row.TransactionCode,
		PostName:             row.PostName,
		RequestTitle:         row.RequestTitle,
		Amount:               int64(math.Round(row.Amount)),
		AmountText:           formatDonorRupiahShort(row.Amount),
		Status:               status,
		StatusLabel:          buildDonorTransactionStatusLabel(status),
		BadgeVariant:         buildDonorTransactionBadgeVariant(status),
		LatestHash:           row.LatestHash,
		ShortLatestHash:      shortenLedgerHash(row.LatestHash),
		VerificationImageURL: row.VerificationImageURL,
		CustodyStepCount:     row.CustodyStepCount,
		ProgressText:         buildDonorTransactionProgressText(status, row.CustodyStepCount),
		DonatedAt:            row.DonatedAt,
		PaidAt:               row.PaidAt,
		VerifiedAt:           row.VerifiedAt,
		ElapsedText:          buildDonorElapsedText(occurredAt),
	}
}

func buildDonorDonationTransactionItems(rows []model.DonorDonationTransactionItemRow) []model.DonorDonationTransactionItem {
	items := make([]model.DonorDonationTransactionItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, model.DonorDonationTransactionItem{
			ItemID:     row.ItemID,
			Name:       row.Name,
			Quantity:   row.Quantity,
			UnitPrice:  row.UnitPrice,
			Subtotal:   row.Subtotal,
			AmountText: formatDonorRupiahShort(row.Subtotal),
		})
	}
	return items
}

func buildDonorDonationTransactionCustodyLogs(rows []model.DonorDonationTransactionCustodyLogRow) []model.DonorDonationTransactionCustodyLogItem {
	logs := make([]model.DonorDonationTransactionCustodyLogItem, 0, len(rows))
	for _, row := range rows {
		logs = append(logs, model.DonorDonationTransactionCustodyLogItem{
			LogsID:           row.LogsID,
			OrderID:          row.OrderID,
			Sequence:         row.Sequence,
			FromActorID:      row.FromActorID,
			ToActorID:        row.ToActorID,
			Latitude:         row.Latitude,
			Longitude:        row.Longitude,
			PrevHash:         row.PrevHash,
			CurrentHash:      row.CurrentHash,
			ShortCurrentHash: shortenLedgerHash(row.CurrentHash),
			CreatedAt:        row.CreatedAt,
			ElapsedText:      buildDonorElapsedText(row.CreatedAt),
		})
	}
	return logs
}

func resolveDonorTransactionStatus(row model.DonorDonationTransactionListRow) string {
	paymentStatus := strings.ToLower(strings.TrimSpace(row.PaymentStatus))
	donationStatus := strings.ToLower(strings.TrimSpace(row.DonationStatus))

	if donationStatus == donationStatusRejected || isFailedMidtransNotification(paymentStatus) {
		return "refund"
	}
	if row.VerifiedAt != nil {
		return "completed"
	}
	if donationStatus == donationStatusApproved && row.LockedOrderID != "" {
		return "locked"
	}
	return "pending"
}

func buildDonorTransactionStatusLabel(status string) string {
	switch status {
	case "completed":
		return "Selesai + Foto"
	case "locked":
		return "Dikirim Kurir"
	case "refund":
		return "Refund Otomatis"
	default:
		return "Menunggu Pembayaran"
	}
}

func buildDonorTransactionBadgeVariant(status string) string {
	switch status {
	case "completed":
		return "success"
	case "locked":
		return "warning"
	case "refund":
		return "muted"
	default:
		return "info"
	}
}

func buildDonorTransactionProgressText(status string, custodyStepCount int) string {
	switch status {
	case "completed":
		return "Bukti foto tersedia"
	case "locked":
		if custodyStepCount <= 0 {
			return "Dana terkunci"
		}
		return "Rantai kustodi berjalan"
	case "refund":
		return "Dana dikembalikan"
	default:
		return "Menunggu pembayaran"
	}
}

func resolveDonorTransactionOccurredAt(row model.DonorDonationTransactionListRow) time.Time {
	if row.VerifiedAt != nil {
		return *row.VerifiedAt
	}
	if row.PaidAt != nil {
		return *row.PaidAt
	}
	return row.DonatedAt
}

func normalizeDonorTransactionStatus(status string) string {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "pending", "locked", "completed", "refund":
		return strings.ToLower(strings.TrimSpace(status))
	default:
		return "all"
	}
}

func normalizeDonorTransactionLimit(limit int) int {
	if limit <= 0 {
		return 20
	}
	if limit > 100 {
		return 100
	}
	return limit
}

func normalizeDonorTransactionOffset(offset int) int {
	if offset < 0 {
		return 0
	}
	return offset
}
