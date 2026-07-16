package service

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/azmiagr/garudahacks-hackathon/entity"
	"github.com/azmiagr/garudahacks-hackathon/internal/repository"
	"github.com/azmiagr/garudahacks-hackathon/model"
	"github.com/azmiagr/garudahacks-hackathon/pkg/database/mariadb"
	apperrors "github.com/azmiagr/garudahacks-hackathon/pkg/errors"
	"gorm.io/gorm"
)

type IStoreDisbursementService interface {
	GetDashboardDisbursement(user *entity.User, req model.StoreDisbursementDashboardRequest) (*model.StoreDisbursementDashboardResponse, error)
}

type StoreDisbursementService struct {
	db                     *gorm.DB
	storeRepository        repository.IStoreRepository
	disbursementRepository repository.IDisbursementRepository
}

func NewStoreDisbursementService(
	storeRepository repository.IStoreRepository,
	disbursementRepository repository.IDisbursementRepository,
) IStoreDisbursementService {
	return &StoreDisbursementService{
		db:                     mariadb.Connection,
		storeRepository:        storeRepository,
		disbursementRepository: disbursementRepository,
	}
}

func (s *StoreDisbursementService) GetDashboardDisbursement(user *entity.User, req model.StoreDisbursementDashboardRequest) (*model.StoreDisbursementDashboardResponse, error) {
	if user == nil {
		return nil, apperrors.Unauthorized("user is required")
	}

	store, err := s.storeRepository.GetStore(s.db, model.GetStoreParam{OwnerID: user.UserID})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.Forbidden("store profile is required")
		}
		return nil, err
	}

	year := normalizeStoreDisbursementRequestYear(req.Year)
	month := normalizeStoreDisbursementRequestMonth(req.Month)
	limit := normalizeStoreDisbursementRequestLimit(req.Limit)
	offset := normalizeStoreDisbursementRequestOffset(req.Offset)

	param := model.StoreDisbursementDashboardParam{
		StoreID: store.StoreID,
		Year:    year,
		Month:   month,
		Limit:   limit,
		Offset:  offset,
	}

	summaryRow, err := s.disbursementRepository.GetStoreDisbursementSummary(s.db, param)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			summaryRow = buildEmptyStoreDisbursementSummaryRow(store)
		} else {
			return nil, err
		}
	}

	historyRows, err := s.disbursementRepository.GetStoreDisbursementHistory(s.db, param)
	if err != nil {
		return nil, err
	}

	totalHistory, err := s.disbursementRepository.CountStoreDisbursementHistory(s.db, param)
	if err != nil {
		return nil, err
	}

	goodnessRow, err := s.disbursementRepository.GetStoreGoodnessTrail(s.db, store.StoreID, year)
	if err != nil {
		return nil, err
	}

	return &model.StoreDisbursementDashboardResponse{
		Summary:       buildStoreDisbursementSummary(*summaryRow),
		History:       buildStoreDisbursementHistoryItems(historyRows),
		GoodnessTrail: buildStoreGoodnessTrail(*goodnessRow, year),
		TotalHistory:  totalHistory,
		Limit:         limit,
		Offset:        offset,
	}, nil
}

func buildEmptyStoreDisbursementSummaryRow(store *entity.Stores) *model.StoreDisbursementSummaryRow {
	return &model.StoreDisbursementSummaryRow{
		StoreID:       store.StoreID,
		StoreName:     store.Name,
		BankName:      store.BankName,
		BankAccountNo: store.BankAccountNo,
	}
}

func buildStoreDisbursementSummary(row model.StoreDisbursementSummaryRow) model.StoreDisbursementSummary {
	medianMinutes := int(math.Round(row.MedianDisbursementMin))
	maskedAccount := maskBankAccount(row.BankAccountNo)

	return model.StoreDisbursementSummary{
		StoreID:                 row.StoreID,
		StoreName:               row.StoreName,
		TotalDisbursedThisMonth: row.TotalDisbursedThisMonth,
		TotalDisbursedText:      formatStoreRupiah(row.TotalDisbursedThisMonth),
		CompletedOrderCount:     row.CompletedOrderCount,
		DisputeCount:            row.DisputeCount,
		BankName:                row.BankName,
		MaskedBankAccount:       maskedAccount,
		MedianDisbursementMin:   medianMinutes,
		MedianDisbursementText:  buildMedianDisbursementText(medianMinutes),
		Subtitle: fmt.Sprintf(
			"%d order selesai · %d sengketa · rekening %s %s",
			row.CompletedOrderCount,
			row.DisputeCount,
			strings.ToUpper(strings.TrimSpace(row.BankName)),
			maskedAccount,
		),
	}
}

func buildStoreDisbursementHistoryItems(rows []model.StoreDisbursementHistoryRow) []model.StoreDisbursementHistoryItem {
	items := make([]model.StoreDisbursementHistoryItem, 0, len(rows))
	for _, row := range rows {
		minutes := int(math.Round(row.MinutesAfterVerification))
		items = append(items, model.StoreDisbursementHistoryItem{
			DisbursementID:               row.DisbursementID,
			OrderID:                      row.OrderID,
			OrderCode:                    row.OrderCode,
			PostName:                     row.PostName,
			Amount:                       row.Amount,
			AmountText:                   formatStoreDisbursementAmount(row.Amount, row.Status),
			Status:                       row.Status,
			StatusLabel:                  mapStoreDisbursementStatusLabel(row.Status),
			BadgeVariant:                 mapStoreDisbursementBadgeVariant(row.Status),
			IdempotencyKey:               row.IdempotencyKey,
			GatewayReference:             row.GatewayReference,
			GatewayAttempt:               row.GatewayAttempt,
			VerificationApprovedAt:       row.VerificationApprovedAt,
			DisbursedAt:                  row.DisbursedAt,
			CreatedAt:                    row.CreatedAt,
			TimelineText:                 buildStoreDisbursementTimelineText(row, minutes),
			MinutesAfterVerification:     minutes,
			MinutesAfterVerificationText: buildMinutesAfterVerificationText(minutes),
		})
	}
	return items
}

