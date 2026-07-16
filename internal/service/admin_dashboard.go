package service

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/azmiagr/garudahacks-hackathon/entity"
	"github.com/azmiagr/garudahacks-hackathon/internal/repository"
	"github.com/azmiagr/garudahacks-hackathon/model"
	"github.com/azmiagr/garudahacks-hackathon/pkg/database/mariadb"
	apperrors "github.com/azmiagr/garudahacks-hackathon/pkg/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IAdminDashboardService interface {
	GetAdminDashboard(user *entity.User) (*model.AdminDashboardHomeResponse, error)
}

type AdminDashboardService struct {
	db                       *gorm.DB
	adminDashboardRepository repository.IAdminDashboardRepository
}

func NewAdminDashboardService(adminDashboardRepository repository.IAdminDashboardRepository) IAdminDashboardService {
	return &AdminDashboardService{
		db:                       mariadb.Connection,
		adminDashboardRepository: adminDashboardRepository,
	}
}

func (s *AdminDashboardService) GetAdminDashboard(user *entity.User) (*model.AdminDashboardHomeResponse, error) {
	if user == nil {
		return nil, apperrors.Unauthorized("user is required")
	}

	if strings.TrimSpace(user.Name) == "" {
		user.Name = "Admin"
	}

	eventRows, err := s.adminDashboardRepository.GetAdminDashboardEvents(s.db, model.AdminDashboardHomeParam{
		UserID: user.UserID,
	})
	if err != nil {
		return nil, err
	}

	orderRows, err := s.adminDashboardRepository.GetAdminDashboardOrders(s.db, model.AdminDashboardHomeParam{
		UserID: user.UserID,
	})
	if err != nil {
		return nil, err
	}

	ordersByPostID := map[string][]model.AdminDashboardOrderPreview{}
	for _, row := range orderRows {
		key := row.PostID.String()
		if len(ordersByPostID[key]) >= 2 {
			continue
		}

		ordersByPostID[key] = append(ordersByPostID[key], buildOrderPreview(row))
	}

	activeEvents := []model.AdminDashboardEventCard{}
	closedEvents := []model.AdminDashboardEventCard{}

	for _, row := range eventRows {
		card := buildEventCard(row, ordersByPostID[row.PostID.String()])

		if card.Status == "closed" {
			closedEvents = append(closedEvents, card)
			continue
		}

		activeEvents = append(activeEvents, card)
	}

	isVerified := user.KYCStatus == "approved"
	verificationText := "ADMIN MENUNGGU VERIFIKASI"
	if isVerified {
		verificationText = "ADMIN TERVERIFIKASI"
	}

	return &model.AdminDashboardHomeResponse{
		GreetingName:     user.Name,
		IsAdminVerified:  isVerified,
		VerificationText: verificationText,
		ActiveEvents:     activeEvents,
		ClosedEvents:     closedEvents,
	}, nil
}

func buildEventCard(row model.AdminDashboardEventRow, orders []model.AdminDashboardOrderPreview) model.AdminDashboardEventCard {
	fundingPercentage := calculateFundingPercentageAdmin(row.FundedAmount, row.FundingTarget)
	status := resolveEventStatus(row)
	statusLabel := strings.ToUpper(status)

	canScanCourierQR := false
	for _, order := range orders {
		if order.Status == "sent" {
			canScanCourierQR = true
			break
		}
	}

	return model.AdminDashboardEventCard{
		PostID:                row.PostID,
		EventCode:             buildEventCode(row.PostID),
		Title:                 row.Title,
		DisasterType:          row.DisasterType,
		Status:                status,
		StatusLabel:           statusLabel,
		ImageURL:              row.ImageURL,
		Address:               row.Address,
		GeofenceRadius:        row.GeofenceRadius,
		AffectedHouseholds:    0,
		StartedAt:             row.StartedAt,
		ElapsedText:           buildElapsedText(row.StartedAt),
		FundingTarget:         row.FundingTarget,
		FundedAmount:          row.FundedAmount,
		FundingPercentage:     fundingPercentage,
		FundingText:           buildFundingText(row.FundedAmount, row.FundingTarget, fundingPercentage),
		OrderCount:            row.OrderCount,
		SummaryText:           buildEventSummaryText(row),
		CanScanCourierQR:      canScanCourierQR,
		CanAddFollowUpRequest: status == "active",
		LatestOrders:          orders,
	}
}

