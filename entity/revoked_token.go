package entity

import "time"

type RevokedToken struct {
	TokenHash string    `json:"token_hash" gorm:"type:varchar(64);primaryKey"`
	ExpiresAt time.Time `json:"expires_at" gorm:"type:datetime;not null;index"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}
