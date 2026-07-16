package repository

import (
	"github.com/azmiagr/garudahacks-hackathon/entity"
	"github.com/azmiagr/garudahacks-hackathon/model"
	"gorm.io/gorm"
)

type IRoleRepository interface {
	GetRole(tx *gorm.DB, param model.GetRoleParam) (*entity.Role, error)
}

type RoleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) IRoleRepository {
	return &RoleRepository{db: db}
}

func (r *RoleRepository) GetRole(tx *gorm.DB, param model.GetRoleParam) (*entity.Role, error) {
	var role entity.Role
	err := tx.Where(&param).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}
