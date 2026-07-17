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
	"github.com/azmiagr/garudahacks-hackathon/pkg/supabase"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const maxDistributionProofPhotoSize = 5 * 1024 * 1024

type IAdminCustodyService interface {
	GetReceiveOrderDetail(user *entity.User, orderID uuid.UUID) (*model.AdminReceiveOrderDetailResponse, error)
	SubmitPostHandoff(user *entity.User, req model.SubmitCustodyHandshakeRequest) (*model.CustodyHandshakeResponse, error)
	CreateSupplementalNeed(user *entity.User, orderID uuid.UUID, req model.CreateSupplementalNeedRequest) (*model.CreateSupplementalNeedResponse, error)
	UploadDistributionProof(user *entity.User, req model.UploadDistributionProofRequest) (*model.UploadDistributionProofResponse, error)
	CompleteDistribution(user *entity.User, orderID uuid.UUID, req model.CompleteDistributionRequest) (*model.CompleteDistributionResponse, error)
}

type AdminCustodyService struct {
	db                             *gorm.DB
	orderRepository                repository.IOrderRepository
	orderItemRepository            repository.IOrderItemRepository
	requestRepository              repository.IRequestRepository
	itemRepository                 repository.IItemRepository
	userRepository                 repository.IUserRepository
	custodyLogRepository           repository.ICustodyLogRepository
	handshakeTokenRepository       repository.ICustodyHandshakeTokenRepository
	deliveryVerificationRepository repository.IDeliveryVerificationRepository
	distributionProofRepository    repository.IDistributionProofRepository
	supplementalNeedRepository     repository.IRequestSupplementalNeedRepository
	pointService                   IPointService
	storage                        supabase.Interface
}

func NewAdminCustodyService(
	orderRepository repository.IOrderRepository,
	orderItemRepository repository.IOrderItemRepository,
	requestRepository repository.IRequestRepository,
	itemRepository repository.IItemRepository,
	userRepository repository.IUserRepository,
	custodyLogRepository repository.ICustodyLogRepository,
	handshakeTokenRepository repository.ICustodyHandshakeTokenRepository,
	deliveryVerificationRepository repository.IDeliveryVerificationRepository,
	distributionProofRepository repository.IDistributionProofRepository,
	supplementalNeedRepository repository.IRequestSupplementalNeedRepository,
	pointService IPointService,
	storage supabase.Interface,
) IAdminCustodyService {
	return &AdminCustodyService{
		db:                             mariadb.Connection,
		orderRepository:                orderRepository,
		orderItemRepository:            orderItemRepository,
		requestRepository:              requestRepository,
		itemRepository:                 itemRepository,
		userRepository:                 userRepository,
		custodyLogRepository:           custodyLogRepository,
		handshakeTokenRepository:       handshakeTokenRepository,
		deliveryVerificationRepository: deliveryVerificationRepository,
		distributionProofRepository:    distributionProofRepository,
		supplementalNeedRepository:     supplementalNeedRepository,
		pointService:                   pointService,
		storage:                        storage,
	}
}

func (s *AdminCustodyService) GetReceiveOrderDetail(user *entity.User, orderID uuid.UUID) (*model.AdminReceiveOrderDetailResponse, error) {
	if user == nil {
		return nil, apperrors.Unauthorized("user is required")
	}

	order, lockContext, err := s.getAdminOwnedOrderContext(s.db, user, orderID)
	if err != nil {
		return nil, err
	}

	return s.buildReceiveOrderDetail(s.db, order, lockContext)
}

