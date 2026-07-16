package model

import "github.com/google/uuid"

type StoreProfileStatsRow struct {
	TotalOrder30Days     int64   `json:"total_order_30_days"`
	AcceptedOrder30Days  int64   `json:"accepted_order_30_days"`
	CancelledOrder30Days int64   `json:"cancelled_order_30_days"`
	ReputationScore      float64 `json:"reputation_score"`
}

type StoreProfileResponse struct {
	StoreID              uuid.UUID `json:"store_id"`
	OwnerID              uuid.UUID `json:"owner_id"`
	Name                 string    `json:"name"`
	Address              string    `json:"address"`
	Latitude             float64   `json:"latitude"`
	Longitude            float64   `json:"longitude"`
	IsOnline             bool      `json:"is_online"`
	StoreStatus          string    `json:"store_status"`
	KYCStatus            string    `json:"kyc_status"`
	KYCLabel             string    `json:"kyc_label"`
	ReputationScore      float64   `json:"reputation_score"`
	BusinessNumber       string    `json:"business_number"`
	NPWP                 string    `json:"npwp"`
	KTPImageURL          string    `json:"ktp_image_url"`
	BankName             string    `json:"bank_name"`
	BankAccountNo        string    `json:"bank_account_no"`
	MaskedBankAccount    string    `json:"masked_bank_account"`
	BankAccountName      string    `json:"bank_account_name"`
	Categories           []string  `json:"categories"`
	CategoriesText       string    `json:"categories_text"`
	TotalOrder30Days     int64     `json:"total_order_30_days"`
	AcceptedOrder30Days  int64     `json:"accepted_order_30_days"`
	CancelledOrder30Days int64     `json:"cancelled_order_30_days"`
	AcceptanceRate30Days float64   `json:"acceptance_rate_30_days"`
	AcceptanceSummary    string    `json:"acceptance_summary"`
	CreatedAt            string    `json:"created_at"`
	UpdatedAt            string    `json:"updated_at"`
}
