package service

import (
	"fmt"
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

const (
	minPostGeofenceRadius = 100
	maxPostGeofenceRadius = 2000
	maxEventPhotoSize     = 5 * 1024 * 1024
)

type IAdminEventService interface {
	CreateEvent(user *entity.User, req model.CreateAdminEventRequest) (*model.CreateAdminEventResponse, error)
}

type AdminEventService struct {
	db                       *gorm.DB
	postRepository           repository.IPostRepository
	disasterReportRepository repository.IDisasterReportRepository
	disasterEventRepository  repository.IDisasterEventRepository
	requestRepository        repository.IRequestRepository
	itemRepository           repository.IItemRepository
	storage                  supabase.Interface
}

func NewAdminEventService(
	postRepository repository.IPostRepository,
	disasterReportRepository repository.IDisasterReportRepository,
	disasterEventRepository repository.IDisasterEventRepository,
	requestRepository repository.IRequestRepository,
	itemRepository repository.IItemRepository,
	storage supabase.Interface,
) IAdminEventService {
	return &AdminEventService{
		db:                       mariadb.Connection,
		postRepository:           postRepository,
		disasterReportRepository: disasterReportRepository,
		disasterEventRepository:  disasterEventRepository,
		requestRepository:        requestRepository,
		itemRepository:           itemRepository,
		storage:                  storage,
	}
}

func (s *AdminEventService) CreateEvent(user *entity.User, req model.CreateAdminEventRequest) (*model.CreateAdminEventResponse, error) {
	if user == nil {
		return nil, apperrors.Unauthorized("user is required")
	}

	err := validateCreateAdminEventRequest(req)
	if err != nil {
		return nil, err
	}

	imageURL, err := supabase.UploadOptionalImage(
		s.storage,
		req.Photo,
		maxEventPhotoSize,
		"event photo size must not exceed 5MB",
	)
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

	disasterEvent, err := s.disasterEventRepository.GetDisasterEventByName(tx, req.DisasterType)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, apperrors.BadRequest("disaster type is not supported")
		}
		return nil, err
	}

	now := time.Now().UTC()
	postID := uuid.New()
	reportID := uuid.New()
	requestID := uuid.New()

	post := &entity.Post{
		PostID:         postID,
		UserID:         user.UserID,
		Name:           strings.TrimSpace(req.Name),
		Description:    strings.TrimSpace(req.Description),
		Address:        strings.TrimSpace(req.Address),
		Latitude:       req.Latitude,
		Longitude:      req.Longitude,
		GeofenceRadius: req.GeofenceRadius,
	}

	err = s.postRepository.CreatePost(tx, post)
	if err != nil {
		return nil, err
	}

	report := &entity.DisasterReport{
		ReportID:        reportID,
		DisasterEventID: disasterEvent.EventID,
		PostID:          postID,
		UserID:          user.UserID,
		ReportTitle:     strings.TrimSpace(req.Name),
		Description:     strings.TrimSpace(req.Description),
		Latitude:        req.Latitude,
		Longitude:       req.Longitude,
		ImageUrl:        imageURL,
		ReportStatus:    "approved",
		ReportedAt:      now,
	}

	err = s.disasterReportRepository.CreateDisasterReport(tx, report)
	if err != nil {
		return nil, err
	}

	items, itemResponse, fundingTarget := buildCreateEventItems(requestID, req.Items)

	request := &entity.Requests{
		RequestID:      requestID,
		ReportID:       reportID,
		CreatedBy:      user.UserID,
		Title:          strings.TrimSpace(req.Name),
		Description:    strings.TrimSpace(req.Description),
		FundingTarget:  fundingTarget,
		FundedAmount:   0,
		ReservedAmount: 0,
		RequestStatus:  "approved",
	}

	err = s.requestRepository.CreateRequest(tx, request)
	if err != nil {
		return nil, err
	}

	err = s.itemRepository.CreateItems(tx, items)
	if err != nil {
		return nil, err
	}

	err = tx.Commit().Error
	if err != nil {
		return nil, err
	}

	committed = true

	return &model.CreateAdminEventResponse{
		PostID:         postID,
		ReportID:       reportID,
		RequestID:      requestID,
		EventCode:      buildPostCode(postID),
		Name:           post.Name,
		DisasterType:   disasterEvent.Name,
		ImageURL:       imageURL,
		Address:        post.Address,
		Latitude:       post.Latitude,
		Longitude:      post.Longitude,
		GeofenceRadius: post.GeofenceRadius,
		FundingTarget:  fundingTarget,
		ReportStatus:   report.ReportStatus,
		RequestStatus:  request.RequestStatus,
		Items:          itemResponse,
	}, nil
}