func (s *AdminCustodyService) SubmitPostHandoff(user *entity.User, req model.SubmitCustodyHandshakeRequest) (*model.CustodyHandshakeResponse, error) {
	if user == nil {
		return nil, apperrors.Unauthorized("user is required")
	}

	method := strings.ToLower(strings.TrimSpace(req.Method))
	if method != entity.HandshakeMethodQR && method != entity.HandshakeMethodPIN {
		return nil, apperrors.BadRequest("handshake method must be qr or pin")
	}

	tx := s.db.Begin()
	defer tx.Rollback()

	if strings.TrimSpace(req.IdempotencyKey) != "" {
		exists, err := s.custodyLogRepository.ExistsIdempotencyKey(tx, req.IdempotencyKey)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, apperrors.Conflict("handshake idempotency key has already been used")
		}
	}

	token, err := getHandshakeTokenForStageUpdate(tx, s.handshakeTokenRepository, req, method, entity.CustodyStageCourierToPost)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.BadRequest("handshake token is invalid or expired")
		}
		return nil, err
	}
	if token.HandoffStage != entity.CustodyStageCourierToPost {
		return nil, apperrors.BadRequest("handshake token stage is not supported")
	}

	order, err := s.orderRepository.GetOrderForUpdate(tx, token.OrderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NotFound("order not found")
		}
		return nil, err
	}
	if order.OrderStatus != entity.OrderStatusInTransit {
		return nil, apperrors.BadRequest("order is not in transit")
	}
	if order.CourierID != token.PresentedBy {
		return nil, apperrors.BadRequest("handshake token does not match order courier")
	}

	lockContext, err := s.requestRepository.GetDonationLockContext(tx, order.RequestID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.BadRequest("request is not linked to a valid disaster report and post")
		}
		return nil, err
	}
	if lockContext.AdminUserID != user.UserID {
		return nil, apperrors.Forbidden("user is not the registered admin for this post")
	}

	capturedAt := req.CapturedAt.UTC()
	if capturedAt.IsZero() {
		capturedAt = time.Now().UTC()
	}

	log, err := appendCustodyLog(tx, s.custodyLogRepository, appendCustodyLogParam{
		OrderID:           order.OrderID,
		HandoffStage:      entity.CustodyStageCourierToPost,
		HandshakeMethod:   method,
		FromActorID:       order.CourierID,
		ToActorID:         user.UserID,
		ScannedBy:         user.UserID,
		Latitude:          req.Latitude,
		Longitude:         req.Longitude,
		IdempotencyKey:    strings.TrimSpace(req.IdempotencyKey),
		CapturedAt:        capturedAt,
		GPSDistanceMeters: calculateDistanceMeters(req.Latitude, req.Longitude, lockContext.Latitude, lockContext.Longitude),
	})
	if err != nil {
		return nil, err
	}

	err = s.handshakeTokenRepository.MarkTokenUsed(tx, token.TokenID, user.UserID, capturedAt)
	if err != nil {
		return nil, err
	}

	order.OrderStatus = entity.OrderStatusDelivered
	order.DeliveredAt = &capturedAt
	order.UpdatedAt = time.Now().UTC()
	if err := s.orderRepository.UpdateOrder(tx, order); err != nil {
		return nil, err
	}

	verification := &entity.DeliveryVerification{
		VerificationID:     uuid.New(),
		OrderID:            order.OrderID,
		SubmittedBy:        order.CourierID,
		VerifiedBy:         user.UserID,
		VerificationStatus: "pending",
		Latitude:           req.Latitude,
		Longitude:          req.Longitude,
		CapturedAt:         capturedAt,
		ReviewedAt:         time.Now().UTC(),
	}
	err = s.deliveryVerificationRepository.CreateDeliveryVerification(tx, verification)
	if err != nil {
		return nil, err
	}

	pointsAwarded, err := s.pointService.AwardCourierDeliveryPoints(tx, order)
	if err != nil {
		return nil, err
	}

	stats, err := s.orderRepository.GetCourierGoodnessStats(tx, order.CourierID)
	if err != nil {
		return nil, err
	}

	err = tx.Commit().Error
	if err != nil {
		return nil, err
	}

	deliveryCount := stats.DeliveryCount
	totalDistanceKm := stats.TotalDistanceKm

	return &model.CustodyHandshakeResponse{
		OrderID:          order.OrderID,
		LogID:            log.LogsID,
		OrderStatus:      order.OrderStatus,
		HandoffStage:     log.HandoffStage,
		HandshakeMethod:  log.HandshakeMethod,
		Sequence:         log.Sequence,
		CurrentHash:      log.CurrentHash,
		ShortCurrentHash: shortenLedgerHash(log.CurrentHash),
		CapturedAt:       capturedAt,
		PointsAwarded:    &pointsAwarded,
		DeliveryCount:    &deliveryCount,
		TotalDistanceKm:  &totalDistanceKm,
	}, nil
}

