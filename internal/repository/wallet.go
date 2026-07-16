package repository

import (
	"github.com/azmiagr/garudahacks-hackathon/entity"
	"github.com/azmiagr/garudahacks-hackathon/model"
	"gorm.io/gorm"
)

type IWalletRepository interface {
	CreateWallet(tx *gorm.DB, wallet *entity.Wallets) error
	GetWallet(tx *gorm.DB, param model.GetWalletParam) (*entity.Wallets, error)
	UpdateWallet(tx *gorm.DB, wallet *entity.Wallets) error
}

type WalletRepository struct {
	db *gorm.DB
}

func NewWalletRepository(db *gorm.DB) IWalletRepository {
	return &WalletRepository{db: db}
}

func (r *WalletRepository) CreateWallet(tx *gorm.DB, wallet *entity.Wallets) error {
	err := tx.Create(wallet).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *WalletRepository) GetWallet(tx *gorm.DB, param model.GetWalletParam) (*entity.Wallets, error) {
	var wallet entity.Wallets
	err := tx.Where(&param).First(&wallet).Error
	if err != nil {
		return nil, err
	}

	return &wallet, nil
}

func (r *WalletRepository) UpdateWallet(tx *gorm.DB, wallet *entity.Wallets) error {
	err := tx.Save(wallet).Error
	if err != nil {
		return err
	}

	return nil
}