func validateCreateAdminEventRequest(req model.CreateAdminEventRequest) error {
	if strings.TrimSpace(req.Name) == "" {
		return apperrors.BadRequest("event name is required")
	}

	if strings.TrimSpace(req.DisasterType) == "" {
		return apperrors.BadRequest("disaster type is required")
	}

	if strings.TrimSpace(req.Address) == "" {
		return apperrors.BadRequest("address is required")
	}

	if req.Photo == nil {
		return apperrors.BadRequest("event photo is required")
	}

	if req.GeofenceRadius < minPostGeofenceRadius || req.GeofenceRadius > maxPostGeofenceRadius {
		return apperrors.BadRequest("geofence radius must be between 100 and 2000 meters")
	}

	if req.Latitude < -90 || req.Latitude > 90 {
		return apperrors.BadRequest("latitude is invalid")
	}

	if req.Longitude < -180 || req.Longitude > 180 {
		return apperrors.BadRequest("longitude is invalid")
	}

	if len(req.Items) == 0 {
		return apperrors.BadRequest("at least one request item is required")
	}

	for _, item := range req.Items {
		if strings.TrimSpace(item.Name) == "" {
			return apperrors.BadRequest("item name is required")
		}

		if item.Price <= 0 {
			return apperrors.BadRequest("item price must be greater than zero")
		}

		if item.QuantityNeeded <= 0 {
			return apperrors.BadRequest("item quantity must be greater than zero")
		}
	}

	return nil
}

func buildCreateEventItems(requestID uuid.UUID, reqItems []model.CreateAdminEventItem) ([]entity.Items, []model.CreateAdminEventItemData, float64) {
	items := make([]entity.Items, 0, len(reqItems))
	responseItems := make([]model.CreateAdminEventItemData, 0, len(reqItems))

	var fundingTarget float64

	for _, reqItem := range reqItems {
		itemID := uuid.New()
		estimatedTotal := reqItem.Price * float64(reqItem.QuantityNeeded)
		estimatedTotal = math.Round(estimatedTotal*100) / 100

		item := entity.Items{
			ItemID:            itemID,
			RequestID:         requestID,
			Name:              strings.TrimSpace(reqItem.Name),
			Description:       strings.TrimSpace(reqItem.Description),
			Price:             reqItem.Price,
			EstimatedTotal:    estimatedTotal,
			QuantityNeeded:    reqItem.QuantityNeeded,
			QuantityFulfilled: 0,
		}

		items = append(items, item)
		responseItems = append(responseItems, model.CreateAdminEventItemData{
			ItemID:         itemID,
			Name:           item.Name,
			Description:    item.Description,
			Price:          item.Price,
			EstimatedTotal: item.EstimatedTotal,
			QuantityNeeded: item.QuantityNeeded,
		})

		fundingTarget += estimatedTotal
	}

	fundingTarget = math.Round(fundingTarget*100) / 100

	return items, responseItems, fundingTarget
}

func buildPostCode(postID uuid.UUID) string {
	value := strings.ReplaceAll(postID.String(), "-", "")
	if len(value) < 4 {
		return "PSK-0000"
	}

	return fmt.Sprintf("PSK-%s", strings.ToUpper(value[:4]))
}
