package entity

import (
	"time"

	"github.com/google/uuid"
)

type Donations struct {
	DonationID          uuid.UUID `json:"donation_id" gorm:"type:varchar(36);primaryKey"`
	RequestID           uuid.UUID `json:"request_id" gorm:"type:varchar(36)"`
	DonatedBy           uuid.UUID `json:"donated_by" gorm:"type:varchar(36)"`
	WalletTransactionID uuid.UUID `json:"wallet_transaction_id" gorm:"type:varchar(36)"`
	DonationAmount      float64   `json:"donation_amount" gorm:"type:decimal(10,2);not null"`
	DonationStatus      string    `json:"donation_status" gorm:"type:enum('pending','approved','rejected');default:'pending'"`
	DonatedAt           time.Time `json:"donated_at" gorm:"autoCreateTime"`
	UpdatedAt           time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
