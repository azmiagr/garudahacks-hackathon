package model

import (
	"time"

	"github.com/google/uuid"
)

type StoreOrderActionRequest struct {
	OrderID uuid.UUID `json:"order_id" binding:"required"`
}

type StoreOrderActionResponse struct {
	OrderID     uuid.UUID `json:"order_id"`
	StoreID     uuid.UUID `json:"store_id"`
	OrderStatus string    `json:"order_status"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type StoreHandoffTokenResponse struct {
	OrderID              uuid.UUID `json:"order_id"`
	TokenID              uuid.UUID `json:"token_id"`
	HandoffStage         string    `json:"handoff_stage"`
	QRPayload            string    `json:"qr_payload"`
	FallbackPIN          string    `json:"fallback_pin"`
	ExpiresAt            time.Time `json:"expires_at"`
	CacheValidUntil      time.Time `json:"cache_valid_until"`
	RefreshInSeconds     int       `json:"refresh_in_seconds"`
	CacheWindowInSeconds int       `json:"cache_window_in_seconds"`
}

type SubmitCustodyHandshakeRequest struct {
	OrderID        uuid.UUID `json:"order_id"`
	Method         string    `json:"method" binding:"required"`
	QRPayload      string    `json:"qr_payload"`
	FallbackPIN    string    `json:"fallback_pin"`
	Latitude       float64   `json:"latitude" binding:"required"`
	Longitude      float64   `json:"longitude" binding:"required"`
	IdempotencyKey string    `json:"idempotency_key"`
	CapturedAt     time.Time `json:"captured_at"`
}

type CustodyHandshakeResponse struct {
	OrderID          uuid.UUID `json:"order_id"`
	LogID            string    `json:"log_id"`
	OrderStatus      string    `json:"order_status"`
	HandoffStage     string    `json:"handoff_stage"`
	HandshakeMethod  string    `json:"handshake_method"`
	Sequence         int       `json:"sequence"`
	CurrentHash      string    `json:"current_hash"`
	ShortCurrentHash string    `json:"short_current_hash"`
	CapturedAt       time.Time `json:"captured_at"`
	PointsAwarded    *int64    `json:"points_awarded,omitempty"`
	DeliveryCount    *int64    `json:"delivery_count,omitempty"`
	TotalDistanceKm  *float64  `json:"total_distance_km,omitempty"`
}
