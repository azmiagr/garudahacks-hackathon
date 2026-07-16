package entity

import (
	"time"

	"github.com/google/uuid"
)

const (
	OrderStatusPending        = "pending"
	OrderStatusBroadcasted    = "broadcasted"
	OrderStatusAccepted       = "accepted"
	OrderStatusPreparing      = "preparing"
	OrderStatusReadyForPickup = "ready_for_pickup"
	OrderStatusPickedUp       = "picked_up"
	OrderStatusInTransit      = "in_transit"
	OrderStatusDelivered      = "delivered"
	OrderStatusCompleted      = "completed"
	OrderStatusCancelled      = "cancelled"
	OrderStatusDisputed       = "disputed"
)

type Orders struct {
	OrderID         uuid.UUID  `json:"order_id" gorm:"type:varchar(36);primaryKey"`
	RequestID       uuid.UUID  `json:"request_id" gorm:"type:varchar(36);index"`
	StoreID         uuid.UUID  `json:"store_id" gorm:"type:varchar(36);index"`
	CourierID       uuid.UUID  `json:"courier_id" gorm:"type:varchar(36);index"`
	OrderCode       string     `json:"order_code" gorm:"type:varchar(50);not null"`
	OrderStatus     string     `json:"order_status" gorm:"type:varchar(40);default:'pending';index"`
	TotalAmount     float64    `json:"total_amount" gorm:"type:decimal(10,2);not null"`
	BroadcastRadius float64    `json:"broadcast_radius" gorm:"type:decimal(10,2)"`
	AcceptedAt      *time.Time `json:"accepted_at" gorm:"type:datetime"`
	ReadyAt         *time.Time `json:"ready_at" gorm:"type:datetime"`
	PickedUpAt      *time.Time `json:"picked_up_at" gorm:"type:datetime"`
	DeliveredAt     *time.Time `json:"delivered_at" gorm:"type:datetime"`
	CompletedAt     *time.Time `json:"completed_at" gorm:"type:datetime"`
	CancelledAt     *time.Time `json:"cancelled_at" gorm:"type:datetime"`
	CreatedAt       time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time  `json:"updated_at" gorm:"autoUpdateTime"`

	DeliveryVerifications []DeliveryVerification  `json:"delivery_verifications" gorm:"foreignKey:OrderID;references:OrderID;constraint:onDelete:CASCADE"`
	CustodyLogs           []CustodyLogs           `json:"custody_logs" gorm:"foreignKey:OrderID;references:OrderID;constraint:onDelete:CASCADE"`
	CustodyTokens         []CustodyHandshakeToken `json:"custody_tokens" gorm:"foreignKey:OrderID;references:OrderID;constraint:onDelete:CASCADE"`
	Disbursements         []Disbursements         `json:"disbursements" gorm:"foreignKey:OrderID;references:OrderID;constraint:onDelete:CASCADE"`
	OrderItems            []OrderItems            `json:"order_items" gorm:"foreignKey:OrderID;references:OrderID;constraint:onDelete:CASCADE"`
}
