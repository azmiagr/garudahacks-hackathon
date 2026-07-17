package service

import (
	"fmt"
	"math"
	"time"

	"github.com/azmiagr/garudahacks-hackathon/entity"
	"github.com/azmiagr/garudahacks-hackathon/internal/repository"
	"github.com/azmiagr/garudahacks-hackathon/model"
	"github.com/azmiagr/garudahacks-hackathon/pkg/database/mariadb"
	apperrors "github.com/azmiagr/garudahacks-hackathon/pkg/errors"
	"gorm.io/gorm"
)

type ICourierGoodnessService interface {
	GetCourierGoodness(user *entity.User, req model.CourierGoodnessRequest) (*model.CourierGoodnessResponse, error)
}

type CourierGoodnessService struct {
	db              *gorm.DB
	orderRepository repository.IOrderRepository
}

func NewCourierGoodnessService(orderRepository repository.IOrderRepository) ICourierGoodnessService {
	return &CourierGoodnessService{
		db:              mariadb.Connection,
		orderRepository: orderRepository,
	}
}

func (s *CourierGoodnessService) GetCourierGoodness(user *entity.User, req model.CourierGoodnessRequest) (*model.CourierGoodnessResponse, error) {
	if user == nil {
		return nil, apperrors.Unauthorized("user is required")
	}

	year := normalizeStoreGoodnessYear(req.Year)
	limit := normalizeStoreGoodnessServiceLimit(req.Limit)
	offset := normalizeStoreGoodnessServiceOffset(req.Offset)
	param := model.CourierGoodnessParam{
		CourierID: user.UserID,
		Year:      year,
		Limit:     limit,
		Offset:    offset,
	}

	stats, err := s.orderRepository.GetCourierGoodnessStats(s.db, user.UserID)
	if err != nil {
		return nil, err
	}

	historyRows, err := s.orderRepository.GetCourierDeliveryHistory(s.db, param)
	if err != nil {
		return nil, err
	}

	totalHistory, err := s.orderRepository.CountCourierDeliveryHistory(s.db, param)
	if err != nil {
		return nil, err
	}

	return &model.CourierGoodnessResponse{
		Certificate:  buildCourierGoodnessCertificate(user, *stats, year),
		History:      buildCourierDeliveryHistoryItems(historyRows),
		TotalHistory: totalHistory,
		Limit:        limit,
		Offset:       offset,
	}, nil
}

func buildCourierGoodnessCertificate(user *entity.User, row model.CourierGoodnessStatsRow, year int) model.CourierGoodnessCertificate {
	reputation := math.Round(row.ReputationScore*10) / 10

	return model.CourierGoodnessCertificate{
		CourierID:         user.UserID,
		CourierName:       user.Name,
		Title:             "Sertifikat Digital",
		PartnerLabel:      "Relawan Pengantaran",
		SinceText:         buildCourierGoodnessSinceText(row.FirstDeliveryAt, year),
		DeliveryCount:     row.DeliveryCount,
		DeliveryCountText: fmt.Sprintf("%d pengantaran selesai", row.DeliveryCount),
		TotalDistanceKm:   row.TotalDistanceKm,
		TotalDistanceText: formatCourierDistance(row.TotalDistanceKm),
		ReputationScore:   reputation,
		ReputationText:    formatReputationText(reputation),
		DisputeCount:      row.DisputeCount,
		DisputeText:       fmt.Sprintf("%d sengketa", row.DisputeCount),
		FirstDeliveryAt:   row.FirstDeliveryAt,
		ShareURL:          fmt.Sprintf("/courier/goodness/%s", user.UserID.String()),
	}
}

func buildCourierDeliveryHistoryItems(rows []model.CourierDeliveryHistoryRow) []model.CourierDeliveryHistoryItem {
	items := make([]model.CourierDeliveryHistoryItem, 0, len(rows))
	for _, row := range rows {
		distanceKm := 0.0
		if row.DeliveryDistanceKm != nil {
			distanceKm = *row.DeliveryDistanceKm
		}

		items = append(items, model.CourierDeliveryHistoryItem{
			OrderID:            row.OrderID,
			OrderCode:          row.OrderCode,
			PostName:           row.PostName,
			DisasterName:       row.DisasterName,
			Title:              fmt.Sprintf("%s %s - #%s", titleCase(row.DisasterName), row.PostName, row.OrderCode),
			ItemCount:          row.ItemCount,
			TotalAmount:        row.TotalAmount,
			DeliveryDistanceKm: distanceKm,
			DeliveredAt:        row.DeliveredAt,
			DeliveredAtText:    formatCourierDeliveredAt(row.DeliveredAt),
		})
	}
	return items
}

func buildCourierGoodnessSinceText(firstDeliveryAt *time.Time, fallbackYear int) string {
	if firstDeliveryAt == nil {
		return fmt.Sprintf("sejak %d", fallbackYear)
	}
	return "sejak " + firstDeliveryAt.Format("Jan 2006")
}

func formatCourierDistance(distanceKm float64) string {
	if distanceKm <= 0 {
		return "0 km ditempuh"
	}
	return fmt.Sprintf("%.1f km ditempuh", distanceKm)
}

func formatCourierDeliveredAt(value *time.Time) string {
	if value == nil {
		return ""
	}
	return formatGoodnessDate(*value)
}
