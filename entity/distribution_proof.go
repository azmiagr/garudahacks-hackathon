package entity

import (
	"time"

	"github.com/google/uuid"
)

type DistributionProof struct {
	ProofID             uuid.UUID `json:"proof_id" gorm:"type:varchar(36);primaryKey"`
	OrderID             uuid.UUID `json:"order_id" gorm:"type:varchar(36);not null;uniqueIndex:idx_distribution_order_item"`
	ItemID              uuid.UUID `json:"item_id" gorm:"type:varchar(36);not null;uniqueIndex:idx_distribution_order_item"`
	SubmittedBy         uuid.UUID `json:"submitted_by" gorm:"type:varchar(36);not null;index"`
	ImageURL            string    `json:"image_url" gorm:"type:varchar(255);not null"`
	RecipientNote       string    `json:"recipient_note" gorm:"type:varchar(150)"`
	DistributedQuantity int       `json:"distributed_quantity" gorm:"type:int;not null;default:0"`
	Latitude            float64   `json:"latitude" gorm:"type:decimal(10,8);not null"`
	Longitude           float64   `json:"longitude" gorm:"type:decimal(11,8);not null"`
	GPSDistanceMeters   float64   `json:"gps_distance_meters" gorm:"type:decimal(10,2);not null;default:0"`
	BlurFaceEnabled     bool      `json:"blur_face_enabled" gorm:"not null;default:true"`
	CapturedFromCamera  bool      `json:"captured_from_camera" gorm:"not null;default:false"`
	CapturedAt          time.Time `json:"captured_at" gorm:"type:datetime"`
	CreatedAt           time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt           time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
