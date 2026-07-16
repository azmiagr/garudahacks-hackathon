package model

import (
	"time"

	"github.com/google/uuid"
)

type DonorDashboardMapParam struct {
	Query        string   `form:"q"`
	DisasterType string   `form:"disaster_type"`
	MinLatitude  *float64 `form:"min_latitude"`
	MinLongitude *float64 `form:"min_longitude"`
	MaxLatitude  *float64 `form:"max_latitude"`
	MaxLongitude *float64 `form:"max_longitude"`
	Limit        int      `form:"limit"`
}

type DonorDashboardMapResponse struct {
	HeatmapPoints []DonorHeatmapPoint `json:"heatmap_points"`
	UrgentPosts   []DonorUrgentPost   `json:"urgent_posts"`
	Legend        []DonorMapLegend    `json:"legend"`
}

type DonorHeatmapPoint struct {
	PostID            uuid.UUID `json:"post_id"`
	Name              string    `json:"name"`
	Latitude          float64   `json:"latitude"`
	Longitude         float64   `json:"longitude"`
	DisasterType      string    `json:"disaster_type"`
	FundingPercentage int       `json:"funding_percentage"`
	UrgencyLevel      string    `json:"urgency_level"`
	Color             string    `json:"color"`
}

type DonorUrgentPost struct {
	PostID            uuid.UUID `json:"post_id"`
	Name              string    `json:"name"`
	Address           string    `json:"address"`
	ImageURL          string    `json:"image_url"`
	FundingTarget     float64   `json:"funding_target"`
	FundedAmount      float64   `json:"funded_amount"`
	FundingPercentage int       `json:"funding_percentage"`
	FundingText       string    `json:"funding_text"`
	ElapsedText       string    `json:"elapsed_text"`
}

type DonorMapLegend struct {
	Label string `json:"label"`
	Level string `json:"level"`
	Color string `json:"color"`
}

type DonorPostDetailParam struct {
	PostID uuid.UUID `uri:"post_id" binding:"required"`
}

type DonorPostDetailRow struct {
	PostID         uuid.UUID `json:"post_id"`
	ReportID       uuid.UUID `json:"report_id"`
	RequestID      uuid.UUID `json:"request_id"`
	Name           string    `json:"name"`
	Address        string    `json:"address"`
	DisasterType   string    `json:"disaster_type"`
	ImageURL       string    `json:"image_url"`
	CreatedAt      time.Time `json:"created_at"`
	ReportedAt     time.Time `json:"reported_at"`
	FundingTarget  float64   `json:"funding_target"`
	FundedAmount   float64   `json:"funded_amount"`
	DonorCount     int64     `json:"donor_count"`
	AdminKYCStatus string    `json:"admin_kyc_status"`
}

type DonorPostDetailItemRow struct {
	ItemID            uuid.UUID `json:"item_id"`
	Name              string    `json:"name"`
	Description       string    `json:"description"`
	Price             float64   `json:"price"`
	EstimatedTotal    float64   `json:"estimated_total"`
	QuantityNeeded    int       `json:"quantity_needed"`
	QuantityFulfilled int       `json:"quantity_fulfilled"`
}

type DonorPostDetailResponse struct {
	PostID                uuid.UUID             `json:"post_id"`
	ReportID              uuid.UUID             `json:"report_id"`
	RequestID             uuid.UUID             `json:"request_id"`
	Name                  string                `json:"name"`
	Address               string                `json:"address"`
	DisasterType          string                `json:"disaster_type"`
	ImageURL              string                `json:"image_url"`
	ElapsedText           string                `json:"elapsed_text"`
	AdminVerified         bool                  `json:"admin_verified"`
	AdminVerificationText string                `json:"admin_verification_text"`
	FundingTarget         float64               `json:"funding_target"`
	FundedAmount          float64               `json:"funded_amount"`
	FundingPercentage     int                   `json:"funding_percentage"`
	FundingText           string                `json:"funding_text"`
	DonorCount            int64                 `json:"donor_count"`
	UrgencyLevel          string                `json:"urgency_level"`
	Items                 []DonorPostDetailItem `json:"items"`
}

type DonorPostDetailItem struct {
	ItemID            uuid.UUID `json:"item_id"`
	Name              string    `json:"name"`
	Description       string    `json:"description"`
	Price             float64   `json:"price"`
	EstimatedTotal    float64   `json:"estimated_total"`
	QuantityNeeded    int       `json:"quantity_needed"`
	QuantityFulfilled int       `json:"quantity_fulfilled"`
	ProgressText      string    `json:"progress_text"`
}
