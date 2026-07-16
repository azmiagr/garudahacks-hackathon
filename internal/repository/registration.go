package repository

import (
	"github.com/azmiagr/garudahacks-hackathon/entity"
	"github.com/azmiagr/garudahacks-hackathon/model"
	"gorm.io/gorm"
)

type IRegistrationRepository interface {
	GetRegistrationSession(tx *gorm.DB, param model.GetRegistrationSessionParam) (*entity.RegistrationSession, error)
	UpsertRegistrationSession(tx *gorm.DB, session *entity.RegistrationSession) error
	UpdateRegistrationSession(tx *gorm.DB, session *entity.RegistrationSession) error
	DeleteRegistrationSession(tx *gorm.DB, session *entity.RegistrationSession) error
}

type RegistrationRepository struct {
	db *gorm.DB
}

func NewRegistrationRepository(db *gorm.DB) IRegistrationRepository {
	return &RegistrationRepository{db: db}
}

func (r *RegistrationRepository) GetRegistrationSession(tx *gorm.DB, param model.GetRegistrationSessionParam) (*entity.RegistrationSession, error) {
	var session entity.RegistrationSession
	err := tx.Where(&param).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *RegistrationRepository) UpsertRegistrationSession(tx *gorm.DB, session *entity.RegistrationSession) error {
	err := tx.Save(session).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *RegistrationRepository) UpdateRegistrationSession(tx *gorm.DB, session *entity.RegistrationSession) error {
	err := tx.Save(session).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *RegistrationRepository) DeleteRegistrationSession(tx *gorm.DB, session *entity.RegistrationSession) error {
	err := tx.Delete(session).Error
	if err != nil {
		return err
	}

	return nil
}
