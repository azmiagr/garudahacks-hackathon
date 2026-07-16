package entity

import "github.com/google/uuid"

type DisasterEvent struct {
	EventID     uuid.UUID `json:"event_id" gorm:"type:varchar(36);primaryKey"`
	Name        string    `json:"name" gorm:"type:varchar(150);not null"`
	Description string    `json:"description" gorm:"type:text"`

	// DisasterReports []DisasterReport `json:"disaster_reports" gorm:"foreignKey:EventID;references:EventID;constraint:onDelete:CASCADE"`
}
