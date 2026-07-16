package repository

import (
	"gorm.io/gorm"
)

type Repository struct {
	UserRepository                 IUserRepository
	PostRepository                 IPostRepository
	DisasterReportRepository       IDisasterReportRepository
	DisasterEventRepository        IDisasterEventRepository
	RequestRepository              IRequestRepository
	DeliveryVerificationRepository IDeliveryVerificationRepository
	DonationRepository             IDonationRepository
	DisbursementRepository         IDisbursementRepository
	CustodyLogRepository           ICustodyLogRepository
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{
		UserRepository:                 NewUserRepository(db),
		PostRepository:                 NewPostRepository(db),
		DisasterReportRepository:       NewDisasterReportRepository(db),
		DisasterEventRepository:        NewDisasterEventRepository(db),
		RequestRepository:              NewRequestRepository(db),
		DeliveryVerificationRepository: NewDeliveryVerificationRepository(db),
		DonationRepository:             NewDonationRepository(db),
		DisbursementRepository:         NewDisbursementRepository(db),
		CustodyLogRepository:           NewCustodyLogRepository(db),
	}
}
