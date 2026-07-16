package repository

import (
	"github.com/azmiagr/garudahacks-hackathon/entity"
	"github.com/azmiagr/garudahacks-hackathon/model"
	"gorm.io/gorm"
)

type IWalletTransactionRepository interface {
	CreateWalletTransaction(tx *gorm.DB, walletTransaction *entity.WalletTransactions) error
	GetWalletTransaction(tx *gorm.DB, param model.GetWalletTransactionParam) (*entity.WalletTransactions, error)
	UpdateWalletTransaction(tx *gorm.DB, walletTransaction *entity.WalletTransactions) error
}

type WalletTransactionRepository struct {
	db *gorm.DB
}

func NewWalletTransactionRepository(db *gorm.DB) IWalletTransactionRepository {
	return &WalletTransactionRepository{db: db}
}

func (r *WalletTransactionRepository) CreateWalletTransaction(tx *gorm.DB, walletTransaction *entity.WalletTransactions) error {
	err := tx.Create(walletTransaction).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *WalletTransactionRepository) GetWalletTransaction(tx *gorm.DB, param model.GetWalletTransactionParam) (*entity.WalletTransactions, error) {
	var walletTransaction entity.WalletTransactions
	err := tx.Where(&param).First(&walletTransaction).Error
	if err != nil {
		return nil, err
	}

	return &walletTransaction, nil
}

func (r *WalletTransactionRepository) UpdateWalletTransaction(tx *gorm.DB, walletTransaction *entity.WalletTransactions) error {
	err := tx.Save(walletTransaction).Error
	if err != nil {
		return err
	}
	return nil
}
