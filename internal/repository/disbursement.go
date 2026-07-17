package repository

import (
	"math"
	"time"

	"github.com/azmiagr/garudahacks-hackathon/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IDisbursementRepository interface {
	GetVerifiedDisbursedTotal(tx *gorm.DB, year int) (float64, error)
	GetMonthlyDisbursements(tx *gorm.DB, year int) ([]model.MonthlyDisbursementRow, error)
	GetMonthlyDisbursementsInRange(tx *gorm.DB, start time.Time, end time.Time) ([]model.MonthlyDisbursementRow, error)
	GetStoreDisbursementSummary(tx *gorm.DB, param model.StoreDisbursementDashboardParam) (*model.StoreDisbursementSummaryRow, error)
	GetStoreDisbursementHistory(tx *gorm.DB, param model.StoreDisbursementDashboardParam) ([]model.StoreDisbursementHistoryRow, error)
	CountStoreDisbursementHistory(tx *gorm.DB, param model.StoreDisbursementDashboardParam) (int64, error)
	GetStoreGoodnessTrail(tx *gorm.DB, storeID uuid.UUID, year int) (*model.StoreGoodnessTrailRow, error)
}

type DisbursementRepository struct {
	db *gorm.DB
}

func NewDisbursementRepository(db *gorm.DB) IDisbursementRepository {
	return &DisbursementRepository{db: db}
}

func (r *DisbursementRepository) GetVerifiedDisbursedTotal(tx *gorm.DB, year int) (float64, error) {
	var total float64

	query := tx.Table("disbursements").
		Select("COALESCE(SUM(amount), 0)").
		Where("status = ?", "success")

	if year > 0 {
		query = query.Where("YEAR(created_at) = ?", year)
	}

	err := query.Scan(&total).Error
	if err != nil {
		return 0, err
	}

	return total, nil
}

