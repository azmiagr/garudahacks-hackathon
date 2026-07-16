package entity

import (
	"time"

	"github.com/google/uuid"
)

type Items struct {
	ItemID            uuid.UUID `json:"item_id" gorm:"type:varchar(36);primaryKey"`
	RequestID         uuid.UUID `json:"request_id" gorm:"type:varchar(36)"`
	Name              string    `json:"name" gorm:"type:varchar(150);not null"`
	Description       string    `json:"description" gorm:"type:text"`
	Price             float64   `json:"price" gorm:"type:decimal(10,2);not null"`
	EstimatedTotal    float64   `json:"estimated_total" gorm:"type:decimal(10,2);not null"`
	QuantityNeeded    int       `json:"quantity_needed" gorm:"type:int;not null"`
	QuantityFulfilled int       `json:"quantity_fulfilled" gorm:"type:int;default:0"`
	CreatedAt         time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt         time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	OrderItems []OrderItems `gorm:"foreignKey:ItemID;references:ItemID;constraint:onDelete:CASCADE"`
}
