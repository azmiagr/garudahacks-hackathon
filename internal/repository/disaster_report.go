package repository

import (
	"strings"

	"github.com/azmiagr/garudahacks-hackathon/model"
	"gorm.io/gorm"
)

type IDisasterReportRepository interface {
	GetLatestByPostIDs(tx *gorm.DB, param model.LatestDisasterReportParam) ([]model.LatestDisasterReportRow, error)
}

type DisasterReportRepository struct {
	db *gorm.DB
}

func NewDisasterReportRepository(db *gorm.DB) IDisasterReportRepository {
	return &DisasterReportRepository{db: db}
}

func (r *DisasterReportRepository) GetLatestByPostIDs(tx *gorm.DB, param model.LatestDisasterReportParam) ([]model.LatestDisasterReportRow, error) {
	var rows []model.LatestDisasterReportRow

	if len(param.PostIDs) == 0 {
		return rows, nil
	}

	latestReportSubquery := tx.Table("disaster_reports").
		Select("post_id, MAX(COALESCE(reported_at, created_at)) AS latest_reported_at").
		Where("post_id IN ?", param.PostIDs).
		Where("report_status IN ?", []string{"pending", "approved"}).
		Group("post_id")

	query := tx.Table("disaster_reports AS dr").
		Select(`
			dr.report_id,
			dr.post_id,
			dr.event_id,
			dr.report_title,
			dr.reported_at,
			dr.created_at
		`).
		Joins(`
			JOIN (?) AS latest_reports
			ON latest_reports.post_id = dr.post_id
			AND latest_reports.latest_reported_at = COALESCE(dr.reported_at, dr.created_at)
		`, latestReportSubquery).
		Where("dr.report_status IN ?", []string{"pending", "approved"})

	if strings.TrimSpace(param.DisasterType) != "" {
		query = query.
			Joins("JOIN disaster_events AS de ON de.event_id = dr.event_id").
			Where("LOWER(de.name) = LOWER(?)", strings.TrimSpace(param.DisasterType))
	}

	err := query.Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	return rows, nil
}
