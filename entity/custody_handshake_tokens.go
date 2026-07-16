package entity

import (
	"time"

	"github.com/google/uuid"
)

const (
	CustodyTokenStatusActive  = "active"
	CustodyTokenStatusUsed    = "used"
	CustodyTokenStatusExpired = "expired"
	CustodyTokenStatusRevoked = "revoked"
)

type CustodyHandshakeToken struct {
	TokenID         uuid.UUID  `json:"token_id" gorm:"type:varchar(36);primaryKey"`
	OrderID         uuid.UUID  `json:"order_id" gorm:"type:varchar(36);index"`
	HandoffStage    string     `json:"handoff_stage" gorm:"type:varchar(40);index"`
	PresentedBy     uuid.UUID  `json:"presented_by" gorm:"type:varchar(36);index"`
	QRPayloadHash   string     `json:"qr_payload_hash" gorm:"type:varchar(255);uniqueIndex"`
	PINHash         string     `json:"pin_hash" gorm:"type:varchar(255);index"`
	Nonce           string     `json:"nonce" gorm:"type:varchar(80);uniqueIndex"`
	Status          string     `json:"status" gorm:"type:varchar(20);default:'active';index"`
	ExpiresAt       time.Time  `json:"expires_at" gorm:"type:datetime;index"`
	CacheValidUntil time.Time  `json:"cache_valid_until" gorm:"type:datetime;index"`
	UsedAt          *time.Time `json:"used_at" gorm:"type:datetime"`
	UsedBy          *uuid.UUID `json:"used_by" gorm:"type:varchar(36);index"`
	CreatedAt       time.Time  `json:"created_at" gorm:"autoCreateTime"`
}
