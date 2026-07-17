package repository

import (
	"github.com/azmiagr/garudahacks-hackathon/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IOrderItemRepository interface {
	CreateOrderItems(tx *gorm.DB, orderItems []entity.OrderItems) error
	CountDistinctItemsByOrder(tx *gorm.DB, orderID uuid.UUID) (int64, error)
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

func (r *OrderItemRepository) CountDistinctItemsByOrder(tx *gorm.DB, orderID uuid.UUID) (int64, error) {
	var count int64
	err := tx.Model(&entity.OrderItems{}).
		Where("order_id = ?", orderID).
		Distinct("item_id").
		Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}
