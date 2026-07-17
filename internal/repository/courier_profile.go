package repository

import (
	"github.com/azmiagr/garudahacks-hackathon/entity"
	"github.com/azmiagr/garudahacks-hackathon/model"
	"gorm.io/gorm"
)

type ICourierProfileRepository interface {
	GetCourierProfile(tx *gorm.DB, param model.GetCourierProfileParam) (*entity.CourierProfile, error)
	CreateCourierProfile(tx *gorm.DB, profile *entity.CourierProfile) error
	UpdateCourierProfilePreferences(tx *gorm.DB, userID string, isAvailable bool, urgentTaskNotificationEnabled bool) error
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

func (r *CourierProfileRepository) UpdateCourierProfilePreferences(tx *gorm.DB, userID string, isAvailable bool, urgentTaskNotificationEnabled bool) error {
	err := tx.Model(&entity.CourierProfile{}).
		Where("user_id = ?", userID).
		Updates(map[string]interface{}{
			"is_available":                     isAvailable,
			"urgent_task_notification_enabled": urgentTaskNotificationEnabled,
		}).Error
	if err != nil {
		return err
	}
	return nil
}
