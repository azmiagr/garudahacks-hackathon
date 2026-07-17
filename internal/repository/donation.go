package repository

import (
	"strings"

	"github.com/azmiagr/garudahacks-hackathon/entity"
	"github.com/azmiagr/garudahacks-hackathon/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IDonationRepository interface {
	GetTransparencySummary(tx *gorm.DB, year int) (*model.DonationTransparencySummaryRow, error)
	CreateDonation(tx *gorm.DB, donation *entity.Donations) error
	GetDonation(tx *gorm.DB, param model.GetDonationParam) (*entity.Donations, error)
	UpdateDonation(tx *gorm.DB, donation *entity.Donations) error
	GetDonorDonationTransactions(tx *gorm.DB, param model.DonorDonationTransactionListParam) ([]model.DonorDonationTransactionListRow, error)
	CountDonorDonationTransactions(tx *gorm.DB, param model.DonorDonationTransactionListParam) (int64, error)
	GetDonorDonationTransactionDetail(tx *gorm.DB, param model.DonorDonationTransactionDetailParam) (*model.DonorDonationTransactionDetailRow, error)
	GetDonorDonationTransactionItems(tx *gorm.DB, param model.DonorDonationTransactionDetailParam) ([]model.DonorDonationTransactionItemRow, error)
	GetDonorDonationTransactionCustodyLogs(tx *gorm.DB, param model.DonorDonationTransactionDetailParam) ([]model.DonorDonationTransactionCustodyLogRow, error)
}

type DonationRepository struct {
	db *gorm.DB
}

func NewDonationRepository(db *gorm.DB) IDonationRepository {
	return &DonationRepository{db: db}
}

func (r *DonationRepository) CreateDonation(tx *gorm.DB, donation *entity.Donations) error {
	err := tx.Create(donation).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *DonationRepository) GetDonation(tx *gorm.DB, param model.GetDonationParam) (*entity.Donations, error) {
	var donation entity.Donations
	err := tx.Where(&param).First(&donation).Error
	if err != nil {
		return nil, err
	}

	return &donation, nil
}

func (r *DonationRepository) UpdateDonation(tx *gorm.DB, donation *entity.Donations) error {
	err := tx.Save(donation).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *DonationRepository) GetTransparencySummary(tx *gorm.DB, year int) (*model.DonationTransparencySummaryRow, error) {
	var row model.DonationTransparencySummaryRow

	query := tx.Table("donations").
		Select(`
			COALESCE(SUM(CASE WHEN donation_status = 'approved' THEN donation_amount ELSE 0 END), 0) AS total_collected,
			COALESCE(SUM(CASE WHEN donation_status = 'rejected' THEN donation_amount ELSE 0 END), 0) AS refund_automatic
		`)

	if year > 0 {
		query = query.Where("YEAR(donated_at) = ?", year)
	}

	if err := query.Scan(&row).Error; err != nil {
		return nil, err
	}

	return &row, nil
}

func (r *DonationRepository) GetDonorDonationTransactions(tx *gorm.DB, param model.DonorDonationTransactionListParam) ([]model.DonorDonationTransactionListRow, error) {
	var rows []model.DonorDonationTransactionListRow

	query := buildDonorDonationTransactionBaseQuery(tx).
		Select(baseDonorDonationTransactionSelect()).
		Where("d.donated_by = ?", param.UserID)

	query = applyDonorDonationTransactionStatusFilter(query, param.Status)

	err := query.
		Order("d.donated_at DESC").
		Limit(normalizeDonorDonationTransactionLimit(param.Limit)).
		Offset(normalizeDonorDonationTransactionOffset(param.Offset)).
		Scan(&rows).Error
	return rows, err
}

func (r *DonationRepository) CountDonorDonationTransactions(tx *gorm.DB, param model.DonorDonationTransactionListParam) (int64, error) {
	var count int64

	query := buildDonorDonationTransactionBaseQuery(tx).
		Where("d.donated_by = ?", param.UserID)

	query = applyDonorDonationTransactionStatusFilter(query, param.Status)
	err := query.Distinct("d.donation_id").Count(&count).Error
	return count, err
}

