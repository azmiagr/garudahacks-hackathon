package model

import (
	"time"

	"github.com/google/uuid"
)

type DonorProfileMetricsRow struct {
	TotalDonatedAmount float64 `json:"total_donated_amount"`
	SupportedPostCount int64   `json:"supported_post_count"`
	ActivePoints       int64   `json:"active_points"`
}

type DonorProfileResponse struct {
	UserID                 uuid.UUID `json:"user_id"`
	Name                   string    `json:"name"`
	Initials               string    `json:"initials"`
	Email                  string    `json:"email"`
	PhoneNumber            string    `json:"phone_number"`
	Role                   string    `json:"role"`
	DisplayRole            string    `json:"display_role"`
	KYCStatus              string    `json:"kyc_status"`
	IsVerified             bool      `json:"is_verified"`
	VerificationText       string    `json:"verification_text"`
	MemberSince            time.Time `json:"member_since"`
	MemberSinceText        string    `json:"member_since_text"`
	Level                  string    `json:"level"`
	TotalDonatedAmount     float64   `json:"total_donated_amount"`
	TotalDonatedAmountText string    `json:"total_donated_amount_text"`
	SupportedPostCount     int64     `json:"supported_post_count"`
	ActivePoints           int64     `json:"active_points"`
}
