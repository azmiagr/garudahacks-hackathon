package repository

import (
	"github.com/azmiagr/garudahacks-hackathon/entity"
	"github.com/azmiagr/garudahacks-hackathon/model"
	"gorm.io/gorm"
)

type IDonorProfileRepository interface {
	GetDonorProfile(tx *gorm.DB, param model.GetDonorProfileParam) (*entity.DonorProfile, error)
	CreateDonorProfile(tx *gorm.DB, profile *entity.DonorProfile) error
	GetDonorProfileMetrics(tx *gorm.DB, userID string) (*model.DonorProfileMetricsRow, error)
}

type DonorProfileRepository struct {
	db *gorm.DB
}

func NewDonorProfileRepository(db *gorm.DB) IDonorProfileRepository {
	return &DonorProfileRepository{db: db}
}

func (r *DonorProfileRepository) GetDonorProfile(tx *gorm.DB, param model.GetDonorProfileParam) (*entity.DonorProfile, error) {
	var profile entity.DonorProfile
	err := tx.Where(&param).First(&profile).Error
	if err != nil {
		return nil, err
	}

	return &profile, nil
}

func (r *DonorProfileRepository) CreateDonorProfile(tx *gorm.DB, profile *entity.DonorProfile) error {
	err := tx.Create(profile).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *DonorProfileRepository) GetDonorProfileMetrics(tx *gorm.DB, userID string) (*model.DonorProfileMetricsRow, error) {
	var row model.DonorProfileMetricsRow

	err := tx.Table("users AS u").
		Select(`
			COALESCE(SUM(CASE WHEN d.donation_status = 'approved' THEN d.donation_amount ELSE 0 END), 0) AS total_donated_amount,
			COUNT(DISTINCT CASE WHEN d.donation_status = 'approved' THEN d.request_id END) AS supported_post_count,
			COALESCE(pa.active_points, 0) AS active_points
		`).
		Joins("LEFT JOIN donations AS d ON d.donated_by = u.user_id").
		Joins("LEFT JOIN point_accounts AS pa ON pa.user_id = u.user_id").
		Where("u.user_id = ?", userID).
		Group("u.user_id, pa.active_points").
		Scan(&row).Error
	if err != nil {
		return nil, err
	}

	return &row, nil
}
