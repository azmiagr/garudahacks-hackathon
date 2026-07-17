package service

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/azmiagr/garudahacks-hackathon/entity"
	"github.com/azmiagr/garudahacks-hackathon/internal/repository"
	"github.com/azmiagr/garudahacks-hackathon/model"
	"github.com/azmiagr/garudahacks-hackathon/pkg/database/mariadb"
	apperrors "github.com/azmiagr/garudahacks-hackathon/pkg/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	pointSourceDonationPayment = "donation_payment"
	pointSourceCourierDelivery = "courier_delivery"
	pointSourceRewardClaim     = "reward_claim"

	pointTransactionEarn   = "earn"
	pointTransactionRedeem = "redeem"

	rewardClaimStatusClaimed = "claimed"

	pointsPerRupiahDivider   = 1000
	pointsPerCourierDelivery = 100
	pointsPerCourierKm       = 10
	defaultPointExpiryDays   = 365
)

type IPointService interface {
	GetDashboard(user *entity.User, param model.PointDashboardParam) (*model.PointDashboardResponse, error)
	GetHistory(user *entity.User, param model.PointHistoryQueryParam) (*model.PointHistoryResponse, error)
	GetRewards(user *entity.User, param model.RewardQueryParam) (*model.RewardListResponse, error)
	ClaimReward(user *entity.User, req model.ClaimRewardRequest) (*model.RewardClaimResponse, error)
	AwardDonationPaymentPoints(tx *gorm.DB, payment *entity.PaymentTransactions) error
	AwardCourierDeliveryPoints(tx *gorm.DB, order *entity.Orders) (int64, error)
}

type PointService struct {
	db              *gorm.DB
	pointRepository repository.IPointRepository
}

func NewPointService(pointRepository repository.IPointRepository) IPointService {
	return &PointService{
		db:              mariadb.Connection,
		pointRepository: pointRepository,
	}
}

func (s *PointService) GetDashboard(user *entity.User, param model.PointDashboardParam) (*model.PointDashboardResponse, error) {
	if user == nil {
		return nil, apperrors.Unauthorized("user is required")
	}

	account, err := s.pointRepository.GetOrCreatePointAccount(s.db, user.UserID)
	if err != nil {
		return nil, err
	}

	rewards, err := s.pointRepository.GetRewards(s.db, model.RewardListParam{
		OnlyActive: true,
		Limit:      normalizePointServiceLimit(param.RewardLimit),
		Offset:     0,
	})
	if err != nil {
		return nil, err
	}

	historyRows, err := s.pointRepository.GetPointHistory(s.db, model.PointHistoryParam{
		UserID: user.UserID,
		Limit:  normalizePointServiceLimit(param.HistoryLimit),
		Offset: 0,
	})
	if err != nil {
		return nil, err
	}

	level, nextLevelPoints := resolvePointLevel(account.ActivePoints)

	return &model.PointDashboardResponse{
		ActivePoints:      account.ActivePoints,
		TotalEarned:       account.TotalEarned,
		TotalRedeemed:     account.TotalRedeemed,
		Level:             level,
		NextLevelPoints:   nextLevelPoints,
		PointsToNextLevel: calculatePointsToNextLevel(account.ActivePoints, nextLevelPoints),
		Rewards:           buildRewardItems(rewards, account.ActivePoints),
		History:           buildPointHistoryItems(historyRows),
	}, nil
}

