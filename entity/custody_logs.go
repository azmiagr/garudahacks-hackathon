package entity

import (
	"time"

	"github.com/google/uuid"
)

const (
	CustodyStageFundLocked            = "fund_locked"
	CustodyStageStoreToCourier        = "store_to_courier"
	CustodyStageCourierToPost         = "courier_to_post"
	CustodyStageDistributionCompleted = "distribution_completed"

	HandshakeMethodSystem = "system"
	HandshakeMethodQR     = "qr"
	HandshakeMethodPIN    = "pin"
)

type CustodyLogs struct {
	LogsID            string     `json:"logs_id" gorm:"type:varchar(36);primaryKey"`
	OrderID           string     `json:"order_id" gorm:"type:varchar(36);uniqueIndex:idx_order_custody_sequence"`
	HandoffStage      string     `json:"handoff_stage" gorm:"type:varchar(40);index"`
	HandshakeMethod   string     `json:"handshake_method" gorm:"type:varchar(20)"`
	FromActorID       uuid.UUID  `json:"from_actor_id" gorm:"type:varchar(36)"`
	ToActorID         uuid.UUID  `json:"to_actor_id" gorm:"type:varchar(36)"`
	ScannedBy         *uuid.UUID `json:"scanned_by" gorm:"type:varchar(36);index"`
	Sequence          int        `json:"sequence" gorm:"type:int;not null;uniqueIndex:idx_order_custody_sequence"`
	Latitude          float64    `json:"latitude" gorm:"type:decimal(10,2);not null"`
	Longitude         float64    `json:"longitude" gorm:"type:decimal(11,8);not null"`
	GPSDistanceMeters *float64   `json:"gps_distance_meters" gorm:"type:decimal(10,2)"`
	IsGPSAnomaly      bool       `json:"is_gps_anomaly" gorm:"default:false"`
	IdempotencyKey    *string    `json:"idempotency_key" gorm:"type:varchar(120);uniqueIndex"`
	PrevHash          string     `json:"prev_hash" gorm:"type:varchar(255);not null"`
	CurrentHash       string     `json:"current_hash" gorm:"type:varchar(255);not null"`
	CapturedAt        *time.Time `json:"captured_at" gorm:"type:datetime"`
	CreatedAt         time.Time  `json:"created_at" gorm:"autoCreateTime"`
}
