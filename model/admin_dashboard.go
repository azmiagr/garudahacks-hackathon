package model

import (
	"time"

	"github.com/google/uuid"
)

type AdminDashboardHomeParam struct {
	UserID uuid.UUID
}

type AdminDashboardHomeResponse struct {
	GreetingName     string                    `json:"greeting_name"`
	IsAdminVerified  bool                      `json:"is_admin_verified"`
	VerificationText string                    `json:"verification_text"`
	ActiveEvents     []AdminDashboardEventCard `json:"active_events"`
	ClosedEvents     []AdminDashboardEventCard `json:"closed_events"`
}

type AdminDashboardEventCard struct {
	PostID                uuid.UUID                    `json:"post_id"`
	EventCode             string                       `json:"event_code"`
	Title                 string                       `json:"title"`
	DisasterType          string                       `json:"disaster_type"`
	Status                string                       `json:"status"`
	StatusLabel           string                       `json:"status_label"`
	ImageURL              string                       `json:"image_url"`
	Address               string                       `json:"address"`
	GeofenceRadius        float64                      `json:"geofence_radius"`
	AffectedHouseholds    int                          `json:"affected_households"`
	StartedAt             time.Time                    `json:"started_at"`
	ElapsedText           string                       `json:"elapsed_text"`
	FundingTarget         float64                      `json:"funding_target"`
	FundedAmount          float64                      `json:"funded_amount"`
	FundingPercentage     float64                      `json:"funding_percentage"`
	FundingText           string                       `json:"funding_text"`
	OrderCount            int64                        `json:"order_count"`
	SummaryText           string                       `json:"summary_text"`
	CanScanCourierQR      bool                         `json:"can_scan_courier_qr"`
	CanAddFollowUpRequest bool                         `json:"can_add_follow_up_request"`
	LatestOrders          []AdminDashboardOrderPreview `json:"latest_orders"`
}

type AdminDashboardOrderPreview struct {
	OrderID      uuid.UUID `json:"order_id"`
	OrderCode    string    `json:"order_code"`
	StoreName    string    `json:"store_name"`
	CourierName  string    `json:"courier_name"`
	Status       string    `json:"status"`
	StatusLabel  string    `json:"status_label"`
	Description  string    `json:"description"`
	BadgeVariant string    `json:"badge_variant"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type AdminDashboardEventRow struct {
	PostID              uuid.UUID `json:"post_id"`
	Title               string    `json:"title"`
	Address             string    `json:"address"`
	GeofenceRadius      float64   `json:"geofence_radius"`
	DisasterType        string    `json:"disaster_type"`
	ImageURL            string    `json:"image_url"`
	StartedAt           time.Time `json:"started_at"`
	FundingTarget       float64   `json:"funding_target"`
	FundedAmount        float64   `json:"funded_amount"`
	OrderCount          int64     `json:"order_count"`
	CompletedOrderCount int64     `json:"completed_order_count"`
}

type AdminDashboardOrderRow struct {
	PostID             uuid.UUID `json:"post_id"`
	OrderID            uuid.UUID `json:"order_id"`
	OrderCode          string    `json:"order_code"`
	OrderStatus        string    `json:"order_status"`
	StoreName          string    `json:"store_name"`
	CourierName        string    `json:"courier_name"`
	UpdatedAt          time.Time `json:"updated_at"`
	VerificationStatus string    `json:"verification_status"`
}
