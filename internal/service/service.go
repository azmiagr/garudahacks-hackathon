package service

import (
	"github.com/azmiagr/garudahacks-hackathon/internal/repository"
	"github.com/azmiagr/garudahacks-hackathon/pkg/bcrypt"
	"github.com/azmiagr/garudahacks-hackathon/pkg/config"
	"github.com/azmiagr/garudahacks-hackathon/pkg/hash"
	"github.com/azmiagr/garudahacks-hackathon/pkg/jwt"
	"github.com/azmiagr/garudahacks-hackathon/pkg/supabase"
)

type Service struct {
	UserService              IUserService
	PublicDashboardService   IPublicDashboardService
	OtpService               IOtpService
	AuthService              IAuthService
	AdminDashboardService    IAdminDashboardService
	AdminEventService        IAdminEventService
	AdminProfileService      IAdminProfileService
	DonorProfileService      IDonorProfileService
	DonorDashboardService    IDonorDashboardService
	DonationPaymentService   IDonationPaymentService
	DonorTransactionService  IDonorTransactionService
	PointService             IPointService
	StoreCustodyService      IStoreCustodyService
	StoreDisbursementService IStoreDisbursementService
	StoreGoodnessService     IStoreGoodnessService
	StoreProfileService      IStoreProfileService
	CourierTaskService       ICourierTaskService
	CourierGoodnessService   ICourierGoodnessService
	AdminCustodyService      IAdminCustodyService
}

func NewService(repository *repository.Repository, bcrypt bcrypt.Interface, jwtAuth jwt.Interface, supabase supabase.Interface, midtransConfig *config.MidtransConfig, hasher hash.Interface) *Service {
	publicDashboardService := NewPublicDashboardService(repository.PostRepository, repository.DisasterReportRepository, repository.DisasterEventRepository, repository.RequestRepository, repository.DeliveryVerificationRepository, repository.DonationRepository, repository.DisbursementRepository, repository.CustodyLogRepository)
	pointService := NewPointService(repository.PointRepository)

	return &Service{
		UserService:              NewUserService(repository.UserRepository, repository.RoleRepository),
		PublicDashboardService:   publicDashboardService,
		OtpService:               NewOtpService(repository.OtpRepository, repository.UserRepository),
		AuthService:              NewAuthService(repository.UserRepository, repository.RoleRepository, repository.RegistrationRepository, repository.AdminPoskoProfileRepository, repository.DonorProfileRepository, repository.RevokedTokenRepository, bcrypt, jwtAuth, hasher, supabase, repository.StoreRepository),
		AdminDashboardService:    NewAdminDashboardService(repository.AdminDashboardRepository),
		AdminEventService:        NewAdminEventService(repository.PostRepository, repository.DisasterReportRepository, repository.DisasterEventRepository, repository.RequestRepository, repository.ItemRepository, supabase),
		AdminProfileService:      NewAdminProfileService(repository.RoleRepository, repository.AdminPoskoProfileRepository),
		DonorProfileService:      NewDonorProfileService(repository.RoleRepository, repository.DonorProfileRepository),
		DonorDashboardService:    NewDonorDashboardService(repository.PostRepository, repository.ItemRepository, publicDashboardService),
		DonorTransactionService:  NewDonorTransactionService(repository.DonationRepository),
		PointService:             pointService,
		StoreCustodyService:      NewStoreCustodyService(repository.OrderRepository, repository.StoreRepository, repository.CustodyLogRepository, repository.CustodyHandshakeTokenRepository, repository.RequestRepository),
		DonationPaymentService:   NewDonationPaymentService(repository.RequestRepository, repository.ItemRepository, repository.WalletRepository, repository.WalletTransactionRepository, repository.DonationRepository, repository.PaymentTransactionRepository, repository.OrderRepository, repository.OrderItemRepository, repository.CustodyLogRepository, pointService, midtransConfig),
		StoreDisbursementService: NewStoreDisbursementService(repository.StoreRepository, repository.DisbursementRepository),
		StoreGoodnessService:     NewStoreGoodnessService(repository.StoreRepository),
		StoreProfileService:      NewStoreProfileService(repository.StoreRepository),
		CourierTaskService:       NewCourierTaskService(repository.OrderRepository, repository.StoreRepository, repository.CustodyLogRepository, repository.CustodyHandshakeTokenRepository),
		CourierGoodnessService:   NewCourierGoodnessService(repository.OrderRepository),
		AdminCustodyService:      NewAdminCustodyService(repository.OrderRepository, repository.RequestRepository, repository.CustodyLogRepository, repository.CustodyHandshakeTokenRepository, repository.DeliveryVerificationRepository, pointService),
	}
}
