package model

import "github.com/google/uuid"

type AdminProfileMetricsRow struct {
	EventCount         int64   `json:"event_count"`
	ManagedAidAmount   float64 `json:"managed_aid_amount"`
	TotalOrderCount    int64   `json:"total_order_count"`
	VerifiedOrderCount int64   `json:"verified_order_count"`
}

type AdminProfileResponse struct {
	UserID                  uuid.UUID `json:"user_id"`
	Name                    string    `json:"name"`
	Initials                string    `json:"initials"`
	Role                    string    `json:"role"`
	DisplayRole             string    `json:"display_role"`
	Affiliation             string    `json:"affiliation"`
	KYCStatus               string    `json:"kyc_status"`
	IsVerified              bool      `json:"is_verified"`
	VerificationText        string    `json:"verification_text"`
	SuccessfulEventsText    string    `json:"successful_events_text"`
	EventCount              int64     `json:"event_count"`
	ManagedAidAmount        float64   `json:"managed_aid_amount"`
	ManagedAidAmountText    string    `json:"managed_aid_amount_text"`
	VerifiedOrderPercentage float64   `json:"verified_order_percentage"`
	VerifiedOrderText       string    `json:"verified_order_text"`
}
