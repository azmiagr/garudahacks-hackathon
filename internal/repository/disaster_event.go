package repository

import (
	"github.com/azmiagr/garudahacks-hackathon/entity"
	"github.com/azmiagr/garudahacks-hackathon/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IDisasterEventRepository interface {
	GetDisasterEvent(tx *gorm.DB, param model.GetDisasterEventParam) (*entity.DisasterEvent, error)
	GetDisasterEventByName(tx *gorm.DB, name string) (*entity.DisasterEvent, error)
	GetEventByIDs(tx *gorm.DB, eventIDs []uuid.UUID) ([]model.DisasterEventRow, error)
}

type DisasterEventRepository struct {
	db *gorm.DB
}

func NewDisasterEventRepository(db *gorm.DB) IDisasterEventRepository {
	return &DisasterEventRepository{db: db}
}

func (r *DisasterEventRepository) GetDisasterEvent(tx *gorm.DB, param model.GetDisasterEventParam) (*entity.DisasterEvent, error) {
	var event entity.DisasterEvent
	err := tx.Where(&param).First(&event).Error
	if err != nil {
		return nil, err
	}

	return &event, nil
}

func (r *DisasterEventRepository) GetEventByIDs(tx *gorm.DB, eventIDs []uuid.UUID) ([]model.DisasterEventRow, error) {
	var rows []model.DisasterEventRow

	if len(eventIDs) == 0 {
		return rows, nil
	}

	err := tx.Table("disaster_events").
		Select("event_id, name").
		Where("event_id IN ?", eventIDs).
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func (r *DisasterEventRepository) GetDisasterEventByName(tx *gorm.DB, name string) (*entity.DisasterEvent, error) {
	var event entity.DisasterEvent
	err := tx.Where("LOWER(name) = LOWER(?)", name).First(&event).Error
	if err != nil {
		return nil, err
	}

	return &event, nil
}
