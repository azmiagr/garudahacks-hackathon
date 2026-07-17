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

	donorDonationsSubquery := tx.Table("donations").
		Select(`
			donated_by AS user_id,
			request_id,
			SUM(donation_amount) AS donor_amount
		`).
		Where("donation_status = ?", "approved").
		Group("donated_by, request_id")

	requestDonationsSubquery := tx.Table("donations").
		Select(`
			request_id,
			SUM(donation_amount) AS request_total_amount
		`).
		Where("donation_status = ?", "approved").
		Group("request_id")

	requestDisbursementsSubquery := tx.Table("orders AS o").
		Select(`
			o.request_id,
			SUM(d.amount) AS disbursed_amount
		`).
		Joins("JOIN disbursements AS d ON d.order_id = o.order_id").
		Where("d.status = ?", "success").
		Group("o.request_id")

	err := tx.Table("users AS u").
		Select(`
			COALESCE(SUM(dd.donor_amount), 0) AS total_donated_amount,
			COALESCE(SUM(
				GREATEST(
					dd.donor_amount - LEAST(
						dd.donor_amount,
						COALESCE((dd.donor_amount / NULLIF(rd.request_total_amount, 0)) * COALESCE(rds.disbursed_amount, 0), 0)
					),
					0
				)
			), 0) AS undistributed_donation_amount,
			COUNT(DISTINCT dd.request_id) AS supported_post_count,
			COALESCE(pa.active_points, 0) AS active_points
		`).
		Joins("LEFT JOIN (?) AS dd ON dd.user_id = u.user_id", donorDonationsSubquery).
		Joins("LEFT JOIN (?) AS rd ON rd.request_id = dd.request_id", requestDonationsSubquery).
		Joins("LEFT JOIN (?) AS rds ON rds.request_id = dd.request_id", requestDisbursementsSubquery).
		Joins("LEFT JOIN point_accounts AS pa ON pa.user_id = u.user_id").
		Where("u.user_id = ?", userID).
		Group("u.user_id, pa.active_points").
		Scan(&row).Error
	if err != nil {
		return nil, err
	}

	return &row, nil
}
