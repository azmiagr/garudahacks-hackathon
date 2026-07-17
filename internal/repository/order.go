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
	ClaimOrderForCourier(tx *gorm.DB, orderID uuid.UUID, courierID uuid.UUID) error
	GetCourierTasks(tx *gorm.DB, param model.CourierTaskListRepositoryParam) ([]model.CourierTaskRow, error)
	GetCourierTaskDetail(tx *gorm.DB, param model.CourierTaskDetailRepositoryParam) (*model.CourierTaskRow, error)
	UpdateCourierLocation(tx *gorm.DB, orderID uuid.UUID, courierID uuid.UUID, lat float64, lng float64, capturedAt time.Time) error
	MarkCourierArrived(tx *gorm.DB, orderID uuid.UUID, courierID uuid.UUID, now time.Time) error
	MarkCourierArrivedAtPost(tx *gorm.DB, orderID uuid.UUID, courierID uuid.UUID, now time.Time) error
	GetCourierGoodnessStats(tx *gorm.DB, courierID uuid.UUID) (*model.CourierGoodnessStatsRow, error)
	GetCourierDeliveryHistory(tx *gorm.DB, param model.CourierGoodnessParam) ([]model.CourierDeliveryHistoryRow, error)
	CountCourierDeliveryHistory(tx *gorm.DB, param model.CourierGoodnessParam) (int64, error)
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

