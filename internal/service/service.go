package service

import (
	"github.com/azmiagr/garudahacks-hackathon/internal/repository"
	"github.com/azmiagr/garudahacks-hackathon/pkg/bcrypt"
	"github.com/azmiagr/garudahacks-hackathon/pkg/jwt"
	"github.com/azmiagr/garudahacks-hackathon/pkg/supabase"
)

type Service struct {
	UserService            IUserService
	PublicDashboardService IPublicDashboardService
	OtpService             IOtpService
	AuthService            IAuthService
	AdminDashboardService  IAdminDashboardService
	AdminEventService      IAdminEventService
}

func NewService(repository *repository.Repository, bcrypt bcrypt.Interface, jwtAuth jwt.Interface, supabase supabase.Interface) *Service {
	return &Service{
		UserService:            NewUserService(repository.UserRepository, repository.RoleRepository),
		PublicDashboardService: NewPublicDashboardService(repository.PostRepository, repository.DisasterReportRepository, repository.DisasterEventRepository, repository.RequestRepository, repository.DeliveryVerificationRepository, repository.DonationRepository, repository.DisbursementRepository, repository.CustodyLogRepository),
		OtpService:             NewOtpService(repository.OtpRepository, repository.UserRepository),
		AuthService:            NewAuthService(repository.UserRepository, repository.RoleRepository, repository.RegistrationRepository, repository.AdminPoskoProfileRepository, bcrypt, jwtAuth),
		AdminDashboardService:  NewAdminDashboardService(repository.AdminDashboardRepository),
		AdminEventService:      NewAdminEventService(repository.PostRepository, repository.DisasterReportRepository, repository.DisasterEventRepository, repository.RequestRepository, repository.ItemRepository, supabase),
	}
}
