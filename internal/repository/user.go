package repository

import (
	"github.com/azmiagr/garudahacks-hackathon/entity"
	"github.com/azmiagr/garudahacks-hackathon/model"
	"gorm.io/gorm"
)

type IUserRepository interface {
	GetUser(tx *gorm.DB, param model.GetUserParam) (*entity.User, error)
	CreateUser(tx *gorm.DB, user *entity.User) error
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetUser(tx *gorm.DB, param model.GetUserParam) (*entity.User, error) {
	var user entity.User
	if err := tx.Where(&param).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) CreateUser(tx *gorm.DB, user *entity.User) error {
	err := tx.Create(user).Error
	if err != nil {
		return err
	}
	return nil
}
