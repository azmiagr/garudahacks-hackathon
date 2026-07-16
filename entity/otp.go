package entity

import (
	"time"

	"github.com/google/uuid"
)

type OtpCode struct {
	OtpID     uuid.UUID `gorm:"type:varchar(36);not null;primaryKey"`
	UserID    uuid.UUID `gorm:"type:varchar(36);not null"`
	Code      string    `gorm:"type:varchar(6);unique"`
	CreatedAt time.Time `gorm:"autoCreateTime;not null"`
	UpdatedAt time.Time `gorm:"autoUpdateTime;not null"`
}
