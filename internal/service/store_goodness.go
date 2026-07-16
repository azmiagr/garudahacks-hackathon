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

type IStoreGoodnessService interface {
	GetStoreGoodness(user *entity.User, req model.StoreGoodnessRequest) (*model.StoreGoodnessResponse, error)
}

type StoreGoodnessService struct {
	db              *gorm.DB
	storeRepository repository.IStoreRepository
}

func NewStoreGoodnessService(storeRepository repository.IStoreRepository) IStoreGoodnessService {
	return &StoreGoodnessService{
		db:              mariadb.Connection,
		storeRepository: storeRepository,
	}
}

func (s *StoreGoodnessService) GetStoreGoodness(user *entity.User, req model.StoreGoodnessRequest) (*model.StoreGoodnessResponse, error) {
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

	year := normalizeStoreGoodnessYear(req.Year)
	limit := normalizeStoreGoodnessServiceLimit(req.Limit)
	offset := normalizeStoreGoodnessServiceOffset(req.Offset)
	param := model.StoreGoodnessParam{
		StoreID: store.StoreID,
		Year:    year,
		Limit:   limit,
		Offset:  offset,
	}

	certificateRow, err := s.storeRepository.GetStoreGoodnessCertificate(s.db, param)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			certificateRow = buildEmptyStoreGoodnessCertificateRow(store)
		} else {
			return nil, err
		}
	}

	historyRows, err := s.storeRepository.GetStoreContributionHistory(s.db, param)
	if err != nil {
		return nil, err
	}

	totalHistory, err := s.storeRepository.CountStoreContributionHistory(s.db, param)
	if err != nil {
		return nil, err
	}

	return &model.StoreGoodnessResponse{
		Certificate:  buildStoreGoodnessCertificate(*certificateRow, year),
		History:      buildStoreContributionHistoryItems(historyRows),
		TotalHistory: totalHistory,
		Limit:        limit,
		Offset:       offset,
	}, nil
}

func buildEmptyStoreGoodnessCertificateRow(store *entity.Stores) *model.StoreGoodnessCertificateRow {
	return &model.StoreGoodnessCertificateRow{
		StoreID:   store.StoreID,
		StoreName: store.Name,
		BankName:  store.BankName,
	}
}

func buildStoreGoodnessCertificate(row model.StoreGoodnessCertificateRow, year int) model.StoreGoodnessCertificate {
	reputation := math.Round(row.ReputationScore*10) / 10

	return model.StoreGoodnessCertificate{
		StoreID:             row.StoreID,
		StoreName:           row.StoreName,
		Title:               "Sertifikat Digital",
		PartnerLabel:        "Mitra Tanggap Bencana",
		SinceText:           buildStoreGoodnessSinceText(row.FirstContributionAt, year),
		VerifiedOrderCount:  row.VerifiedOrderCount,
		VerifiedOrderText:   fmt.Sprintf("%d order selesai", row.VerifiedOrderCount),
		VerifiedAmountTotal: row.VerifiedAmountTotal,
		VerifiedAmountText:  formatGoodnessRupiahShort(row.VerifiedAmountTotal),
		ReputationScore:     reputation,
		ReputationText:      formatReputationText(reputation),
		DisputeCount:        row.DisputeCount,
		DisputeText:         fmt.Sprintf("%d sengketa", row.DisputeCount),
		FirstContributionAt: row.FirstContributionAt,
		ShareURL:            fmt.Sprintf("/store/goodness/%s", row.StoreID.String()),
	}
}

func buildStoreContributionHistoryItems(rows []model.StoreContributionHistoryRow) []model.StoreContributionHistoryItem {
	items := make([]model.StoreContributionHistoryItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, model.StoreContributionHistoryItem{
			OrderID:         row.OrderID,
			OrderCode:       row.OrderCode,
			PostName:        row.PostName,
			DisasterName:    row.DisasterName,
			Title:           fmt.Sprintf("%s %s · #%s", titleCase(row.DisasterName), row.PostName, row.OrderCode),
			ItemCount:       row.ItemCount,
			TotalAmount:     row.TotalAmount,
			TotalAmountText: formatGoodnessRupiah(row.TotalAmount),
			VerifiedAt:      row.VerifiedAt,
			VerifiedAtText:  formatGoodnessDate(row.VerifiedAt),
			LatestHash:      row.LatestHash,
			ShortLatestHash: shortenLedgerHash(row.LatestHash),
		})
	}
	return items
}

func buildStoreGoodnessSinceText(firstContributionAt *time.Time, fallbackYear int) string {
	if firstContributionAt == nil {
		return fmt.Sprintf("sejak %d", fallbackYear)
	}
	return "sejak " + firstContributionAt.Format("Jan 2006")
}

func formatReputationText(score float64) string {
	if score <= 0 {
		return "belum ada reputasi"
	}
	return fmt.Sprintf("%.1f reputasi", score)
}

func formatGoodnessRupiah(amount float64) string {
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

func formatGoodnessRupiahShort(amount float64) string {
	if amount >= 1000000000 {
		return fmt.Sprintf("Rp%.1f M", amount/1000000000)
	}
	if amount >= 1000000 {
		return fmt.Sprintf("Rp%.0fjt", amount/1000000)
	}
	if amount >= 1000 {
		return fmt.Sprintf("Rp%.0frb", amount/1000)
	}
	return formatGoodnessRupiah(amount)
}

func formatGoodnessDate(value time.Time) string {
	now := time.Now()
	yearNow, monthNow, dayNow := now.Date()
	yearValue, monthValue, dayValue := value.Date()
	if yearNow == yearValue && monthNow == monthValue && dayNow == dayValue {
		return "Hari ini"
	}
	return value.Format("2 Jan")
}

func titleCase(value string) string {
	words := strings.Fields(strings.ToLower(strings.TrimSpace(value)))
	for i, word := range words {
		if word == "" {
			continue
		}
		words[i] = strings.ToUpper(word[:1]) + word[1:]
	}
	return strings.Join(words, " ")
}

func normalizeStoreGoodnessYear(year int) int {
	if year <= 0 {
		return time.Now().Year()
	}
	return year
}

func normalizeStoreGoodnessServiceLimit(limit int) int {
	if limit <= 0 {
		return 20
	}
	if limit > 100 {
		return 100
	}
	return limit
}

func normalizeStoreGoodnessServiceOffset(offset int) int {
	if offset < 0 {
		return 0
	}
	return offset
}
