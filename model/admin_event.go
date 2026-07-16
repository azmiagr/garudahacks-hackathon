package model

import (
	"mime/multipart"

	"github.com/google/uuid"
)

type GetPostParam struct {
	PostID uuid.UUID `json:"post_id"`
	UserID uuid.UUID `json:"user_id"`
}

type GetDisasterEventParam struct {
	EventID uuid.UUID `json:"event_id"`
	Name    string    `json:"name"`
}

type GetDisasterReportParam struct {
	ReportID uuid.UUID `json:"report_id"`
	PostID   uuid.UUID `json:"post_id"`
	UserID   uuid.UUID `json:"user_id"`
}

type GetRequestParam struct {
	RequestID uuid.UUID `json:"request_id"`
	ReportID  uuid.UUID `json:"report_id"`
	CreatedBy uuid.UUID `json:"created_by"`
}

type GetItemParam struct {
	ItemID    uuid.UUID `json:"item_id"`
	RequestID uuid.UUID `json:"request_id"`
}

type CreateAdminEventRequest struct {
	Name           string                 `json:"name" binding:"required"`
	Description    string                 `json:"description"`
	DisasterType   string                 `json:"disaster_type" binding:"required"`
	Address        string                 `json:"address" binding:"required"`
	Latitude       float64                `json:"latitude" binding:"required"`
	Longitude      float64                `json:"longitude" binding:"required"`
	GeofenceRadius float64                `json:"geofence_radius" binding:"required"`
	Photo          *multipart.FileHeader  `json:"-"`
	Items          []CreateAdminEventItem `json:"items" binding:"required,min=1"`
}

type CreateAdminEventItem struct {
	Name           string  `json:"name" binding:"required"`
	Description    string  `json:"description"`
	Price          float64 `json:"price" binding:"required"`
	QuantityNeeded int     `json:"quantity_needed" binding:"required"`
}

type CreateAdminEventResponse struct {
	PostID         uuid.UUID                  `json:"post_id"`
	ReportID       uuid.UUID                  `json:"report_id"`
	RequestID      uuid.UUID                  `json:"request_id"`
	EventCode      string                     `json:"event_code"`
	Name           string                     `json:"name"`
	DisasterType   string                     `json:"disaster_type"`
	ImageURL       string                     `json:"image_url"`
	Address        string                     `json:"address"`
	Latitude       float64                    `json:"latitude"`
	Longitude      float64                    `json:"longitude"`
	GeofenceRadius float64                    `json:"geofence_radius"`
	FundingTarget  float64                    `json:"funding_target"`
	ReportStatus   string                     `json:"report_status"`
	RequestStatus  string                     `json:"request_status"`
	Items          []CreateAdminEventItemData `json:"items"`
}

type CreateAdminEventItemData struct {
	ItemID         uuid.UUID `json:"item_id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	Price          float64   `json:"price"`
	EstimatedTotal float64   `json:"estimated_total"`
	QuantityNeeded int       `json:"quantity_needed"`
}
