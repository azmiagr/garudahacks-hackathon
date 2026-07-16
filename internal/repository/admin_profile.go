package repository

import (
	"github.com/azmiagr/garudahacks-hackathon/entity"
	"github.com/azmiagr/garudahacks-hackathon/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IAdminPoskoProfileRepository interface {
	GetAdminPoskoProfile(tx *gorm.DB, param model.GetAdminPoskoProfileParam) (*entity.AdminProfile, error)
	GetAdminProfileMetrics(tx *gorm.DB, userID uuid.UUID) (*model.AdminProfileMetricsRow, error)
	CreateAdminPoskoProfile(tx *gorm.DB, profile *entity.AdminProfile) error
}

type AdminPoskoProfileRepository struct {
	db *gorm.DB
}

func NewAdminPoskoProfileRepository(db *gorm.DB) IAdminPoskoProfileRepository {
	return &AdminPoskoProfileRepository{db: db}
}

func (r *AdminPoskoProfileRepository) GetAdminPoskoProfile(tx *gorm.DB, param model.GetAdminPoskoProfileParam) (*entity.AdminProfile, error) {
	var profile entity.AdminProfile
	err := tx.Where(&param).First(&profile).Error
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (r *AdminPoskoProfileRepository) GetAdminProfileMetrics(tx *gorm.DB, userID uuid.UUID) (*model.AdminProfileMetricsRow, error) {
	var row model.AdminProfileMetricsRow

	err := tx.Table("posts").
		Select("COUNT(post_id) AS event_count").
		Where("user_id = ?", userID).
		Scan(&row).
		Error
	if err != nil {
		return nil, err
	}

	err = tx.Table("requests AS req").
		Select("COALESCE(SUM(req.funded_amount), 0) AS managed_aid_amount").
		Joins("JOIN disaster_reports AS dr ON dr.report_id = req.report_id").
		Joins("JOIN posts AS p ON p.post_id = dr.post_id").
		Where("p.user_id = ?", userID).
		Where("req.request_status IN ?", []string{"pending", "approved"}).
		Scan(&row).
		Error
	if err != nil {
		return nil, err
	}

	err = tx.Table("orders AS o").
		Select(`
			COUNT(DISTINCT o.order_id) AS total_order_count,
			COUNT(DISTINCT CASE
				WHEN dv.verification_status = 'approved' THEN o.order_id
			END) AS verified_order_count
		`).
		Joins("JOIN requests AS req ON req.request_id = o.request_id").
		Joins("JOIN disaster_reports AS dr ON dr.report_id = req.report_id").
		Joins("JOIN posts AS p ON p.post_id = dr.post_id").
		Joins("LEFT JOIN delivery_verifications AS dv ON dv.order_id = o.order_id").
		Where("p.user_id = ?", userID).
		Scan(&row).
		Error
	if err != nil {
		return nil, err
	}

	return &row, nil
}

func (r *AdminPoskoProfileRepository) CreateAdminPoskoProfile(tx *gorm.DB, profile *entity.AdminProfile) error {
	err := tx.Create(profile).Error
	if err != nil {
		return err
	}
	return nil
}
