package model

import (
	"time"

	"github.com/google/uuid"
)

type StoreOrderListParam struct {
	Status string `form:"status"` // available, mine, accepted, ready, in_transit
	Limit  int    `form:"limit"`
	Offset int    `form:"offset"`
}

type StoreOrderListRepositoryParam struct {
	StoreID uuid.UUID
	Status  string
	Limit   int
	Offset  int
}

type StoreOrderDetailRepositoryParam struct {
	OrderID uuid.UUID
	StoreID uuid.UUID
}

type StoreOrderRow struct {
	OrderID       uuid.UUID  `json:"order_id"`
	RequestID     uuid.UUID  `json:"request_id"`
	OrderCode     string     `json:"order_code"`
	OrderStatus   string     `json:"order_status"`
	TotalAmount   float64    `json:"total_amount"`
	StoreID       uuid.UUID  `json:"store_id"`
	CourierID     uuid.UUID  `json:"courier_id"`
	RequestTitle  string     `json:"request_title"`
	PostName      string     `json:"post_name"`
	PostAddress   string     `json:"post_address"`
	PostLatitude  float64    `json:"post_latitude"`
	PostLongitude float64    `json:"post_longitude"`
	StoreName     string     `json:"store_name"`
	CourierName   string     `json:"courier_name"`
	AcceptedAt    *time.Time `json:"accepted_at"`
	ReadyAt       *time.Time `json:"ready_at"`
	PickedUpAt    *time.Time `json:"picked_up_at"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type StoreOrderItemRow struct {
	ItemID    uuid.UUID `json:"item_id"`
	Name      string    `json:"name"`
	Quantity  int       `json:"quantity"`
	Unit      int       `json:"unit"`
	UnitPrice float64   `json:"unit_price"`
	Subtotal  float64   `json:"subtotal"`
}

type StoreOrderListResponse struct {
	Items  []StoreOrderListItem `json:"items"`
	Limit  int                  `json:"limit"`
	Offset int                  `json:"offset"`
}

type StoreOrderListItem struct {
	OrderID      uuid.UUID `json:"order_id"`
	OrderCode    string    `json:"order_code"`
	OrderStatus  string    `json:"order_status"`
	TotalAmount  float64   `json:"total_amount"`
	RequestTitle string    `json:"request_title"`
	PostName     string    `json:"post_name"`
	PostAddress  string    `json:"post_address"`
	StoreName    string    `json:"store_name"`
	CourierName  string    `json:"courier_name"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type StoreOrderDetailResponse struct {
	StoreOrderListItem
	RequestID     uuid.UUID            `json:"request_id"`
	StoreID       uuid.UUID            `json:"store_id"`
	CourierID     uuid.UUID            `json:"courier_id"`
	PostLatitude  float64              `json:"post_latitude"`
	PostLongitude float64              `json:"post_longitude"`
	AcceptedAt    *time.Time           `json:"accepted_at"`
	ReadyAt       *time.Time           `json:"ready_at"`
	PickedUpAt    *time.Time           `json:"picked_up_at"`
	CreatedAt     time.Time            `json:"created_at"`
	Items         []StoreOrderItemItem `json:"items"`
}

type StoreOrderItemItem struct {
	ItemID    uuid.UUID `json:"item_id"`
	Name      string    `json:"name"`
	Quantity  int       `json:"quantity"`
	Unit      int       `json:"unit"`
	UnitPrice float64   `json:"unit_price"`
	Subtotal  float64   `json:"subtotal"`
}
