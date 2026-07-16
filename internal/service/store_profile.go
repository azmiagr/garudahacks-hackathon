package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"strings"

	"github.com/azmiagr/garudahacks-hackathon/entity"
	"github.com/azmiagr/garudahacks-hackathon/internal/repository"
	"github.com/azmiagr/garudahacks-hackathon/model"
	"github.com/azmiagr/garudahacks-hackathon/pkg/database/mariadb"
	apperrors "github.com/azmiagr/garudahacks-hackathon/pkg/errors"
	"gorm.io/gorm"
)

type IStoreProfileService interface {
	GetStoreProfile(user *entity.User) (*model.StoreProfileResponse, error)
}

type StoreProfileService struct {
	db              *gorm.DB
	storeRepository repository.IStoreRepository
}

func NewStoreProfileService(storeRepository repository.IStoreRepository) IStoreProfileService {
	return &StoreProfileService{
		db:              mariadb.Connection,
		storeRepository: storeRepository,
	}
}

func (s *StoreProfileService) GetStoreProfile(user *entity.User) (*model.StoreProfileResponse, error) {
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

	stats, err := s.storeRepository.GetStoreProfileStats(s.db, store.StoreID)
	if err != nil {
		return nil, err
	}

	categories := parseStoreCategories(store.CategoriesJSON)
	acceptanceRate := calculateStoreAcceptanceRate(stats.AcceptedOrder30Days, stats.TotalOrder30Days)

	return &model.StoreProfileResponse{
		StoreID:              store.StoreID,
		OwnerID:              store.OwnerID,
		Name:                 store.Name,
		Address:              store.Address,
		Latitude:             store.Latitude,
		Longitude:            store.Longitude,
		IsOnline:             user.Status == "active",
		StoreStatus:          mapStoreOnlineStatus(user.Status),
		KYCStatus:            user.KYCStatus,
		KYCLabel:             mapStoreKYCLabel(user.KYCStatus),
		ReputationScore:      math.Round(stats.ReputationScore*10) / 10,
		BusinessNumber:       store.BusinessNumber,
		NPWP:                 store.NPWP,
		KTPImageURL:          store.KTPImageURL,
		BankName:             store.BankName,
		BankAccountNo:        store.BankAccountNo,
		MaskedBankAccount:    maskBankAccount(store.BankAccountNo),
		BankAccountName:      store.BankAccountName,
		Categories:           categories,
		CategoriesText:       buildStoreCategoriesText(categories),
		TotalOrder30Days:     stats.TotalOrder30Days,
		AcceptedOrder30Days:  stats.AcceptedOrder30Days,
		CancelledOrder30Days: stats.CancelledOrder30Days,
		AcceptanceRate30Days: acceptanceRate,
		AcceptanceSummary:    buildStoreAcceptanceSummary(acceptanceRate, stats.CancelledOrder30Days),
		CreatedAt:            store.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:            store.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

func parseStoreCategories(raw string) []string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return []string{}
	}

	var categories []string
	if err := json.Unmarshal([]byte(raw), &categories); err != nil {
		return []string{}
	}

	return categories
}

func buildStoreCategoriesText(categories []string) string {
	if len(categories) == 0 {
		return ""
	}
	if len(categories) <= 2 {
		return strings.Join(categories, " · ")
	}
	return strings.Join(categories[:2], " · ") + " · +tambah"
}

func calculateStoreAcceptanceRate(accepted int64, total int64) float64 {
	if total <= 0 {
		return 0
	}
	return math.Round((float64(accepted)/float64(total))*10000) / 100
}

func buildStoreAcceptanceSummary(acceptanceRate float64, cancelledCount int64) string {
	return fmt.Sprintf("Tingkat penerimaan order 30 hari: %.0f%% · pembatalan: %d", acceptanceRate, cancelledCount)
}

func mapStoreOnlineStatus(status string) string {
	if status == "active" {
		return "Online"
	}
	return "Offline"
}

func mapStoreKYCLabel(status string) string {
	switch status {
	case "approved":
		return "KYC TERVERIFIKASI"
	case "rejected":
		return "KYC DITOLAK"
	default:
		return "KYC MENUNGGU"
	}
}