func (r *DonationRepository) GetDonorDonationTransactionDetail(tx *gorm.DB, param model.DonorDonationTransactionDetailParam) (*model.DonorDonationTransactionDetailRow, error) {
	var row model.DonorDonationTransactionDetailRow

	err := buildDonorDonationTransactionBaseQuery(tx).
		Select(baseDonorDonationTransactionSelect()+`,
			p.address AS post_address,
			p.latitude,
			p.longitude,
			req.funding_target,
			req.funded_amount,
			COUNT(DISTINCT CASE WHEN d_all.donation_status = 'approved' THEN d_all.donated_by END) AS donor_count,
			COUNT(DISTINCT oi.order_item_id) AS total_item_count
		`).
		Joins("LEFT JOIN donations AS d_all ON d_all.request_id = req.request_id").
		Joins("LEFT JOIN order_items AS oi ON oi.order_id = o.order_id").
		Where("d.donated_by = ? AND d.donation_id = ?", param.UserID, param.DonationID).
		Group(donorDonationTransactionGroupColumns() + ", p.address, p.latitude, p.longitude, req.funding_target, req.funded_amount").
		Scan(&row).Error

	if err != nil {
		return nil, err
	}
	if row.DonationID == uuid.Nil {
		return nil, gorm.ErrRecordNotFound
	}
	return &row, nil
}

func (r *DonationRepository) GetDonorDonationTransactionItems(tx *gorm.DB, param model.DonorDonationTransactionDetailParam) ([]model.DonorDonationTransactionItemRow, error) {
	var rows []model.DonorDonationTransactionItemRow
	err := tx.Table("donations AS d").
		Select("i.item_id, i.name, oi.quantity, oi.unit_price, oi.subtotal").
		Joins("JOIN orders AS o ON o.order_code = "+donationTransactionCodeSQL("d.donation_id")).
		Joins("JOIN order_items AS oi ON oi.order_id = o.order_id").
		Joins("JOIN items AS i ON i.item_id = oi.item_id").
		Where("d.donated_by = ? AND d.donation_id = ?", param.UserID, param.DonationID).
		Order("oi.created_at ASC").
		Scan(&rows).Error
	return rows, err
}

func (r *DonationRepository) GetDonorDonationTransactionCustodyLogs(tx *gorm.DB, param model.DonorDonationTransactionDetailParam) ([]model.DonorDonationTransactionCustodyLogRow, error) {
	var rows []model.DonorDonationTransactionCustodyLogRow
	err := tx.Table("donations AS d").
		Select("cl.logs_id, cl.order_id, cl.sequence, cl.from_actor_id, cl.to_actor_id, cl.latitude, cl.longitude, cl.prev_hash, cl.current_hash, cl.created_at").
		Joins("JOIN orders AS o ON o.order_code = "+donationTransactionCodeSQL("d.donation_id")).
		Joins("JOIN custody_logs AS cl ON cl.order_id = o.order_id").
		Where("d.donated_by = ? AND d.donation_id = ?", param.UserID, param.DonationID).
		Order("cl.sequence ASC").
		Scan(&rows).Error
	return rows, err
}

