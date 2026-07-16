package repository

import (
	"strings"

	"github.com/azmiagr/garudahacks-hackathon/entity"
	"github.com/azmiagr/garudahacks-hackathon/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IStoreRepository interface {
	CreateStore(tx *gorm.DB, store *entity.Stores) error
	GetStore(tx *gorm.DB, param model.GetStoreParam) (*entity.Stores, error)
	GetStoreProfileStats(tx *gorm.DB, storeID uuid.UUID) (*model.StoreProfileStatsRow, error)
	GetStoreGoodnessCertificate(tx *gorm.DB, param model.StoreGoodnessParam) (*model.StoreGoodnessCertificateRow, error)
	GetStoreContributionHistory(tx *gorm.DB, param model.StoreGoodnessParam) ([]model.StoreContributionHistoryRow, error)
	CountStoreContributionHistory(tx *gorm.DB, param model.StoreGoodnessParam) (int64, error)
}

type StoreRepository struct {
	db *gorm.DB
}

func NewStoreRepository(db *gorm.DB) IStoreRepository {
	return &StoreRepository{db: db}
}

func (r *StoreRepository) CreateStore(tx *gorm.DB, store *entity.Stores) error {
	err := tx.Create(store).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *StoreRepository) GetStore(tx *gorm.DB, param model.GetStoreParam) (*entity.Stores, error) {
	var store entity.Stores
	query := tx.Model(&entity.Stores{})

	if param.StoreID != uuid.Nil {
		query = query.Where("store_id = ?", param.StoreID)
	}
	if param.OwnerID != uuid.Nil {
		query = query.Where("owner_id = ?", param.OwnerID).
			Where("store_id <> ?", uuid.Nil)
	}
	if strings.TrimSpace(param.BusinessNumber) != "" {
		query = query.Where("business_number = ?", strings.TrimSpace(param.BusinessNumber))
	}

	err := query.First(&store).Error
	if err != nil {
		return nil, err
	}

	return &store, nil
}

func (r *StoreRepository) GetStoreProfileStats(tx *gorm.DB, storeID uuid.UUID) (*model.StoreProfileStatsRow, error) {
	var row model.StoreProfileStatsRow
	err := tx.Table("orders AS o").
		Select(`
			COUNT(DISTINCT CASE
				WHEN o.created_at >= DATE_SUB(NOW(), INTERVAL 30 DAY)
				THEN o.order_id
			END) AS total_order_30_days,
			COUNT(DISTINCT CASE
				WHEN o.created_at >= DATE_SUB(NOW(), INTERVAL 30 DAY)
					AND o.order_status NOT IN ('cancelled', 'disputed', 'rejected')
				THEN o.order_id
			END) AS accepted_order_30_days,
			COUNT(DISTINCT CASE
				WHEN o.created_at >= DATE_SUB(NOW(), INTERVAL 30 DAY)
					AND o.order_status IN ('cancelled', 'rejected')
				THEN o.order_id
			END) AS cancelled_order_30_days,
			CASE
				WHEN COUNT(DISTINCT dv.order_id) = 0 THEN 0
				ELSE ROUND(4.5 + LEAST(COUNT(DISTINCT dv.order_id), 50) / 125, 1)
			END AS reputation_score
		`).
		Joins("LEFT JOIN delivery_verifications AS dv ON dv.order_id = o.order_id AND dv.verification_status = 'approved'").
		Where("o.store_id = ?", storeID).
		Scan(&row).Error
	if err != nil {
		return nil, err
	}

	return &row, nil
}

