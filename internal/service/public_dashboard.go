package service

import (
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

type IPublicDashboardService interface {
	GetPublicMap(param model.PublicDashboardParam) (*model.PublicDashboardMapResponse, error)
	GetSummary(param model.PublicDashboardParam) (*model.PublicDashboardSummaryResponse, error)
	GetDistributions(param model.PublicDistributionParam) (*model.PublicDistributionResponse, error)
	GetTransparency(param model.PublicTransparencyParam) (*model.PublicTransparencyResponse, error)
}

type PublicDashboardService struct {
	db                       *gorm.DB
	postRepo                 repository.IPostRepository
	disasterReportRepo       repository.IDisasterReportRepository
	disasterEventRepo        repository.IDisasterEventRepository
	requestRepo              repository.IRequestRepository
	deliveryVerificationRepo repository.IDeliveryVerificationRepository
	donationRepo             repository.IDonationRepository
	disbursementRepo         repository.IDisbursementRepository
	custodyLogRepo           repository.ICustodyLogRepository
}

func NewPublicDashboardService(
	postRepo repository.IPostRepository,
	disasterReportRepo repository.IDisasterReportRepository,
	disasterEventRepo repository.IDisasterEventRepository,
	requestRepo repository.IRequestRepository,
	deliveryVerificationRepo repository.IDeliveryVerificationRepository,
	donationRepo repository.IDonationRepository,
	disbursementRepo repository.IDisbursementRepository,
	custodyLogRepo repository.ICustodyLogRepository,
) IPublicDashboardService {
	return &PublicDashboardService{
		db:                       mariadb.Connection,
		postRepo:                 postRepo,
		disasterReportRepo:       disasterReportRepo,
		disasterEventRepo:        disasterEventRepo,
		requestRepo:              requestRepo,
		deliveryVerificationRepo: deliveryVerificationRepo,
		donationRepo:             donationRepo,
		disbursementRepo:         disbursementRepo,
		custodyLogRepo:           custodyLogRepo,
	}
}

func (s *PublicDashboardService) GetPublicMap(param model.PublicDashboardParam) (*model.PublicDashboardMapResponse, error) {
	items, err := s.buildDashboardItems(param)
	if err != nil {
		return nil, err
	}

	return &model.PublicDashboardMapResponse{
		Items: items,
	}, nil
}

func (s *PublicDashboardService) GetSummary(param model.PublicDashboardParam) (*model.PublicDashboardSummaryResponse, error) {
	items, err := s.buildDashboardItems(param)
	if err != nil {
		return nil, err
	}

	sort.SliceStable(items, func(i, j int) bool {
		if items[i].FundingPercentage == items[j].FundingPercentage {
			return items[i].LatestReportedAt.After(items[j].LatestReportedAt)
		}

		return items[i].FundingPercentage < items[j].FundingPercentage
	})

	var totalTarget float64
	var totalFunded float64

	for _, item := range items {
		totalTarget += item.FundingTarget
		totalFunded += item.FundedAmount
	}

	summaryItems := make([]model.PublicDashboardSummaryItem, 0, len(items))
	for _, item := range items {
		summaryItems = append(summaryItems, model.PublicDashboardSummaryItem{
			PostID:            item.PostID,
			Name:              item.Name,
			Address:           item.Address,
			DisasterEvent:     item.DisasterEvent,
			LatestReportTitle: item.LatestReportTitle,
			FundingTarget:     item.FundingTarget,
			FundedAmount:      item.FundedAmount,
			FundingPercentage: item.FundingPercentage,
			UrgencyLevel:      item.UrgencyLevel,
			RequestCount:      item.RequestCount,
		})
	}

	return &model.PublicDashboardSummaryResponse{
		ActivePoskoCount:  int64(len(items)),
		TotalTarget:       totalTarget,
		TotalFunded:       totalFunded,
		FundingPercentage: calculateFundingPercentage(totalFunded, totalTarget),
		Items:             summaryItems,
	}, nil
}

func (s *PublicDashboardService) GetTransparency(param model.PublicTransparencyParam) (*model.PublicTransparencyResponse, error) {
	year := param.Year

	donationSummary, err := s.donationRepo.GetTransparencySummary(s.db, year)
	if err != nil {
		return nil, apperrors.InternalServer("failed to get donation transparency summary")
	}

	totalDisbursed, err := s.disbursementRepo.GetVerifiedDisbursedTotal(s.db, year)
	if err != nil {
		return nil, apperrors.InternalServer("failed to get verified disbursed total")
	}

	monthlyRows, err := s.disbursementRepo.GetMonthlyDisbursements(s.db, year)
	if err != nil {
		return nil, apperrors.InternalServer("failed to get monthly disbursements")
	}

	allocationRows, err := s.requestRepo.GetAllocationByDisaster(s.db, year)
	if err != nil {
		return nil, apperrors.InternalServer("failed to get allocation by disaster")
	}

	fulfillmentRate, err := s.deliveryVerificationRepo.GetVerifiedFulfillmentRate(s.db, year)
	if err != nil {
		return nil, apperrors.InternalServer("failed to get verified fulfillment rate")
	}

	ledgerRows, err := s.custodyLogRepo.GetLatestPublicLedger(s.db, year, param.Limit)
	if err != nil {
		return nil, apperrors.InternalServer("failed to get latest public ledger")
	}

	return &model.PublicTransparencyResponse{
		Summary: model.PublicTransparencySummary{
			TotalDonationCollected:  donationSummary.TotalCollected,
			TotalDisbursedVerified:  totalDisbursed,
			RefundAutomatic:         donationSummary.RefundAutomatic,
			VerifiedFulfillmentRate: calculateRatePercentage(fulfillmentRate.VerifiedOrders, fulfillmentRate.TotalOrders),
		},
		MonthlyDisbursements: buildMonthlyDisbursementItems(monthlyRows),
		AllocationByDisaster: buildDisasterAllocationItems(allocationRows),
		LatestLedger:         buildPublicLedgerItems(ledgerRows),
	}, nil
}

func calculateRatePercentage(part int64, total int64) int {
	if total <= 0 {
		return 0
	}

	return calculateFundingPercentage(float64(part), float64(total))
}

func buildMonthlyDisbursementItems(rows []model.MonthlyDisbursementRow) []model.MonthlyDisbursementItem {
	monthNames := map[int]string{
		1:  "Jan",
		2:  "Feb",
		3:  "Mar",
		4:  "Apr",
		5:  "Mei",
		6:  "Jun",
		7:  "Jul",
		8:  "Agu",
		9:  "Sep",
		10: "Okt",
		11: "Nov",
		12: "Des",
	}

	items := make([]model.MonthlyDisbursementItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, model.MonthlyDisbursementItem{
			Month: monthNames[row.Month],
			Total: row.Total,
		})
	}

	return items
}

