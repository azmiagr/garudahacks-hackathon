package model

import (
	"time"

	"github.com/google/uuid"
)

type StoreDisbursementDashboardParam struct {
	StoreID uuid.UUID
	Year    int
	Month   int
	Limit   int
	Offset  int
}

type StoreDisbursementSummaryRow struct {
	StoreID                 uuid.UUID `json:"store_id"`
	StoreName               string    `json:"store_name"`
	BankName                string    `json:"bank_name"`
	BankAccountNo           string    `json:"bank_account_no"`
	TotalDisbursedThisMonth float64   `json:"total_disbursed_this_month"`
	CompletedOrderCount     int64     `json:"completed_order_count"`
	DisputeCount            int64     `json:"dispute_count"`
	MedianDisbursementMin   float64   `json:"median_disbursement_min"`
}

type StoreDisbursementHistoryRow struct {
	DisbursementID           uuid.UUID  `json:"disbursement_id"`
	OrderID                  uuid.UUID  `json:"order_id"`
	OrderCode                string     `json:"order_code"`
	PostName                 string     `json:"post_name"`
	Amount                   float64    `json:"amount"`
	Status                   string     `json:"status"`
	IdempotencyKey           string     `json:"idempotency_key"`
	GatewayReference         string     `json:"gateway_reference"`
	GatewayAttempt           int        `json:"gateway_attempt"`
	VerificationApprovedAt   *time.Time `json:"verification_approved_at"`
	DisbursedAt              *time.Time `json:"disbursed_at"`
	CreatedAt                time.Time  `json:"created_at"`
	UpdatedAt                time.Time  `json:"updated_at"`
	MinutesAfterVerification float64    `json:"minutes_after_verification"`
}

type StoreGoodnessTrailRow struct {
	StoreID             uuid.UUID  `json:"store_id"`
	VerifiedOrderCount  int64      `json:"verified_order_count"`
	VerifiedAmountTotal float64    `json:"verified_amount_total"`
	FirstContributionAt *time.Time `json:"first_contribution_at"`
	LastContributionAt  *time.Time `json:"last_contribution_at"`
}

type StoreDisbursementDashboardRequest struct {
	Year   int `form:"year"`
	Month  int `form:"month"`
	Limit  int `form:"limit"`
	Offset int `form:"offset"`
}

type StoreDisbursementDashboardResponse struct {
	Summary       StoreDisbursementSummary       `json:"summary"`
	History       []StoreDisbursementHistoryItem `json:"history"`
	GoodnessTrail StoreGoodnessTrail             `json:"goodness_trail"`
	TotalHistory  int64                          `json:"total_history"`
	Limit         int                            `json:"limit"`
	Offset        int                            `json:"offset"`
}

type StoreDisbursementSummary struct {
	StoreID                 uuid.UUID `json:"store_id"`
	StoreName               string    `json:"store_name"`
	TotalDisbursedThisMonth float64   `json:"total_disbursed_this_month"`
	TotalDisbursedText      string    `json:"total_disbursed_text"`
	CompletedOrderCount     int64     `json:"completed_order_count"`
	DisputeCount            int64     `json:"dispute_count"`
	BankName                string    `json:"bank_name"`
	MaskedBankAccount       string    `json:"masked_bank_account"`
	MedianDisbursementMin   int       `json:"median_disbursement_min"`
	MedianDisbursementText  string    `json:"median_disbursement_text"`
	Subtitle                string    `json:"subtitle"`
}

type StoreDisbursementHistoryItem struct {
	DisbursementID               uuid.UUID  `json:"disbursement_id"`
	OrderID                      uuid.UUID  `json:"order_id"`
	OrderCode                    string     `json:"order_code"`
	PostName                     string     `json:"post_name"`
	Amount                       float64    `json:"amount"`
	AmountText                   string     `json:"amount_text"`
	Status                       string     `json:"status"`
	StatusLabel                  string     `json:"status_label"`
	BadgeVariant                 string     `json:"badge_variant"`
	IdempotencyKey               string     `json:"idempotency_key"`
	GatewayReference             string     `json:"gateway_reference"`
	GatewayAttempt               int        `json:"gateway_attempt"`
	VerificationApprovedAt       *time.Time `json:"verification_approved_at"`
	DisbursedAt                  *time.Time `json:"disbursed_at"`
	CreatedAt                    time.Time  `json:"created_at"`
	TimelineText                 string     `json:"timeline_text"`
	MinutesAfterVerification     int        `json:"minutes_after_verification"`
	MinutesAfterVerificationText string     `json:"minutes_after_verification_text"`
}

type StoreGoodnessTrail struct {
	StoreID             uuid.UUID  `json:"store_id"`
	VerifiedOrderCount  int64      `json:"verified_order_count"`
	VerifiedAmountTotal float64    `json:"verified_amount_total"`
	VerifiedAmountText  string     `json:"verified_amount_text"`
	Year                int        `json:"year"`
	SummaryText         string     `json:"summary_text"`
	FirstContributionAt *time.Time `json:"first_contribution_at"`
	LastContributionAt  *time.Time `json:"last_contribution_at"`
}
