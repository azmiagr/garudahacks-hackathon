package model

import (
	"time"

	"github.com/google/uuid"
)

type CourierGoodnessParam struct {
	CourierID uuid.UUID
	Year      int
	Limit     int
	Offset    int
}

type CourierGoodnessStatsRow struct {
	DeliveryCount   int64      `json:"delivery_count"`
	TotalDistanceKm float64    `json:"total_distance_km"`
	DisputeCount    int64      `json:"dispute_count"`
	FirstDeliveryAt *time.Time `json:"first_delivery_at"`
	ReputationScore float64    `json:"reputation_score"`
}

type CourierDeliveryHistoryRow struct {
	OrderID            uuid.UUID  `json:"order_id"`
	OrderCode          string     `json:"order_code"`
	PostName           string     `json:"post_name"`
	DisasterName       string     `json:"disaster_name"`
	ItemCount          int64      `json:"item_count"`
	TotalAmount        float64    `json:"total_amount"`
	DeliveryDistanceKm *float64   `json:"delivery_distance_km"`
	DeliveredAt        *time.Time `json:"delivered_at"`
}

type CourierGoodnessRequest struct {
	Year   int `form:"year"`
	Limit  int `form:"limit"`
	Offset int `form:"offset"`
}

type CourierGoodnessResponse struct {
	Certificate  CourierGoodnessCertificate   `json:"certificate"`
	History      []CourierDeliveryHistoryItem `json:"history"`
	TotalHistory int64                        `json:"total_history"`
	Limit        int                          `json:"limit"`
	Offset       int                          `json:"offset"`
}

type CourierGoodnessCertificate struct {
	CourierID         uuid.UUID  `json:"courier_id"`
	CourierName       string     `json:"courier_name"`
	Title             string     `json:"title"`
	PartnerLabel      string     `json:"partner_label"`
	SinceText         string     `json:"since_text"`
	DeliveryCount     int64      `json:"delivery_count"`
	DeliveryCountText string     `json:"delivery_count_text"`
	TotalDistanceKm   float64    `json:"total_distance_km"`
	TotalDistanceText string     `json:"total_distance_text"`
	ReputationScore   float64    `json:"reputation_score"`
	ReputationText    string     `json:"reputation_text"`
	DisputeCount      int64      `json:"dispute_count"`
	DisputeText       string     `json:"dispute_text"`
	FirstDeliveryAt   *time.Time `json:"first_delivery_at"`
	ShareURL          string     `json:"share_url"`
}

type CourierDeliveryHistoryItem struct {
	OrderID            uuid.UUID  `json:"order_id"`
	OrderCode          string     `json:"order_code"`
	PostName           string     `json:"post_name"`
	DisasterName       string     `json:"disaster_name"`
	Title              string     `json:"title"`
	ItemCount          int64      `json:"item_count"`
	TotalAmount        float64    `json:"total_amount"`
	DeliveryDistanceKm float64    `json:"delivery_distance_km"`
	DeliveredAt        *time.Time `json:"delivered_at"`
	DeliveredAtText    string     `json:"delivered_at_text"`
}
