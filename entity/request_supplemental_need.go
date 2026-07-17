package entity

import (
	"time"

	"github.com/google/uuid"
)

type RequestSupplementalNeed struct {
	SupplementalID        uuid.UUID `json:"supplemental_id" gorm:"type:varchar(36);primaryKey"`
	RequestID             uuid.UUID `json:"request_id" gorm:"type:varchar(36);not null;index"`
	OrderID               uuid.UUID `json:"order_id" gorm:"type:varchar(36);not null;index"`
	CreatedBy             uuid.UUID `json:"created_by" gorm:"type:varchar(36);not null;index"`
	Reason                string    `json:"reason" gorm:"type:text;not null"`
	ReservedAmountApplied float64   `json:"reserved_amount_applied" gorm:"type:decimal(10,2);not null;default:0"`
	AdditionalTarget      float64   `json:"additional_target" gorm:"type:decimal(10,2);not null;default:0"`
	CreatedAt             time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt             time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