func (s *AdminCustodyService) CreateSupplementalNeed(user *entity.User, orderID uuid.UUID, req model.CreateSupplementalNeedRequest) (*model.CreateSupplementalNeedResponse, error) {
	if user == nil {
		return nil, apperrors.Unauthorized("user is required")
	}
	if strings.TrimSpace(req.Reason) == "" {
		return nil, apperrors.BadRequest("reason is required")
	}
	if len(req.Items) == 0 {
		return nil, apperrors.BadRequest("at least one supplemental item is required")
	}

	tx := s.db.Begin()
	defer tx.Rollback()

	order, _, err := s.getAdminOwnedOrderContext(tx, user, orderID)
	if err != nil {
		return nil, err
	}
	if order.OrderStatus == entity.OrderStatusCompleted || order.OrderStatus == entity.OrderStatusCancelled {
		return nil, apperrors.BadRequest("supplemental needs cannot be added for this order status")
	}

	request, err := s.requestRepository.GetRequest(tx, model.GetRequestParam{
		RequestID: order.RequestID,
	})
	if err != nil {
		return nil, err
	}

	itemResponse, additionalTarget, err := s.applySupplementalItems(tx, order.RequestID, req.Items)
	if err != nil {
		return nil, err
	}

	reservedApplied := normalizeReservedAmount(req.ReservedAmountApplied, request.ReservedAmount, additionalTarget)
	netAdditionalTarget := math.Round((additionalTarget-reservedApplied)*100) / 100
	if netAdditionalTarget > 0 {
		err = s.requestRepository.IncrementFundingTarget(tx, order.RequestID, netAdditionalTarget)
		if err != nil {
			return nil, err
		}
		request.FundingTarget += netAdditionalTarget
	}

	supplemental := &entity.RequestSupplementalNeed{
		SupplementalID:        uuid.New(),
		RequestID:             order.RequestID,
		OrderID:               order.OrderID,
		CreatedBy:             user.UserID,
		Reason:                strings.TrimSpace(req.Reason),
		ReservedAmountApplied: reservedApplied,
		AdditionalTarget:      additionalTarget,
	}
	err = s.supplementalNeedRepository.CreateSupplementalNeed(tx, supplemental)
	if err != nil {
		return nil, err
	}

	err = tx.Commit().Error
	if err != nil {
		return nil, err
	}

	return &model.CreateSupplementalNeedResponse{
		SupplementalID:        supplemental.SupplementalID,
		RequestID:             order.RequestID,
		OrderID:               order.OrderID,
		Reason:                supplemental.Reason,
		ReservedAmountApplied: reservedApplied,
		AdditionalTarget:      additionalTarget,
		NewFundingTarget:      request.FundingTarget,
		Items:                 itemResponse,
	}, nil
}

func (s *AdminCustodyService) UploadDistributionProof(user *entity.User, req model.UploadDistributionProofRequest) (*model.UploadDistributionProofResponse, error) {
	if user == nil {
		return nil, apperrors.Unauthorized("user is required")
	}
	if req.Photo == nil {
		return nil, apperrors.BadRequest("photo is required")
	}
	if req.Photo.Size > maxDistributionProofPhotoSize {
		return nil, apperrors.BadRequest("photo size must not exceed 5MB")
	}
	if !req.CapturedFromCamera {
		return nil, apperrors.BadRequest("photo must be captured from live camera")
	}

	imageURL, err := s.storage.UploadFile(req.Photo)
	if err != nil {
		return nil, err
	}

	committed := false
	defer func() {
		if !committed && imageURL != "" {
			_ = supabase.DeleteFileIfPresent(s.storage, imageURL)
		}
	}()

	tx := s.db.Begin()
	defer tx.Rollback()

	order, lockContext, err := s.getAdminOwnedOrderContext(tx, user, req.OrderID)
	if err != nil {
		return nil, err
	}
	if order.OrderStatus != entity.OrderStatusDelivered {
		return nil, apperrors.BadRequest("order must be received before uploading distribution proof")
	}

	item, err := s.itemRepository.GetItem(tx, model.GetItemParam{
		ItemID:    req.ItemID,
		RequestID: order.RequestID,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.BadRequest("item does not belong to this order request")
		}
		return nil, err
	}

	if !s.orderHasItem(tx, order.OrderID, item.ItemID) {
		return nil, apperrors.BadRequest("item is not part of this order")
	}

	capturedAt := req.CapturedAt.UTC()
	if capturedAt.IsZero() {
		capturedAt = time.Now().UTC()
	}

	proof := &entity.DistributionProof{
		ProofID:             uuid.New(),
		OrderID:             order.OrderID,
		ItemID:              item.ItemID,
		SubmittedBy:         user.UserID,
		ImageURL:            imageURL,
		RecipientNote:       strings.TrimSpace(req.RecipientNote),
		DistributedQuantity: req.DistributedQuantity,
		Latitude:            req.Latitude,
		Longitude:           req.Longitude,
		GPSDistanceMeters:   calculateDistanceMeters(req.Latitude, req.Longitude, lockContext.Latitude, lockContext.Longitude),
		BlurFaceEnabled:     req.BlurFaceEnabled,
		CapturedFromCamera:  true,
		CapturedAt:          capturedAt,
	}

	err = s.distributionProofRepository.CreateDistributionProof(tx, proof)
	if err != nil {
		return nil, err
	}

	requiredCount, uploadedCount, err := s.getDistributionProofProgress(tx, order.OrderID)
	if err != nil {
		return nil, err
	}

	err = tx.Commit().Error
	if err != nil {
		return nil, err
	}
	committed = true

	return &model.UploadDistributionProofResponse{
		Proof:              buildDistributionProofData(*proof),
		RequiredPhotoCount: int(requiredCount),
		UploadedPhotoCount: int(uploadedCount),
		ReadyToComplete:    requiredCount > 0 && uploadedCount >= requiredCount,
	}, nil
}

