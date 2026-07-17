package repository

import (
	"github.com/azmiagr/garudahacks-hackathon/entity"
	"github.com/azmiagr/garudahacks-hackathon/model"
	"gorm.io/gorm"
)

type ICourierProfileRepository interface {
	GetCourierProfile(tx *gorm.DB, param model.GetCourierProfileParam) (*entity.CourierProfile, error)
	CreateCourierProfile(tx *gorm.DB, profile *entity.CourierProfile) error
}

type CourierProfileRepository struct {
	db *gorm.DB
}

func NewCourierProfileRepository(db *gorm.DB) ICourierProfileRepository {
	return &CourierProfileRepository{db: db}
}

func (r *CourierProfileRepository) GetCourierProfile(tx *gorm.DB, param model.GetCourierProfileParam) (*entity.CourierProfile, error) {
	var profile entity.CourierProfile
	err := tx.Where(&param).First(&profile).Error
	if err != nil {
		return nil, err
	}

	return &profile, nil
}

func (r *CourierProfileRepository) CreateCourierProfile(tx *gorm.DB, profile *entity.CourierProfile) error {
	err := tx.Create(profile).Error
	if err != nil {
		return err
	}

	return nil
}
