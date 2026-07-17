package service

import (
	"errors"
	"sort"
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

const assumedAverageSpeedKmh = 25

type ICourierTaskService interface {
	GetCourierTasks(user *entity.User, param model.CourierTaskListParam) (*model.CourierTaskListResponse, error)
	GetCourierTaskDetail(user *entity.User, orderID uuid.UUID) (*model.CourierTaskDetailResponse, error)
	ClaimTask(user *entity.User, orderID uuid.UUID, req model.CourierTaskClaimRequest) (*model.CourierTaskActionResponse, error)
	UpdateLocation(user *entity.User, orderID uuid.UUID, req model.CourierLocationPingRequest) (*model.CourierLocationPingResponse, error)
	MarkArrived(user *entity.User, orderID uuid.UUID) (*model.CourierArrivedResponse, error)
}

type CourierTaskService struct {
	db              *gorm.DB
	orderRepository repository.IOrderRepository
	storeRepository repository.IStoreRepository
}

func NewCourierTaskService(orderRepository repository.IOrderRepository, storeRepository repository.IStoreRepository) ICourierTaskService {
	return &CourierTaskService{
		db:              mariadb.Connection,
		orderRepository: orderRepository,
		storeRepository: storeRepository,
	}
}

func (s *CourierTaskService) GetCourierTasks(user *entity.User, param model.CourierTaskListParam) (*model.CourierTaskListResponse, error) {
	if user == nil {
		return nil, apperrors.Unauthorized("user is required")
	}

	status := normalizeCourierTaskStatus(param.Status)
	limit := normalizeStoreOrderServiceLimit(param.Limit)
	offset := normalizeStoreOrderServiceOffset(param.Offset)

	rows, err := s.orderRepository.GetCourierTasks(s.db, model.CourierTaskListRepositoryParam{
		CourierID: user.UserID,
		Status:    status,
	})
	if err != nil {
		return nil, err
	}

	items := buildCourierTaskListItems(rows, param)
	items = paginateCourierTaskItems(items, limit, offset)

	return &model.CourierTaskListResponse{
		Items:  items,
		Limit:  limit,
		Offset: offset,
	}, nil
}

func (s *CourierTaskService) GetCourierTaskDetail(user *entity.User, orderID uuid.UUID) (*model.CourierTaskDetailResponse, error) {
	if user == nil {
		return nil, apperrors.Unauthorized("user is required")
	}

	row, err := s.orderRepository.GetCourierTaskDetail(s.db, model.CourierTaskDetailRepositoryParam{
		OrderID:   orderID,
		CourierID: user.UserID,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NotFound("task not found")
		}
		return nil, err
	}

	return buildCourierTaskDetailResponse(*row), nil
}

func (s *CourierTaskService) ClaimTask(user *entity.User, orderID uuid.UUID, req model.CourierTaskClaimRequest) (*model.CourierTaskActionResponse, error) {
	if user == nil {
		return nil, apperrors.Unauthorized("user is required")
	}

	tx := s.db.Begin()
	defer tx.Rollback()

	err := s.orderRepository.ClaimOrderForCourier(tx, orderID, user.UserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.Conflict("task is no longer available")
		}
		return nil, err
	}

	order, err := s.orderRepository.GetOrder(tx, orderID)
	if err != nil {
		return nil, err
	}

	if req.HasCoords {
		store, err := s.storeRepository.GetStore(tx, model.GetStoreParam{StoreID: order.StoreID})
		if err != nil {
			return nil, err
		}

		now := time.Now().UTC()
		distanceKm := calculateDistanceMeters(req.Latitude, req.Longitude, store.Latitude, store.Longitude) / 1000
		travelMinutes := estimateTravelMinutes(distanceKm)
		deadline := now.Add(2 * time.Duration(travelMinutes*float64(time.Minute)))

		lat := req.Latitude
		lng := req.Longitude
		order.CourierLatitude = &lat
		order.CourierLongitude = &lng
		order.CourierLocationUpdatedAt = &now
		order.PickupDeadlineAt = &deadline
		order.UpdatedAt = now
		if err := s.orderRepository.UpdateOrder(tx, order); err != nil {
			return nil, err
		}
	}

	err = tx.Commit().Error
	if err != nil {
		return nil, err
	}

	return &model.CourierTaskActionResponse{
		OrderID:     order.OrderID,
		CourierID:   order.CourierID,
		OrderStatus: order.OrderStatus,
		UpdatedAt:   order.UpdatedAt,
	}, nil
}

func (s *CourierTaskService) UpdateLocation(user *entity.User, orderID uuid.UUID, req model.CourierLocationPingRequest) (*model.CourierLocationPingResponse, error) {
	if user == nil {
		return nil, apperrors.Unauthorized("user is required")
	}

	tx := s.db.Begin()
	defer tx.Rollback()

	capturedAt := time.Now().UTC()
	if req.CapturedAt != nil && !req.CapturedAt.IsZero() {
		capturedAt = req.CapturedAt.UTC()
	}

	err := s.orderRepository.UpdateCourierLocation(tx, orderID, user.UserID, req.Latitude, req.Longitude, capturedAt)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.Conflict("task is not active for this courier")
		}
		return nil, err
	}

	order, err := s.orderRepository.GetOrder(tx, orderID)
	if err != nil {
		return nil, err
	}

	store, err := s.storeRepository.GetStore(tx, model.GetStoreParam{StoreID: order.StoreID})
	if err != nil {
		return nil, err
	}

	err = tx.Commit().Error
	if err != nil {
		return nil, err
	}

	distanceKm := calculateDistanceMeters(req.Latitude, req.Longitude, store.Latitude, store.Longitude) / 1000
	etaMinutes := estimateTravelMinutes(distanceKm)

	return &model.CourierLocationPingResponse{
		OrderID:                  orderID,
		CourierLatitude:          req.Latitude,
		CourierLongitude:         req.Longitude,
		CourierLocationUpdatedAt: capturedAt,
		PickupDistanceKm:         distanceKm,
		EtaMinutes:               etaMinutes,
	}, nil
}

