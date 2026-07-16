package repository

import (
	"github.com/azmiagr/garudahacks-hackathon/model"
	"gorm.io/gorm"
)

type IDisbursementRepository interface {
	GetVerifiedDisbursedTotal(tx *gorm.DB, year int) (float64, error)
	GetMonthlyDisbursements(tx *gorm.DB, year int) ([]model.MonthlyDisbursementRow, error)
}

type DisbursementRepository struct {
	db *gorm.DB
}

func NewDisbursementRepository(db *gorm.DB) IDisbursementRepository {
	return &DisbursementRepository{db: db}
}

func (r *DisbursementRepository) GetVerifiedDisbursedTotal(tx *gorm.DB, year int) (float64, error) {
	var total float64

	query := tx.Table("disbursements").
		Select("COALESCE(SUM(amount), 0)").
		Where("status = ?", "success")

	if year > 0 {
		query = query.Where("YEAR(created_at) = ?", year)
	}

	err := query.Scan(&total).Error
	if err != nil {
		return 0, err
	}

	return total, nil
}

func (r *DisbursementRepository) GetMonthlyDisbursements(tx *gorm.DB, year int) ([]model.MonthlyDisbursementRow, error) {
	var rows []model.MonthlyDisbursementRow

	query := tx.Table("disbursements").
		Select(`
			MONTH(created_at) AS month,
			COALESCE(SUM(amount), 0) AS total
		`).
		Where("status = ?", "success").
		Group("MONTH(created_at)").
		Order("MONTH(created_at) ASC")

	if year > 0 {
		query = query.Where("YEAR(created_at) = ?", year)
	}

	err := query.Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	return rows, nil
}