func (s *PointService) GetHistory(user *entity.User, param model.PointHistoryQueryParam) (*model.PointHistoryResponse, error) {
	if user == nil {
		return nil, apperrors.Unauthorized("user is required")
	}

	limit := normalizePointServiceLimit(param.Limit)
	offset := normalizePointServiceOffset(param.Offset)

	rows, err := s.pointRepository.GetPointHistory(s.db, model.PointHistoryParam{
		UserID: user.UserID,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	return &model.PointHistoryResponse{
		Items:  buildPointHistoryItems(rows),
		Limit:  limit,
		Offset: offset,
	}, nil
}

func (s *PointService) GetRewards(user *entity.User, param model.RewardQueryParam) (*model.RewardListResponse, error) {
	if user == nil {
		return nil, apperrors.Unauthorized("user is required")
	}

	account, err := s.pointRepository.GetOrCreatePointAccount(s.db, user.UserID)
	if err != nil {
		return nil, err
	}

	limit := normalizePointServiceLimit(param.Limit)
	offset := normalizePointServiceOffset(param.Offset)

	rewards, err := s.pointRepository.GetRewards(s.db, model.RewardListParam{
		OnlyActive: true,
		Limit:      limit,
		Offset:     offset,
	})
	if err != nil {
		return nil, err
	}

	return &model.RewardListResponse{
		Items:  buildRewardItems(rewards, account.ActivePoints),
		Limit:  limit,
		Offset: offset,
	}, nil
}

func (s *PointService) ClaimReward(user *entity.User, req model.ClaimRewardRequest) (*model.RewardClaimResponse, error) {
	if user == nil {
		return nil, apperrors.Unauthorized("user is required")
	}
	if req.RewardID == uuid.Nil {
		return nil, apperrors.BadRequest("reward_id is required")
	}

	tx := s.db.Begin()
	defer tx.Rollback()

	account, err := s.pointRepository.GetOrCreatePointAccount(tx, user.UserID)
	if err != nil {
		return nil, err
	}

	active := true
	reward, err := s.pointRepository.GetReward(tx, model.GetRewardParam{
		RewardID: req.RewardID,
		IsActive: &active,
	})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NotFound("reward not found")
		}
		return nil, err
	}

	if reward.Stock <= 0 {
		return nil, apperrors.BadRequest("reward is out of stock")
	}
	if account.ActivePoints < reward.PointsCost {
		return nil, apperrors.BadRequest("not enough points")
	}

	now := time.Now().UTC()
	claim := &entity.RewardClaim{
		RewardClaimID:  uuid.New(),
		PointAccountID: account.PointAccountID,
		UserID:         user.UserID,
		RewardID:       reward.RewardID,
		PointsSpent:    reward.PointsCost,
		ClaimStatus:    rewardClaimStatusClaimed,
		ClaimedAt:      now,
		UpdatedAt:      now,
	}

	err = s.pointRepository.CreateRewardClaim(tx, claim)
	if err != nil {
		return nil, err
	}

	err = s.pointRepository.RedeemPoints(tx, account.PointAccountID, reward.PointsCost)
	if err != nil {
		return nil, err
	}

	err = s.pointRepository.UpdateRewardStock(tx, reward.RewardID, -1)
	if err != nil {
		return nil, err
	}

	pointTransaction := &entity.PointTransaction{
		PointTransactionID: uuid.New(),
		PointAccountID:     account.PointAccountID,
		UserID:             user.UserID,
		RewardClaimID:      &claim.RewardClaimID,
		Points:             -reward.PointsCost,
		TransactionType:    pointTransactionRedeem,
		SourceType:         pointSourceRewardClaim,
		SourceID:           claim.RewardClaimID.String(),
		Description:        fmt.Sprintf("Tukar %s", reward.Name),
		CreatedAt:          now,
	}
	err = s.pointRepository.CreatePointTransaction(tx, pointTransaction)
	if err != nil {
		return nil, err
	}

	err = tx.Commit().Error
	if err != nil {
		return nil, err
	}

	return &model.RewardClaimResponse{
		RewardClaimID:   claim.RewardClaimID,
		RewardID:        reward.RewardID,
		RewardName:      reward.Name,
		PointsSpent:     reward.PointsCost,
		RemainingPoints: account.ActivePoints - reward.PointsCost,
		ClaimStatus:     claim.ClaimStatus,
		ClaimedAt:       claim.ClaimedAt,
	}, nil
}

func (s *PointService) AwardDonationPaymentPoints(tx *gorm.DB, payment *entity.PaymentTransactions) error {
	if payment == nil {
		return nil
	}

	points := calculateDonationPoints(payment.Amount)
	if points <= 0 {
		return nil
	}

	sourceID := payment.DonationID.String()
	exists, err := s.pointRepository.HasPointTransactionSource(tx, pointSourceDonationPayment, sourceID)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	account, err := s.pointRepository.GetOrCreatePointAccount(tx, payment.UserID)
	if err != nil {
		return err
	}

	expiresAt := time.Now().UTC().AddDate(0, 0, defaultPointExpiryDays)
	pointTransaction := &entity.PointTransaction{
		PointTransactionID: uuid.New(),
		PointAccountID:     account.PointAccountID,
		UserID:             payment.UserID,
		DonationID:         &payment.DonationID,
		Points:             points,
		TransactionType:    pointTransactionEarn,
		SourceType:         pointSourceDonationPayment,
		SourceID:           sourceID,
		Description:        fmt.Sprintf("Poin donasi %s", payment.OrderID),
		ExpiresAt:          &expiresAt,
		CreatedAt:          time.Now().UTC(),
	}

	err = s.pointRepository.CreatePointTransaction(tx, pointTransaction)
	if err != nil {
		return err
	}

	return s.pointRepository.AddEarnedPoints(tx, account.PointAccountID, points)
}

