package model

import (
	"time"

	"github.com/google/uuid"
)

type DonorDonationTransactionListParam struct {
	UserID uuid.UUID
	Status string // "", "all", "pending", "locked", "preparing", "ready", "shipping", "completed", "refund"
	Limit  int
	Offset int
}

type DonorDonationTransactionDetailParam struct {
	UserID     uuid.UUID
	DonationID uuid.UUID
}

type DonorDonationTransactionListRow struct {
	DonationID           uuid.UUID  `json:"donation_id"`
	PaymentTransactionID uuid.UUID  `json:"payment_transaction_id"`
	PaymentOrderID       string     `json:"payment_order_id"`
	RequestID            uuid.UUID  `json:"request_id"`
	LockedOrderID        string     `json:"locked_order_id"`
	TransactionCode      string     `json:"transaction_code"`
	PostName             string     `json:"post_name"`
	RequestTitle         string     `json:"request_title"`
	Amount               float64    `json:"amount"`
	DonationStatus       string     `json:"donation_status"`
	PaymentStatus        string     `json:"payment_status"`
	OrderStatus          string     `json:"order_status"`
	LatestHash           string     `json:"latest_hash"`
	VerificationImageURL string     `json:"verification_image_url"`
	CustodyStepCount     int        `json:"custody_step_count"`
	DonatedAt            time.Time  `json:"donated_at"`
	PaidAt               *time.Time `json:"paid_at"`
	VerifiedAt           *time.Time `json:"verified_at"`
}

type DonorDonationTransactionDetailRow struct {
	DonorDonationTransactionListRow
	PostAddress    string  `json:"post_address"`
	Latitude       float64 `json:"latitude"`
	Longitude      float64 `json:"longitude"`
	FundingTarget  float64 `json:"funding_target"`
	FundedAmount   float64 `json:"funded_amount"`
	DonorCount     int64   `json:"donor_count"`
	TotalItemCount int64   `json:"total_item_count"`
}

type DonorDonationTransactionItemRow struct {
	ItemID    uuid.UUID `json:"item_id"`
	Name      string    `json:"name"`
	Quantity  int       `json:"quantity"`
	UnitPrice float64   `json:"unit_price"`
	Subtotal  float64   `json:"subtotal"`
}

type DonorDonationTransactionCustodyLogRow struct {
	LogsID      string    `json:"logs_id"`
	OrderID     string    `json:"order_id"`
	Sequence    int       `json:"sequence"`
	FromActorID uuid.UUID `json:"from_actor_id"`
	ToActorID   uuid.UUID `json:"to_actor_id"`
	Latitude    float64   `json:"latitude"`
	Longitude   float64   `json:"longitude"`
	PrevHash    string    `json:"prev_hash"`
	CurrentHash string    `json:"current_hash"`
	CreatedAt   time.Time `json:"created_at"`
}

type DonorDonationTransactionParam struct {
	Status string `form:"status"` // all, pending, locked, preparing, ready, shipping, completed, refund
	Limit  int    `form:"limit"`
	Offset int    `form:"offset"`
}

type DonorDonationTransactionListResponse struct {
	Items  []DonorDonationTransactionListItem `json:"items"`
	Total  int64                              `json:"total"`
	Limit  int                                `json:"limit"`
	Offset int                                `json:"offset"`
}

type DonorDonationTransactionListItem struct {
	DonationID           uuid.UUID  `json:"donation_id"`
	PaymentTransactionID uuid.UUID  `json:"payment_transaction_id"`
	PaymentOrderID       string     `json:"payment_order_id"`
	RequestID            uuid.UUID  `json:"request_id"`
	LockedOrderID        string     `json:"locked_order_id"`
	TransactionCode      string     `json:"transaction_code"`
	PostName             string     `json:"post_name"`
	RequestTitle         string     `json:"request_title"`
	Amount               int64      `json:"amount"`
	AmountText           string     `json:"amount_text"`
	Status               string     `json:"status"`
	StatusLabel          string     `json:"status_label"`
	BadgeVariant         string     `json:"badge_variant"`
	LatestHash           string     `json:"latest_hash"`
	ShortLatestHash      string     `json:"short_latest_hash"`
	VerificationImageURL string     `json:"verification_image_url"`
	CustodyStepCount     int        `json:"custody_step_count"`
	ProgressText         string     `json:"progress_text"`
	DonatedAt            time.Time  `json:"donated_at"`
	PaidAt               *time.Time `json:"paid_at"`
	VerifiedAt           *time.Time `json:"verified_at"`
	ElapsedText          string     `json:"elapsed_text"`
}

type DonorDonationTransactionDetailResponse struct {
	DonorDonationTransactionListItem
	PostAddress       string                                   `json:"post_address"`
	Latitude          float64                                  `json:"latitude"`
	Longitude         float64                                  `json:"longitude"`
	FundingTarget     float64                                  `json:"funding_target"`
	FundedAmount      float64                                  `json:"funded_amount"`
	FundingPercentage int                                      `json:"funding_percentage"`
	FundingText       string                                   `json:"funding_text"`
	DonorCount        int64                                    `json:"donor_count"`
	TotalItemCount    int64                                    `json:"total_item_count"`
	Items             []DonorDonationTransactionItem           `json:"items"`
	CustodyLogs       []DonorDonationTransactionCustodyLogItem `json:"custody_logs"`
}

type DonorDonationTransactionItem struct {
	ItemID     uuid.UUID `json:"item_id"`
	Name       string    `json:"name"`
	Quantity   int       `json:"quantity"`
	UnitPrice  float64   `json:"unit_price"`
	Subtotal   float64   `json:"subtotal"`
	AmountText string    `json:"amount_text"`
}

type DonorDonationTransactionCustodyLogItem struct {
	LogsID           string    `json:"logs_id"`
	OrderID          string    `json:"order_id"`
	Sequence         int       `json:"sequence"`
	FromActorID      uuid.UUID `json:"from_actor_id"`
	ToActorID        uuid.UUID `json:"to_actor_id"`
	Latitude         float64   `json:"latitude"`
	Longitude        float64   `json:"longitude"`
	PrevHash         string    `json:"prev_hash"`
	CurrentHash      string    `json:"current_hash"`
	ShortCurrentHash string    `json:"short_current_hash"`
	CreatedAt        time.Time `json:"created_at"`
	ElapsedText      string    `json:"elapsed_text"`
}
