package repository

import (
	"github.com/azmiagr/garudahacks-hackathon/model"
	"gorm.io/gorm"
)

type IAdminDashboardRepository interface {
	GetAdminDashboardEvents(tx *gorm.DB, param model.AdminDashboardHomeParam) ([]model.AdminDashboardEventRow, error)
	GetAdminDashboardOrders(tx *gorm.DB, param model.AdminDashboardHomeParam) ([]model.AdminDashboardOrderRow, error)
}

type AdminDashboardRepository struct {
	db *gorm.DB
}

func NewAdminDashboardRepository(db *gorm.DB) IAdminDashboardRepository {
	return &AdminDashboardRepository{db: db}
}

func (r *AdminDashboardRepository) GetAdminDashboardEvents(tx *gorm.DB, param model.AdminDashboardHomeParam) ([]model.AdminDashboardEventRow, error) {
	var rows []model.AdminDashboardEventRow

	latestReportSubquery := tx.Table("disaster_reports").
		Select("post_id, MAX(COALESCE(reported_at, created_at)) AS latest_reported_at").
		Group("post_id")

	fundingSubquery := tx.Table("disaster_reports AS dr").
		Select(`
			dr.post_id,
			COALESCE(SUM(req.funding_target), 0) AS funding_target,
			COALESCE(SUM(req.funded_amount), 0) AS funded_amount
		`).
		Joins("LEFT JOIN requests AS req ON req.report_id = dr.report_id").
		Where("req.request_status IN ?", []string{"pending", "approved"}).
		Group("dr.post_id")

	orderSubquery := tx.Table("disaster_reports AS dr").
		Select(`
			dr.post_id,
			COUNT(DISTINCT o.order_id) AS order_count,
			COUNT(DISTINCT CASE
				WHEN dv.verification_status = 'approved' THEN o.order_id
			END) AS completed_order_count
		`).
		Joins("LEFT JOIN requests AS req ON req.report_id = dr.report_id").
		Joins("LEFT JOIN orders AS o ON o.request_id = req.request_id").
		Joins("LEFT JOIN delivery_verifications AS dv ON dv.order_id = o.order_id").
		Group("dr.post_id")

	err := tx.Table("posts AS p").
		Select(`
			p.post_id,
			p.name AS title,
			p.address,
			p.geofence_radius,
			COALESCE(de.name, '') AS disaster_type,
			COALESCE(dr.image_url, '') AS image_url,
			COALESCE(dr.reported_at, dr.created_at, p.created_at) AS started_at,
			COALESCE(funding.funding_target, 0) AS funding_target,
			COALESCE(funding.funded_amount, 0) AS funded_amount,
			COALESCE(orders.order_count, 0) AS order_count,
			COALESCE(orders.completed_order_count, 0) AS completed_order_count
		`).
		Joins(`
			LEFT JOIN (?) AS latest_reports
			ON latest_reports.post_id = p.post_id
		`, latestReportSubquery).
		Joins(`
			LEFT JOIN disaster_reports AS dr
			ON dr.post_id = latest_reports.post_id
			AND COALESCE(dr.reported_at, dr.created_at) = latest_reports.latest_reported_at
		`).
		Joins("LEFT JOIN disaster_events AS de ON de.event_id = dr.event_id").
		Joins("LEFT JOIN (?) AS funding ON funding.post_id = p.post_id", fundingSubquery).
		Joins("LEFT JOIN (?) AS orders ON orders.post_id = p.post_id", orderSubquery).
		Where("p.user_id = ?", param.UserID).
		Order("p.created_at DESC").
		Scan(&rows).
		Error

	if err != nil {
		return nil, err
	}

	return rows, nil
}

func (r *AdminDashboardRepository) GetAdminDashboardOrders(tx *gorm.DB, param model.AdminDashboardHomeParam) ([]model.AdminDashboardOrderRow, error) {
	var rows []model.AdminDashboardOrderRow

	latestVerificationSubquery := tx.Table("delivery_verifications").
		Select("order_id, MAX(created_at) AS latest_created_at").
		Group("order_id")

	err := tx.Table("orders AS o").
		Select(`
			dr.post_id,
			o.order_id,
			o.order_code,
			o.order_status,
			COALESCE(s.name, '') AS store_name,
			COALESCE(courier.name, '') AS courier_name,
			o.updated_at,
			COALESCE(dv.verification_status, '') AS verification_status
		`).
		Joins("JOIN requests AS req ON req.request_id = o.request_id").
		Joins("JOIN disaster_reports AS dr ON dr.report_id = req.report_id").
		Joins("JOIN posts AS p ON p.post_id = dr.post_id").
		Joins("LEFT JOIN stores AS s ON s.store_id = o.store_id").
		Joins("LEFT JOIN users AS courier ON courier.user_id = o.courier_id").
		Joins("LEFT JOIN (?) AS latest_dv ON latest_dv.order_id = o.order_id", latestVerificationSubquery).
		Joins(`
			LEFT JOIN delivery_verifications AS dv
			ON dv.order_id = latest_dv.order_id
			AND dv.created_at = latest_dv.latest_created_at
		`).
		Where("p.user_id = ?", param.UserID).
		Order("o.updated_at DESC").
		Limit(50).
		Scan(&rows).
		Error

	if err != nil {
		return nil, err
	}

	return rows, nil
}
