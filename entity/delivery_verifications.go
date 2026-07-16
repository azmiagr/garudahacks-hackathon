package entity

import (
	"time"

	"github.com/google/uuid"
)

type DeliveryVerification struct {
	VerificationID     uuid.UUID `json:"verification_id" gorm:"type:varchar(36);primaryKey"`
	OrderID            uuid.UUID `json:"order_id" gorm:"type:varchar(36)"`
	SubmittedBy        uuid.UUID `json:"submitted_by" gorm:"type:varchar(36)"`
	VerifiedBy         uuid.UUID `json:"verified_by" gorm:"type:varchar(36)"`
	VerificationStatus string    `json:"verification_status" gorm:"type:enum('pending','approved','rejected');default:'pending'"`
	ImageURL           string    `json:"image_url" gorm:"type:varchar(255)"`
	Latitude           float64   `json:"latitude" gorm:"type:decimal(10,8);not null"`
	Longitude          float64   `json:"longitude" gorm:"type:decimal(11,8);not null"`
	CapturedAt         time.Time `json:"captured_at" gorm:"type:datetime"`
	ReviewedAt         time.Time `json:"reviewed_at" gorm:"type:datetime"`
	CreatedAt          time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt          time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
