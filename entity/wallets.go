package entity

import (
	"time"

	"github.com/google/uuid"
)

type Wallets struct {
	WalletID        uuid.UUID `json:"wallet_id" gorm:"type:varchar(36);primaryKey"`
	UserID          uuid.UUID `json:"user_id" gorm:"type:varchar(36)"`
	Balance         float64   `json:"balance" gorm:"type:decimal(10,2);not null"`
	LockedBalance   float64   `json:"locked_balance" gorm:"type:decimal(10,2);not null"`
	ReservedBalance float64   `json:"reserved_balance" gorm:"type:decimal(10,2);not null"`
	CreatedAt       time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	WalletTransactions []WalletTransactions `json:"wallet_transactions" gorm:"foreignKey:WalletID;references:WalletID;constraint:onDelete:CASCADE"`
}
