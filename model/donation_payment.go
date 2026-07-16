package model

import (
	"time"

	"github.com/google/uuid"
)

type DonationLockContextRow struct {
	RequestID   uuid.UUID `json:"request_id"`
	PostID      uuid.UUID `json:"post_id"`
	AdminUserID uuid.UUID `json:"admin_user_id"`
	PostName    string    `json:"post_name"`
	Latitude    float64   `json:"latitude"`
	Longitude   float64   `json:"longitude"`
}

type CreateDonationPaymentRequest struct {
	RequestID     uuid.UUID `json:"request_id" binding:"required"`
	Amount        int64     `json:"amount" binding:"required,min=10000"`
	PaymentMethod string    `json:"payment_method" binding:"required"` // qris, virtual_account
	Bank          string    `json:"bank"`                              // bca, bni, bri, mandiri, permata
	Autonomous    bool      `json:"autonomous"`
}

type CreateDonationPaymentResponse struct {
	OrderID              string     `json:"order_id"`
	DonationID           uuid.UUID  `json:"donation_id"`
	PaymentTransactionID uuid.UUID  `json:"payment_transaction_id"`
	RequestID            uuid.UUID  `json:"request_id"`
	Amount               int64      `json:"amount"`
	PaymentMethod        string     `json:"payment_method"`
	PaymentChannel       string     `json:"payment_channel"`
	TransactionStatus    string     `json:"transaction_status"`
	QRString             string     `json:"qr_string"`
	QRURL                string     `json:"qr_url"`
	VANumber             string     `json:"va_number"`
	VABank               string     `json:"va_bank"`
	PermataVANumber      string     `json:"permata_va_number"`
	ExpiredAt            *time.Time `json:"expired_at"`
}

type MidtransNotificationRequest struct {
	TransactionStatus string             `json:"transaction_status"`
	TransactionID     string             `json:"transaction_id"`
	StatusCode        string             `json:"status_code"`
	SignatureKey      string             `json:"signature_key"`
	PaymentType       string             `json:"payment_type"`
	OrderID           string             `json:"order_id"`
	GrossAmount       string             `json:"gross_amount"`
	FraudStatus       string             `json:"fraud_status"`
	SettlementTime    string             `json:"settlement_time"`
	Bank              string             `json:"bank"`
	PermataVANumber   string             `json:"permata_va_number"`
	VANumbers         []MidtransVANumber `json:"va_numbers"`
}

type MidtransVANumber struct {
	Bank     string `json:"bank"`
	VANumber string `json:"va_number"`
}

type DonationLockStatusResponse struct {
	PaymentOrderID  string    `json:"payment_order_id"`
	DonationID      uuid.UUID `json:"donation_id"`
	LockedOrderID   uuid.UUID `json:"locked_order_id"`
	TransactionCode string    `json:"transaction_code"`
	PostName        string    `json:"post_name"`
	Amount          int64     `json:"amount"`
	FundStatus      string    `json:"fund_status"` // LOCKED
	AllocationText  string    `json:"allocation_text"`
	LedgerHash      string    `json:"ledger_hash"`
	ShortLedgerHash string    `json:"short_ledger_hash"`
	Processed       bool      `json:"processed"`
}
