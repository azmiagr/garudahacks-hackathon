package repository

import (
	"github.com/azmiagr/garudahacks-hackathon/entity"
	"github.com/azmiagr/garudahacks-hackathon/model"
	"gorm.io/gorm"
)

type IAdminPoskoProfileRepository interface {
	GetAdminPoskoProfile(tx *gorm.DB, param model.GetAdminPoskoProfileParam) (*entity.AdminProfile, error)
	CreateAdminPoskoProfile(tx *gorm.DB, profile *entity.AdminProfile) error
}

type AdminPoskoProfileRepository struct {
	db *gorm.DB
}

func NewAdminPoskoProfileRepository(db *gorm.DB) IAdminPoskoProfileRepository {
	return &AdminPoskoProfileRepository{db: db}
}

func (r *AdminPoskoProfileRepository) GetAdminPoskoProfile(tx *gorm.DB, param model.GetAdminPoskoProfileParam) (*entity.AdminProfile, error) {
	var profile entity.AdminProfile
	err := tx.Where(&param).First(&profile).Error
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (r *AdminPoskoProfileRepository) CreateAdminPoskoProfile(tx *gorm.DB, profile *entity.AdminProfile) error {
	err := tx.Create(profile).Error
	if err != nil {
		return err
	}
	return nil
}
