package model

import (
	"time"

	"github.com/google/uuid"
)

type GetPointAccountParam struct {
	PointAccountID uuid.UUID
	UserID         uuid.UUID
}

type GetRewardParam struct {
	RewardID uuid.UUID
	IsActive *bool
}

type PointSummaryRow struct {
	UserID        uuid.UUID `json:"user_id"`
	ActivePoints  int64     `json:"active_points"`
	TotalEarned   int64     `json:"total_earned"`
	TotalRedeemed int64     `json:"total_redeemed"`
}

type PointHistoryParam struct {
	UserID uuid.UUID
	Limit  int
	Offset int
}

type RewardListParam struct {
	OnlyActive bool
	Limit      int
	Offset     int
}

type PointHistoryRow struct {
	PointTransactionID uuid.UUID  `json:"point_transaction_id"`
	DonationID         *uuid.UUID `json:"donation_id"`
	RewardClaimID      *uuid.UUID `json:"reward_claim_id"`
	Points             int64      `json:"points"`
	TransactionType    string     `json:"transaction_type"`
	SourceType         string     `json:"source_type"`
	SourceID           string     `json:"source_id"`
	Description        string     `json:"description"`
	ExpiresAt          *time.Time `json:"expires_at"`
	CreatedAt          time.Time  `json:"created_at"`
}

type PointDashboardParam struct {
	HistoryLimit int `form:"history_limit"`
	RewardLimit  int `form:"reward_limit"`
}

type PointHistoryQueryParam struct {
	Limit  int `form:"limit"`
	Offset int `form:"offset"`
}

type RewardQueryParam struct {
	Limit  int `form:"limit"`
	Offset int `form:"offset"`
}

type ClaimRewardRequest struct {
	RewardID uuid.UUID `json:"reward_id" binding:"required"`
}

type PointDashboardResponse struct {
	ActivePoints      int64              `json:"active_points"`
	TotalEarned       int64              `json:"total_earned"`
	TotalRedeemed     int64              `json:"total_redeemed"`
	Level             string             `json:"level"`
	NextLevelPoints   int64              `json:"next_level_points"`
	PointsToNextLevel int64              `json:"points_to_next_level"`
	Rewards           []RewardItem       `json:"rewards"`
	History           []PointHistoryItem `json:"history"`
}

type PointHistoryResponse struct {
	Items  []PointHistoryItem `json:"items"`
	Limit  int                `json:"limit"`
	Offset int                `json:"offset"`
}

type RewardListResponse struct {
	Items  []RewardItem `json:"items"`
	Limit  int          `json:"limit"`
	Offset int          `json:"offset"`
}

type RewardClaimResponse struct {
	RewardClaimID   uuid.UUID `json:"reward_claim_id"`
	RewardID        uuid.UUID `json:"reward_id"`
	RewardName      string    `json:"reward_name"`
	PointsSpent     int64     `json:"points_spent"`
	RemainingPoints int64     `json:"remaining_points"`
	ClaimStatus     string    `json:"claim_status"`
	ClaimedAt       time.Time `json:"claimed_at"`
}

type RewardItem struct {
	RewardID     uuid.UUID `json:"reward_id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	RewardType   string    `json:"reward_type"`
	PointsCost   int64     `json:"points_cost"`
	Stock        int       `json:"stock"`
	IsActive     bool      `json:"is_active"`
	ValidityDays int       `json:"validity_days"`
	CanClaim     bool      `json:"can_claim"`
}

type PointHistoryItem struct {
	PointTransactionID uuid.UUID  `json:"point_transaction_id"`
	DonationID         *uuid.UUID `json:"donation_id"`
	RewardClaimID      *uuid.UUID `json:"reward_claim_id"`
	Points             int64      `json:"points"`
	TransactionType    string     `json:"transaction_type"`
	SourceType         string     `json:"source_type"`
	SourceID           string     `json:"source_id"`
	Description        string     `json:"description"`
	ExpiresAt          *time.Time `json:"expires_at"`
	CreatedAt          time.Time  `json:"created_at"`
	PointsText         string     `json:"points_text"`
}
