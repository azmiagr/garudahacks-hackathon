package entity

import (
	"time"

	"github.com/google/uuid"
)

type Stores struct {
	StoreID        uuid.UUID `json:"store_id" gorm:"type:varchar(36);primaryKey"`
	OwnerID        uuid.UUID `json:"owner_id" gorm:"type:varchar(36)"`
	Name           string    `json:"name" gorm:"type:varchar(255);not null"`
	BusinessNumber string    `json:"business_number" gorm:"type:varchar(255);not null;uniqueIndex"`
	Address        string    `json:"address" gorm:"type:text;not null"`
	Latitude       float64   `json:"latitude" gorm:"type:decimal(10,2);not null"`
	Longitude      float64   `json:"longitude" gorm:"type:decimal(10,2);not null"`
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	Orders        []Orders        `gorm:"foreignKey:StoreID;references:StoreID;constraint:onDelete:CASCADE"`
	Disbursements []Disbursements `gorm:"foreignKey:StoreID;references:StoreID;constraint:onDelete:CASCADE"`
}
