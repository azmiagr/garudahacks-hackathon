package model

import "github.com/google/uuid"

type CourierProfileResponse struct {
	UserID                        uuid.UUID `json:"user_id"`
	ProfileID                     uuid.UUID `json:"profile_id"`
	Name                          string    `json:"name"`
	Initials                      string    `json:"initials"`
	Email                         string    `json:"email"`
	Role                          string    `json:"role"`
	DisplayRole                   string    `json:"display_role"`
	KYCStatus                     string    `json:"kyc_status"`
	IsVerified                    bool      `json:"is_verified"`
	VerificationText              string    `json:"verification_text"`
	OperationalArea               string    `json:"operational_area"`
	OperationRadiusKM             int       `json:"operation_radius_km"`
	OperationAreaText             string    `json:"operation_area_text"`
	VehicleType                   string    `json:"vehicle_type"`
	VehicleCapacityKG             int       `json:"vehicle_capacity_kg"`
	VehicleText                   string    `json:"vehicle_text"`
	WaiverAccepted                bool      `json:"waiver_accepted"`
	IsAvailable                   bool      `json:"is_available"`
	UrgentTaskNotificationEnabled bool      `json:"urgent_task_notification_enabled"`
	ReputationScore               float64   `json:"reputation_score"`
	ReputationText                string    `json:"reputation_text"`
	ActivePoints                  int64     `json:"active_points"`
	ActivePointsText              string    `json:"active_points_text"`
	DeliveryCount                 int64     `json:"delivery_count"`
	TotalDistanceKm               float64   `json:"total_distance_km"`
}

type UpdateCourierProfilePreferencesRequest struct {
	IsAvailable                   *bool `json:"is_available"`
	UrgentTaskNotificationEnabled *bool `json:"urgent_task_notification_enabled"`
}

type UpdateCourierProfilePreferencesResponse struct {
	IsAvailable                   bool `json:"is_available"`
	UrgentTaskNotificationEnabled bool `json:"urgent_task_notification_enabled"`
}
