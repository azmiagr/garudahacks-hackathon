package repository

import (
	"time"

	"github.com/azmiagr/garudahacks-hackathon/entity"
	"github.com/azmiagr/garudahacks-hackathon/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type IOrderRepository interface {
	CreateOrder(tx *gorm.DB, order *entity.Orders) error
	GetOrder(tx *gorm.DB, orderID uuid.UUID) (*entity.Orders, error)
	GetOrderForUpdate(tx *gorm.DB, orderID uuid.UUID) (*entity.Orders, error)
	AcceptOrderForStore(tx *gorm.DB, orderID uuid.UUID, storeID uuid.UUID, now time.Time) error
	MarkReadyForPickup(tx *gorm.DB, orderID uuid.UUID, storeID uuid.UUID, now time.Time) error
	AssignCourier(tx *gorm.DB, orderID uuid.UUID, courierID uuid.UUID) error
	UpdateOrder(tx *gorm.DB, order *entity.Orders) error
	GetStoreOrders(tx *gorm.DB, param model.StoreOrderListRepositoryParam) ([]model.StoreOrderRow, error)
	GetStoreOrderDetail(tx *gorm.DB, param model.StoreOrderDetailRepositoryParam) (*model.StoreOrderRow, error)
	GetStoreOrderItems(tx *gorm.DB, orderID uuid.UUID) ([]model.StoreOrderItemRow, error)
}

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) IOrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) CreateOrder(tx *gorm.DB, order *entity.Orders) error {
	err := tx.Create(order).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *OrderRepository) GetOrder(tx *gorm.DB, orderID uuid.UUID) (*entity.Orders, error) {
	var order entity.Orders
	err := tx.Where("order_id = ?", orderID).First(&order).Error
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (r *OrderRepository) GetOrderForUpdate(tx *gorm.DB, orderID uuid.UUID) (*entity.Orders, error) {
	var order entity.Orders
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("order_id = ?", orderID).
		First(&order).Error
	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (r *OrderRepository) AcceptOrderForStore(tx *gorm.DB, orderID uuid.UUID, storeID uuid.UUID, now time.Time) error {
	result := tx.Model(&entity.Orders{}).
		Where("order_id = ?", orderID).
		Where("order_status IN ?", []string{entity.OrderStatusPending, entity.OrderStatusBroadcasted}).
		Where("(store_id = ? OR store_id = ?)", uuid.Nil, storeID).
		Updates(map[string]interface{}{
			"store_id":     storeID,
			"order_status": entity.OrderStatusAccepted,
			"accepted_at":  now,
			"updated_at":   now,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *OrderRepository) MarkReadyForPickup(tx *gorm.DB, orderID uuid.UUID, storeID uuid.UUID, now time.Time) error {
	result := tx.Model(&entity.Orders{}).
		Where("order_id = ? AND store_id = ?", orderID, storeID).
		Where("order_status IN ?", []string{entity.OrderStatusAccepted, entity.OrderStatusPreparing}).
		Updates(map[string]interface{}{
			"order_status": entity.OrderStatusReadyForPickup,
			"ready_at":     now,
			"updated_at":   now,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *OrderRepository) AssignCourier(tx *gorm.DB, orderID uuid.UUID, courierID uuid.UUID) error {
	result := tx.Model(&entity.Orders{}).
		Where("order_id = ?", orderID).
		Update("courier_id", courierID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *OrderRepository) UpdateOrder(tx *gorm.DB, order *entity.Orders) error {
	return tx.Save(order).Error
}

func (r *OrderRepository) GetStoreOrders(tx *gorm.DB, param model.StoreOrderListRepositoryParam) ([]model.StoreOrderRow, error) {
	var rows []model.StoreOrderRow
	err := applyStoreOrderFilter(buildStoreOrderBaseQuery(tx), param).
		Order("o.updated_at DESC").
		Limit(normalizeStoreOrderLimit(param.Limit)).
		Offset(normalizeStoreOrderOffset(param.Offset)).
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func (r *OrderRepository) GetStoreOrderDetail(tx *gorm.DB, param model.StoreOrderDetailRepositoryParam) (*model.StoreOrderRow, error) {
	var row model.StoreOrderRow
	err := buildStoreOrderBaseQuery(tx).
		Where("o.order_id = ?", param.OrderID).
		Where("(o.store_id = ? OR o.store_id = ?)", param.StoreID, uuid.Nil).
		Scan(&row).Error
	if err != nil {
		return nil, err
	}
	if row.OrderID == uuid.Nil {
		return nil, gorm.ErrRecordNotFound
	}

	return &row, nil
}

func (r *OrderRepository) GetStoreOrderItems(tx *gorm.DB, orderID uuid.UUID) ([]model.StoreOrderItemRow, error) {
	var rows []model.StoreOrderItemRow
	err := tx.Table("order_items AS oi").
		Select(`
			i.item_id,
			i.name,
			oi.quantity,
			oi.unit,
			oi.unit_price,
			oi.subtotal
		`).
		Joins("JOIN items AS i ON i.item_id = oi.item_id").
		Where("oi.order_id = ?", orderID).
		Order("oi.created_at ASC").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func buildStoreOrderBaseQuery(tx *gorm.DB) *gorm.DB {
	return tx.Table("orders AS o").
		Select(`
			o.order_id,
			o.request_id,
			o.order_code,
			o.order_status,
			o.total_amount,
			o.store_id,
			o.courier_id,
			req.title AS request_title,
			p.name AS post_name,
			p.address AS post_address,
			p.latitude AS post_latitude,
			p.longitude AS post_longitude,
			COALESCE(s.name, '') AS store_name,
			COALESCE(courier.name, '') AS courier_name,
			o.accepted_at,
			o.ready_at,
			o.picked_up_at,
			o.created_at,
			o.updated_at
		`).
		Joins("JOIN requests AS req ON req.request_id = o.request_id").
		Joins("JOIN disaster_reports AS dr ON dr.report_id = req.report_id").
		Joins("JOIN posts AS p ON p.post_id = dr.post_id").
		Joins("LEFT JOIN stores AS s ON s.store_id = o.store_id").
		Joins("LEFT JOIN users AS courier ON courier.user_id = o.courier_id")
}

func applyStoreOrderFilter(query *gorm.DB, param model.StoreOrderListRepositoryParam) *gorm.DB {
	switch param.Status {
	case "mine":
		return query.Where("o.store_id = ?", param.StoreID)
	case "accepted":
		return query.Where("o.store_id = ? AND o.order_status IN ?", param.StoreID, []string{
			entity.OrderStatusAccepted,
			entity.OrderStatusPreparing,
		})
	case "ready":
		return query.Where("o.store_id = ? AND o.order_status = ?", param.StoreID, entity.OrderStatusReadyForPickup)
	case "in_transit":
		return query.Where("o.store_id = ? AND o.order_status IN ?", param.StoreID, []string{
			entity.OrderStatusPickedUp,
			entity.OrderStatusInTransit,
		})
	default:
		return query.Where("o.order_status IN ?", []string{
			entity.OrderStatusPending,
			entity.OrderStatusBroadcasted,
		}).Where("o.store_id = ?", uuid.Nil)
	}
}

func normalizeStoreOrderLimit(limit int) int {
	if limit <= 0 {
		return 20
	}
	if limit > 100 {
		return 100
	}
	return limit
}

func normalizeStoreOrderOffset(offset int) int {
	if offset < 0 {
		return 0
	}
	return offset
}
