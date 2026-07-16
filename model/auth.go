package model

import "github.com/google/uuid"

type GetRegistrationSessionParam struct {
	RegistrationID uuid.UUID `json:"registration_id"`
	Email          string    `json:"email"`
}

type GetRoleParam struct {
	RoleID   uuid.UUID `json:"role_id"`
	RoleName string    `json:"role_name"`
}

type GetAdminPoskoProfileParam struct {
	ProfileID uuid.UUID `json:"profile_id"`
	UserID    uuid.UUID `json:"user_id"`
	NIK       string    `json:"nik"`
}

type RequestAdminRegisterOtpRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type RequestAdminRegisterOtpResponse struct {
	RegistrationID     uuid.UUID `json:"registration_id"`
	Email              string    `json:"email"`
	OtpExpiresInSecond int       `json:"otp_expires_in_seconds"`
}

type VerifyAdminRegisterOtpRequest struct {
	RegistrationID uuid.UUID `json:"registration_id" binding:"required"`
	OtpCode        string    `json:"otp_code" binding:"required,len=6"`
}

type VerifyAdminRegisterOtpResponse struct {
	RegistrationID uuid.UUID `json:"registration_id"`
	OtpVerified    bool      `json:"otp_verified"`
}

type SetAdminRegisterPasswordRequest struct {
	RegistrationID  uuid.UUID `json:"registration_id" binding:"required"`
	Password        string    `json:"password" binding:"required,min=8"`
	ConfirmPassword string    `json:"confirm_password" binding:"required"`
}

type SetAdminRegisterPasswordResponse struct {
	RegistrationID  uuid.UUID `json:"registration_id"`
	PasswordCreated bool      `json:"password_created"`
}

type CompleteAdminRegisterRequest struct {
	RegistrationID uuid.UUID `json:"registration_id" binding:"required"`
	FullName       string    `json:"full_name" binding:"required"`
	NIK            string    `json:"nik" binding:"required,len=16,numeric"`
	Affiliation    string    `json:"affiliation" binding:"required"`
}

type CompleteAdminRegisterResponse struct {
	Token string               `json:"token"`
	User  RegisterUserResponse `json:"user"`
}

type RegisterUserResponse struct {
	UserID      uuid.UUID `json:"user_id"`
	Role        string    `json:"role"`
	DisplayRole string    `json:"display_role"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Status      string    `json:"status"`
	KYCStatus   string    `json:"kyc_status"`
}
