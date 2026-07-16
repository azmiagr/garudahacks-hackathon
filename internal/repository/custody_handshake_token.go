package repository

import (
	"time"

	"github.com/azmiagr/garudahacks-hackathon/entity"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ICustodyHandshakeTokenRepository interface {
	CreateToken(tx *gorm.DB, token *entity.CustodyHandshakeToken) error
	GetActiveTokenByQRHashForUpdate(tx *gorm.DB, qrPayloadHash string, now time.Time) (*entity.CustodyHandshakeToken, error)
	GetActiveTokenByPINHashForUpdate(tx *gorm.DB, orderID uuid.UUID, handoffStage string, pinHash string, now time.Time) (*entity.CustodyHandshakeToken, error)
	MarkTokenUsed(tx *gorm.DB, tokenID uuid.UUID, usedBy uuid.UUID, usedAt time.Time) error
	ExpireTokens(tx *gorm.DB, now time.Time) error
}

type CustodyHandshakeTokenRepository struct {
	db *gorm.DB
}

func NewCustodyHandshakeTokenRepository(db *gorm.DB) ICustodyHandshakeTokenRepository {
	return &CustodyHandshakeTokenRepository{db: db}
}

func (r *CustodyHandshakeTokenRepository) CreateToken(tx *gorm.DB, token *entity.CustodyHandshakeToken) error {
	return tx.Create(token).Error
}

func (r *CustodyHandshakeTokenRepository) GetActiveTokenByQRHashForUpdate(tx *gorm.DB, qrPayloadHash string, now time.Time) (*entity.CustodyHandshakeToken, error) {
	var token entity.CustodyHandshakeToken
	err := activeTokenBaseQuery(tx, now).
		Where("qr_payload_hash = ?", qrPayloadHash).
		First(&token).Error
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (r *CustodyHandshakeTokenRepository) GetActiveTokenByPINHashForUpdate(tx *gorm.DB, orderID uuid.UUID, handoffStage string, pinHash string, now time.Time) (*entity.CustodyHandshakeToken, error) {
	var token entity.CustodyHandshakeToken
	err := activeTokenBaseQuery(tx, now).
		Where("order_id = ? AND handoff_stage = ? AND pin_hash = ?", orderID, handoffStage, pinHash).
		First(&token).Error
	if err != nil {
		return nil, err
	}

	return &token, nil
}

func (r *CustodyHandshakeTokenRepository) MarkTokenUsed(tx *gorm.DB, tokenID uuid.UUID, usedBy uuid.UUID, usedAt time.Time) error {
	result := tx.Model(&entity.CustodyHandshakeToken{}).
		Where("token_id = ? AND status = ?", tokenID, entity.CustodyTokenStatusActive).
		Updates(map[string]interface{}{
			"status":  entity.CustodyTokenStatusUsed,
			"used_by": usedBy,
			"used_at": usedAt,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (r *CustodyHandshakeTokenRepository) ExpireTokens(tx *gorm.DB, now time.Time) error {
	err := tx.Model(&entity.CustodyHandshakeToken{}).
		Where("status = ? AND cache_valid_until < ?", entity.CustodyTokenStatusActive, now).
		Update("status", entity.CustodyTokenStatusExpired).Error
	if err != nil {
		return err
	}

	return nil
}

func activeTokenBaseQuery(tx *gorm.DB, now time.Time) *gorm.DB {
	return tx.Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("status = ?", entity.CustodyTokenStatusActive).
		Where("cache_valid_until >= ?", now)
}
