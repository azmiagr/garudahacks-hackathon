package repository

import (
	"github.com/azmiagr/garudahacks-hackathon/entity"
	"github.com/azmiagr/garudahacks-hackathon/model"
	"gorm.io/gorm"
)

type IDonationRepository interface {
	GetTransparencySummary(tx *gorm.DB, year int) (*model.DonationTransparencySummaryRow, error)
	CreateDonation(tx *gorm.DB, donation *entity.Donations) error
	GetDonation(tx *gorm.DB, param model.GetDonationParam) (*entity.Donations, error)
	UpdateDonation(tx *gorm.DB, donation *entity.Donations) error
}

type DonationRepository struct {
	db *gorm.DB
}

func NewDonationRepository(db *gorm.DB) IDonationRepository {
	return &DonationRepository{db: db}
}

func (r *DonationRepository) CreateDonation(tx *gorm.DB, donation *entity.Donations) error {
	err := tx.Create(donation).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *DonationRepository) GetDonation(tx *gorm.DB, param model.GetDonationParam) (*entity.Donations, error) {
	var donation entity.Donations
	err := tx.Where(&param).First(&donation).Error
	if err != nil {
		return nil, err
	}

	return &donation, nil
}

func (r *DonationRepository) UpdateDonation(tx *gorm.DB, donation *entity.Donations) error {
	err := tx.Save(donation).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *DonationRepository) GetTransparencySummary(tx *gorm.DB, year int) (*model.DonationTransparencySummaryRow, error) {
	var row model.DonationTransparencySummaryRow

	query := tx.Table("donations").
		Select(`
			COALESCE(SUM(CASE WHEN donation_status = 'approved' THEN donation_amount ELSE 0 END), 0) AS total_collected,
			COALESCE(SUM(CASE WHEN donation_status = 'rejected' THEN donation_amount ELSE 0 END), 0) AS refund_automatic
		`)

	if year > 0 {
		query = query.Where("YEAR(donated_at) = ?", year)
	}

	if err := query.Scan(&row).Error; err != nil {
		return nil, err
	}

	return &row, nil
}