func buildDonorDonationTransactionBaseQuery(tx *gorm.DB) *gorm.DB {
	custodyCountSubquery := tx.Table("custody_logs").
		Select("order_id, COUNT(*) AS custody_step_count").
		Group("order_id")

	latestCustodySubquery := tx.Table("custody_logs").
		Select("order_id, MAX(sequence) AS latest_sequence").
		Group("order_id")

	verificationProofSubquery := tx.Table("delivery_verifications").
		Select(`
			order_id,
			MAX(captured_at) AS verified_at,
			SUBSTRING_INDEX(GROUP_CONCAT(image_url ORDER BY captured_at DESC SEPARATOR '||'), '||', 1) AS verification_image_url
		`).
		Where("verification_status = ?", "approved").
		Group("order_id")

	return tx.Table("donations AS d").
		Joins("JOIN payment_transactions AS pt ON pt.donation_id = d.donation_id").
		Joins("JOIN requests AS req ON req.request_id = d.request_id").
		Joins("JOIN disaster_reports AS dr ON dr.report_id = req.report_id").
		Joins("JOIN posts AS p ON p.post_id = dr.post_id").
		Joins("LEFT JOIN orders AS o ON o.order_code = "+donationTransactionCodeSQL("d.donation_id")).
		Joins("LEFT JOIN (?) AS cc ON cc.order_id = o.order_id", custodyCountSubquery).
		Joins("LEFT JOIN (?) AS lc ON lc.order_id = o.order_id", latestCustodySubquery).
		Joins("LEFT JOIN custody_logs AS cl ON cl.order_id = lc.order_id AND cl.sequence = lc.latest_sequence").
		Joins("LEFT JOIN (?) AS vp ON vp.order_id = o.order_id", verificationProofSubquery)
}

func baseDonorDonationTransactionSelect() string {
	return `
		d.donation_id,
		pt.payment_transaction_id,
		pt.order_id AS payment_order_id,
		d.request_id,
		COALESCE(o.order_id, '') AS locked_order_id,
		COALESCE(o.order_code, '') AS transaction_code,
		p.name AS post_name,
		req.title AS request_title,
		d.donation_amount AS amount,
		d.donation_status,
		pt.transaction_status AS payment_status,
		COALESCE(o.order_status, '') AS order_status,
		COALESCE(cl.current_hash, '') AS latest_hash,
		COALESCE(vp.verification_image_url, '') AS verification_image_url,
		COALESCE(cc.custody_step_count, 0) AS custody_step_count,
		d.donated_at,
		pt.paid_at,
		vp.verified_at
	`
}

func donorDonationTransactionGroupColumns() string {
	return `
		d.donation_id,
		pt.payment_transaction_id,
		pt.order_id,
		d.request_id,
		o.order_id,
		o.order_code,
		p.name,
		req.title,
		d.donation_amount,
		d.donation_status,
		pt.transaction_status,
		o.order_status,
		cl.current_hash,
		vp.verification_image_url,
		cc.custody_step_count,
		d.donated_at,
		pt.paid_at,
		vp.verified_at
	`
}

func applyDonorDonationTransactionStatusFilter(query *gorm.DB, status string) *gorm.DB {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "", "all":
		return query
	case "pending":
		return query.Where("d.donation_status = ?", "pending")
	case "locked":
		return query.Where("d.donation_status = ? AND o.order_id IS NOT NULL AND vp.verified_at IS NULL AND (o.order_status = ? OR o.order_status IS NULL)", "approved", "pending")
	case "preparing":
		return query.Where("d.donation_status = ? AND o.order_status IN ? AND vp.verified_at IS NULL", "approved", []string{"accepted", "preparing"})
	case "ready":
		return query.Where("d.donation_status = ? AND o.order_status = ? AND vp.verified_at IS NULL", "approved", "ready_for_pickup")
	case "shipping":
		return query.Where("d.donation_status = ? AND o.order_status IN ? AND vp.verified_at IS NULL", "approved", []string{"picked_up", "in_transit", "delivered"})
	case "completed":
		return query.Where("vp.verified_at IS NOT NULL")
	case "refund":
		return query.Where(
			"d.donation_status = ? OR pt.transaction_status IN ?",
			"rejected",
			[]string{"expire", "cancel", "deny", "failure"},
		)
	default:
		return query
	}
}

func normalizeDonorDonationTransactionLimit(limit int) int {
	if limit <= 0 {
		return 20
	}
	if limit > 100 {
		return 100
	}
	return limit
}

func normalizeDonorDonationTransactionOffset(offset int) int {
	if offset < 0 {
		return 0
	}
	return offset
}

func donationTransactionCodeSQL(column string) string {
	cleanColumn := "REPLACE(UPPER(" + column + "), '-', '')"
	return "CONCAT('DN-', SUBSTRING(" + cleanColumn + ", 1, 4), '-', SUBSTRING(" + cleanColumn + ", 5, 5))"
}