func buildStoreGoodnessTrail(row model.StoreGoodnessTrailRow, year int) model.StoreGoodnessTrail {
	return model.StoreGoodnessTrail{
		StoreID:             row.StoreID,
		VerifiedOrderCount:  row.VerifiedOrderCount,
		VerifiedAmountTotal: row.VerifiedAmountTotal,
		VerifiedAmountText:  formatStoreRupiahShort(row.VerifiedAmountTotal),
		Year:                year,
		SummaryText: fmt.Sprintf(
			"%d order kebencanaan terverifikasi · %s sejak %d",
			row.VerifiedOrderCount,
			formatStoreRupiahShort(row.VerifiedAmountTotal),
			year,
		),
		FirstContributionAt: row.FirstContributionAt,
		LastContributionAt:  row.LastContributionAt,
	}
}

func maskBankAccount(account string) string {
	account = strings.TrimSpace(account)
	if account == "" {
		return "••••"
	}
	if len(account) <= 4 {
		return "••••" + account
	}
	return "••••" + account[len(account)-4:]
}

func buildMedianDisbursementText(minutes int) string {
	if minutes <= 0 {
		return "Median cair belum tersedia"
	}
	return fmt.Sprintf("Median cair %d mnt", minutes)
}

func formatStoreDisbursementAmount(amount float64, status string) string {
	prefix := "+"
	if strings.ToLower(strings.TrimSpace(status)) == "failed" {
		prefix = ""
	}
	return prefix + formatStoreRupiah(amount)
}

func formatStoreRupiah(amount float64) string {
	value := int64(math.Round(amount))
	raw := fmt.Sprintf("%d", value)
	parts := []string{}
	for len(raw) > 3 {
		parts = append([]string{raw[len(raw)-3:]}, parts...)
		raw = raw[:len(raw)-3]
	}
	parts = append([]string{raw}, parts...)
	return "Rp" + strings.Join(parts, ".")
}

func formatStoreRupiahShort(amount float64) string {
	if amount >= 1000000000 {
		return fmt.Sprintf("Rp%.1f M", amount/1000000000)
	}
	if amount >= 1000000 {
		return fmt.Sprintf("Rp%.0f jt", amount/1000000)
	}
	if amount >= 1000 {
		return fmt.Sprintf("Rp%.0f rb", amount/1000)
	}
	return formatStoreRupiah(amount)
}

func mapStoreDisbursementStatusLabel(status string) string {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "success":
		return "CAIR"
	case "failed":
		return "RETRY"
	default:
		return "PENDING"
	}
}

func mapStoreDisbursementBadgeVariant(status string) string {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "success":
		return "success"
	case "failed":
		return "warning"
	default:
		return "neutral"
	}
}

func buildStoreDisbursementTimelineText(row model.StoreDisbursementHistoryRow, minutes int) string {
	switch strings.ToLower(strings.TrimSpace(row.Status)) {
	case "success":
		if row.DisbursedAt != nil {
			return fmt.Sprintf("%s · %s setelah verifikasi tiba", formatStoreDateTime(*row.DisbursedAt), buildMinutesAfterVerificationText(minutes))
		}
		return fmt.Sprintf("%s · cair", formatStoreDateTime(row.UpdatedAt))
	case "failed":
		if row.GatewayAttempt > 0 {
			return fmt.Sprintf("%s · antre retry gateway (percobaan %d/5)", formatStoreDate(row.UpdatedAt), row.GatewayAttempt)
		}
		return fmt.Sprintf("%s · antre retry gateway", formatStoreDate(row.UpdatedAt))
	default:
		return fmt.Sprintf("%s · menunggu gateway", formatStoreDateTime(row.CreatedAt))
	}
}

func buildMinutesAfterVerificationText(minutes int) string {
	if minutes <= 0 {
		return "belum tersedia"
	}
	if minutes < 60 {
		return fmt.Sprintf("%d mnt", minutes)
	}
	hours := minutes / 60
	remaining := minutes % 60
	if remaining == 0 {
		return fmt.Sprintf("%d jam", hours)
	}
	return fmt.Sprintf("%d jam %d mnt", hours, remaining)
}

func formatStoreDateTime(value time.Time) string {
	now := time.Now()
	if sameDate(now, value) {
		return "Hari ini " + value.Format("15:04")
	}
	return value.Format("02 Jan 15:04")
}

func formatStoreDate(value time.Time) string {
	now := time.Now()
	if sameDate(now, value) {
		return "Hari ini"
	}
	return value.Format("02 Jan")
}

func sameDate(a time.Time, b time.Time) bool {
	ay, am, ad := a.Date()
	by, bm, bd := b.Date()
	return ay == by && am == bm && ad == bd
}

func normalizeStoreDisbursementRequestYear(year int) int {
	if year <= 0 {
		return time.Now().Year()
	}
	return year
}

func normalizeStoreDisbursementRequestMonth(month int) int {
	if month < 1 || month > 12 {
		return int(time.Now().Month())
	}
	return month
}

func normalizeStoreDisbursementRequestLimit(limit int) int {
	if limit <= 0 {
		return 20
	}
	if limit > 100 {
		return 100
	}
	return limit
}

func normalizeStoreDisbursementRequestOffset(offset int) int {
	if offset < 0 {
		return 0
	}
	return offset
}