func buildDisasterAllocationItems(rows []model.DisasterAllocationRow) []model.DisasterAllocationItem {
	var total float64
	for _, row := range rows {
		total += row.TotalAmount
	}

	items := make([]model.DisasterAllocationItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, model.DisasterAllocationItem{
			DisasterEvent: row.DisasterEvent,
			TotalAmount:   row.TotalAmount,
			Percentage:    calculateFundingPercentage(row.TotalAmount, total),
		})
	}

	return items
}

func buildPublicLedgerItems(rows []model.PublicLedgerRow) []model.PublicLedgerItem {
	items := make([]model.PublicLedgerItem, 0, len(rows))

	for _, row := range rows {
		items = append(items, model.PublicLedgerItem{
			OccurredAt: row.OccurredAt,
			Event:      row.Event,
			PostName:   row.PostName,
			ValueLabel: row.ValueLabel,
			Hash:       shortenAuditHash(row.Hash),
		})
	}

	return items
}

func (s *PublicDashboardService) buildDashboardItems(param model.PublicDashboardParam) ([]model.PublicDashboardMapItem, error) {
	posts, err := s.postRepo.GetPublicMapPosts(s.db, model.PublicMapPostParam{
		Query:        param.Query,
		MinLatitude:  param.MinLatitude,
		MinLongitude: param.MinLongitude,
		MaxLatitude:  param.MaxLatitude,
		MaxLongitude: param.MaxLongitude,
		Limit:        param.Limit,
	})
	if err != nil {
		return nil, apperrors.InternalServer("failed to get public map post")
	}

	if len(posts) == 0 {
		return []model.PublicDashboardMapItem{}, nil
	}

	postIDs := collectPostIDs(posts)

	reports, err := s.disasterReportRepo.GetLatestByPostIDs(s.db, model.LatestDisasterReportParam{
		PostIDs:      postIDs,
		DisasterType: param.DisasterType,
	})
	if err != nil {
		return nil, apperrors.InternalServer("failed to get latest disaster report")
	}

	if len(reports) == 0 {
		return []model.PublicDashboardMapItem{}, nil
	}

	reportByPostID := make(map[uuid.UUID]model.LatestDisasterReportRow, len(reports))
	eventIDs := make([]uuid.UUID, 0, len(reports))
	reportIDs := make([]uuid.UUID, 0, len(reports))

	for _, report := range reports {
		reportByPostID[report.PostID] = report
		eventIDs = append(eventIDs, report.EventID)
		reportIDs = append(reportIDs, report.ReportID)
	}

	events, err := s.disasterEventRepo.GetEventByIDs(s.db, uniqueUUIDs(eventIDs))
	if err != nil {
		return nil, apperrors.InternalServer("failed to get disaster event")
	}

	eventNameByID := make(map[uuid.UUID]string, len(events))
	for _, event := range events {
		eventNameByID[event.EventID] = event.Name
	}

	fundingRows, err := s.requestRepo.GetFundingSummaryByReportIDs(s.db, model.RequestFundingSummaryParam{
		ReportIDs: uniqueUUIDs(reportIDs),
	})
	if err != nil {
		return nil, apperrors.InternalServer("failed to get funding summary")
	}

	fundingByReportID := make(map[uuid.UUID]model.RequestFundingSummaryRow, len(fundingRows))
	for _, funding := range fundingRows {
		fundingByReportID[funding.ReportID] = funding
	}

	items := make([]model.PublicDashboardMapItem, 0, len(posts))
	for _, post := range posts {
		report, ok := reportByPostID[post.PostID]
		if !ok {
			continue
		}

		funding := fundingByReportID[report.ReportID]
		if funding.RequestCount == 0 {
			continue
		}

		percentage := calculateFundingPercentage(funding.FundedAmount, funding.FundingTarget)

		items = append(items, model.PublicDashboardMapItem{
			PostID:            post.PostID,
			Name:              post.Name,
			Address:           post.Address,
			Latitude:          post.Latitude,
			Longitude:         post.Longitude,
			DisasterEvent:     eventNameByID[report.EventID],
			LatestReportTitle: report.ReportTitle,
			FundingTarget:     funding.FundingTarget,
			FundedAmount:      funding.FundedAmount,
			FundingPercentage: percentage,
			UrgencyLevel:      resolveUrgencyLevel(percentage),
			RequestCount:      funding.RequestCount,
			LatestReportedAt:  resolveReportTime(report),
		})
	}

	return items, nil
}

