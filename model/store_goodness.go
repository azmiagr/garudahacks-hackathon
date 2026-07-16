package model

import (
	"time"

	"github.com/google/uuid"
)

type StoreGoodnessParam struct {
	StoreID uuid.UUID
	Year    int
	Limit   int
	Offset  int
}

type StoreGoodnessCertificateRow struct {
	StoreID             uuid.UUID  `json:"store_id"`
	StoreName           string     `json:"store_name"`
	BankName            string     `json:"bank_name"`
	VerifiedOrderCount  int64      `json:"verified_order_count"`
	VerifiedAmountTotal float64    `json:"verified_amount_total"`
	ReputationScore     float64    `json:"reputation_score"`
	DisputeCount        int64      `json:"dispute_count"`
	FirstContributionAt *time.Time `json:"first_contribution_at"`
}

type StoreContributionHistoryRow struct {
	OrderID      uuid.UUID `json:"order_id"`
	OrderCode    string    `json:"order_code"`
	PostName     string    `json:"post_name"`
	DisasterName string    `json:"disaster_name"`
	ItemCount    int64     `json:"item_count"`
	TotalAmount  float64   `json:"total_amount"`
	VerifiedAt   time.Time `json:"verified_at"`
	LatestHash   string    `json:"latest_hash"`
}

type StoreGoodnessRequest struct {
	Year   int `form:"year"`
	Limit  int `form:"limit"`
	Offset int `form:"offset"`
}

type StoreGoodnessResponse struct {
	Certificate  StoreGoodnessCertificate       `json:"certificate"`
	History      []StoreContributionHistoryItem `json:"history"`
	TotalHistory int64                          `json:"total_history"`
	Limit        int                            `json:"limit"`
	Offset       int                            `json:"offset"`
}

type StoreGoodnessCertificate struct {
	StoreID             uuid.UUID  `json:"store_id"`
	StoreName           string     `json:"store_name"`
	Title               string     `json:"title"`
	PartnerLabel        string     `json:"partner_label"`
	SinceText           string     `json:"since_text"`
	VerifiedOrderCount  int64      `json:"verified_order_count"`
	VerifiedOrderText   string     `json:"verified_order_text"`
	VerifiedAmountTotal float64    `json:"verified_amount_total"`
	VerifiedAmountText  string     `json:"verified_amount_text"`
	ReputationScore     float64    `json:"reputation_score"`
	ReputationText      string     `json:"reputation_text"`
	DisputeCount        int64      `json:"dispute_count"`
	DisputeText         string     `json:"dispute_text"`
	FirstContributionAt *time.Time `json:"first_contribution_at"`
	ShareURL            string     `json:"share_url"`
}

type StoreContributionHistoryItem struct {
	OrderID         uuid.UUID `json:"order_id"`
	OrderCode       string    `json:"order_code"`
	PostName        string    `json:"post_name"`
	DisasterName    string    `json:"disaster_name"`
	Title           string    `json:"title"`
	ItemCount       int64     `json:"item_count"`
	TotalAmount     float64   `json:"total_amount"`
	TotalAmountText string    `json:"total_amount_text"`
	VerifiedAt      time.Time `json:"verified_at"`
	VerifiedAtText  string    `json:"verified_at_text"`
	LatestHash      string    `json:"latest_hash"`
	ShortLatestHash string    `json:"short_latest_hash"`
}
