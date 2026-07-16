package repository

import (
	"github.com/azmiagr/garudahacks-hackathon/entity"
	"github.com/azmiagr/garudahacks-hackathon/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ICustodyLogRepository interface {
	GetLatestPublicLedger(tx *gorm.DB, year int, limit int) ([]model.PublicLedgerRow, error)
	GetLatestCustodyLog(tx *gorm.DB) (*entity.CustodyLogs, error)
	GetLatestCustodyLogByOrderForUpdate(tx *gorm.DB, orderID uuid.UUID) (*entity.CustodyLogs, error)
	ListCustodyLogsByOrder(tx *gorm.DB, orderID uuid.UUID) ([]entity.CustodyLogs, error)
	CreateCustodyLog(tx *gorm.DB, log *entity.CustodyLogs) error
	ExistsIdempotencyKey(tx *gorm.DB, key string) (bool, error)
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

func (r *CustodyLogRepository) GetLatestCustodyLog(tx *gorm.DB) (*entity.CustodyLogs, error) {
	var log entity.CustodyLogs
	err := tx.Order("sequence DESC").First(&log).Error
	if err != nil {
		return nil, err
	}
	return &log, nil
}

func (r *CustodyLogRepository) GetLatestCustodyLogByOrderForUpdate(tx *gorm.DB, orderID uuid.UUID) (*entity.CustodyLogs, error) {
	var log entity.CustodyLogs
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("order_id = ?", orderID.String()).
		Order("sequence DESC").
		First(&log).Error
	if err != nil {
		return nil, err
	}

	return &log, nil
}

func (r *CustodyLogRepository) ListCustodyLogsByOrder(tx *gorm.DB, orderID uuid.UUID) ([]entity.CustodyLogs, error) {
	var logs []entity.CustodyLogs
	err := tx.Where("order_id = ?", orderID.String()).
		Order("sequence ASC").
		Find(&logs).Error
	if err != nil {
		return nil, err
	}

	return logs, nil
}

func (r *CustodyLogRepository) CreateCustodyLog(tx *gorm.DB, log *entity.CustodyLogs) error {
	err := tx.Create(log).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *CustodyLogRepository) ExistsIdempotencyKey(tx *gorm.DB, key string) (bool, error) {
	if key == "" {
		return false, nil
	}

	var count int64
	err := tx.Model(&entity.CustodyLogs{}).
		Where("idempotency_key = ?", key).
		Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
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
