package repository

import (
	"errors"
	"fmt"

	"github.com/azmiagr/garudahacks-hackathon/entity"
	"github.com/azmiagr/garudahacks-hackathon/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type IPointRepository interface {
	GetOrCreatePointAccount(tx *gorm.DB, userID uuid.UUID) (*entity.PointAccount, error)
	GetPointAccount(tx *gorm.DB, param model.GetPointAccountParam) (*entity.PointAccount, error)
	GetPointSummary(tx *gorm.DB, userID uuid.UUID) (*model.PointSummaryRow, error)
	CreatePointTransaction(tx *gorm.DB, pointTransaction *entity.PointTransaction) error
	HasPointTransactionSource(tx *gorm.DB, sourceType string, sourceID string) (bool, error)
	AddEarnedPoints(tx *gorm.DB, accountID uuid.UUID, points int64) error
	RedeemPoints(tx *gorm.DB, accountID uuid.UUID, points int64) error
	GetPointHistory(tx *gorm.DB, param model.PointHistoryParam) ([]model.PointHistoryRow, error)
	CreateReward(tx *gorm.DB, reward *entity.Reward) error
	GetReward(tx *gorm.DB, param model.GetRewardParam) (*entity.Reward, error)
	GetRewards(tx *gorm.DB, param model.RewardListParam) ([]entity.Reward, error)
	CreateRewardClaim(tx *gorm.DB, claim *entity.RewardClaim) error
	UpdateRewardStock(tx *gorm.DB, rewardID uuid.UUID, delta int) error
}

type PointRepository struct {
	db *gorm.DB
}

func NewPointRepository(db *gorm.DB) IPointRepository {
	return &PointRepository{db: db}
}

func (r *PointRepository) GetOrCreatePointAccount(tx *gorm.DB, userID uuid.UUID) (*entity.PointAccount, error) {
	var account entity.PointAccount
	err := tx.Where("user_id = ?", userID).First(&account).Error
	if err == nil {
		return &account, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	account = entity.PointAccount{
		PointAccountID: uuid.New(),
		UserID:         userID,
		ActivePoints:   0,
		TotalEarned:    0,
		TotalRedeemed:  0,
	}
	err = tx.Create(&account).Error
	if err != nil {
		return nil, err
	}

	return &account, nil
}

func (r *PointRepository) GetPointAccount(tx *gorm.DB, param model.GetPointAccountParam) (*entity.PointAccount, error) {
	var account entity.PointAccount
	err := tx.Where(&param).First(&account).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *PointRepository) GetPointSummary(tx *gorm.DB, userID uuid.UUID) (*model.PointSummaryRow, error) {
	var row model.PointSummaryRow
	err := tx.Table("point_accounts").
		Select("user_id, active_points, total_earned, total_redeemed").
		Where("user_id = ?", userID).
		Scan(&row).Error
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *PointRepository) CreatePointTransaction(tx *gorm.DB, pointTransaction *entity.PointTransaction) error {
	return tx.Clauses(clause.OnConflict{DoNothing: true}).Create(pointTransaction).Error
}

func (r *PointRepository) HasPointTransactionSource(tx *gorm.DB, sourceType string, sourceID string) (bool, error) {
	var count int64
	err := tx.Model(&entity.PointTransaction{}).
		Where("source_type = ? AND source_id = ?", sourceType, sourceID).
		Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *PointRepository) AddEarnedPoints(tx *gorm.DB, accountID uuid.UUID, points int64) error {
	err := tx.Model(&entity.PointAccount{}).
		Where("point_account_id = ?", accountID).
		Updates(map[string]interface{}{
			"active_points": gorm.Expr("active_points + ?", points),
			"total_earned":  gorm.Expr("total_earned + ?", points),
		}).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *PointRepository) RedeemPoints(tx *gorm.DB, accountID uuid.UUID, points int64) error {
	result := tx.Model(&entity.PointAccount{}).
		Where("point_account_id = ? AND active_points >= ?", accountID, points).
		Updates(map[string]interface{}{
			"active_points":  gorm.Expr("active_points - ?", points),
			"total_redeemed": gorm.Expr("total_redeemed + ?", points),
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("not enough points")
	}

	return nil
}

func (r *PointRepository) GetPointHistory(tx *gorm.DB, param model.PointHistoryParam) ([]model.PointHistoryRow, error) {
	var rows []model.PointHistoryRow
	err := tx.Table("point_transactions").
		Select("point_transaction_id, donation_id, reward_claim_id, points, transaction_type, source_type, source_id, description, expires_at, created_at").
		Where("user_id = ?", param.UserID).
		Order("created_at DESC").
		Limit(normalizePointLimit(param.Limit)).
		Offset(normalizePointOffset(param.Offset)).
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func (r *PointRepository) CreateReward(tx *gorm.DB, reward *entity.Reward) error {
	err := tx.Create(reward).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *PointRepository) GetReward(tx *gorm.DB, param model.GetRewardParam) (*entity.Reward, error) {
	var reward entity.Reward
	query := tx.Where("reward_id = ?", param.RewardID)
	if param.IsActive != nil {
		query = query.Where("is_active = ?", *param.IsActive)
	}
	err := query.First(&reward).Error
	if err != nil {
		return nil, err
	}

	return &reward, nil
}

func (r *PointRepository) GetRewards(tx *gorm.DB, param model.RewardListParam) ([]entity.Reward, error) {
	var rewards []entity.Reward
	query := tx.Model(&entity.Reward{})
	if param.OnlyActive {
		query = query.Where("is_active = ?", true)
	}
	err := query.Order("points_cost ASC").
		Limit(normalizePointLimit(param.Limit)).
		Offset(normalizePointOffset(param.Offset)).
		Find(&rewards).Error
	if err != nil {
		return nil, err
	}

	return rewards, nil
}

func (r *PointRepository) UpdateRewardStock(tx *gorm.DB, rewardID uuid.UUID, delta int) error {
	query := tx.Model(&entity.Reward{}).Where("reward_id = ?", rewardID)
	if delta < 0 {
		query = query.Where("stock >= ?", -delta)
	}

	result := query.Update("stock", gorm.Expr("stock + ?", delta))
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("reward stock is not enough")
	}

	return nil
}

func (r *PointRepository) CreateRewardClaim(tx *gorm.DB, claim *entity.RewardClaim) error {
	err := tx.Create(claim).Error
	if err != nil {
		return err
	}

	return nil
}

func normalizePointLimit(limit int) int {
	if limit <= 0 {
		return 20
	}
	if limit > 100 {
		return 100
	}
	return limit
}

func normalizePointOffset(offset int) int {
	if offset < 0 {
		return 0
	}
	return offset
}
