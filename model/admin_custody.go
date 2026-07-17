package model

import (
	"mime/multipart"
	"time"

	"github.com/google/uuid"
)

type AdminReceiveOrderDetailResponse struct {
	OrderID            uuid.UUID                    `json:"order_id"`
	OrderCode          string                       `json:"order_code"`
	OrderStatus        string                       `json:"order_status"`
	RequestID          uuid.UUID                    `json:"request_id"`
	RequestTitle       string                       `json:"request_title"`
	PostName           string                       `json:"post_name"`
	CourierID          uuid.UUID                    `json:"courier_id"`
	CourierName        string                       `json:"courier_name"`
	DeliveredAt        *time.Time                   `json:"delivered_at"`
	CompletedAt        *time.Time                   `json:"completed_at"`
	Items              []AdminReceiveOrderItem      `json:"items"`
	Proofs             []AdminDistributionProofData `json:"proofs"`
	RequiredPhotoCount int                          `json:"required_photo_count"`
	UploadedPhotoCount int                          `json:"uploaded_photo_count"`
}

type AdminReceiveOrderItem struct {
	ItemID    uuid.UUID `json:"item_id"`
	Name      string    `json:"name"`
	Quantity  int       `json:"quantity"`
	Unit      int       `json:"unit"`
	UnitPrice float64   `json:"unit_price"`
	Subtotal  float64   `json:"subtotal"`
	HasProof  bool      `json:"has_proof"`
}

type CreateSupplementalNeedRequest struct {
	Reason                string                       `json:"reason" binding:"required"`
	ReservedAmountApplied float64                      `json:"reserved_amount_applied"`
	Items                 []CreateSupplementalNeedItem `json:"items" binding:"required,min=1"`
}

type CreateSupplementalNeedItem struct {
	ItemID         uuid.UUID `json:"item_id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	Price          float64   `json:"price" binding:"required"`
	QuantityNeeded int       `json:"quantity_needed" binding:"required"`
}

type CreateSupplementalNeedResponse struct {
	SupplementalID        uuid.UUID                  `json:"supplemental_id"`
	RequestID             uuid.UUID                  `json:"request_id"`
	OrderID               uuid.UUID                  `json:"order_id"`
	Reason                string                     `json:"reason"`
	ReservedAmountApplied float64                    `json:"reserved_amount_applied"`
	AdditionalTarget      float64                    `json:"additional_target"`
	NewFundingTarget      float64                    `json:"new_funding_target"`
	Items                 []CreateAdminEventItemData `json:"items"`
}

type UploadDistributionProofRequest struct {
	OrderID             uuid.UUID             `form:"-" json:"-"`
	ItemID              uuid.UUID             `form:"item_id" json:"item_id" binding:"required"`
	Photo               *multipart.FileHeader `form:"photo" json:"-"`
	RecipientNote       string                `form:"recipient_note" json:"recipient_note"`
	DistributedQuantity int                   `form:"distributed_quantity" json:"distributed_quantity"`
	Latitude            float64               `form:"latitude" json:"latitude" binding:"required"`
	Longitude           float64               `form:"longitude" json:"longitude" binding:"required"`
	BlurFaceEnabled     bool                  `form:"blur_face_enabled" json:"blur_face_enabled"`
	CapturedFromCamera  bool                  `form:"captured_from_camera" json:"captured_from_camera"`
	CapturedAt          time.Time             `form:"captured_at" json:"captured_at"`
}

type UploadDistributionProofResponse struct {
	Proof              AdminDistributionProofData `json:"proof"`
	RequiredPhotoCount int                        `json:"required_photo_count"`
	UploadedPhotoCount int                        `json:"uploaded_photo_count"`
	ReadyToComplete    bool                       `json:"ready_to_complete"`
}

type AdminDistributionProofData struct {
	ProofID             uuid.UUID `json:"proof_id"`
	OrderID             uuid.UUID `json:"order_id"`
	ItemID              uuid.UUID `json:"item_id"`
	ImageURL            string    `json:"image_url"`
	RecipientNote       string    `json:"recipient_note"`
	DistributedQuantity int       `json:"distributed_quantity"`
	Latitude            float64   `json:"latitude"`
	Longitude           float64   `json:"longitude"`
	GPSDistanceMeters   float64   `json:"gps_distance_meters"`
	BlurFaceEnabled     bool      `json:"blur_face_enabled"`
	CapturedFromCamera  bool      `json:"captured_from_camera"`
	CapturedAt          time.Time `json:"captured_at"`
}

type CompleteDistributionRequest struct {
	IdempotencyKey string    `json:"idempotency_key"`
	Latitude       float64   `json:"latitude" binding:"required"`
	Longitude      float64   `json:"longitude" binding:"required"`
	CapturedAt     time.Time `json:"captured_at"`
}

type CompleteDistributionResponse struct {
	OrderID            uuid.UUID `json:"order_id"`
	OrderStatus        string    `json:"order_status"`
	RequiredPhotoCount int       `json:"required_photo_count"`
	UploadedPhotoCount int       `json:"uploaded_photo_count"`
	FinalHash          string    `json:"final_hash"`
	ShortFinalHash     string    `json:"short_final_hash"`
	CompletedAt        time.Time `json:"completed_at"`
}
