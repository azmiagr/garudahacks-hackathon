package repository

import (
	"github.com/azmiagr/garudahacks-hackathon/entity"
	"github.com/azmiagr/garudahacks-hackathon/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IRequestRepository interface {
	CreateRequest(tx *gorm.DB, request *entity.Requests) error
	GetRequest(tx *gorm.DB, param model.GetRequestParam) (*entity.Requests, error)
	GetFundingSummaryByReportIDs(tx *gorm.DB, param model.RequestFundingSummaryParam) ([]model.RequestFundingSummaryRow, error)
	GetAllocationByDisaster(tx *gorm.DB, year int) ([]model.DisasterAllocationRow, error)
	IncrementFundedAmount(tx *gorm.DB, requestID uuid.UUID, amount float64) error
	GetDonationLockContext(tx *gorm.DB, requestID uuid.UUID) (*model.DonationLockContextRow, error)
}

type RequestRepository struct {
	db *gorm.DB
}

func NewRequestRepository(db *gorm.DB) IRequestRepository {
	return &RequestRepository{db: db}
}

func (r *RequestRepository) CreateRequest(tx *gorm.DB, request *entity.Requests) error {
	err := tx.Create(request).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *RequestRepository) GetRequest(tx *gorm.DB, param model.GetRequestParam) (*entity.Requests, error) {
	var request entity.Requests
	err := tx.Where(&param).First(&request).Error
	if err != nil {
		return nil, err
	}

	return &request, nil
}

func (r *RequestRepository) GetFundingSummaryByReportIDs(tx *gorm.DB, param model.RequestFundingSummaryParam) ([]model.RequestFundingSummaryRow, error) {
	var rows []model.RequestFundingSummaryRow

	if len(param.ReportIDs) == 0 {
		return rows, nil
	}

	err := tx.Table("requests").
		Select(`
			report_id,
			COALESCE(SUM(funding_target), 0) AS funding_target,
			COALESCE(SUM(funded_amount), 0) AS funded_amount,
			COALESCE(SUM(reserved_amount), 0) AS reserved_amount,
			COUNT(request_id) AS request_count
		`).
		Where("report_id IN ?", param.ReportIDs).
		Where("request_status IN ?", []string{"pending", "approved"}).
		Group("report_id").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func (r *RequestRepository) GetAllocationByDisaster(tx *gorm.DB, year int) ([]model.DisasterAllocationRow, error) {
	var rows []model.DisasterAllocationRow

	query := tx.Table("requests AS req").
		Select(`
			de.name AS disaster_event,
			COALESCE(SUM(req.funded_amount), 0) AS total_amount
		`).
		Joins("JOIN disaster_reports AS dr ON dr.report_id = req.report_id").
		Joins("JOIN disaster_events AS de ON de.event_id = dr.event_id").
		Where("req.request_status IN ?", []string{"pending", "approved"}).
		Group("de.name").
		Order("total_amount DESC")

	if year > 0 {
		query = query.Where("YEAR(req.created_at) = ?", year)
	}

	err := query.Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func (r *RequestRepository) IncrementFundedAmount(tx *gorm.DB, requestID uuid.UUID, amount float64) error {
	err := tx.Model(&entity.Requests{}).
		Where("request_id = ?", requestID).
		Update("funded_amount", gorm.Expr("funded_amount + ?", amount)).
		Error
	if err != nil {
		return err
	}

	return nil
}

func (r *RequestRepository) GetDonationLockContext(tx *gorm.DB, requestID uuid.UUID) (*model.DonationLockContextRow, error) {
	var row model.DonationLockContextRow

	err := tx.Table("requests AS req").
		Select(`
			req.request_id,
			p.post_id,
			p.user_id AS admin_user_id,
			p.name AS post_name,
			p.latitude,
			p.longitude
		`).
		Joins("JOIN disaster_reports AS dr ON dr.report_id = req.report_id").
		Joins("JOIN posts AS p ON p.post_id = dr.post_id").
		Where("req.request_id = ?", requestID).
		Scan(&row).Error
	if err != nil {
		return nil, err
	}

	if row.RequestID == uuid.Nil {
		return nil, gorm.ErrRecordNotFound
	}

	return &row, nil
}
