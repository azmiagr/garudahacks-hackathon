package repository

import (
	"github.com/azmiagr/garudahacks-hackathon/entity"
	"github.com/azmiagr/garudahacks-hackathon/model"
	"gorm.io/gorm"
)

type IItemRepository interface {
	CreateItem(tx *gorm.DB, item *entity.Items) error
	CreateItems(tx *gorm.DB, items []entity.Items) error
	GetItem(tx *gorm.DB, param model.GetItemParam) (*entity.Items, error)
	GetItemsByRequestID(tx *gorm.DB, param model.GetItemParam) ([]entity.Items, error)
}

type ItemRepository struct {
	db *gorm.DB
}

func NewItemRepository(db *gorm.DB) IItemRepository {
	return &ItemRepository{db: db}
}

func (r *ItemRepository) CreateItem(tx *gorm.DB, item *entity.Items) error {
	err := tx.Create(item).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *ItemRepository) CreateItems(tx *gorm.DB, items []entity.Items) error {
	if len(items) == 0 {
		return nil
	}

	err := tx.Create(&items).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *ItemRepository) GetItem(tx *gorm.DB, param model.GetItemParam) (*entity.Items, error) {
	var item entity.Items
	err := tx.Where(&param).First(&item).Error
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (r *ItemRepository) GetItemsByRequestID(tx *gorm.DB, param model.GetItemParam) ([]entity.Items, error) {
	var items []entity.Items
	err := tx.Where("request_id = ?", param.RequestID).Find(&items).Error
	if err != nil {
		return nil, err
	}

	return items, nil
}
