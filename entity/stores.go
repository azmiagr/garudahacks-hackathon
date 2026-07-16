package entity

import (
	"time"

	"github.com/google/uuid"
)

type Stores struct {
	StoreID         uuid.UUID `json:"store_id" gorm:"type:varchar(36);primaryKey"`
	OwnerID         uuid.UUID `json:"owner_id" gorm:"type:varchar(36)"`
	Name            string    `json:"name" gorm:"type:varchar(255);not null"`
	BusinessNumber  string    `json:"business_number" gorm:"type:varchar(255);not null;uniqueIndex"` // NIB
	NPWP            string    `json:"npwp" gorm:"type:varchar(40)"`
	KTPImageURL     string    `json:"ktp_image_url" gorm:"type:varchar(255)"`
	BankName        string    `json:"bank_name" gorm:"type:varchar(40)"`
	BankAccountNo   string    `json:"bank_account_no" gorm:"type:varchar(80)"`
	BankAccountName string    `json:"bank_account_name" gorm:"type:varchar(150)"`
	CategoriesJSON  string    `json:"categories_json" gorm:"type:json"`
	Address         string    `json:"address" gorm:"type:text;not null"`
	Latitude        float64   `json:"latitude" gorm:"type:decimal(10,8);not null"`
	Longitude       float64   `json:"longitude" gorm:"type:decimal(11,8);not null"`
	CreatedAt       time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	Orders        []Orders        `gorm:"foreignKey:StoreID;references:StoreID;constraint:onDelete:CASCADE"`
	Disbursements []Disbursements `gorm:"foreignKey:StoreID;references:StoreID;constraint:onDelete:CASCADE"`
}
