package entity

import (
	"time"

	"github.com/google/uuid"
)

type CourierProfile struct {
	ProfileID         uuid.UUID `json:"profile_id" gorm:"type:varchar(36);primaryKey"`
	UserID            uuid.UUID `json:"user_id" gorm:"type:varchar(36);not null;uniqueIndex"`
	NIK               string    `json:"-" gorm:"type:varchar(64);not null;uniqueIndex"`
	VehicleType       string    `json:"vehicle_type" gorm:"type:varchar(40);not null"`
	VehicleCapacityKG int       `json:"vehicle_capacity_kg" gorm:"not null;default:0"`
	OperationalArea   string    `json:"operational_area" gorm:"type:varchar(150);not null"`
	OperationRadiusKM int       `json:"operation_radius_km" gorm:"not null;default:0"`
	WaiverAccepted    bool      `json:"waiver_accepted" gorm:"not null;default:false"`
	WaiverAcceptedAt  time.Time `json:"waiver_accepted_at"`
	CreatedAt         time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt         time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
