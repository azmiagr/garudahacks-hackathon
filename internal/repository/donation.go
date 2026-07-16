package repository

import (
	"github.com/azmiagr/garudahacks-hackathon/model"
	"gorm.io/gorm"
)

type IDonationRepository interface {
	GetTransparencySummary(tx *gorm.DB, year int) (*model.DonationTransparencySummaryRow, error)
}

type DonationRepository struct {
	db *gorm.DB
}

func NewDonationRepository(db *gorm.DB) IDonationRepository {
	return &DonationRepository{db: db}
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

	err := query.Scan(&row).Error
	if err != nil {
		return nil, err
	}

	return &row, nil
}
