package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/azmiagr/garudahacks-hackathon/entity"
	"github.com/azmiagr/garudahacks-hackathon/internal/repository"
	"github.com/azmiagr/garudahacks-hackathon/model"
	"github.com/azmiagr/garudahacks-hackathon/pkg/database/mariadb"
	apperrors "github.com/azmiagr/garudahacks-hackathon/pkg/errors"
	"gorm.io/gorm"
)

type IDonorProfileService interface {
	GetProfile(user *entity.User) (*model.DonorProfileResponse, error)
}

type DonorProfileService struct {
	db                     *gorm.DB
	roleRepository         repository.IRoleRepository
	donorProfileRepository repository.IDonorProfileRepository
}

func NewDonorProfileService(roleRepository repository.IRoleRepository, donorProfileRepository repository.IDonorProfileRepository) IDonorProfileService {
	return &DonorProfileService{
		db:                     mariadb.Connection,
		roleRepository:         roleRepository,
		donorProfileRepository: donorProfileRepository,
	}
}

func (s *DonorProfileService) GetProfile(user *entity.User) (*model.DonorProfileResponse, error) {
	if user == nil {
		return nil, apperrors.Unauthorized("user is required")
	}

	role, err := s.roleRepository.GetRole(s.db, model.GetRoleParam{RoleID: user.RoleID})
	if err != nil {
		return nil, err
	}

	donorProfile, err := s.donorProfileRepository.GetDonorProfile(s.db, model.GetDonorProfileParam{UserID: user.UserID})
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	metrics, err := s.donorProfileRepository.GetDonorProfileMetrics(s.db, user.UserID.String())
	if err != nil {
		return nil, err
	}

	phoneNumber := ""
	memberSince := user.CreatedAt
	if donorProfile != nil {
		phoneNumber = donorProfile.PhoneNumber
		if !donorProfile.CreatedAt.IsZero() {
			memberSince = donorProfile.CreatedAt
		}
	}

	isVerified := user.KYCStatus == "approved"

	return &model.DonorProfileResponse{
		UserID:                          user.UserID,
		Name:                            user.Name,
		Initials:                        buildInitials(user.Name),
		Email:                           user.Email,
		PhoneNumber:                     phoneNumber,
		Role:                            role.RoleName,
		DisplayRole:                     resolveDisplayRole(role.RoleName),
		KYCStatus:                       user.KYCStatus,
		IsVerified:                      isVerified,
		VerificationText:                buildDonorProfileVerificationText(isVerified),
		MemberSince:                     memberSince,
		MemberSinceText:                 buildDonorMemberSinceText(memberSince),
		Level:                           resolveDonorProfileLevel(metrics.ActivePoints),
		TotalDonatedAmount:              metrics.TotalDonatedAmount,
		TotalDonatedAmountText:          formatRupiahShortProfile(metrics.TotalDonatedAmount),
		UndistributedDonationAmount:     metrics.UndistributedDonationAmount,
		UndistributedDonationAmountText: formatRupiahShortProfile(metrics.UndistributedDonationAmount),
		SupportedPostCount:              metrics.SupportedPostCount,
		ActivePoints:                    metrics.ActivePoints,
	}, nil
}

func buildDonorProfileVerificationText(isVerified bool) string {
	if isVerified {
		return "TERVERIFIKASI"
	}

	return "MENUNGGU VERIFIKASI"
}

func buildDonorMemberSinceText(memberSince time.Time) string {
	if memberSince.IsZero() {
		return ""
	}

	return fmt.Sprintf("Donatur sejak %s", memberSince.Format("Jan 2006"))
}

func resolveDonorProfileLevel(activePoints int64) string {
	switch {
	case activePoints >= 5000:
		return "Legenda Kebaikan"
	case activePoints >= 2500:
		return "Penerang Nusantara"
	case activePoints >= 1000:
		return "Sahabat Posko"
	default:
		return "Relawan Kebaikan"
	}
}
