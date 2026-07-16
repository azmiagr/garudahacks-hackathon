package service

import (
	"github.com/azmiagr/garudahacks-hackathon/internal/repository"
	"github.com/azmiagr/garudahacks-hackathon/pkg/bcrypt"
	"github.com/azmiagr/garudahacks-hackathon/pkg/jwt"
)

type Service struct {
	UserService            IUserService
	PublicDashboardService IPublicDashboardService
	OtpService             IOtpService
	AuthService            IAuthService
	AdminDashboardService  IAdminDashboardService
}

func NewService(repository *repository.Repository, bcrypt bcrypt.Interface, jwtAuth jwt.Interface) *Service {
	return &Service{
		UserService:            NewUserService(repository.UserRepository, repository.RoleRepository),
		PublicDashboardService: NewPublicDashboardService(repository.PostRepository, repository.DisasterReportRepository, repository.DisasterEventRepository, repository.RequestRepository, repository.DeliveryVerificationRepository, repository.DonationRepository, repository.DisbursementRepository, repository.CustodyLogRepository),
		OtpService:             NewOtpService(repository.OtpRepository, repository.UserRepository),
		AuthService:            NewAuthService(repository.UserRepository, repository.RoleRepository, repository.RegistrationRepository, repository.AdminPoskoProfileRepository, bcrypt, jwtAuth),
		AdminDashboardService:  NewAdminDashboardService(repository.AdminDashboardRepository),
	}
}
