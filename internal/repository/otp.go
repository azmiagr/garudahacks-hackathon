package repository

import (
	"github.com/azmiagr/garudahacks-hackathon/entity"
	"github.com/azmiagr/garudahacks-hackathon/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IOtpRepository interface {
	GetOtp(tx *gorm.DB, param model.GetOtp) (*entity.OtpCode, error)
	CreateOtp(tx *gorm.DB, otp *entity.OtpCode) error
	UpdateOtp(tx *gorm.DB, otp *entity.OtpCode) error
	DeleteOtp(tx *gorm.DB, otp *entity.OtpCode) error
	DeleteOtpByUserID(tx *gorm.DB, userID uuid.UUID) error
}

type OtpRepository struct {
	db *gorm.DB
}

func NewOtpRepository(db *gorm.DB) IOtpRepository {
	return &OtpRepository{
		db: db,
	}
}

func (o *OtpRepository) GetOtp(tx *gorm.DB, param model.GetOtp) (*entity.OtpCode, error) {
	var otp *entity.OtpCode
	err := tx.Debug().Where(&param).First(&otp).Error
	if err != nil {
		return nil, err
	}

	return otp, nil
}

func (o *OtpRepository) CreateOtp(tx *gorm.DB, otp *entity.OtpCode) error {
	err := tx.Debug().Create(otp).Error
	if err != nil {
		return err
	}

	return nil
}

func (o *OtpRepository) UpdateOtp(tx *gorm.DB, otp *entity.OtpCode) error {
	err := tx.Debug().Updates(otp).Error
	if err != nil {
		return err
	}

	return nil
}

func (o *OtpRepository) DeleteOtp(tx *gorm.DB, otp *entity.OtpCode) error {
	err := tx.Debug().Delete(otp).Error
	if err != nil {
		return err
	}

	return nil
}

func (o *OtpRepository) DeleteOtpByUserID(tx *gorm.DB, userID uuid.UUID) error {
	err := tx.Debug().Where("user_id = ?", userID).Delete(&entity.OtpCode{}).Error
	if err != nil {
		return err
	}

	return nil
}