func buildOrderPreview(row model.AdminDashboardOrderRow) model.AdminDashboardOrderPreview {
	status := resolveOrderStatus(row)
	statusLabel := mapOrderStatusLabel(status)

	return model.AdminDashboardOrderPreview{
		OrderID:      row.OrderID,
		OrderCode:    row.OrderCode,
		StoreName:    row.StoreName,
		CourierName:  row.CourierName,
		Status:       status,
		StatusLabel:  statusLabel,
		Description:  buildOrderDescription(status, row),
		BadgeVariant: mapOrderBadgeVariant(status),
		UpdatedAt:    row.UpdatedAt,
	}
}

func resolveEventStatus(row model.AdminDashboardEventRow) string {
	if row.OrderCount > 0 && row.CompletedOrderCount >= row.OrderCount {
		return "closed"
	}

	return "active"
}

func resolveOrderStatus(row model.AdminDashboardOrderRow) string {
	if row.VerificationStatus == "approved" {
		return "completed"
	}

	if row.OrderStatus == "approved" {
		return "sent"
	}

	if row.OrderStatus == "rejected" {
		return "cancelled"
	}

	return "pending"
}

func mapOrderStatusLabel(status string) string {
	switch status {
	case "completed":
		return "SELESAI"
	case "sent":
		return "DIKIRIM"
	case "cancelled":
		return "BATAL"
	default:
		return "MENUNGGU"
	}
}

func mapOrderBadgeVariant(status string) string {
	switch status {
	case "completed":
		return "success"
	case "sent":
		return "warning"
	case "cancelled":
		return "danger"
	default:
		return "neutral"
	}
}

func buildOrderDescription(status string, row model.AdminDashboardOrderRow) string {
	switch status {
	case "completed":
		return "Selesai - foto bukti tayang publik"
	case "sent":
		if strings.TrimSpace(row.CourierName) == "" {
			return "Kurir dalam perjalanan"
		}
		return fmt.Sprintf("Kurir %s dalam perjalanan", row.CourierName)
	case "cancelled":
		return "Order dibatalkan"
	default:
		return "Menunggu toko atau kurir"
	}
}

func calculateFundingPercentageAdmin(fundedAmount, fundingTarget float64) float64 {
	if fundingTarget <= 0 {
		return 0
	}

	percentage := (fundedAmount / fundingTarget) * 100
	if percentage > 100 {
		return 100
	}

	return math.Round(percentage*100) / 100
}

func buildFundingText(fundedAmount, fundingTarget, percentage float64) string {
	return fmt.Sprintf("%s / %s - %.0f%%", formatRupiahShort(fundedAmount), formatRupiahShort(fundingTarget), percentage)
}

func buildEventSummaryText(row model.AdminDashboardEventRow) string {
	if resolveEventStatus(row) == "closed" {
		return fmt.Sprintf("Selesai - %d order - laporan otomatis tersedia", row.OrderCount)
	}

	return fmt.Sprintf("%s · radius %.0f m", buildEventCode(row.PostID), row.GeofenceRadius)
}

func buildEventCode(postID uuid.UUID) string {
	value := strings.ReplaceAll(postID.String(), "-", "")
	if len(value) < 4 {
		return "PSK-0000"
	}

	return "PSK-" + strings.ToUpper(value[:4])
}

func buildElapsedText(startedAt time.Time) string {
	if startedAt.IsZero() {
		return ""
	}

	duration := time.Since(startedAt)
	if duration < time.Hour {
		minutes := int(duration.Minutes())
		if minutes <= 0 {
			minutes = 1
		}
		return fmt.Sprintf("%d menit berjalan", minutes)
	}

	if duration < 24*time.Hour {
		return fmt.Sprintf("%d jam berjalan", int(duration.Hours()))
	}

	return fmt.Sprintf("%d hari berjalan", int(duration.Hours()/24))
}

func formatRupiahShort(amount float64) string {
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
