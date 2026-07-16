package service

import (
	"errors"
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/azmiagr/garudahacks-hackathon/internal/repository"
	"github.com/azmiagr/garudahacks-hackathon/model"
	"github.com/azmiagr/garudahacks-hackathon/pkg/database/mariadb"
	apperrors "github.com/azmiagr/garudahacks-hackathon/pkg/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IDonorDashboardService interface {
	GetDonorMap(param model.DonorDashboardMapParam) (*model.DonorDashboardMapResponse, error)
	GetPostDetail(postID uuid.UUID) (*model.DonorPostDetailResponse, error)
}

type DonorDashboardService struct {
	db                     *gorm.DB
	postRepository         repository.IPostRepository
	itemRepository         repository.IItemRepository
	publicDashboardService IPublicDashboardService
}

func NewDonorDashboardService(
	postRepository repository.IPostRepository,
	itemRepository repository.IItemRepository,
	publicDashboardService IPublicDashboardService,
) IDonorDashboardService {
	return &DonorDashboardService{
		db:                     mariadb.Connection,
		postRepository:         postRepository,
		itemRepository:         itemRepository,
		publicDashboardService: publicDashboardService,
	}
}

func (s *DonorDashboardService) GetDonorMap(param model.DonorDashboardMapParam) (*model.DonorDashboardMapResponse, error) {
	publicMap, err := s.publicDashboardService.GetPublicMap(model.PublicDashboardParam{
		Query:        param.Query,
		DisasterType: param.DisasterType,
		MinLatitude:  param.MinLatitude,
		MinLongitude: param.MinLongitude,
		MaxLatitude:  param.MaxLatitude,
		MaxLongitude: param.MaxLongitude,
		Limit:        param.Limit,
	})
	if err != nil {
		return nil, err
	}

	heatmapPoints := make([]model.DonorHeatmapPoint, 0, len(publicMap.Items))
	urgentItems := make([]model.DonorUrgentPost, 0, len(publicMap.Items))

	for _, item := range publicMap.Items {
		heatmapPoints = append(heatmapPoints, model.DonorHeatmapPoint{
			PostID:            item.PostID,
			Name:              item.Name,
			Latitude:          item.Latitude,
			Longitude:         item.Longitude,
			DisasterType:      item.DisasterEvent,
			FundingPercentage: item.FundingPercentage,
			UrgencyLevel:      item.UrgencyLevel,
			Color:             donorUrgencyColor(item.FundingPercentage),
		})

		urgentItems = append(urgentItems, model.DonorUrgentPost{
			PostID:            item.PostID,
			Name:              item.Name,
			Address:           item.Address,
			ImageURL:          item.ImageURL,
			FundingTarget:     item.FundingTarget,
			FundedAmount:      item.FundedAmount,
			FundingPercentage: item.FundingPercentage,
			FundingText:       buildDonorFundingText(item.FundedAmount, item.FundingTarget, item.FundingPercentage),
			ElapsedText:       buildDonorElapsedText(item.LatestReportedAt),
		})
	}

	sort.SliceStable(urgentItems, func(i, j int) bool {
		if urgentItems[i].FundingPercentage == urgentItems[j].FundingPercentage {
			return urgentItems[i].FundingTarget > urgentItems[j].FundingTarget
		}

		return urgentItems[i].FundingPercentage < urgentItems[j].FundingPercentage
	})

	if len(urgentItems) > 5 {
		urgentItems = urgentItems[:5]
	}

	return &model.DonorDashboardMapResponse{
		HeatmapPoints: heatmapPoints,
		UrgentPosts:   urgentItems,
		Legend:        buildDonorMapLegend(),
	}, nil
}

func (s *DonorDashboardService) GetPostDetail(postID uuid.UUID) (*model.DonorPostDetailResponse, error) {
	row, err := s.postRepository.GetDonorPostDetail(s.db, postID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.NotFound("post not found")
		}
		return nil, err
	}

	itemRows, err := s.itemRepository.GetDonorPostDetailItems(s.db, row.RequestID)
	if err != nil {
		return nil, err
	}

	percentage := calculateDonorFundingPercentage(row.FundedAmount, row.FundingTarget)
	adminVerified := row.AdminKYCStatus == "approved"

	return &model.DonorPostDetailResponse{
		PostID:                row.PostID,
		ReportID:              row.ReportID,
		RequestID:             row.RequestID,
		Name:                  row.Name,
		Address:               row.Address,
		DisasterType:          row.DisasterType,
		ImageURL:              row.ImageURL,
		ElapsedText:           buildDonorElapsedText(resolveDonorReportedAt(row)),
		AdminVerified:         adminVerified,
		AdminVerificationText: buildDonorAdminVerificationText(adminVerified),
		FundingTarget:         row.FundingTarget,
		FundedAmount:          row.FundedAmount,
		FundingPercentage:     percentage,
		FundingText:           buildDonorFundingText(row.FundedAmount, row.FundingTarget, percentage),
		DonorCount:            row.DonorCount,
		UrgencyLevel:          resolveDonorUrgencyLevel(percentage),
		Items:                 buildDonorPostDetailItems(itemRows),
	}, nil
}

