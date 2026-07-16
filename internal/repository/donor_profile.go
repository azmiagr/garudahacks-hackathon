package repository

import (
	"github.com/azmiagr/garudahacks-hackathon/entity"
	"github.com/azmiagr/garudahacks-hackathon/model"
	"gorm.io/gorm"
)

type IDonorProfileRepository interface {
	GetDonorProfile(tx *gorm.DB, param model.GetDonorProfileParam) (*entity.DonorProfile, error)
	CreateDonorProfile(tx *gorm.DB, profile *entity.DonorProfile) error
}

type DonorProfileRepository struct {
	db *gorm.DB
}

func NewDonorProfileRepository(db *gorm.DB) IDonorProfileRepository {
	return &DonorProfileRepository{db: db}
}

func (r *DonorProfileRepository) GetDonorProfile(tx *gorm.DB, param model.GetDonorProfileParam) (*entity.DonorProfile, error) {
	var profile entity.DonorProfile
	err := tx.Where(&param).First(&profile).Error
	if err != nil {
		return nil, err
	}

	return &profile, nil
}

func (r *DonorProfileRepository) CreateDonorProfile(tx *gorm.DB, profile *entity.DonorProfile) error {
	err := tx.Create(profile).Error
	if err != nil {
		return err
	}

	return nil
}
