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

type IAdminProfileService interface {
	GetProfile(user *entity.User) (*model.AdminProfileResponse, error)
}

type AdminProfileService struct {
	db                     *gorm.DB
	roleRepository         repository.IRoleRepository
	adminProfileRepository repository.IAdminPoskoProfileRepository
}

func NewAdminProfileService(roleRepository repository.IRoleRepository, adminProfileRepository repository.IAdminPoskoProfileRepository) IAdminProfileService {
	return &AdminProfileService{
		db:                     mariadb.Connection,
		roleRepository:         roleRepository,
		adminProfileRepository: adminProfileRepository,
	}
}

func (s *AdminProfileService) GetProfile(user *entity.User) (*model.AdminProfileResponse, error) {
	if user == nil {
		return nil, apperrors.Unauthorized("user is required")
	}

	role, err := s.roleRepository.GetRole(s.db, model.GetRoleParam{
		RoleID: user.RoleID,
	})
	if err != nil {
		return nil, err
	}

	adminProfile, err := s.adminProfileRepository.GetAdminPoskoProfile(s.db, model.GetAdminPoskoProfileParam{
		UserID: user.UserID,
	})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	metrics, err := s.adminProfileRepository.GetAdminProfileMetrics(s.db, user.UserID)
	if err != nil {
		return nil, err
	}

	affiliation := ""
	if adminProfile != nil {
		affiliation = adminProfile.Affiliation
	}

	verifiedPercentage := calculateVerifiedOrderPercentage(metrics.VerifiedOrderCount, metrics.TotalOrderCount)
	isVerified := user.KYCStatus == "approved"

	return &model.AdminProfileResponse{
		UserID:                  user.UserID,
		Name:                    user.Name,
		Initials:                buildInitials(user.Name),
		Role:                    role.RoleName,
		DisplayRole:             resolveDisplayRole(role.RoleName),
		Affiliation:             affiliation,
		KYCStatus:               user.KYCStatus,
		IsVerified:              isVerified,
		VerificationText:        buildAdminVerificationText(isVerified),
		SuccessfulEventsText:    fmt.Sprintf("%d EVENTS SUKSES", metrics.EventCount),
		EventCount:              metrics.EventCount,
		ManagedAidAmount:        metrics.ManagedAidAmount,
		ManagedAidAmountText:    formatRupiahShortProfile(metrics.ManagedAidAmount),
		VerifiedOrderPercentage: verifiedPercentage,
		VerifiedOrderText:       fmt.Sprintf("%.0f%%", verifiedPercentage),
	}, nil
}

func calculateVerifiedOrderPercentage(verifiedOrderCount, totalOrderCount int64) float64 {
	if totalOrderCount <= 0 {
		return 0
	}

	percentage := (float64(verifiedOrderCount) / float64(totalOrderCount)) * 100
	return math.Round(percentage*100) / 100
}

func buildAdminVerificationText(isVerified bool) string {
	if isVerified {
		return "TERVERIFIKASI"
	}

	return "MENUNGGU VERIFIKASI"
}

func buildInitials(name string) string {
	parts := strings.Fields(strings.TrimSpace(name))
	if len(parts) == 0 {
		return "A"
	}

	initials := ""
	for _, part := range parts {
		if part == "" {
			continue
		}
		initials += strings.ToUpper(string([]rune(part)[0]))
		if len([]rune(initials)) >= 2 {
			break
		}
	}

	return initials
}

func formatRupiahShortProfile(amount float64) string {
	switch {
	case amount >= 1000000000:
		return fmt.Sprintf("Rp%.1fM", amount/1000000000)
	case amount >= 1000000:
		return fmt.Sprintf("Rp%.0fjt", amount/1000000)
	case amount >= 1000:
		return fmt.Sprintf("Rp%.0frb", amount/1000)
	default:
		return fmt.Sprintf("Rp%.0f", amount)
	}
}