func buildDonorPostDetailItems(rows []model.DonorPostDetailItemRow) []model.DonorPostDetailItem {
	items := make([]model.DonorPostDetailItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, model.DonorPostDetailItem{
			ItemID:            row.ItemID,
			Name:              row.Name,
			Description:       row.Description,
			Price:             row.Price,
			EstimatedTotal:    row.EstimatedTotal,
			QuantityNeeded:    row.QuantityNeeded,
			QuantityFulfilled: row.QuantityFulfilled,
			ProgressText:      fmt.Sprintf("%d/%d", row.QuantityFulfilled, row.QuantityNeeded),
		})
	}

	return items
}

func buildDonorMapLegend() []model.DonorMapLegend {
	return []model.DonorMapLegend{
		{Label: "0-20% terdana", Level: "critical", Color: "#8B1E23"},
		{Label: "21-50%", Level: "high", Color: "#E53935"},
		{Label: "51-80%", Level: "medium", Color: "#FF7A1A"},
		{Label: "81-99%", Level: "low", Color: "#F4C20D"},
		{Label: "Terdanai penuh", Level: "funded", Color: "#66BB6A"},
	}
}

func donorUrgencyColor(percentage int) string {
	switch {
	case percentage >= 100:
		return "#66BB6A"
	case percentage >= 81:
		return "#F4C20D"
	case percentage >= 51:
		return "#FF7A1A"
	case percentage >= 21:
		return "#E53935"
	default:
		return "#8B1E23"
	}
}

func resolveDonorUrgencyLevel(percentage int) string {
	switch {
	case percentage >= 100:
		return "funded"
	case percentage >= 81:
		return "low"
	case percentage >= 51:
		return "medium"
	case percentage >= 21:
		return "high"
	default:
		return "critical"
	}
}

func calculateDonorFundingPercentage(fundedAmount float64, fundingTarget float64) int {
	if fundingTarget <= 0 {
		return 0
	}

	percentage := int(math.Round((fundedAmount / fundingTarget) * 100))
	if percentage < 0 {
		return 0
	}

	if percentage > 100 {
		return 100
	}

	return percentage
}

func buildDonorFundingText(fundedAmount float64, fundingTarget float64, percentage int) string {
	return fmt.Sprintf("%s / %s - %d%%", formatDonorRupiahShort(fundedAmount), formatDonorRupiahShort(fundingTarget), percentage)
}

func formatDonorRupiahShort(amount float64) string {
	switch {
	case amount >= 1000000000:
		return fmt.Sprintf("Rp%.1fM", amount/1000000000)
	case amount >= 1000000:
		return fmt.Sprintf("Rp%.1fjt", amount/1000000)
	case amount >= 1000:
		return fmt.Sprintf("Rp%.0frb", amount/1000)
	default:
		return fmt.Sprintf("Rp%.0f", amount)
	}
}

func buildDonorElapsedText(startedAt time.Time) string {
	if startedAt.IsZero() {
		return ""
	}

	duration := time.Since(startedAt)
	if duration < time.Hour {
		minutes := int(duration.Minutes())
		if minutes <= 0 {
			minutes = 1
		}
		return fmt.Sprintf("%d menit lalu", minutes)
	}

	if duration < 24*time.Hour {
		return fmt.Sprintf("%d jam lalu", int(duration.Hours()))
	}

	return fmt.Sprintf("%d hari lalu", int(duration.Hours()/24))
}

func buildDonorAdminVerificationText(isVerified bool) string {
	if isVerified {
		return "Admin terverifikasi"
	}

	return "Admin menunggu verifikasi"
}

func resolveDonorReportedAt(row *model.DonorPostDetailRow) time.Time {
	if !row.ReportedAt.IsZero() {
		return row.ReportedAt
	}

	return row.CreatedAt
}
