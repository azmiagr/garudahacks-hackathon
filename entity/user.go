package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	UserID    uuid.UUID `json:"user_id" gorm:"type:varchar(36);primaryKey"`
	RoleID    uuid.UUID `json:"role_id" gorm:"type:varchar(36)"`
	Name      string    `json:"name" gorm:"type:varchar(150);not null"`
	Email     string    `json:"email" gorm:"type:varchar(150);not null;uniqueIndex"`
	Password  string    `json:"password" gorm:"type:varchar(255);not null"`
	Status    string    `json:"status" gorm:"type:enum('active','inactive');default:'inactive'"`
	KYCStatus string    `json:"kyc_status" gorm:"type:enum('pending','approved','rejected');default:'pending'"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	AdminProfile          AdminProfile           `json:"admin_profile" gorm:"foreignKey:UserID;references:UserID;constraint:onDelete:CASCADE"`
	DonorProfile          DonorProfile           `json:"donor_profile" gorm:"foreignKey:UserID;references:UserID;constraint:onDelete:CASCADE"`
	OtpCodes              []OtpCode              `json:"otp_codes" gorm:"foreignKey:UserID;references:UserID;constraint:onDelete:CASCADE"`
	Posts                 []Post                 `json:"posts" gorm:"foreignKey:UserID;references:UserID;constraint:onDelete:CASCADE"`
	DisasterReports       []DisasterReport       `json:"disaster_reports" gorm:"-"`
	Requests              []Requests             `json:"requests" gorm:"foreignKey:CreatedBy;references:UserID;constraint:onDelete:CASCADE"`
	Donations             []Donations            `json:"donations" gorm:"foreignKey:DonatedBy;references:UserID;constraint:onDelete:CASCADE"`
	Orders                []Orders               `json:"orders" gorm:"foreignKey:CourierID;references:UserID;constraint:onDelete:CASCADE"`
	DeliveryVerifications []DeliveryVerification `json:"delivery_verifications" gorm:"foreignKey:VerifiedBy;references:UserID;constraint:onDelete:CASCADE"`
	DeliverySubmissions   []DeliveryVerification `json:"delivery_submissions" gorm:"foreignKey:SubmittedBy;references:UserID;constraint:onDelete:CASCADE"`
	SentCustodyLogs       []CustodyLogs          `json:"sent_custody_logs" gorm:"foreignKey:FromActorID;references:UserID;constraint:onDelete:CASCADE"`
	ReceivedCustodyLogs   []CustodyLogs          `json:"received_custody_logs" gorm:"foreignKey:ToActorID;references:UserID;constraint:onDelete:CASCADE"`
	Disbursements         []Disbursements        `json:"disbursements" gorm:"foreignKey:HeldBy;references:UserID;constraint:onDelete:CASCADE"`
	Stores                []Stores               `json:"stores" gorm:"foreignKey:OwnerID;references:UserID;constraint:onDelete:CASCADE"`
	Wallet                Wallets                `json:"wallet" gorm:"foreignKey:UserID;references:UserID;constraint:onDelete:CASCADE"`
}
