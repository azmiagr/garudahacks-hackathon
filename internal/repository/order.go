package repository

import (
	"github.com/azmiagr/garudahacks-hackathon/entity"
	"gorm.io/gorm"
)

type IOrderRepository interface {
	CreateOrder(tx *gorm.DB, order *entity.Orders) error
}

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) IOrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) CreateOrder(tx *gorm.DB, order *entity.Orders) error {
	err := tx.Create(order).Error
	if err != nil {
		return err
	}

	return nil
}
