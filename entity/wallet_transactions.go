package entity

import (
	"time"

	"github.com/google/uuid"
)

type WalletTransactions struct {
	WalletTransactionID uuid.UUID `json:"wallet_transaction_id" gorm:"type:varchar(36);primaryKey"`
	WalletID            uuid.UUID `json:"wallet_id" gorm:"type:varchar(36)"`
	Amount              float64   `json:"amount" gorm:"type:decimal(10,2);not null"`
	BalanceBefore       float64   `json:"balance_before" gorm:"type:decimal(10,2);not null"`
	BalanceAfter        float64   `json:"balance_after" gorm:"type:decimal(10,2);not null"`
	TransactionType     string    `json:"transaction_type" gorm:"type:enum('deposit','withdrawal');default:'deposit'"`
	TransactionStatus   string    `json:"transaction_status" gorm:"type:enum('pending','approved','rejected');default:'pending'"`
	CreatedAt           time.Time `json:"created_at" gorm:"autoCreateTime"`

	Donations []Donations `gorm:"foreignKey:WalletTransactionID;references:WalletTransactionID;constraint:onDelete:CASCADE"`
}
