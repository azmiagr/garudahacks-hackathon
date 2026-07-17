package repository

import (
	"github.com/azmiagr/garudahacks-hackathon/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type IDistributionProofRepository interface {
	CreateDistributionProof(tx *gorm.DB, proof *entity.DistributionProof) error
	GetDistributionProof(tx *gorm.DB, orderID uuid.UUID, itemID uuid.UUID) (*entity.DistributionProof, error)
	ListDistributionProofsByOrder(tx *gorm.DB, orderID uuid.UUID) ([]entity.DistributionProof, error)
	CountDistributionProofsByOrder(tx *gorm.DB, orderID uuid.UUID) (int64, error)
}

type DistributionProofRepository struct {
	db *gorm.DB
}

func NewDistributionProofRepository(db *gorm.DB) IDistributionProofRepository {
	return &DistributionProofRepository{db: db}
}

func (r *DistributionProofRepository) CreateDistributionProof(tx *gorm.DB, proof *entity.DistributionProof) error {
	err := tx.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "order_id"},
			{Name: "item_id"},
		},
		DoUpdates: clause.AssignmentColumns([]string{
			"submitted_by",
			"image_url",
			"recipient_note",
			"distributed_quantity",
			"latitude",
			"longitude",
			"gps_distance_meters",
			"blur_face_enabled",
			"captured_from_camera",
			"captured_at",
			"updated_at",
		}),
	}).Create(proof).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *DistributionProofRepository) GetDistributionProof(tx *gorm.DB, orderID uuid.UUID, itemID uuid.UUID) (*entity.DistributionProof, error) {
	var proof entity.DistributionProof
	err := tx.Where("order_id = ? AND item_id = ?", orderID, itemID).First(&proof).Error
	if err != nil {
		return nil, err
	}

	return &proof, nil
}

func (r *DistributionProofRepository) ListDistributionProofsByOrder(tx *gorm.DB, orderID uuid.UUID) ([]entity.DistributionProof, error) {
	var proofs []entity.DistributionProof
	err := tx.Where("order_id = ?", orderID).Find(&proofs).Error
	if err != nil {
		return nil, err
	}

	return proofs, nil
}

func (r *DistributionProofRepository) CountDistributionProofsByOrder(tx *gorm.DB, orderID uuid.UUID) (int64, error) {
	var count int64
	err := tx.Model(&entity.DistributionProof{}).
		Where("order_id = ?", orderID).
		Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}
