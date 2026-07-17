package model

import (
	"mime/multipart"

	"github.com/google/uuid"
)

type GetRegistrationSessionParam struct {
	RegistrationID uuid.UUID `json:"registration_id"`
	Email          string    `json:"email"`
	RoleName       string    `json:"role_name"`
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

type GetDonorProfileParam struct {
	ProfileID   uuid.UUID `json:"profile_id"`
	UserID      uuid.UUID `json:"user_id"`
	PhoneNumber string    `json:"phone_number"`
}

type GetCourierProfileParam struct {
	ProfileID uuid.UUID `json:"profile_id"`
	UserID    uuid.UUID `json:"user_id"`
	NIK       string    `json:"nik"`
}

type RequestAdminRegisterOtpRequest struct {
	Email string `json:"email" binding:"required,email"`
	Role  string `json:"role"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type LogoutResponse struct {
	LoggedOut bool `json:"logged_out"`
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

type CompleteDonorRegisterRequest struct {
	RegistrationID      uuid.UUID `json:"registration_id" binding:"required"`
	FullName            string    `json:"full_name" binding:"required"`
	PhoneNumber         string    `json:"phone_number" binding:"required"`
	DonationPreferences []string  `json:"donation_preferences"`
	ConsentAccepted     bool      `json:"consent_accepted" binding:"required"`
}

type CompleteDonorRegisterResponse struct {
	Token string               `json:"token"`
	User  RegisterUserResponse `json:"user"`
}

type CompleteCourierRegisterRequest struct {
	RegistrationID    uuid.UUID `json:"registration_id" binding:"required"`
	FullName          string    `json:"full_name" binding:"required"`
	NIK               string    `json:"nik" binding:"required,len=16,numeric"`
	VehicleType       string    `json:"vehicle_type" binding:"required"`
	VehicleCapacityKG int       `json:"vehicle_capacity_kg"`
	OperationalArea   string    `json:"operational_area" binding:"required"`
	OperationRadiusKM int       `json:"operation_radius_km" binding:"required"`
	WaiverAccepted    bool      `json:"waiver_accepted" binding:"required"`
}

type CompleteCourierRegisterResponse struct {
	Token   string                  `json:"token"`
	User    RegisterUserResponse    `json:"user"`
	Courier CourierRegisterResponse `json:"courier"`
}

type CourierRegisterResponse struct {
	ProfileID         uuid.UUID `json:"profile_id"`
	UserID            uuid.UUID `json:"user_id"`
	VehicleType       string    `json:"vehicle_type"`
	VehicleCapacityKG int       `json:"vehicle_capacity_kg"`
	OperationalArea   string    `json:"operational_area"`
	OperationRadiusKM int       `json:"operation_radius_km"`
	WaiverAccepted    bool      `json:"waiver_accepted"`
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

type CompleteStoreRegisterRequest struct {
	RegistrationID  uuid.UUID             `form:"registration_id" json:"registration_id" binding:"required"`
	StoreName       string                `form:"store_name" json:"store_name" binding:"required"`
	OwnerName       string                `form:"owner_name" json:"owner_name" binding:"required"`
	NIB             string                `form:"nib" json:"nib" binding:"required"`
	NPWP            string                `form:"npwp" json:"npwp"`
	KTPImage        *multipart.FileHeader `form:"ktp_image" json:"-"`
	BankName        string                `form:"bank_name" json:"bank_name"`
	BankAccountNo   string                `form:"bank_account_no" json:"bank_account_no"`
	BankAccountName string                `form:"bank_account_name" json:"bank_account_name"`
	Categories      []string              `form:"categories" json:"categories"`
	CategoriesJSON  string                `form:"categories_json" json:"categories_json"`
	Address         string                `form:"address" json:"address" binding:"required"`
	Latitude        float64               `form:"latitude" json:"latitude" binding:"required"`
	Longitude       float64               `form:"longitude" json:"longitude" binding:"required"`
}

type CompleteStoreRegisterResponse struct {
	Token string                `json:"token"`
	User  RegisterUserResponse  `json:"user"`
	Store StoreRegisterResponse `json:"store"`
}

type StoreRegisterResponse struct {
	StoreID        uuid.UUID `json:"store_id"`
	OwnerID        uuid.UUID `json:"owner_id"`
	Name           string    `json:"name"`
	BusinessNumber string    `json:"business_number"`
	Address        string    `json:"address"`
	Latitude       float64   `json:"latitude"`
	Longitude      float64   `json:"longitude"`
}

type GetStoreParam struct {
	StoreID        uuid.UUID `json:"store_id"`
	OwnerID        uuid.UUID `json:"owner_id"`
	BusinessNumber string    `json:"business_number"`
}
