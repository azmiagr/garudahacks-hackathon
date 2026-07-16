package repository

import (
	"github.com/azmiagr/garudahacks-hackathon/entity"
	"github.com/azmiagr/garudahacks-hackathon/model"
	"gorm.io/gorm"
)

type IPaymentTransactionRepository interface {
	CreatePaymentTransaction(tx *gorm.DB, payment *entity.PaymentTransactions) error
	GetPaymentTransaction(tx *gorm.DB, param model.GetPaymentTransactionParam) (*entity.PaymentTransactions, error)
	UpdatePaymentTransaction(tx *gorm.DB, payment *entity.PaymentTransactions) error
}

type PaymentTransactionRepository struct {
	db *gorm.DB
}

func NewPaymentTransactionRepository(db *gorm.DB) IPaymentTransactionRepository {
	return &PaymentTransactionRepository{db: db}
}

func (r *PaymentTransactionRepository) CreatePaymentTransaction(tx *gorm.DB, payment *entity.PaymentTransactions) error {
	err := tx.Create(payment).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *PaymentTransactionRepository) GetPaymentTransaction(tx *gorm.DB, param model.GetPaymentTransactionParam) (*entity.PaymentTransactions, error) {
	var payment entity.PaymentTransactions
	err := tx.Where(&param).First(&payment).Error
	if err != nil {
		return nil, err
	}

	return &payment, nil
}

func (r *PaymentTransactionRepository) UpdatePaymentTransaction(tx *gorm.DB, payment *entity.PaymentTransactions) error {
	err := tx.Save(payment).Error
	if err != nil {
		return err
	}

	return nil
}