func (s *AdminCustodyService) CompleteDistribution(user *entity.User, orderID uuid.UUID, req model.CompleteDistributionRequest) (*model.CompleteDistributionResponse, error) {
	if user == nil {
		return nil, apperrors.Unauthorized("user is required")
	}

	tx := s.db.Begin()
	defer tx.Rollback()

	if strings.TrimSpace(req.IdempotencyKey) != "" {
		exists, err := s.custodyLogRepository.ExistsIdempotencyKey(tx, req.IdempotencyKey)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, apperrors.Conflict("completion idempotency key has already been used")
		}
	}

	order, lockContext, err := s.getAdminOwnedOrderContextForUpdate(tx, user, orderID)
	if err != nil {
		return nil, err
	}
	if order.OrderStatus != entity.OrderStatusDelivered {
		return nil, apperrors.BadRequest("order must be received and delivered before completion")
	}

	requiredCount, uploadedCount, err := s.getDistributionProofProgress(tx, order.OrderID)
	if err != nil {
		return nil, err
	}
	if requiredCount == 0 {
		return nil, apperrors.BadRequest("order has no items to verify")
	}
	if uploadedCount < requiredCount {
		return nil, apperrors.BadRequest("all item distribution photos must be uploaded before completion")
	}

	proofs, err := s.distributionProofRepository.ListDistributionProofsByOrder(tx, order.OrderID)
	if err != nil {
		return nil, err
	}
	if len(proofs) == 0 {
		return nil, apperrors.BadRequest("distribution proof is required")
	}

	capturedAt := req.CapturedAt.UTC()
	if capturedAt.IsZero() {
		capturedAt = time.Now().UTC()
	}

	log, err := appendCustodyLog(tx, s.custodyLogRepository, appendCustodyLogParam{
		OrderID:           order.OrderID,
		HandoffStage:      entity.CustodyStageDistributionCompleted,
		HandshakeMethod:   entity.HandshakeMethodSystem,
		FromActorID:       user.UserID,
		ToActorID:         user.UserID,
		ScannedBy:         user.UserID,
		Latitude:          req.Latitude,
		Longitude:         req.Longitude,
		IdempotencyKey:    strings.TrimSpace(req.IdempotencyKey),
		CapturedAt:        capturedAt,
		GPSDistanceMeters: calculateDistanceMeters(req.Latitude, req.Longitude, lockContext.Latitude, lockContext.Longitude),
	})
	if err != nil {
		return nil, err
	}

	order.OrderStatus = entity.OrderStatusCompleted
	order.CompletedAt = &capturedAt
	order.UpdatedAt = time.Now().UTC()
	if err := s.orderRepository.UpdateOrder(tx, order); err != nil {
		return nil, err
	}

	verification := &entity.DeliveryVerification{
		VerificationID:     uuid.New(),
		OrderID:            order.OrderID,
		SubmittedBy:        user.UserID,
		VerifiedBy:         user.UserID,
		VerificationStatus: "approved",
		ImageURL:           proofs[0].ImageURL,
		Latitude:           req.Latitude,
		Longitude:          req.Longitude,
		CapturedAt:         capturedAt,
		ReviewedAt:         time.Now().UTC(),
	}
	err = s.deliveryVerificationRepository.CreateDeliveryVerification(tx, verification)
	if err != nil {
		return nil, err
	}

	err = tx.Commit().Error
	if err != nil {
		return nil, err
	}

	return &model.CompleteDistributionResponse{
		OrderID:            order.OrderID,
		OrderStatus:        order.OrderStatus,
		RequiredPhotoCount: int(requiredCount),
		UploadedPhotoCount: int(uploadedCount),
		FinalHash:          log.CurrentHash,
		ShortFinalHash:     shortenLedgerHash(log.CurrentHash),
		CompletedAt:        capturedAt,
	}, nil
}

