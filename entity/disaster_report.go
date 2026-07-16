package entity

import (
	"time"

	"github.com/google/uuid"
)

type DisasterReport struct {
	ReportID        uuid.UUID `json:"report_id" gorm:"type:varchar(36);primaryKey"`
	DisasterEventID uuid.UUID `json:"event_id" gorm:"column:event_id;type:varchar(36)"`
	PostID          uuid.UUID `json:"post_id" gorm:"type:varchar(36)"`
	UserID          uuid.UUID `json:"user_id" gorm:"type:varchar(36)"`
	ReportTitle     string    `json:"report_title" gorm:"type:varchar(150);not null"`
	Description     string    `json:"description" gorm:"type:text"`
	Latitude        float64   `json:"latitude" gorm:"type:decimal(10,8);not null"`
	Longitude       float64   `json:"longitude" gorm:"type:decimal(11,8);not null"`
	ImageUrl        string    `json:"image_url" gorm:"type:varchar(255)"`
	ReportStatus    string    `json:"report_status" gorm:"type:enum('pending','approved','rejected');default:'pending'"`
	ReportedAt      time.Time `json:"reported_at" gorm:"autoCreateTime"`
	CreatedAt       time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	DisasterEvent DisasterEvent `json:"disaster_event" gorm:"foreignKey:DisasterEventID;references:EventID;constraint:OnDelete:CASCADE"`
}
