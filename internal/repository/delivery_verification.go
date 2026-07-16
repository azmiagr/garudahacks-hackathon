package repository

import (
	"strings"

	"github.com/azmiagr/garudahacks-hackathon/entity"
	"github.com/azmiagr/garudahacks-hackathon/model"
	"gorm.io/gorm"
)

type IDeliveryVerificationRepository interface {
	GetPublicDistributionProofs(tx *gorm.DB, param model.PublicDistributionParam) ([]model.PublicDistributionProofRow, error)
	GetVerifiedFulfillmentRate(tx *gorm.DB, year int) (*model.VerifiedFulfillmentRateRow, error)
}

type DeliveryVerificationRepository struct {
	db *gorm.DB
}

func NewDeliveryVerificationRepository(db *gorm.DB) IDeliveryVerificationRepository {
	return &DeliveryVerificationRepository{db: db}
}

func (r *DeliveryVerificationRepository) GetPublicDistributionProofs(tx *gorm.DB, param model.PublicDistributionParam) ([]model.PublicDistributionProofRow, error) {
	var rows []model.PublicDistributionProofRow

	latestCustodySubquery := tx.Table("custody_logs").
		Select("order_id, MAX(sequence) AS latest_sequence").
		Group("order_id")

	query := tx.Table("delivery_verifications AS dv").
		Select(`
			dv.verification_id,
			dv.order_id,
			p.post_id,
			p.name AS post_name,
			req.title AS request_title,
			de.name AS disaster_event,
			dv.image_url,
			dv.verification_status,
			dv.latitude,
			dv.longitude,
			dv.captured_at,
			o.total_amount,
			COUNT(DISTINCT d.donated_by) AS donor_count,
			COALESCE(cl.current_hash, '') AS current_hash
		`).
		Joins("JOIN orders AS o ON o.order_id = dv.order_id").
		Joins("JOIN requests AS req ON req.request_id = o.request_id").
		Joins("JOIN disaster_reports AS dr ON dr.report_id = req.report_id").
		Joins("JOIN posts AS p ON p.post_id = dr.post_id").
		Joins("JOIN disaster_events AS de ON de.event_id = dr.event_id").
		Joins("LEFT JOIN donations AS d ON d.request_id = req.request_id AND d.donation_status = 'approved'").
		Joins("LEFT JOIN (?) AS latest_cl ON latest_cl.order_id = o.order_id", latestCustodySubquery).
		Joins("LEFT JOIN custody_logs AS cl ON cl.order_id = latest_cl.order_id AND cl.sequence = latest_cl.latest_sequence").
		Where("dv.verification_status = ?", "approved").
		Group(`
			dv.verification_id,
			dv.order_id,
			p.post_id,
			p.name,
			req.title,
			de.name,
			dv.image_url,
			dv.verification_status,
			dv.latitude,
			dv.longitude,
			dv.captured_at,
			o.total_amount,
			cl.current_hash
		`)

	disasterFilter := strings.TrimSpace(param.DisasterType)
	if disasterFilter == "" {
		switch strings.TrimSpace(param.Filter) {
		case "banjir", "gempa", "longsor", "erupsi":
			disasterFilter = strings.TrimSpace(param.Filter)
		}
	}

	if disasterFilter != "" {
		query = query.Where("LOWER(de.name) = LOWER(?)", disasterFilter)
	}

	switch strings.TrimSpace(param.Filter) {
	case "largest_amount":
		query = query.Order("o.total_amount DESC")
	default:
		query = query.Order("dv.captured_at DESC")
	}

	err := query.
		Limit(normalizePublicDistributionLimit(param.Limit)).
		Scan(&rows).
		Error
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func (r *DeliveryVerificationRepository) GetVerifiedFulfillmentRate(tx *gorm.DB, year int) (*model.VerifiedFulfillmentRateRow, error) {
	var row model.VerifiedFulfillmentRateRow

	query := tx.Table("orders AS o").
		Select(`
			COUNT(DISTINCT o.order_id) AS total_orders,
			COUNT(DISTINCT CASE
				WHEN dv.verification_status = 'approved' THEN o.order_id
			END) AS verified_orders
		`).
		Joins("LEFT JOIN delivery_verifications AS dv ON dv.order_id = o.order_id").
		Where("o.order_status IN ?", []string{
			entity.OrderStatusAccepted,
			entity.OrderStatusPreparing,
			entity.OrderStatusReadyForPickup,
			entity.OrderStatusPickedUp,
			entity.OrderStatusInTransit,
			entity.OrderStatusDelivered,
			entity.OrderStatusCompleted,
		})

	if year > 0 {
		query = query.Where("YEAR(o.created_at) = ?", year)
	}

	err := query.Scan(&row).Error
	if err != nil {
		return nil, err
	}

	return &row, nil
}

func normalizePublicDistributionLimit(limit int) int {
	if limit <= 0 {
		return 6
	}

	if limit > 50 {
		return 50
	}

	return limit
}