func (s *AdminCustodyService) getAdminOwnedOrderContext(tx *gorm.DB, user *entity.User, orderID uuid.UUID) (*entity.Orders, *model.DonationLockContextRow, error) {
	order, err := s.orderRepository.GetOrder(tx, orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, apperrors.NotFound("order not found")
		}
		return nil, nil, err
	}

	lockContext, err := s.getAdminOwnedLockContext(tx, user, order.RequestID)
	if err != nil {
		return nil, nil, err
	}

	return order, lockContext, nil
}

func (s *AdminCustodyService) getAdminOwnedOrderContextForUpdate(tx *gorm.DB, user *entity.User, orderID uuid.UUID) (*entity.Orders, *model.DonationLockContextRow, error) {
	order, err := s.orderRepository.GetOrderForUpdate(tx, orderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, apperrors.NotFound("order not found")
		}
		return nil, nil, err
	}

	lockContext, err := s.getAdminOwnedLockContext(tx, user, order.RequestID)
	if err != nil {
		return nil, nil, err
	}

	return order, lockContext, nil
}

func (s *AdminCustodyService) getAdminOwnedLockContext(tx *gorm.DB, user *entity.User, requestID uuid.UUID) (*model.DonationLockContextRow, error) {
	lockContext, err := s.requestRepository.GetDonationLockContext(tx, requestID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.BadRequest("request is not linked to a valid disaster report and post")
		}
		return nil, err
	}
	if lockContext.AdminUserID != user.UserID {
		return nil, apperrors.Forbidden("user is not the registered admin for this post")
	}

	return lockContext, nil
}

func (s *AdminCustodyService) buildReceiveOrderDetail(tx *gorm.DB, order *entity.Orders, lockContext *model.DonationLockContextRow) (*model.AdminReceiveOrderDetailResponse, error) {
	request, err := s.requestRepository.GetRequest(tx, model.GetRequestParam{
		RequestID: order.RequestID,
	})
	if err != nil {
		return nil, err
	}

	items, err := s.orderRepository.GetStoreOrderItems(tx, order.OrderID)
	if err != nil {
		return nil, err
	}

	proofs, err := s.distributionProofRepository.ListDistributionProofsByOrder(tx, order.OrderID)
	if err != nil {
		return nil, err
	}

	proofByItem := make(map[uuid.UUID]entity.DistributionProof, len(proofs))
	proofResponse := make([]model.AdminDistributionProofData, 0, len(proofs))
	for _, proof := range proofs {
		proofByItem[proof.ItemID] = proof
		proofResponse = append(proofResponse, buildDistributionProofData(proof))
	}

	itemResponse := make([]model.AdminReceiveOrderItem, 0, len(items))
	for _, item := range items {
		_, hasProof := proofByItem[item.ItemID]
		itemResponse = append(itemResponse, model.AdminReceiveOrderItem{
			ItemID:    item.ItemID,
			Name:      item.Name,
			Quantity:  item.Quantity,
			Unit:      item.Unit,
			UnitPrice: item.UnitPrice,
			Subtotal:  item.Subtotal,
			HasProof:  hasProof,
		})
	}

	courierName := ""
	if order.CourierID != uuid.Nil {
		courier, err := s.userRepository.GetUser(tx, model.GetUserParam{
			UserID: order.CourierID,
		})
		if err == nil && courier != nil {
			courierName = courier.Name
		}
	}

	return &model.AdminReceiveOrderDetailResponse{
		OrderID:            order.OrderID,
		OrderCode:          order.OrderCode,
		OrderStatus:        order.OrderStatus,
		RequestID:          order.RequestID,
		RequestTitle:       request.Title,
		PostName:           lockContext.PostName,
		CourierID:          order.CourierID,
		CourierName:        courierName,
		DeliveredAt:        order.DeliveredAt,
		CompletedAt:        order.CompletedAt,
		Items:              itemResponse,
		Proofs:             proofResponse,
		RequiredPhotoCount: len(items),
		UploadedPhotoCount: len(proofs),
	}, nil
}

