package entity

import (
	"time"

	"github.com/google/uuid"
)

type Requests struct {
	RequestID      uuid.UUID `json:"request_id" gorm:"type:varchar(36);primaryKey"`
	ReportID       uuid.UUID `json:"report_id" gorm:"type:varchar(36)"`
	CreatedBy      uuid.UUID `json:"created_by" gorm:"type:varchar(36)"`
	Title          string    `json:"title" gorm:"type:varchar(150);not null"`
	Description    string    `json:"description" gorm:"type:text"`
	FundingTarget  float64   `json:"funding_target" gorm:"type:decimal(10,2);not null"`
	FundedAmount   float64   `json:"funded_amount" gorm:"type:decimal(10,2);not null"`
	ReservedAmount float64   `json:"reserved_amount" gorm:"type:decimal(10,2);not null"`
	RequestStatus  string    `json:"request_status" gorm:"type:enum('pending','approved','rejected');default:'pending'"`
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	Items     []Items     `gorm:"foreignKey:RequestID;references:RequestID;constraint:onDelete:CASCADE"`
	Donations []Donations `gorm:"foreignKey:RequestID;references:RequestID;constraint:onDelete:CASCADE"`
	Orders    []Orders    `gorm:"foreignKey:RequestID;references:RequestID;constraint:onDelete:CASCADE"`
}
