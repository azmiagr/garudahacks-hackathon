package entity

import (
	"time"

	"github.com/google/uuid"
)

type Post struct {
	PostID         uuid.UUID `json:"post_id" gorm:"type:varchar(36);primaryKey"`
	UserID         uuid.UUID `json:"user_id" gorm:"type:varchar(36)"`
	Name           string    `json:"name" gorm:"type:varchar(150);not null"`
	Description    string    `json:"description" gorm:"type:text"`
	Address        string    `json:"address" gorm:"type:varchar(255);not null"`
	PhoneNumber    string    `json:"phone_number" gorm:"type:varchar(30)"`
	Latitude       float64   `json:"latitude" gorm:"type:decimal(10,8);not null"`
	Longitude      float64   `json:"longitude" gorm:"type:decimal(11,8);not null"`
	GeofenceRadius float64   `json:"geofence_radius" gorm:"type:decimal(10,2)"`
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	DisasterReports []DisasterReport `json:"disaster_reports" gorm:"foreignKey:PostID;references:PostID;constraint:OnDelete:CASCADE"`
}
