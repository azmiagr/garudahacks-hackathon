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

func normalizePublicLimit(limit int) int {
	if limit <= 0 {
		return 50
	}

	if limit > 500 {
		return 500
	}

	return limit
}
