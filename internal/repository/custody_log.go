package repository

import (
	"github.com/azmiagr/garudahacks-hackathon/model"
	"gorm.io/gorm"
)

type ICustodyLogRepository interface {
	GetLatestPublicLedger(tx *gorm.DB, year int, limit int) ([]model.PublicLedgerRow, error)
}

type CustodyLogRepository struct {
	db *gorm.DB
}

func NewCustodyLogRepository(db *gorm.DB) ICustodyLogRepository {
	return &CustodyLogRepository{db: db}
}

func (r *CustodyLogRepository) GetLatestPublicLedger(tx *gorm.DB, year int, limit int) ([]model.PublicLedgerRow, error) {
	var rows []model.PublicLedgerRow

	query := tx.Table("custody_logs AS cl").
		Select(`
			cl.created_at AS occurred_at,
			'Kustodi kurir -> posko (handshake QR)' AS event,
			p.name AS post_name,
			CONCAT(COUNT(oi.order_item_id), ' item') AS value_label,
			cl.current_hash AS hash
		`).
		Joins("JOIN orders AS o ON o.order_id = cl.order_id").
		Joins("JOIN requests AS req ON req.request_id = o.request_id").
		Joins("JOIN disaster_reports AS dr ON dr.report_id = req.report_id").
		Joins("JOIN posts AS p ON p.post_id = dr.post_id").
		Joins("LEFT JOIN order_items AS oi ON oi.order_id = o.order_id").
		Group("cl.logs_id, cl.created_at, p.name, cl.current_hash").
		Order("cl.created_at DESC").
		Limit(normalizeLedgerLimit(limit))

	if year > 0 {
		query = query.Where("YEAR(cl.created_at) = ?", year)
	}

	err := query.Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func normalizeLedgerLimit(limit int) int {
	if limit <= 0 {
		return 4
	}

	if limit > 100 {
		return 100
	}

	return limit
}
