package repository

import (
	"github.com/azmiagr/garudahacks-hackathon/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IRequestRepository interface {
	GetFundingSummaryByReportIDs(tx *gorm.DB, param model.RequestFundingSummaryParam) ([]model.RequestFundingSummaryRow, error)
	GetAllocationByDisaster(tx *gorm.DB, year int) ([]model.DisasterAllocationRow, error)
}

type RequestRepository struct {
	db *gorm.DB
}

func NewRequestRepository(db *gorm.DB) IRequestRepository {
	return &RequestRepository{db: db}
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

func UUIDsFromReports(reports []model.LatestDisasterReportRow) []uuid.UUID {
	ids := make([]uuid.UUID, 0, len(reports))

	for _, report := range reports {
		ids = append(ids, report.ReportID)
	}

	return ids
}