func (s *CourierTaskService) MarkArrived(user *entity.User, orderID uuid.UUID) (*model.CourierArrivedResponse, error) {
	if user == nil {
		return nil, apperrors.Unauthorized("user is required")
	}

	tx := s.db.Begin()
	defer tx.Rollback()

	now := time.Now().UTC()
	err := s.orderRepository.MarkCourierArrived(tx, orderID, user.UserID, now)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.Conflict("task is not active for this courier")
		}
		return nil, err
	}

	order, err := s.orderRepository.GetOrder(tx, orderID)
	if err != nil {
		return nil, err
	}

	err = tx.Commit().Error
	if err != nil {
		return nil, err
	}

	return &model.CourierArrivedResponse{
		OrderID:     order.OrderID,
		OrderStatus: order.OrderStatus,
		ArrivedAt:   now,
	}, nil
}

func estimateTravelMinutes(distanceKm float64) float64 {
	return distanceKm / assumedAverageSpeedKmh * 60
}

func normalizeCourierTaskStatus(status string) string {
	if strings.ToLower(strings.TrimSpace(status)) == "mine" {
		return "mine"
	}
	return "available"
}

func buildCourierTaskListItems(rows []model.CourierTaskRow, param model.CourierTaskListParam) []model.CourierTaskListItem {
	items := make([]model.CourierTaskListItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, buildCourierTaskListItem(row, param))
	}

	if param.HasCoords {
		sortCourierTaskItemsByDistance(items)
		if param.RadiusKm > 0 {
			items = filterCourierTaskItemsByRadius(items, param.RadiusKm)
		}
	}

	return items
}

func buildCourierTaskListItem(row model.CourierTaskRow, param model.CourierTaskListParam) model.CourierTaskListItem {
	item := model.CourierTaskListItem{
		OrderID:       row.OrderID,
		OrderCode:     row.OrderCode,
		OrderStatus:   row.OrderStatus,
		TotalAmount:   row.TotalAmount,
		RequestTitle:  row.RequestTitle,
		EventName:     row.EventName,
		StoreName:     row.StoreName,
		StoreAddress:  row.StoreAddress,
		PostName:      row.PostName,
		PostAddress:   row.PostAddress,
		ItemCount:     row.ItemCount,
		TotalQuantity: row.TotalQuantity,
		UpdatedAt:     row.UpdatedAt,
	}

	if param.HasCoords {
		pickupKm := calculateDistanceMeters(param.Latitude, param.Longitude, row.StoreLatitude, row.StoreLongitude) / 1000
		dropoffKm := calculateDistanceMeters(row.StoreLatitude, row.StoreLongitude, row.PostLatitude, row.PostLongitude) / 1000
		totalKm := pickupKm + dropoffKm
		item.PickupDistanceKm = &pickupKm
		item.DropoffDistanceKm = &dropoffKm
		item.TotalDistanceKm = &totalKm
	}

	return item
}

func sortCourierTaskItemsByDistance(items []model.CourierTaskListItem) {
	sort.Slice(items, func(i, j int) bool {
		return *items[i].TotalDistanceKm < *items[j].TotalDistanceKm
	})
}

func filterCourierTaskItemsByRadius(items []model.CourierTaskListItem, radiusKm float64) []model.CourierTaskListItem {
	filtered := make([]model.CourierTaskListItem, 0, len(items))
	for _, item := range items {
		if *item.PickupDistanceKm <= radiusKm {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

func paginateCourierTaskItems(items []model.CourierTaskListItem, limit int, offset int) []model.CourierTaskListItem {
	if offset >= len(items) {
		return []model.CourierTaskListItem{}
	}
	end := offset + limit
	if end > len(items) {
		end = len(items)
	}
	return items[offset:end]
}

func buildCourierTaskDetailResponse(row model.CourierTaskRow) *model.CourierTaskDetailResponse {
	var etaMinutes *float64
	if row.CourierLatitude != nil && row.CourierLongitude != nil {
		distanceKm := calculateDistanceMeters(*row.CourierLatitude, *row.CourierLongitude, row.StoreLatitude, row.StoreLongitude) / 1000
		minutes := estimateTravelMinutes(distanceKm)
		etaMinutes = &minutes
	}

	return &model.CourierTaskDetailResponse{
		CourierTaskListItem:      buildCourierTaskListItem(row, model.CourierTaskListParam{}),
		RequestID:                row.RequestID,
		StoreID:                  row.StoreID,
		CourierID:                row.CourierID,
		CourierName:              row.CourierName,
		StoreLatitude:            row.StoreLatitude,
		StoreLongitude:           row.StoreLongitude,
		StorePhoneNumber:         row.StorePhoneNumber,
		PostLatitude:             row.PostLatitude,
		PostLongitude:            row.PostLongitude,
		PostContactName:          row.PostContactName,
		CourierLatitude:          row.CourierLatitude,
		CourierLongitude:         row.CourierLongitude,
		CourierLocationUpdatedAt: row.CourierLocationUpdatedAt,
		ArrivedAt:                row.ArrivedAt,
		PickupDeadlineAt:         row.PickupDeadlineAt,
		EtaMinutes:               etaMinutes,
		AcceptedAt:               row.AcceptedAt,
		ReadyAt:                  row.ReadyAt,
		PickedUpAt:               row.PickedUpAt,
		CreatedAt:                row.CreatedAt,
	}
}
