package entity

import (
	"time"

	"github.com/google/uuid"
)

type Orders struct {
	OrderID         uuid.UUID `json:"order_id" gorm:"type:varchar(36);primaryKey"`
	RequestID       uuid.UUID `json:"request_id" gorm:"type:varchar(36)"`
	StoreID         uuid.UUID `json:"store_id" gorm:"type:varchar(36)"`
	CourierID       uuid.UUID `json:"courier_id" gorm:"type:varchar(36)"`
	OrderCode       string    `json:"order_code" gorm:"type:varchar(50);not null"`
	OrderStatus     string    `json:"order_status" gorm:"type:enum('pending','approved','rejected');default:'pending'"`
	TotalAmount     float64   `json:"total_amount" gorm:"type:decimal(10,2);not null"`
	BroadcastRadius float64   `json:"broadcast_radius" gorm:"type:decimal(10,2)"`
	CreatedAt       time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	DeliveryVerifications []DeliveryVerification `json:"delivery_verifications" gorm:"foreignKey:OrderID;references:OrderID;constraint:onDelete:CASCADE"`
	CustodyLogs           []CustodyLogs          `json:"custody_logs" gorm:"foreignKey:OrderID;references:OrderID;constraint:onDelete:CASCADE"`
	Disbursements         []Disbursements        `json:"disbursements" gorm:"foreignKey:OrderID;references:OrderID;constraint:onDelete:CASCADE"`
	OrderItems            []OrderItems           `json:"order_items" gorm:"foreignKey:OrderID;references:OrderID;constraint:onDelete:CASCADE"`
}