func (r *OrderRepository) ClaimOrderForCourier(tx *gorm.DB, orderID uuid.UUID, courierID uuid.UUID) error {
	result := tx.Model(&entity.Orders{}).
		Where("order_id = ?", orderID).
		Where("order_status = ?", entity.OrderStatusReadyForPickup).
		Where("courier_id = ?", uuid.Nil).
		Update("courier_id", courierID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *OrderRepository) GetCourierTasks(tx *gorm.DB, param model.CourierTaskListRepositoryParam) ([]model.CourierTaskRow, error) {
	var rows []model.CourierTaskRow
	err := applyCourierTaskFilter(buildCourierTaskBaseQuery(tx), param).
		Order("o.updated_at DESC").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func (r *OrderRepository) GetCourierTaskDetail(tx *gorm.DB, param model.CourierTaskDetailRepositoryParam) (*model.CourierTaskRow, error) {
	var row model.CourierTaskRow
	err := buildCourierTaskBaseQuery(tx).
		Where("o.order_id = ?", param.OrderID).
		Where("(o.courier_id = ? OR o.courier_id = ?)", param.CourierID, uuid.Nil).
		Scan(&row).Error
	if err != nil {
		return nil, err
	}
	if row.OrderID == uuid.Nil {
		return nil, gorm.ErrRecordNotFound
	}

	return &row, nil
}

func (r *OrderRepository) UpdateCourierLocation(tx *gorm.DB, orderID uuid.UUID, courierID uuid.UUID, lat float64, lng float64, capturedAt time.Time) error {
	result := tx.Model(&entity.Orders{}).
		Where("order_id = ?", orderID).
		Where("courier_id = ?", courierID).
		Where("order_status IN ?", []string{entity.OrderStatusReadyForPickup, entity.OrderStatusInTransit}).
		Updates(map[string]interface{}{
			"courier_latitude":            lat,
			"courier_longitude":           lng,
			"courier_location_updated_at": capturedAt,
			"updated_at":                  capturedAt,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *OrderRepository) MarkCourierArrived(tx *gorm.DB, orderID uuid.UUID, courierID uuid.UUID, now time.Time) error {
	result := tx.Model(&entity.Orders{}).
		Where("order_id = ?", orderID).
		Where("courier_id = ?", courierID).
		Where("order_status = ?", entity.OrderStatusReadyForPickup).
		Updates(map[string]interface{}{
			"arrived_at": now,
			"updated_at": now,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *OrderRepository) MarkCourierArrivedAtPost(tx *gorm.DB, orderID uuid.UUID, courierID uuid.UUID, now time.Time) error {
	result := tx.Model(&entity.Orders{}).
		Where("order_id = ?", orderID).
		Where("courier_id = ?", courierID).
		Where("order_status = ?", entity.OrderStatusInTransit).
		Updates(map[string]interface{}{
			"arrived_at_post_at": now,
			"updated_at":         now,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *OrderRepository) GetCourierGoodnessStats(tx *gorm.DB, courierID uuid.UUID) (*model.CourierGoodnessStatsRow, error) {
	var row model.CourierGoodnessStatsRow

	err := tx.Table("orders AS o").
		Select(`
			COUNT(DISTINCT CASE WHEN o.order_status IN ('delivered', 'completed') THEN o.order_id END) AS delivery_count,
			COALESCE(SUM(CASE WHEN o.order_status IN ('delivered', 'completed') THEN o.delivery_distance_km ELSE 0 END), 0) AS total_distance_km,
			COUNT(DISTINCT CASE WHEN o.order_status = 'disputed' THEN o.order_id END) AS dispute_count,
			MIN(CASE WHEN o.order_status IN ('delivered', 'completed') THEN o.delivered_at END) AS first_delivery_at,
			CASE
				WHEN COUNT(DISTINCT CASE WHEN o.order_status IN ('delivered', 'completed') THEN o.order_id END) = 0 THEN 0
				ELSE ROUND(4.5 + LEAST(COUNT(DISTINCT CASE WHEN o.order_status IN ('delivered', 'completed') THEN o.order_id END), 50) / 125, 1)
			END AS reputation_score
		`).
		Where("o.courier_id = ?", courierID).
		Scan(&row).Error
	if err != nil {
		return nil, err
	}

	return &row, nil
}

func (r *OrderRepository) GetCourierDeliveryHistory(tx *gorm.DB, param model.CourierGoodnessParam) ([]model.CourierDeliveryHistoryRow, error) {
	var rows []model.CourierDeliveryHistoryRow

	err := buildCourierDeliveryHistoryQuery(tx, param).
		Order("o.delivered_at DESC").
		Limit(normalizeStoreGoodnessLimit(param.Limit)).
		Offset(normalizeStoreGoodnessOffset(param.Offset)).
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func (r *OrderRepository) CountCourierDeliveryHistory(tx *gorm.DB, param model.CourierGoodnessParam) (int64, error) {
	var count int64

	err := buildCourierDeliveryHistoryQuery(tx, param).
		Distinct("o.order_id").
		Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}

func buildCourierDeliveryHistoryQuery(tx *gorm.DB, param model.CourierGoodnessParam) *gorm.DB {
	query := tx.Table("orders AS o").
		Select(`
			o.order_id,
			o.order_code,
			p.name AS post_name,
			de.name AS disaster_name,
			COALESCE(item_summary.item_count, 0) AS item_count,
			o.total_amount,
			o.delivery_distance_km,
			o.delivered_at
		`).
		Joins("JOIN requests AS req ON req.request_id = o.request_id").
		Joins("JOIN disaster_reports AS dr ON dr.report_id = req.report_id").
		Joins("JOIN disaster_events AS de ON de.event_id = dr.event_id").
		Joins("JOIN posts AS p ON p.post_id = dr.post_id").
		Joins("LEFT JOIN (?) AS item_summary ON item_summary.order_id = o.order_id", buildCourierTaskItemSummarySubquery(tx)).
		Where("o.courier_id = ?", param.CourierID).
		Where("o.order_status IN ?", []string{entity.OrderStatusDelivered, entity.OrderStatusCompleted})

	if param.Year > 0 {
		query = query.Where("YEAR(o.delivered_at) = ?", param.Year)
	}

	return query
}

func buildCourierTaskItemSummarySubquery(tx *gorm.DB) *gorm.DB {
	return tx.Table("order_items").
		Select("order_id, COUNT(order_item_id) AS item_count, COALESCE(SUM(quantity), 0) AS total_quantity").
		Group("order_id")
}

func buildCourierTaskBaseQuery(tx *gorm.DB) *gorm.DB {
	itemSummary := buildCourierTaskItemSummarySubquery(tx)

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
			de.name AS event_name,
			COALESCE(s.name, '') AS store_name,
			COALESCE(s.address, '') AS store_address,
			COALESCE(s.latitude, 0) AS store_latitude,
			COALESCE(s.longitude, 0) AS store_longitude,
			COALESCE(s.phone_number, '') AS store_phone_number,
			p.name AS post_name,
			p.address AS post_address,
			p.latitude AS post_latitude,
			p.longitude AS post_longitude,
			COALESCE(p.phone_number, '') AS post_phone_number,
			COALESCE(post_owner.name, '') AS post_contact_name,
			COALESCE(courier.name, '') AS courier_name,
			COALESCE(item_summary.item_count, 0) AS item_count,
			COALESCE(item_summary.total_quantity, 0) AS total_quantity,
			o.courier_latitude,
			o.courier_longitude,
			o.courier_location_updated_at,
			o.arrived_at,
			o.pickup_deadline_at,
			o.delivery_deadline_at,
			o.arrived_at_post_at,
			o.accepted_at,
			o.ready_at,
			o.picked_up_at,
			o.delivered_at,
			o.created_at,
			o.updated_at
		`).
		Joins("JOIN requests AS req ON req.request_id = o.request_id").
		Joins("JOIN disaster_reports AS dr ON dr.report_id = req.report_id").
		Joins("JOIN disaster_events AS de ON de.event_id = dr.event_id").
		Joins("JOIN posts AS p ON p.post_id = dr.post_id").
		Joins("LEFT JOIN users AS post_owner ON post_owner.user_id = p.user_id").
		Joins("LEFT JOIN stores AS s ON s.store_id = o.store_id").
		Joins("LEFT JOIN users AS courier ON courier.user_id = o.courier_id").
		Joins("LEFT JOIN (?) AS item_summary ON item_summary.order_id = o.order_id", itemSummary)
}

func applyCourierTaskFilter(query *gorm.DB, param model.CourierTaskListRepositoryParam) *gorm.DB {
	switch param.Status {
	case "mine":
		return query.Where("o.courier_id = ?", param.CourierID).
			Where("o.order_status IN ?", []string{
				entity.OrderStatusReadyForPickup,
				entity.OrderStatusPickedUp,
				entity.OrderStatusInTransit,
			})
	default:
		return query.Where("o.order_status = ?", entity.OrderStatusReadyForPickup).
			Where("o.courier_id = ?", uuid.Nil)
	}
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
	case "all":
		return query.Where(`
			(
				o.order_status IN ?
				AND o.store_id = ?
			)
			OR
			(
				o.store_id = ?
				AND o.order_status IN ?
			)
		`,
			[]string{entity.OrderStatusPending, entity.OrderStatusBroadcasted},
			uuid.Nil,
			param.StoreID,
			[]string{
				entity.OrderStatusAccepted,
				entity.OrderStatusPreparing,
				entity.OrderStatusReadyForPickup,
				entity.OrderStatusPickedUp,
				entity.OrderStatusInTransit,
				entity.OrderStatusDelivered,
				entity.OrderStatusCompleted,
			},
		)
	case "available":
		return query.Where("o.order_status IN ?", []string{
			entity.OrderStatusPending,
			entity.OrderStatusBroadcasted,
		}).Where("o.store_id = ?", uuid.Nil)
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
		return query
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
