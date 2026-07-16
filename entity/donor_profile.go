package entity

import (
	"time"

	"github.com/google/uuid"
)

type DonorProfile struct {
	ProfileID         uuid.UUID `json:"profile_id" gorm:"type:varchar(36);primaryKey"`
	UserID            uuid.UUID `json:"user_id" gorm:"type:varchar(36);not null;uniqueIndex"`
	PhoneNumber       string    `json:"phone_number" gorm:"type:varchar(30);not null;uniqueIndex"`
	PreferenceJSON    string    `json:"preference_json" gorm:"type:json"`
	ConsentAccepted   bool      `json:"consent_accepted" gorm:"not null;default:false"`
	ConsentAcceptedAt time.Time `json:"consent_accepted_at"`
	CreatedAt         time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt         time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
