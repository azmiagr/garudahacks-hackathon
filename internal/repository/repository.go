package repository

import (
	"gorm.io/gorm"
)

type Repository struct {
	UserRepository                 IUserRepository
	RoleRepository                 IRoleRepository
	RegistrationRepository         IRegistrationRepository
	AdminPoskoProfileRepository    IAdminPoskoProfileRepository
	DonorProfileRepository         IDonorProfileRepository
	AdminDashboardRepository       IAdminDashboardRepository
	PostRepository                 IPostRepository
	DisasterReportRepository       IDisasterReportRepository
	DisasterEventRepository        IDisasterEventRepository
	RequestRepository              IRequestRepository
	ItemRepository                 IItemRepository
	DeliveryVerificationRepository IDeliveryVerificationRepository
	DonationRepository             IDonationRepository
	DisbursementRepository         IDisbursementRepository
	CustodyLogRepository           ICustodyLogRepository
	OtpRepository                  IOtpRepository
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		UserRepository:                 NewUserRepository(db),
		RoleRepository:                 NewRoleRepository(db),
		RegistrationRepository:         NewRegistrationRepository(db),
		AdminPoskoProfileRepository:    NewAdminPoskoProfileRepository(db),
		DonorProfileRepository:         NewDonorProfileRepository(db),
		AdminDashboardRepository:       NewAdminDashboardRepository(db),
		PostRepository:                 NewPostRepository(db),
		DisasterReportRepository:       NewDisasterReportRepository(db),
		DisasterEventRepository:        NewDisasterEventRepository(db),
		RequestRepository:              NewRequestRepository(db),
		ItemRepository:                 NewItemRepository(db),
		DeliveryVerificationRepository: NewDeliveryVerificationRepository(db),
		DonationRepository:             NewDonationRepository(db),
		DisbursementRepository:         NewDisbursementRepository(db),
		CustodyLogRepository:           NewCustodyLogRepository(db),
		OtpRepository:                  NewOtpRepository(db),
	}
}
