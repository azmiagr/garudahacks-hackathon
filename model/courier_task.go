package model

import (
	"time"

	"github.com/google/uuid"
)

type CourierTaskListParam struct {
	Status    string  `form:"status"`
	Latitude  float64 `form:"lat"`
	Longitude float64 `form:"lng"`
	HasCoords bool
	RadiusKm  float64 `form:"radius_km"`
	Limit     int     `form:"limit"`
	Offset    int     `form:"offset"`
}

type CourierTaskListRepositoryParam struct {
	CourierID uuid.UUID
	Status    string
}

type CourierTaskDetailRepositoryParam struct {
	OrderID   uuid.UUID
	CourierID uuid.UUID
}

type CourierTaskRow struct {
	OrderID                  uuid.UUID  `json:"order_id"`
	RequestID                uuid.UUID  `json:"request_id"`
	OrderCode                string     `json:"order_code"`
	OrderStatus              string     `json:"order_status"`
	TotalAmount              float64    `json:"total_amount"`
	StoreID                  uuid.UUID  `json:"store_id"`
	CourierID                uuid.UUID  `json:"courier_id"`
	RequestTitle             string     `json:"request_title"`
	EventName                string     `json:"event_name"`
	StoreName                string     `json:"store_name"`
	StoreAddress             string     `json:"store_address"`
	StoreLatitude            float64    `json:"store_latitude"`
	StoreLongitude           float64    `json:"store_longitude"`
	StorePhoneNumber         string     `json:"store_phone_number"`
	PostName                 string     `json:"post_name"`
	PostAddress              string     `json:"post_address"`
	PostLatitude             float64    `json:"post_latitude"`
	PostLongitude            float64    `json:"post_longitude"`
	PostPhoneNumber          string     `json:"post_phone_number"`
	PostContactName          string     `json:"post_contact_name"`
	CourierName              string     `json:"courier_name"`
	ItemCount                int        `json:"item_count"`
	TotalQuantity            int        `json:"total_quantity"`
	CourierLatitude          *float64   `json:"courier_latitude"`
	CourierLongitude         *float64   `json:"courier_longitude"`
	CourierLocationUpdatedAt *time.Time `json:"courier_location_updated_at"`
	ArrivedAt                *time.Time `json:"arrived_at"`
	PickupDeadlineAt         *time.Time `json:"pickup_deadline_at"`
	DeliveryDeadlineAt       *time.Time `json:"delivery_deadline_at"`
	ArrivedAtPostAt          *time.Time `json:"arrived_at_post_at"`
	AcceptedAt               *time.Time `json:"accepted_at"`
	ReadyAt                  *time.Time `json:"ready_at"`
	PickedUpAt               *time.Time `json:"picked_up_at"`
	DeliveredAt              *time.Time `json:"delivered_at"`
	CreatedAt                time.Time  `json:"created_at"`
	UpdatedAt                time.Time  `json:"updated_at"`
}

type CourierTaskListResponse struct {
	Items  []CourierTaskListItem `json:"items"`
	Limit  int                   `json:"limit"`
	Offset int                   `json:"offset"`
}

type CourierTaskListItem struct {
	OrderID           uuid.UUID `json:"order_id"`
	OrderCode         string    `json:"order_code"`
	OrderStatus       string    `json:"order_status"`
	TotalAmount       float64   `json:"total_amount"`
	RequestTitle      string    `json:"request_title"`
	EventName         string    `json:"event_name"`
	StoreName         string    `json:"store_name"`
	StoreAddress      string    `json:"store_address"`
	PostName          string    `json:"post_name"`
	PostAddress       string    `json:"post_address"`
	ItemCount         int       `json:"item_count"`
	TotalQuantity     int       `json:"total_quantity"`
	PickupDistanceKm  *float64  `json:"pickup_distance_km"`
	DropoffDistanceKm *float64  `json:"dropoff_distance_km"`
	TotalDistanceKm   *float64  `json:"total_distance_km"`
	UpdatedAt         time.Time `json:"updated_at"`
}

type CourierTaskDetailResponse struct {
	CourierTaskListItem
	RequestID                uuid.UUID            `json:"request_id"`
	StoreID                  uuid.UUID            `json:"store_id"`
	CourierID                uuid.UUID            `json:"courier_id"`
	CourierName              string               `json:"courier_name"`
	StoreLatitude            float64              `json:"store_latitude"`
	StoreLongitude           float64              `json:"store_longitude"`
	StorePhoneNumber         string               `json:"store_phone_number"`
	PostLatitude             float64              `json:"post_latitude"`
	PostLongitude            float64              `json:"post_longitude"`
	PostPhoneNumber          string               `json:"post_phone_number"`
	PostContactName          string               `json:"post_contact_name"`
	CourierLatitude          *float64             `json:"courier_latitude"`
	CourierLongitude         *float64             `json:"courier_longitude"`
	CourierLocationUpdatedAt *time.Time           `json:"courier_location_updated_at"`
	ArrivedAt                *time.Time           `json:"arrived_at"`
	ArrivedAtPostAt          *time.Time           `json:"arrived_at_post_at"`
	PickupDeadlineAt         *time.Time           `json:"pickup_deadline_at"`
	DeliveryDeadlineAt       *time.Time           `json:"delivery_deadline_at"`
	EtaMinutes               *float64             `json:"eta_minutes"`
	LatestCustodyStage       string               `json:"latest_custody_stage"`
	LatestCustodyHash        string               `json:"latest_custody_hash"`
	LatestCustodyShortHash   string               `json:"latest_custody_short_hash"`
	LatestCustodyCapturedAt  *time.Time           `json:"latest_custody_captured_at"`
	AcceptedAt               *time.Time           `json:"accepted_at"`
	ReadyAt                  *time.Time           `json:"ready_at"`
	PickedUpAt               *time.Time           `json:"picked_up_at"`
	DeliveredAt              *time.Time           `json:"delivered_at"`
	CreatedAt                time.Time            `json:"created_at"`
	Items                    []StoreOrderItemItem `json:"items"`
}

type CourierTaskActionResponse struct {
	OrderID     uuid.UUID `json:"order_id"`
	CourierID   uuid.UUID `json:"courier_id"`
	OrderStatus string    `json:"order_status"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CourierTaskClaimRequest struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lng"`
	HasCoords bool    `json:"-"`
}

type CourierLocationPingRequest struct {
	Latitude   float64    `json:"lat" binding:"required"`
	Longitude  float64    `json:"lng" binding:"required"`
	CapturedAt *time.Time `json:"captured_at"`
}

type CourierLocationPingResponse struct {
	OrderID                  uuid.UUID `json:"order_id"`
	CourierLatitude          float64   `json:"courier_latitude"`
	CourierLongitude         float64   `json:"courier_longitude"`
	CourierLocationUpdatedAt time.Time `json:"courier_location_updated_at"`
	PickupDistanceKm         float64   `json:"pickup_distance_km"`
	EtaMinutes               float64   `json:"eta_minutes"`
}

type CourierArrivedResponse struct {
	OrderID     uuid.UUID `json:"order_id"`
	OrderStatus string    `json:"order_status"`
	ArrivedAt   time.Time `json:"arrived_at"`
}

type CourierArrivedAtPostResponse struct {
	OrderID         uuid.UUID `json:"order_id"`
	OrderStatus     string    `json:"order_status"`
	ArrivedAtPostAt time.Time `json:"arrived_at_post_at"`
}
