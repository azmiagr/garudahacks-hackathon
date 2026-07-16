package entity

import (
	"time"

	"github.com/google/uuid"
)

type PointAccount struct {
	PointAccountID uuid.UUID `json:"point_account_id" gorm:"type:varchar(36);primaryKey"`
	UserID         uuid.UUID `json:"user_id" gorm:"type:varchar(36);not null;uniqueIndex"`
	ActivePoints   int64     `json:"active_points" gorm:"type:bigint;not null;default:0"`
	TotalEarned    int64     `json:"total_earned" gorm:"type:bigint;not null;default:0"`
	TotalRedeemed  int64     `json:"total_redeemed" gorm:"type:bigint;not null;default:0"`
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	PointTransactions []PointTransaction `json:"point_transactions" gorm:"foreignKey:PointAccountID;references:PointAccountID;constraint:onDelete:CASCADE"`
	RewardClaims      []RewardClaim      `json:"reward_claims" gorm:"foreignKey:PointAccountID;references:PointAccountID;constraint:onDelete:CASCADE"`
}

type PointTransaction struct {
	PointTransactionID uuid.UUID  `json:"point_transaction_id" gorm:"type:varchar(36);primaryKey"`
	PointAccountID     uuid.UUID  `json:"point_account_id" gorm:"type:varchar(36);not null;index"`
	UserID             uuid.UUID  `json:"user_id" gorm:"type:varchar(36);not null;index"`
	DonationID         *uuid.UUID `json:"donation_id" gorm:"type:varchar(36);index"`
	RewardClaimID      *uuid.UUID `json:"reward_claim_id" gorm:"type:varchar(36);index"`
	Points             int64      `json:"points" gorm:"type:bigint;not null"`
	TransactionType    string     `json:"transaction_type" gorm:"type:enum('earn','redeem','adjustment');not null"`
	SourceType         string     `json:"source_type" gorm:"type:varchar(40);not null;uniqueIndex:idx_point_source"`
	SourceID           string     `json:"source_id" gorm:"type:varchar(80);not null;uniqueIndex:idx_point_source"`
	Description        string     `json:"description" gorm:"type:varchar(255)"`
	ExpiresAt          *time.Time `json:"expires_at"`
	CreatedAt          time.Time  `json:"created_at" gorm:"autoCreateTime"`
}

type Reward struct {
	RewardID     uuid.UUID `json:"reward_id" gorm:"type:varchar(36);primaryKey"`
	Name         string    `json:"name" gorm:"type:varchar(120);not null"`
	Description  string    `json:"description" gorm:"type:varchar(255)"`
	RewardType   string    `json:"reward_type" gorm:"type:enum('pulsa','voucher','donation');not null"`
	PointsCost   int64     `json:"points_cost" gorm:"type:bigint;not null"`
	Stock        int       `json:"stock" gorm:"type:int;not null;default:0"`
	IsActive     bool      `json:"is_active" gorm:"not null;default:true"`
	ValidityDays int       `json:"validity_days" gorm:"type:int;not null;default:0"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	RewardClaims []RewardClaim `json:"reward_claims" gorm:"foreignKey:RewardID;references:RewardID;constraint:onDelete:CASCADE"`
}

type RewardClaim struct {
	RewardClaimID  uuid.UUID `json:"reward_claim_id" gorm:"type:varchar(36);primaryKey"`
	PointAccountID uuid.UUID `json:"point_account_id" gorm:"type:varchar(36);not null;index"`
	UserID         uuid.UUID `json:"user_id" gorm:"type:varchar(36);not null;index"`
	RewardID       uuid.UUID `json:"reward_id" gorm:"type:varchar(36);not null;index"`
	PointsSpent    int64     `json:"points_spent" gorm:"type:bigint;not null"`
	ClaimStatus    string    `json:"claim_status" gorm:"type:enum('claimed','processing','fulfilled','rejected');not null;default:'claimed'"`
	ClaimedAt      time.Time `json:"claimed_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