func (r *DisbursementRepository) GetMonthlyDisbursements(tx *gorm.DB, year int) ([]model.MonthlyDisbursementRow, error) {
	var rows []model.MonthlyDisbursementRow

	query := tx.Table("disbursements").
		Select(`
			MONTH(created_at) AS month,
			COALESCE(SUM(amount), 0) AS total
		`).
		Where("status = ?", "success").
		Group("MONTH(created_at)").
		Order("MONTH(created_at) ASC")

	if year > 0 {
		query = query.Where("YEAR(created_at) = ?", year)
	}

	err := query.Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func (r *DisbursementRepository) GetMonthlyDisbursementsInRange(tx *gorm.DB, start time.Time, end time.Time) ([]model.MonthlyDisbursementRow, error) {
	var rows []model.MonthlyDisbursementRow

	disbursedAtExpr := "COALESCE(disbursed_at, updated_at, created_at)"
	err := tx.Table("disbursements").
		Select(`
			YEAR(`+disbursedAtExpr+`) AS year,
			MONTH(`+disbursedAtExpr+`) AS month,
			COALESCE(SUM(amount), 0) AS total
		`).
		Where("status = ?", "success").
		Where(disbursedAtExpr+" >= ?", start).
		Where(disbursedAtExpr+" < ?", end).
		Group("YEAR(" + disbursedAtExpr + "), MONTH(" + disbursedAtExpr + ")").
		Order("YEAR(" + disbursedAtExpr + ") ASC, MONTH(" + disbursedAtExpr + ") ASC").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func (r *DisbursementRepository) GetStoreDisbursementSummary(tx *gorm.DB, param model.StoreDisbursementDashboardParam) (*model.StoreDisbursementSummaryRow, error) {
	var row model.StoreDisbursementSummaryRow

	month := normalizeStoreDisbursementMonth(param.Month)
	year := normalizeStoreDisbursementYear(param.Year)

	err := tx.Table("stores AS s").
		Select(`
			s.store_id,
			s.name AS store_name,
			s.bank_name,
			s.bank_account_no,
			COALESCE(SUM(CASE
				WHEN d.status = 'success'
					AND YEAR(COALESCE(d.disbursed_at, d.updated_at)) = ?
					AND MONTH(COALESCE(d.disbursed_at, d.updated_at)) = ?
				THEN d.amount ELSE 0 END), 0) AS total_disbursed_this_month,
			COUNT(DISTINCT CASE
				WHEN d.status = 'success'
					AND YEAR(COALESCE(d.disbursed_at, d.updated_at)) = ?
					AND MONTH(COALESCE(d.disbursed_at, d.updated_at)) = ?
				THEN d.order_id END) AS completed_order_count,
			COUNT(DISTINCT CASE
				WHEN o.order_status = 'disputed'
				THEN o.order_id END) AS dispute_count,
			COALESCE(AVG(CASE
				WHEN d.status = 'success' AND dv.reviewed_at IS NOT NULL
				THEN TIMESTAMPDIFF(MINUTE, dv.reviewed_at, COALESCE(d.disbursed_at, d.updated_at))
			END), 0) AS median_disbursement_min
		`, year, month, year, month).
		Joins("LEFT JOIN disbursements AS d ON d.store_id = s.store_id").
		Joins("LEFT JOIN orders AS o ON o.order_id = d.order_id").
		Joins("LEFT JOIN delivery_verifications AS dv ON dv.order_id = o.order_id AND dv.verification_status = 'approved'").
		Where("s.store_id = ?", param.StoreID).
		Group("s.store_id, s.name, s.bank_name, s.bank_account_no").
		Scan(&row).Error
	if err != nil {
		return nil, err
	}
	if row.StoreID == uuid.Nil {
		return nil, gorm.ErrRecordNotFound
	}

	return &row, nil
}

func (r *DisbursementRepository) GetStoreDisbursementHistory(tx *gorm.DB, param model.StoreDisbursementDashboardParam) ([]model.StoreDisbursementHistoryRow, error) {
	var rows []model.StoreDisbursementHistoryRow

	err := buildStoreDisbursementHistoryQuery(tx, param).
		Order("d.created_at DESC").
		Limit(normalizeStoreDisbursementLimit(param.Limit)).
		Offset(normalizeStoreDisbursementOffset(param.Offset)).
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func (r *DisbursementRepository) CountStoreDisbursementHistory(tx *gorm.DB, param model.StoreDisbursementDashboardParam) (int64, error) {
	var count int64

	err := buildStoreDisbursementHistoryQuery(tx, param).
		Distinct("d.disbursement_id").
		Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}

func buildStoreDisbursementHistoryQuery(tx *gorm.DB, param model.StoreDisbursementDashboardParam) *gorm.DB {
	query := tx.Table("disbursements AS d").
		Select(`
			d.disbursement_id,
			d.order_id,
			o.order_code,
			p.name AS post_name,
			d.amount,
			d.status,
			d.idempotency_key,
			COALESCE(d.gateway_ref, '') AS gateway_reference,
			COALESCE(d.gateway_attempt, 0) AS gateway_attempt,
			dv.reviewed_at AS verification_approved_at,
			d.disbursed_at,
			d.created_at,
			d.updated_at,
			COALESCE(TIMESTAMPDIFF(MINUTE, dv.reviewed_at, COALESCE(d.disbursed_at, d.updated_at)), 0) AS minutes_after_verification
		`).
		Joins("JOIN orders AS o ON o.order_id = d.order_id").
		Joins("JOIN requests AS req ON req.request_id = o.request_id").
		Joins("JOIN disaster_reports AS dr ON dr.report_id = req.report_id").
		Joins("JOIN posts AS p ON p.post_id = dr.post_id").
		Joins("LEFT JOIN delivery_verifications AS dv ON dv.order_id = o.order_id AND dv.verification_status = 'approved'").
		Where("d.store_id = ?", param.StoreID)

	if param.Year > 0 {
		query = query.Where("YEAR(d.created_at) = ?", param.Year)
	}
	if param.Month > 0 {
		query = query.Where("MONTH(d.created_at) = ?", param.Month)
	}

	return query
}
func (r *DisbursementRepository) GetStoreGoodnessTrail(tx *gorm.DB, storeID uuid.UUID, year int) (*model.StoreGoodnessTrailRow, error) {
	var row model.StoreGoodnessTrailRow

	query := tx.Table("orders AS o").
		Select(`
			o.store_id,
			COUNT(DISTINCT o.order_id) AS verified_order_count,
			COALESCE(SUM(o.total_amount), 0) AS verified_amount_total,
			MIN(dv.reviewed_at) AS first_contribution_at,
			MAX(dv.reviewed_at) AS last_contribution_at
		`).
		Joins("JOIN delivery_verifications AS dv ON dv.order_id = o.order_id AND dv.verification_status = 'approved'").
		Where("o.store_id = ?", storeID).
		Group("o.store_id")

	if year > 0 {
		query = query.Where("YEAR(dv.reviewed_at) = ?", year)
	}

	err := query.Scan(&row).Error
	if err != nil {
		return nil, err
	}
	if row.StoreID == uuid.Nil {
		row.StoreID = storeID
	}

	return &row, nil
}

func normalizeStoreDisbursementYear(year int) int {
	if year <= 0 {
		return time.Now().Year()
	}
	return year
}

func normalizeStoreDisbursementMonth(month int) int {
	if month < 1 || month > 12 {
		return int(time.Now().Month())
	}
	return month
}

func normalizeStoreDisbursementLimit(limit int) int {
	if limit <= 0 {
		return 20
	}
	return int(math.Min(float64(limit), 100))
}

func normalizeStoreDisbursementOffset(offset int) int {
	if offset < 0 {
		return 0
	}
	return offset
}