func (s *PublicDashboardService) GetDistributions(param model.PublicDistributionParam) (*model.PublicDistributionResponse, error) {
	rows, err := s.deliveryVerificationRepo.GetPublicDistributionProofs(s.db, param)
	if err != nil {
		return nil, apperrors.InternalServer("failed to get public distribution proofs")
	}

	items := make([]model.PublicDistributionItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, model.PublicDistributionItem{
			VerificationID: row.VerificationID,
			OrderID:        row.OrderID,
			PostID:         row.PostID,
			Title:          row.RequestTitle + " — " + row.PostName,
			PostName:       row.PostName,
			DisasterEvent:  row.DisasterEvent,
			ImageURL:       row.ImageURL,
			GPSValid:       row.VerificationStatus == "approved",
			Latitude:       row.Latitude,
			Longitude:      row.Longitude,
			CapturedAt:     row.CapturedAt,
			TotalAmount:    row.TotalAmount,
			DonorCount:     row.DonorCount,
			AuditHash:      shortenAuditHash(row.CurrentHash),
		})
	}

	return &model.PublicDistributionResponse{
		Items: items,
	}, nil
}

func shortenAuditHash(hash string) string {
	if len(hash) <= 10 {
		return hash
	}

	return hash[:6] + "..." + hash[len(hash)-4:]
}

func collectPostIDs(posts []model.PublicMapPostRow) []uuid.UUID {
	ids := make([]uuid.UUID, 0, len(posts))

	for _, post := range posts {
		ids = append(ids, post.PostID)
	}

	return uniqueUUIDs(ids)
}

func uniqueUUIDs(ids []uuid.UUID) []uuid.UUID {
	seen := make(map[uuid.UUID]struct{}, len(ids))
	unique := make([]uuid.UUID, 0, len(ids))

	for _, id := range ids {
		if id == uuid.Nil {
			continue
		}

		if _, exists := seen[id]; exists {
			continue
		}

		seen[id] = struct{}{}
		unique = append(unique, id)
	}

	return unique
}

func calculateFundingPercentage(fundedAmount float64, fundingTarget float64) int {
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

func resolveUrgencyLevel(percentage int) string {
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

func resolveReportTime(report model.LatestDisasterReportRow) time.Time {
	if !report.ReportedAt.IsZero() {
		return report.ReportedAt
	}

	return report.CreatedAt
}
