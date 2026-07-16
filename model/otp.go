package model

import "github.com/google/uuid"

type GetOtp struct {
	OtpID  uuid.UUID `json:"otp_id"`
	UserID uuid.UUID `json:"user_id"`
	Code   string    `json:"code"`
}
