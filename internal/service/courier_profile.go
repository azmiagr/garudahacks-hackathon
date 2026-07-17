package service

import (
	"errors"
	"fmt"
	"math"
	"strings"

	"github.com/azmiagr/garudahacks-hackathon/entity"
	"github.com/azmiagr/garudahacks-hackathon/internal/repository"
	"github.com/azmiagr/garudahacks-hackathon/model"
	"github.com/azmiagr/garudahacks-hackathon/pkg/database/mariadb"
	apperrors "github.com/azmiagr/garudahacks-hackathon/pkg/errors"
	"gorm.io/gorm"
)

type ICourierProfileService interface {
	GetProfile(user *entity.User) (*model.CourierProfileResponse, error)
	UpdatePreferences(user *entity.User, req model.UpdateCourierProfilePreferencesRequest) (*model.UpdateCourierProfilePreferencesResponse, error)
}

type CourierProfileService struct {
	db                       *gorm.DB
	roleRepository           repository.IRoleRepository
	courierProfileRepository repository.ICourierProfileRepository
	orderRepository          repository.IOrderRepository
	pointRepository          repository.IPointRepository
}

func NewCourierProfileService(
	roleRepository repository.IRoleRepository,
	courierProfileRepository repository.ICourierProfileRepository,
	orderRepository repository.IOrderRepository,
	pointRepository repository.IPointRepository,
) ICourierProfileService {
	return &CourierProfileService{
		db:                       mariadb.Connection,
		roleRepository:           roleRepository,
		courierProfileRepository: courierProfileRepository,
		orderRepository:          orderRepository,
		pointRepository:          pointRepository,
	}
}

func (s *CourierProfileService) GetProfile(user *entity.User) (*model.CourierProfileResponse, error) {
	if user == nil {
		return nil, apperrors.Unauthorized("user is required")
	}

	role, err := s.roleRepository.GetRole(s.db, model.GetRoleParam{RoleID: user.RoleID})
	if err != nil {
		return nil, err
	}

	profile, err := s.courierProfileRepository.GetCourierProfile(s.db, model.GetCourierProfileParam{UserID: user.UserID})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NotFound("courier profile not found")
		}
		return nil, err
	}

	stats, err := s.orderRepository.GetCourierGoodnessStats(s.db, user.UserID)
	if err != nil {
		return nil, err
	}

	points, err := s.pointRepository.GetPointSummary(s.db, user.UserID)
	if err != nil {
		return nil, err
	}

	reputation := math.Round(stats.ReputationScore*10) / 10
	isVerified := user.KYCStatus == "approved"

	return &model.CourierProfileResponse{
		UserID:                        user.UserID,
		ProfileID:                     profile.ProfileID,
		Name:                          user.Name,
		Initials:                      buildInitials(user.Name),
		Email:                         user.Email,
		Role:                          role.RoleName,
		DisplayRole:                   resolveDisplayRole(role.RoleName),
		KYCStatus:                     user.KYCStatus,
		IsVerified:                    isVerified,
		VerificationText:              buildCourierVerificationText(isVerified),
		OperationalArea:               profile.OperationalArea,
		OperationRadiusKM:             profile.OperationRadiusKM,
		OperationAreaText:             buildCourierOperationAreaText(profile.OperationalArea, profile.OperationRadiusKM),
		VehicleType:                   profile.VehicleType,
		VehicleCapacityKG:             profile.VehicleCapacityKG,
		VehicleText:                   buildCourierVehicleText(profile.VehicleType, profile.VehicleCapacityKG),
		WaiverAccepted:                profile.WaiverAccepted,
		IsAvailable:                   profile.IsAvailable,
		UrgentTaskNotificationEnabled: profile.UrgentTaskNotificationEnabled,
		ReputationScore:               reputation,
		ReputationText:                formatReputationText(reputation),
		ActivePoints:                  points.ActivePoints,
		ActivePointsText:              fmt.Sprintf("%d poin aktif", points.ActivePoints),
		DeliveryCount:                 stats.DeliveryCount,
		TotalDistanceKm:               stats.TotalDistanceKm,
	}, nil
}

func (s *CourierProfileService) UpdatePreferences(user *entity.User, req model.UpdateCourierProfilePreferencesRequest) (*model.UpdateCourierProfilePreferencesResponse, error) {
	if user == nil {
		return nil, apperrors.Unauthorized("user is required")
	}

	profile, err := s.courierProfileRepository.GetCourierProfile(s.db, model.GetCourierProfileParam{UserID: user.UserID})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NotFound("courier profile not found")
		}
		return nil, err
	}

	isAvailable := profile.IsAvailable
	if req.IsAvailable != nil {
		isAvailable = *req.IsAvailable
	}

	urgentNotification := profile.UrgentTaskNotificationEnabled
	if req.UrgentTaskNotificationEnabled != nil {
		urgentNotification = *req.UrgentTaskNotificationEnabled
	}

	err = s.courierProfileRepository.UpdateCourierProfilePreferences(s.db, user.UserID.String(), isAvailable, urgentNotification)
	if err != nil {
		return nil, err
	}

	return &model.UpdateCourierProfilePreferencesResponse{
		IsAvailable:                   isAvailable,
		UrgentTaskNotificationEnabled: urgentNotification,
	}, nil
}

func buildCourierVerificationText(isVerified bool) string {
	if isVerified {
		return "NIK TERVERIFIKASI"
	}

	return "MENUNGGU VERIFIKASI"
}

func buildCourierVehicleText(vehicleType string, capacityKG int) string {
	vehicleType = strings.TrimSpace(vehicleType)
	if vehicleType == "" {
		return ""
	}

	if capacityKG <= 0 {
		return titleCase(vehicleType)
	}

	return fmt.Sprintf("%s + box - kapasitas <=%d kg", titleCase(vehicleType), capacityKG)
}

func buildCourierOperationAreaText(area string, radiusKM int) string {
	area = strings.TrimSpace(area)
	if area == "" {
		return ""
	}

	if radiusKM <= 0 {
		return area
	}

	return fmt.Sprintf("%s + radius %d km", area, radiusKM)
}
