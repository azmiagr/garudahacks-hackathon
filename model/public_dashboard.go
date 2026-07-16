package model

import (
	"time"

	"github.com/google/uuid"
)

type PublicMapPostParam struct {
	Query        string
	MinLatitude  *float64
	MinLongitude *float64
	MaxLatitude  *float64
	MaxLongitude *float64
	Limit        int
}

type PublicMapPostRow struct {
	PostID    uuid.UUID `json:"post_id"`
	Name      string    `json:"name"`
	Address   string    `json:"address"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
}

type LatestDisasterReportParam struct {
	PostIDs      []uuid.UUID
	DisasterType string
}

type LatestDisasterReportRow struct {
	ReportID    uuid.UUID `json:"report_id"`
	PostID      uuid.UUID `json:"post_id"`
	EventID     uuid.UUID `json:"event_id"`
	ReportTitle string    `json:"report_title"`
	ImageURL    string    `json:"image_url"`
	ReportedAt  time.Time `json:"reported_at"`
	CreatedAt   time.Time `json:"created_at"`
}

type DisasterEventRow struct {
	EventID uuid.UUID `json:"event_id"`
	Name    string    `json:"name"`
}

type RequestFundingSummaryParam struct {
	ReportIDs []uuid.UUID
}

type RequestFundingSummaryRow struct {
	ReportID       uuid.UUID `json:"report_id"`
	FundingTarget  float64   `json:"funding_target"`
	FundedAmount   float64   `json:"funded_amount"`
	ReservedAmount float64   `json:"reserved_amount"`
	RequestCount   int64     `json:"request_count"`
}

type PublicDashboardParam struct {
	Query        string   `form:"q"`
	DisasterType string   `form:"disaster_type"`
	MinLatitude  *float64 `form:"min_latitude"`
	MinLongitude *float64 `form:"min_longitude"`
	MaxLatitude  *float64 `form:"max_latitude"`
	MaxLongitude *float64 `form:"max_longitude"`
	Limit        int      `form:"limit"`
}

type PublicDashboardMapResponse struct {
	Items []PublicDashboardMapItem `json:"items"`
}

type PublicDashboardMapItem struct {
	PostID            uuid.UUID `json:"post_id"`
	Name              string    `json:"name"`
	Address           string    `json:"address"`
	Latitude          float64   `json:"latitude"`
	Longitude         float64   `json:"longitude"`
	DisasterEvent     string    `json:"disaster_event"`
	LatestReportTitle string    `json:"latest_report_title"`
	ImageURL          string    `json:"image_url"`
	FundingTarget     float64   `json:"funding_target"`
	FundedAmount      float64   `json:"funded_amount"`
	FundingPercentage int       `json:"funding_percentage"`
	UrgencyLevel      string    `json:"urgency_level"`
	RequestCount      int64     `json:"request_count"`
	LatestReportedAt  time.Time `json:"latest_reported_at"`
}

type PublicDashboardSummaryResponse struct {
	ActivePoskoCount  int64                        `json:"active_posko_count"`
	TotalTarget       float64                      `json:"total_target"`
	TotalFunded       float64                      `json:"total_funded"`
	FundingPercentage int                          `json:"funding_percentage"`
	Items             []PublicDashboardSummaryItem `json:"items"`
}

type PublicDashboardSummaryItem struct {
	PostID            uuid.UUID `json:"post_id"`
	Name              string    `json:"name"`
	Address           string    `json:"address"`
	DisasterEvent     string    `json:"disaster_event"`
	LatestReportTitle string    `json:"latest_report_title"`
	FundingTarget     float64   `json:"funding_target"`
	FundedAmount      float64   `json:"funded_amount"`
	FundingPercentage int       `json:"funding_percentage"`
	UrgencyLevel      string    `json:"urgency_level"`
	RequestCount      int64     `json:"request_count"`
}

type PublicDistributionParam struct {
	Filter       string `form:"filter"`
	DisasterType string `form:"disaster_type"`
	Sort         string `form:"sort"`
	Limit        int    `form:"limit"`
}

type PublicDistributionProofRow struct {
	VerificationID     uuid.UUID `json:"verification_id"`
	OrderID            uuid.UUID `json:"order_id"`
	PostID             uuid.UUID `json:"post_id"`
	PostName           string    `json:"post_name"`
	RequestTitle       string    `json:"request_title"`
	DisasterEvent      string    `json:"disaster_event"`
	ImageURL           string    `json:"image_url"`
	VerificationStatus string    `json:"verification_status"`
	Latitude           float64   `json:"latitude"`
	Longitude          float64   `json:"longitude"`
	CapturedAt         time.Time `json:"captured_at"`
	TotalAmount        float64   `json:"total_amount"`
	DonorCount         int64     `json:"donor_count"`
	CurrentHash        string    `json:"current_hash"`
}

type PublicDistributionResponse struct {
	Items []PublicDistributionItem `json:"items"`
}

type PublicDistributionItem struct {
	VerificationID uuid.UUID `json:"verification_id"`
	OrderID        uuid.UUID `json:"order_id"`
	PostID         uuid.UUID `json:"post_id"`
	Title          string    `json:"title"`
	PostName       string    `json:"post_name"`
	DisasterEvent  string    `json:"disaster_event"`
	ImageURL       string    `json:"image_url"`
	GPSValid       bool      `json:"gps_valid"`
	Latitude       float64   `json:"latitude"`
	Longitude      float64   `json:"longitude"`
	CapturedAt     time.Time `json:"captured_at"`
	TotalAmount    float64   `json:"total_amount"`
	DonorCount     int64     `json:"donor_count"`
	AuditHash      string    `json:"audit_hash"`
}

type PublicTransparencyParam struct {
	Year  int `form:"year"`
	Limit int `form:"limit"`
}

type PublicTransparencyResponse struct {
	Summary              PublicTransparencySummary `json:"summary"`
	MonthlyDisbursements []MonthlyDisbursementItem `json:"monthly_disbursements"`
	AllocationByDisaster []DisasterAllocationItem  `json:"allocation_by_disaster"`
	LatestLedger         []PublicLedgerItem        `json:"latest_ledger"`
}

type PublicTransparencySummary struct {
	TotalDonationCollected  float64 `json:"total_donation_collected"`
	TotalDisbursedVerified  float64 `json:"total_disbursed_verified"`
	RefundAutomatic         float64 `json:"refund_automatic"`
	VerifiedFulfillmentRate int     `json:"verified_fulfillment_rate"`
}

type MonthlyDisbursementRow struct {
	Month int     `json:"month"`
	Total float64 `json:"total"`
}

type MonthlyDisbursementItem struct {
	Month string  `json:"month"`
	Total float64 `json:"total"`
}

type DisasterAllocationRow struct {
	DisasterEvent string  `json:"disaster_event"`
	TotalAmount   float64 `json:"total_amount"`
}

type DisasterAllocationItem struct {
	DisasterEvent string  `json:"disaster_event"`
	TotalAmount   float64 `json:"total_amount"`
	Percentage    int     `json:"percentage"`
}

type PublicLedgerRow struct {
	OccurredAt time.Time `json:"occurred_at"`
	Event      string    `json:"event"`
	PostName   string    `json:"post_name"`
	ValueLabel string    `json:"value_label"`
	Hash       string    `json:"hash"`
}

type PublicLedgerItem struct {
	OccurredAt time.Time `json:"occurred_at"`
	Event      string    `json:"event"`
	PostName   string    `json:"post_name"`
	ValueLabel string    `json:"value_label"`
	Hash       string    `json:"hash"`
}

type DonationTransparencySummaryRow struct {
	TotalCollected  float64 `json:"total_collected"`
	RefundAutomatic float64 `json:"refund_automatic"`
}

type VerifiedFulfillmentRateRow struct {
	TotalOrders    int64 `json:"total_orders"`
	VerifiedOrders int64 `json:"verified_orders"`
}
