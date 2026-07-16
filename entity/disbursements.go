package entity

import (
	"time"

	"github.com/google/uuid"
)

type Disbursements struct {
	DisbursementID uuid.UUID `json:"disbursement_id" gorm:"type:varchar(36);primaryKey"`
	OrderID        uuid.UUID `json:"order_id" gorm:"type:varchar(36)"`
	StoreID        uuid.UUID `json:"store_id" gorm:"type:varchar(36)"`
	HeldBy         uuid.UUID `json:"held_by" gorm:"type:varchar(36)"`
	Amount         float64   `json:"amount" gorm:"type:decimal(10,2);not null"`
	IdempotencyKey string    `json:"idempotency_key" gorm:"type:varchar(255);not null;uniqueIndex"`
	Status         string    `json:"status" gorm:"type:enum('pending','success','failed');default:'pending'"`
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
