package service

import (
	"errors"
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

type IAdminCustodyService interface {
	SubmitPostHandoff(user *entity.User, req model.SubmitCustodyHandshakeRequest) (*model.CustodyHandshakeResponse, error)
}

type AdminCustodyService struct {
	db                             *gorm.DB
	orderRepository                repository.IOrderRepository
	requestRepository              repository.IRequestRepository
	custodyLogRepository           repository.ICustodyLogRepository
	handshakeTokenRepository       repository.ICustodyHandshakeTokenRepository
	deliveryVerificationRepository repository.IDeliveryVerificationRepository
	pointService                   IPointService
}

func NewAdminCustodyService(
	orderRepository repository.IOrderRepository,
	requestRepository repository.IRequestRepository,
	custodyLogRepository repository.ICustodyLogRepository,
	handshakeTokenRepository repository.ICustodyHandshakeTokenRepository,
	deliveryVerificationRepository repository.IDeliveryVerificationRepository,
	pointService IPointService,
) IAdminCustodyService {
	return &AdminCustodyService{
		db:                             mariadb.Connection,
		orderRepository:                orderRepository,
		requestRepository:              requestRepository,
		custodyLogRepository:           custodyLogRepository,
		handshakeTokenRepository:       handshakeTokenRepository,
		deliveryVerificationRepository: deliveryVerificationRepository,
		pointService:                   pointService,
	}
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
		VerificationStatus: "approved",
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