func (s *AdminCustodyService) applySupplementalItems(tx *gorm.DB, requestID uuid.UUID, reqItems []model.CreateSupplementalNeedItem) ([]model.CreateAdminEventItemData, float64, error) {
	responseItems := make([]model.CreateAdminEventItemData, 0, len(reqItems))

	var additionalTarget float64
	for _, reqItem := range reqItems {
		name := strings.TrimSpace(reqItem.Name)
		if reqItem.ItemID == uuid.Nil && name == "" {
			return nil, 0, apperrors.BadRequest("item name is required")
		}
		if reqItem.Price <= 0 {
			return nil, 0, apperrors.BadRequest("item price must be greater than zero")
		}
		if reqItem.QuantityNeeded <= 0 {
			return nil, 0, apperrors.BadRequest("item quantity must be greater than zero")
		}

		estimatedTotal := math.Round(reqItem.Price*float64(reqItem.QuantityNeeded)*100) / 100
		additionalTarget += estimatedTotal

		if reqItem.ItemID != uuid.Nil {
			item, err := s.itemRepository.GetItem(tx, model.GetItemParam{
				ItemID:    reqItem.ItemID,
				RequestID: requestID,
			})
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return nil, 0, apperrors.BadRequest("item does not belong to this request")
				}
				return nil, 0, err
			}

			err = s.itemRepository.IncrementQuantityNeeded(tx, item.ItemID, reqItem.QuantityNeeded, estimatedTotal)
			if err != nil {
				return nil, 0, err
			}

			item.QuantityNeeded += reqItem.QuantityNeeded
			item.EstimatedTotal += estimatedTotal
			responseItems = append(responseItems, model.CreateAdminEventItemData{
				ItemID:         item.ItemID,
				Name:           item.Name,
				Description:    item.Description,
				Price:          item.Price,
				EstimatedTotal: item.EstimatedTotal,
				QuantityNeeded: item.QuantityNeeded,
			})
			continue
		}

		item := entity.Items{
			ItemID:            uuid.New(),
			RequestID:         requestID,
			Name:              name,
			Description:       strings.TrimSpace(reqItem.Description),
			Price:             reqItem.Price,
			EstimatedTotal:    estimatedTotal,
			QuantityNeeded:    reqItem.QuantityNeeded,
			QuantityFulfilled: 0,
		}
		err := s.itemRepository.CreateItem(tx, &item)
		if err != nil {
			return nil, 0, err
		}

		responseItems = append(responseItems, model.CreateAdminEventItemData{
			ItemID:         item.ItemID,
			Name:           item.Name,
			Description:    item.Description,
			Price:          item.Price,
			EstimatedTotal: item.EstimatedTotal,
			QuantityNeeded: item.QuantityNeeded,
		})
	}

	additionalTarget = math.Round(additionalTarget*100) / 100
	return responseItems, additionalTarget, nil
}

func (s *AdminCustodyService) orderHasItem(tx *gorm.DB, orderID uuid.UUID, itemID uuid.UUID) bool {
	var count int64
	err := tx.Model(&entity.OrderItems{}).
		Where("order_id = ? AND item_id = ?", orderID, itemID).
		Count(&count).Error
	return err == nil && count > 0
}

func (s *AdminCustodyService) getDistributionProofProgress(tx *gorm.DB, orderID uuid.UUID) (int64, int64, error) {
	requiredCount, err := s.orderItemRepository.CountDistinctItemsByOrder(tx, orderID)
	if err != nil {
		return 0, 0, err
	}

	uploadedCount, err := s.distributionProofRepository.CountDistributionProofsByOrder(tx, orderID)
	if err != nil {
		return 0, 0, err
	}

	return requiredCount, uploadedCount, nil
}

func normalizeReservedAmount(requested float64, available float64, additionalTarget float64) float64 {
	if requested <= 0 || available <= 0 || additionalTarget <= 0 {
		return 0
	}

	if requested > available {
		requested = available
	}
	if requested > additionalTarget {
		requested = additionalTarget
	}

	return math.Round(requested*100) / 100
}

func buildDistributionProofData(proof entity.DistributionProof) model.AdminDistributionProofData {
	return model.AdminDistributionProofData{
		ProofID:             proof.ProofID,
		OrderID:             proof.OrderID,
		ItemID:              proof.ItemID,
		ImageURL:            proof.ImageURL,
		RecipientNote:       proof.RecipientNote,
		DistributedQuantity: proof.DistributedQuantity,
		Latitude:            proof.Latitude,
		Longitude:           proof.Longitude,
		GPSDistanceMeters:   proof.GPSDistanceMeters,
		BlurFaceEnabled:     proof.BlurFaceEnabled,
		CapturedFromCamera:  proof.CapturedFromCamera,
		CapturedAt:          proof.CapturedAt,
	}
}
