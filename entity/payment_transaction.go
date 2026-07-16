package entity

import (
	"time"

	"github.com/google/uuid"
)

type PaymentTransactions struct {
	PaymentTransactionID   uuid.UUID  `json:"payment_transaction_id" gorm:"type:varchar(36);primaryKey"`
	OrderID                string     `json:"order_id" gorm:"type:varchar(80);not null;uniqueIndex"`
	UserID                 uuid.UUID  `json:"user_id" gorm:"type:varchar(36);not null;index"`
	RequestID              uuid.UUID  `json:"request_id" gorm:"type:varchar(36);not null;index"`
	DonationID             uuid.UUID  `json:"donation_id" gorm:"type:varchar(36);not null;index"`
	WalletTransactionID    uuid.UUID  `json:"wallet_transaction_id" gorm:"type:varchar(36);not null;index"`
	Amount                 float64    `json:"amount" gorm:"type:decimal(10,2);not null"`
	PaymentMethod          string     `json:"payment_method" gorm:"type:varchar(40);not null"`
	PaymentChannel         string     `json:"payment_channel" gorm:"type:varchar(40)"`
	TransactionStatus      string     `json:"transaction_status" gorm:"type:varchar(40);not null;default:'pending'"`
	FraudStatus            string     `json:"fraud_status" gorm:"type:varchar(40)"`
	MidtransStatusCode     string     `json:"midtrans_status_code" gorm:"type:varchar(10)"`
	QRString               string     `json:"qr_string" gorm:"type:text"`
	QRURL                  string     `json:"qr_url" gorm:"type:text"`
	VANumber               string     `json:"va_number" gorm:"type:varchar(80)"`
	VABank                 string     `json:"va_bank" gorm:"type:varchar(40)"`
	PermataVANumber        string     `json:"permata_va_number" gorm:"type:varchar(80)"`
	RawChargeResponse      *string    `json:"raw_charge_response" gorm:"type:json"`
	RawNotificationPayload *string    `json:"raw_notification_payload" gorm:"type:json"`
	PaidAt                 *time.Time `json:"paid_at"`
	ExpiredAt              *time.Time `json:"expired_at"`
	CreatedAt              time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt              time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}
