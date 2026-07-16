package entity

import (
	"time"

	"github.com/google/uuid"
)

type CustodyLogs struct {
	LogsID      string    `json:"logs_id" gorm:"type:varchar(36);primaryKey"`
	OrderID     string    `json:"order_id" gorm:"type:varchar(36)"`
	FromActorID uuid.UUID `json:"from_actor_id" gorm:"type:varchar(36)"`
	ToActorID   uuid.UUID `json:"to_actor_id" gorm:"type:varchar(36)"`
	Sequence    int       `json:"sequence" gorm:"type:int;not null;uniqueIndex"`
	Latitude    float64   `json:"latitude" gorm:"type:decimal(10,2);not null"`
	Longitude   float64   `json:"longitude" gorm:"type:decimal(11,8);not null"`
	PrevHash    string    `json:"prev_hash" gorm:"type:varchar(255);not null"`
	CurrentHash string    `json:"current_hash" gorm:"type:varchar(255);not null"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
}
