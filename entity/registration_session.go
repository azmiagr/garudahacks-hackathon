package entity

import (
	"time"

	"github.com/google/uuid"
)

type RegistrationSession struct {
	RegistrationID uuid.UUID  `json:"registration_id" gorm:"type:varchar(36);primaryKey"`
	Email          string     `json:"email" gorm:"type:varchar(150);not null;uniqueIndex:idx_registration_email_role"`
	RoleName       string     `json:"role_name" gorm:"type:varchar(50);not null;uniqueIndex:idx_registration_email_role"`
	OtpCode        string     `json:"-" gorm:"type:varchar(6);not null"`
	OtpExpiresAt   time.Time  `json:"otp_expires_at" gorm:"not null"`
	OtpVerifiedAt  *time.Time `json:"otp_verified_at"`
	PasswordHash   string     `json:"-" gorm:"type:varchar(255)"`
	ExpiresAt      time.Time  `json:"expires_at" gorm:"not null"`
	CreatedAt      time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}