func (s *PointService) AwardCourierDeliveryPoints(tx *gorm.DB, order *entity.Orders) (int64, error) {
	if order == nil || order.CourierID == uuid.Nil {
		return 0, nil
	}

	points := calculateCourierDeliveryPoints(order.DeliveryDistanceKm)
	if points <= 0 {
		return 0, nil
	}

	sourceID := order.OrderID.String()
	exists, err := s.pointRepository.HasPointTransactionSource(tx, pointSourceCourierDelivery, sourceID)
	if err != nil {
		return 0, err
	}
	if exists {
		return 0, nil
	}

	account, err := s.pointRepository.GetOrCreatePointAccount(tx, order.CourierID)
	if err != nil {
		return 0, err
	}

	expiresAt := time.Now().UTC().AddDate(0, 0, defaultPointExpiryDays)
	pointTransaction := &entity.PointTransaction{
		PointTransactionID: uuid.New(),
		PointAccountID:     account.PointAccountID,
		UserID:             order.CourierID,
		Points:             points,
		TransactionType:    pointTransactionEarn,
		SourceType:         pointSourceCourierDelivery,
		SourceID:           sourceID,
		Description:        fmt.Sprintf("Poin pengantaran #%s", order.OrderCode),
		ExpiresAt:          &expiresAt,
		CreatedAt:          time.Now().UTC(),
	}

	err = s.pointRepository.CreatePointTransaction(tx, pointTransaction)
	if err != nil {
		return 0, err
	}

	if err := s.pointRepository.AddEarnedPoints(tx, account.PointAccountID, points); err != nil {
		return 0, err
	}

	return points, nil
}

func calculateDonationPoints(amount float64) int64 {
	return int64(math.Floor(amount / pointsPerRupiahDivider))
}

func calculateCourierDeliveryPoints(distanceKm *float64) int64 {
	points := int64(pointsPerCourierDelivery)
	if distanceKm != nil && *distanceKm > 0 {
		points += int64(math.Floor(*distanceKm * pointsPerCourierKm))
	}
	return points
}

func buildRewardItems(rewards []entity.Reward, activePoints int64) []model.RewardItem {
	items := make([]model.RewardItem, 0, len(rewards))
	for _, reward := range rewards {
		items = append(items, model.RewardItem{
			RewardID:     reward.RewardID,
			Name:         reward.Name,
			Description:  reward.Description,
			RewardType:   reward.RewardType,
			PointsCost:   reward.PointsCost,
			Stock:        reward.Stock,
			IsActive:     reward.IsActive,
			ValidityDays: reward.ValidityDays,
			CanClaim:     reward.IsActive && reward.Stock > 0 && activePoints >= reward.PointsCost,
		})
	}
	return items
}

func buildPointHistoryItems(rows []model.PointHistoryRow) []model.PointHistoryItem {
	items := make([]model.PointHistoryItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, model.PointHistoryItem{
			PointTransactionID: row.PointTransactionID,
			DonationID:         row.DonationID,
			RewardClaimID:      row.RewardClaimID,
			Points:             row.Points,
			TransactionType:    row.TransactionType,
			SourceType:         row.SourceType,
			SourceID:           row.SourceID,
			Description:        row.Description,
			ExpiresAt:          row.ExpiresAt,
			CreatedAt:          row.CreatedAt,
			PointsText:         formatPointDelta(row.Points),
		})
	}
	return items
}

func formatPointDelta(points int64) string {
	if points > 0 {
		return fmt.Sprintf("+%d", points)
	}
	return fmt.Sprintf("%d", points)
}

func resolvePointLevel(activePoints int64) (string, int64) {
	switch {
	case activePoints >= 5000:
		return "Legenda Kebaikan", 0
	case activePoints >= 2500:
		return "Penerang Nusantara", 5000
	case activePoints >= 1000:
		return "Sahabat Posko", 2500
	default:
		return "Relawan Kebaikan", 1000
	}
}

func calculatePointsToNextLevel(activePoints int64, nextLevelPoints int64) int64 {
	if nextLevelPoints <= 0 || activePoints >= nextLevelPoints {
		return 0
	}
	return nextLevelPoints - activePoints
}

func normalizePointServiceLimit(limit int) int {
	if limit <= 0 {
		return 20
	}
	if limit > 100 {
		return 100
	}
	return limit
}

func normalizePointServiceOffset(offset int) int {
	if offset < 0 {
		return 0
	}
	return offset
}