func (r *StoreRepository) GetStoreGoodnessCertificate(tx *gorm.DB, param model.StoreGoodnessParam) (*model.StoreGoodnessCertificateRow, error) {
	var row model.StoreGoodnessCertificateRow
	verifiedOrders := verifiedStoreOrdersSubquery(tx, param.Year)

	err := tx.Table("stores AS s").
		Select(`
			s.store_id,
			s.name AS store_name,
			s.bank_name,
			COUNT(vo.order_id) AS verified_order_count,
			COALESCE(SUM(vo.total_amount), 0) AS verified_amount_total,
			CASE
				WHEN COUNT(vo.order_id) = 0 THEN 0
				ELSE ROUND(4.5 + LEAST(COUNT(vo.order_id), 50) / 125, 1)
			END AS reputation_score,
			COUNT(DISTINCT CASE WHEN o.order_status = ? THEN o.order_id END) AS dispute_count,
			MIN(vo.verified_at) AS first_contribution_at
		`, entity.OrderStatusDisputed).
		Joins("LEFT JOIN (?) AS vo ON vo.store_id = s.store_id", verifiedOrders).
		Joins("LEFT JOIN orders AS o ON o.store_id = s.store_id").
		Where("s.store_id = ?", param.StoreID).
		Group("s.store_id, s.name, s.bank_name").
		Scan(&row).Error
	if err != nil {
		return nil, err
	}
	if row.StoreID == uuid.Nil {
		return nil, gorm.ErrRecordNotFound
	}

	return &row, nil
}

func (r *StoreRepository) GetStoreContributionHistory(tx *gorm.DB, param model.StoreGoodnessParam) ([]model.StoreContributionHistoryRow, error) {
	var rows []model.StoreContributionHistoryRow

	err := buildStoreContributionHistoryQuery(tx, param).
		Order("verified_at DESC").
		Limit(normalizeStoreGoodnessLimit(param.Limit)).
		Offset(normalizeStoreGoodnessOffset(param.Offset)).
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func (r *StoreRepository) CountStoreContributionHistory(tx *gorm.DB, param model.StoreGoodnessParam) (int64, error) {
	var count int64

	err := buildStoreContributionHistoryQuery(tx, param).
		Distinct("o.order_id").
		Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}

func buildStoreContributionHistoryQuery(tx *gorm.DB, param model.StoreGoodnessParam) *gorm.DB {
	latestCustodySubquery := tx.Table("custody_logs").
		Select("order_id, MAX(sequence) AS latest_sequence").
		Group("order_id")

	query := tx.Table("orders AS o").
		Select(`
			o.order_id,
			o.order_code,
			p.name AS post_name,
			de.name AS disaster_name,
			COUNT(DISTINCT oi.order_item_id) AS item_count,
			o.total_amount,
			MAX(dv.reviewed_at) AS verified_at,
			COALESCE(cl.current_hash, '') AS latest_hash
		`).
		Joins("JOIN delivery_verifications AS dv ON dv.order_id = o.order_id AND dv.verification_status = 'approved'").
		Joins("JOIN requests AS req ON req.request_id = o.request_id").
		Joins("JOIN disaster_reports AS dr ON dr.report_id = req.report_id").
		Joins("JOIN disaster_events AS de ON de.event_id = dr.event_id").
		Joins("JOIN posts AS p ON p.post_id = dr.post_id").
		Joins("LEFT JOIN order_items AS oi ON oi.order_id = o.order_id").
		Joins("LEFT JOIN (?) AS latest_cl ON latest_cl.order_id = o.order_id", latestCustodySubquery).
		Joins("LEFT JOIN custody_logs AS cl ON cl.order_id = latest_cl.order_id AND cl.sequence = latest_cl.latest_sequence").
		Where("o.store_id = ?", param.StoreID).
		Group(`
			o.order_id,
			o.order_code,
			p.name,
			de.name,
			o.total_amount,
			cl.current_hash
		`)

	if param.Year > 0 {
		query = query.Where("YEAR(dv.reviewed_at) = ?", param.Year)
	}

	return query
}

func verifiedStoreOrdersSubquery(tx *gorm.DB, year int) *gorm.DB {
	query := tx.Table("orders AS o").
		Select(`
			o.order_id,
			o.store_id,
			o.order_code,
			o.total_amount,
			o.order_status,
			MAX(dv.reviewed_at) AS verified_at
		`).
		Joins("JOIN delivery_verifications AS dv ON dv.order_id = o.order_id AND dv.verification_status = 'approved'").
		Group("o.order_id, o.store_id, o.order_code, o.total_amount, o.order_status")

	if year > 0 {
		query = query.Where("YEAR(dv.reviewed_at) = ?", year)
	}

	return query
}

func normalizeStoreGoodnessLimit(limit int) int {
	if limit <= 0 {
		return 20
	}
	if limit > 100 {
		return 100
	}
	return limit
}

func normalizeStoreGoodnessOffset(offset int) int {
	if offset < 0 {
		return 0
	}
	return offset
}
