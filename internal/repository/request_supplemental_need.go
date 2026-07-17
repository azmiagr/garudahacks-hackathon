package repository

import (
	"github.com/azmiagr/garudahacks-hackathon/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IRequestSupplementalNeedRepository interface {
	CreateSupplementalNeed(tx *gorm.DB, need *entity.RequestSupplementalNeed) error
	ListSupplementalNeedsByOrder(tx *gorm.DB, orderID uuid.UUID) ([]entity.RequestSupplementalNeed, error)
}

type RequestSupplementalNeedRepository struct {
	db *gorm.DB
}

func NewRequestSupplementalNeedRepository(db *gorm.DB) IRequestSupplementalNeedRepository {
	return &RequestSupplementalNeedRepository{db: db}
}

func (r *RequestSupplementalNeedRepository) CreateSupplementalNeed(tx *gorm.DB, need *entity.RequestSupplementalNeed) error {
	return tx.Create(need).Error
}

func (r *RequestSupplementalNeedRepository) ListSupplementalNeedsByOrder(tx *gorm.DB, orderID uuid.UUID) ([]entity.RequestSupplementalNeed, error) {
	var needs []entity.RequestSupplementalNeed
	err := tx.Where("order_id = ?", orderID).
		Order("created_at DESC").
		Find(&needs).Error
	if err != nil {
		return nil, err
	}

	return needs, nil
}
