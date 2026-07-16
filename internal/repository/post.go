package repository

import (
	"strings"

	"github.com/azmiagr/garudahacks-hackathon/entity"
	"github.com/azmiagr/garudahacks-hackathon/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IPostRepository interface {
	CreatePost(tx *gorm.DB, post *entity.Post) error
	GetPost(tx *gorm.DB, param model.GetPostParam) (*entity.Post, error)
	GetPostByIDs(tx *gorm.DB, postIDs []uuid.UUID) ([]*entity.Post, error)
	GetPublicMapPosts(tx *gorm.DB, param model.PublicMapPostParam) ([]model.PublicMapPostRow, error)
	GetDonorPostDetail(tx *gorm.DB, postID uuid.UUID) (*model.DonorPostDetailRow, error)
}

type PostRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) IPostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) CreatePost(tx *gorm.DB, post *entity.Post) error {
	err := tx.Create(post).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *PostRepository) GetPost(tx *gorm.DB, param model.GetPostParam) (*entity.Post, error) {
	var post entity.Post
	err := tx.Where(&param).First(&post).Error
	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (r *PostRepository) GetPostByIDs(tx *gorm.DB, postIDs []uuid.UUID) ([]*entity.Post, error) {
	var posts []*entity.Post

	if len(postIDs) == 0 {
		return posts, nil
	}

	err := tx.Where("post_id IN ?", postIDs).
		Find(&posts).
		Error
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *PostRepository) GetPublicMapPosts(tx *gorm.DB, param model.PublicMapPostParam) ([]model.PublicMapPostRow, error) {
	var rows []model.PublicMapPostRow

	query := tx.Table("posts").
		Select("post_id, name, address, latitude, longitude")

	if param.MinLatitude != nil && param.MaxLatitude != nil {
		query = query.Where("latitude BETWEEN ? AND ?", *param.MinLatitude, *param.MaxLatitude)
	}

	if param.MinLongitude != nil && param.MaxLongitude != nil {
		query = query.Where("longitude BETWEEN ? AND ?", *param.MinLongitude, *param.MaxLongitude)
	}

	if strings.TrimSpace(param.Query) != "" {
		keyword := "%" + strings.TrimSpace(param.Query) + "%"
		query = query.Where("(name LIKE ? OR address LIKE ?)", keyword, keyword)
	}

	err := query.
		Order("created_at DESC").
		Limit(normalizePublicLimit(param.Limit)).
		Scan(&rows).
		Error
	if err != nil {
		return nil, err
	}

	return rows, nil
}

func (r *PostRepository) GetDonorPostDetail(tx *gorm.DB, postID uuid.UUID) (*model.DonorPostDetailRow, error) {
	var row model.DonorPostDetailRow

	latestReportSubquery := tx.Table("disaster_reports").
		Select("post_id, MAX(COALESCE(reported_at, created_at)) AS latest_reported_at").
		Where("report_status = ?", "approved").
		Group("post_id")

	err := tx.Table("posts AS p").
		Select(`
			p.post_id,
			dr.report_id,
			req.request_id,
			p.name,
			p.address,
			de.name AS disaster_type,
			COALESCE(dr.image_url, '') AS image_url,
			p.created_at,
			COALESCE(dr.reported_at, dr.created_at) AS reported_at,
			COALESCE(req.funding_target, 0) AS funding_target,
			COALESCE(req.funded_amount, 0) AS funded_amount,
			COUNT(DISTINCT CASE WHEN d.donation_status = 'approved' THEN d.donated_by END) AS donor_count,
			u.kyc_status AS admin_kyc_status
		`).
		Joins("JOIN (?) AS latest_reports ON latest_reports.post_id = p.post_id", latestReportSubquery).
		Joins("JOIN disaster_reports AS dr ON dr.post_id = latest_reports.post_id AND COALESCE(dr.reported_at, dr.created_at) = latest_reports.latest_reported_at").
		Joins("JOIN disaster_events AS de ON de.event_id = dr.event_id").
		Joins("JOIN requests AS req ON req.report_id = dr.report_id AND req.request_status = 'approved'").
		Joins("JOIN users AS u ON u.user_id = p.user_id").
		Joins("LEFT JOIN donations AS d ON d.request_id = req.request_id").
		Where("p.post_id = ?", postID).
		Group("p.post_id, dr.report_id, req.request_id, p.name, p.address, de.name, dr.image_url, p.created_at, dr.reported_at, dr.created_at, req.funding_target, req.funded_amount, u.kyc_status").
		Scan(&row).Error
	if err != nil {
		return nil, err
	}

	if row.PostID == uuid.Nil {
		return nil, gorm.ErrRecordNotFound
	}

	return &row, nil
}

func normalizePublicLimit(limit int) int {
	if limit <= 0 {
		return 50
	}

	if limit > 500 {
		return 500
	}

	return limit
}
