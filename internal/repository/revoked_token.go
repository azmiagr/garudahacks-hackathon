package repository

import (
	"time"

	"github.com/azmiagr/garudahacks-hackathon/entity"
	"gorm.io/gorm"
)

type IRevokedTokenRepository interface {
	CreateRevokedToken(tx *gorm.DB, token *entity.RevokedToken) error
	ExistsActiveTokenHash(tx *gorm.DB, tokenHash string, now time.Time) (bool, error)
	DeleteExpiredTokens(tx *gorm.DB, now time.Time) error
}

type RevokedTokenRepository struct {
	db *gorm.DB
}

func NewRevokedTokenRepository(db *gorm.DB) IRevokedTokenRepository {
	return &RevokedTokenRepository{db: db}
}

func (r *RevokedTokenRepository) CreateRevokedToken(tx *gorm.DB, token *entity.RevokedToken) error {
	err := tx.FirstOrCreate(token, entity.RevokedToken{TokenHash: token.TokenHash}).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *RevokedTokenRepository) ExistsActiveTokenHash(tx *gorm.DB, tokenHash string, now time.Time) (bool, error) {
	var count int64
	err := tx.Model(&entity.RevokedToken{}).
		Where("token_hash = ? AND expires_at > ?", tokenHash, now).
		Count(&count).
		Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *RevokedTokenRepository) DeleteExpiredTokens(tx *gorm.DB, now time.Time) error {
	err := tx.Where("expires_at <= ?", now).Delete(&entity.RevokedToken{}).Error
	if err != nil {
		return err
	}

	return nil
}
