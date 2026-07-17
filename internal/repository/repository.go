package repository

import (
	"gorm.io/gorm"
)

type Repository struct {
	UserRepository                    IUserRepository
	RoleRepository                    IRoleRepository
	RegistrationRepository            IRegistrationRepository
	AdminPoskoProfileRepository       IAdminPoskoProfileRepository
	DonorProfileRepository            IDonorProfileRepository
	CourierProfileRepository          ICourierProfileRepository
	AdminDashboardRepository          IAdminDashboardRepository
	PostRepository                    IPostRepository
	DisasterReportRepository          IDisasterReportRepository
	DisasterEventRepository           IDisasterEventRepository
	RequestRepository                 IRequestRepository
	ItemRepository                    IItemRepository
	DeliveryVerificationRepository    IDeliveryVerificationRepository
	DistributionProofRepository       IDistributionProofRepository
	DonationRepository                IDonationRepository
	DisbursementRepository            IDisbursementRepository
	RequestSupplementalNeedRepository IRequestSupplementalNeedRepository
	CustodyLogRepository              ICustodyLogRepository
	OtpRepository                     IOtpRepository
	WalletRepository                  IWalletRepository
	WalletTransactionRepository       IWalletTransactionRepository
	PaymentTransactionRepository      IPaymentTransactionRepository
	OrderRepository                   IOrderRepository
	OrderItemRepository               IOrderItemRepository
	PointRepository                   IPointRepository
	StoreRepository                   IStoreRepository
	CustodyHandshakeTokenRepository   ICustodyHandshakeTokenRepository
	RevokedTokenRepository            IRevokedTokenRepository
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		UserRepository:                    NewUserRepository(db),
		RoleRepository:                    NewRoleRepository(db),
		RegistrationRepository:            NewRegistrationRepository(db),
		AdminPoskoProfileRepository:       NewAdminPoskoProfileRepository(db),
		DonorProfileRepository:            NewDonorProfileRepository(db),
		CourierProfileRepository:          NewCourierProfileRepository(db),
		AdminDashboardRepository:          NewAdminDashboardRepository(db),
		PostRepository:                    NewPostRepository(db),
		DisasterReportRepository:          NewDisasterReportRepository(db),
		DisasterEventRepository:           NewDisasterEventRepository(db),
		RequestRepository:                 NewRequestRepository(db),
		ItemRepository:                    NewItemRepository(db),
		DeliveryVerificationRepository:    NewDeliveryVerificationRepository(db),
		DistributionProofRepository:       NewDistributionProofRepository(db),
		DonationRepository:                NewDonationRepository(db),
		DisbursementRepository:            NewDisbursementRepository(db),
		RequestSupplementalNeedRepository: NewRequestSupplementalNeedRepository(db),
		CustodyLogRepository:              NewCustodyLogRepository(db),
		OtpRepository:                     NewOtpRepository(db),
		WalletRepository:                  NewWalletRepository(db),
		WalletTransactionRepository:       NewWalletTransactionRepository(db),
		PaymentTransactionRepository:      NewPaymentTransactionRepository(db),
		OrderRepository:                   NewOrderRepository(db),
		OrderItemRepository:               NewOrderItemRepository(db),
		PointRepository:                   NewPointRepository(db),
		StoreRepository:                   NewStoreRepository(db),
		CustodyHandshakeTokenRepository:   NewCustodyHandshakeTokenRepository(db),
		RevokedTokenRepository:            NewRevokedTokenRepository(db),
	}
}
