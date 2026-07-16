package service

import (
	"github.com/azmiagr/garudahacks-hackathon/internal/repository"
	"github.com/azmiagr/garudahacks-hackathon/pkg/bcrypt"
	"github.com/azmiagr/garudahacks-hackathon/pkg/config"
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
	AdminProfileService    IAdminProfileService
	DonorDashboardService  IDonorDashboardService
	DonationPaymentService IDonationPaymentService
}

func NewService(repository *repository.Repository, bcrypt bcrypt.Interface, jwtAuth jwt.Interface, supabase supabase.Interface, midtransConfig *config.MidtransConfig) *Service {
	publicDashboardService := NewPublicDashboardService(repository.PostRepository, repository.DisasterReportRepository, repository.DisasterEventRepository, repository.RequestRepository, repository.DeliveryVerificationRepository, repository.DonationRepository, repository.DisbursementRepository, repository.CustodyLogRepository)

	return &Service{
		UserService:            NewUserService(repository.UserRepository, repository.RoleRepository),
		PublicDashboardService: publicDashboardService,
		OtpService:             NewOtpService(repository.OtpRepository, repository.UserRepository),
		AuthService:            NewAuthService(repository.UserRepository, repository.RoleRepository, repository.RegistrationRepository, repository.AdminPoskoProfileRepository, repository.DonorProfileRepository, bcrypt, jwtAuth),
		AdminDashboardService:  NewAdminDashboardService(repository.AdminDashboardRepository),
		AdminEventService:      NewAdminEventService(repository.PostRepository, repository.DisasterReportRepository, repository.DisasterEventRepository, repository.RequestRepository, repository.ItemRepository, supabase),
		AdminProfileService:    NewAdminProfileService(repository.RoleRepository, repository.AdminPoskoProfileRepository),
		DonorDashboardService:  NewDonorDashboardService(repository.PostRepository, repository.ItemRepository, publicDashboardService),
		DonationPaymentService: NewDonationPaymentService(
			repository.RequestRepository,
			repository.ItemRepository,
			repository.WalletRepository,
			repository.WalletTransactionRepository,
			repository.DonationRepository,
			repository.PaymentTransactionRepository,
			repository.OrderRepository,
			repository.OrderItemRepository,
			repository.CustodyLogRepository,
			midtransConfig,
		),
	}
}
