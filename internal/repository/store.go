package repository

import (
	"github.com/azmiagr/garudahacks-hackathon/entity"
	"github.com/azmiagr/garudahacks-hackathon/model"
	"gorm.io/gorm"
)

type IStoreRepository interface {
	CreateStore(tx *gorm.DB, store *entity.Stores) error
	GetStore(tx *gorm.DB, param model.GetStoreParam) (*entity.Stores, error)
}

type StoreRepository struct {
	db *gorm.DB
}

func NewStoreRepository(db *gorm.DB) IStoreRepository {
	return &StoreRepository{db: db}
}

func (r *StoreRepository) CreateStore(tx *gorm.DB, store *entity.Stores) error {
	err := tx.Create(store).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *StoreRepository) GetStore(tx *gorm.DB, param model.GetStoreParam) (*entity.Stores, error) {
	var store entity.Stores
	err := tx.Where(&param).First(&store).Error
	if err != nil {
		return nil, err
	}

	return &store, nil
}
