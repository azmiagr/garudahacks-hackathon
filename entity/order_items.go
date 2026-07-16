package entity

import (
	"time"

	"github.com/google/uuid"
)

type OrderItems struct {
	OrderItemID uuid.UUID `json:"order_item_id" gorm:"type:varchar(36);primaryKey"`
	OrderID     uuid.UUID `json:"order_id" gorm:"type:varchar(36)"`
	ItemID      uuid.UUID `json:"item_id" gorm:"type:varchar(36)"`
	Quantity    int       `json:"quantity" gorm:"type:int;not null"`
	Unit        int       `json:"unit" gorm:"type:int;not null"`
	UnitPrice   float64   `json:"unit_price" gorm:"type:decimal(10,2);not null"`
	Subtotal    float64   `json:"subtotal" gorm:"type:decimal(10,2);not null"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
