package repository

import (
	"github.com/azmiagr/garudahacks-hackathon/entity"
	"gorm.io/gorm"
)

type IOrderItemRepository interface {
	CreateOrderItems(tx *gorm.DB, orderItems []entity.OrderItems) error
}

type OrderItemRepository struct {
	db *gorm.DB
}

func NewOrderItemRepository(db *gorm.DB) IOrderItemRepository {
	return &OrderItemRepository{db: db}
}

func (r *OrderItemRepository) CreateOrderItems(tx *gorm.DB, orderItems []entity.OrderItems) error {
	if len(orderItems) == 0 {
		return nil
	}

	err := tx.Create(&orderItems).Error
	if err != nil {
		return err
	}
	return nil
}
