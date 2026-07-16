package entity

import (
	"time"

	"github.com/google/uuid"
)

type AdminProfile struct {
	ProfileID   uuid.UUID `json:"profile_id" gorm:"type:varchar(36);primaryKey"`
	UserID      uuid.UUID `json:"user_id" gorm:"type:varchar(36);not null;uniqueIndex"`
	NIK         string    `json:"-" gorm:"type:varchar(64);not null;uniqueIndex"`
	Affiliation string    `json:"affiliation" gorm:"type:varchar(150);not null"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
