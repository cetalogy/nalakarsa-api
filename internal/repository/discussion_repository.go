package repository

import (
	"errors"

	"nalakarsa/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DiscussionRepository interface {
	Create(disc *model.Discussion) error
	GetByID(id uuid.UUID) (*model.Discussion, error)
	List(search, category, role, sort string, page, limit int) ([]model.Discussion, int64, error)
	Update(disc *model.Discussion) error
	Delete(id uuid.UUID) error
	CreateComment(comment *model.Comment) error
	GetCommentByID(id uuid.UUID) (*model.Comment, error)
	DeleteComment(id uuid.UUID) error
}

type pgDiscussionRepository struct {
	db *gorm.DB
}

func NewDiscussionRepository(db *gorm.DB) DiscussionRepository {
	return &pgDiscussionRepository{db: db}
}

func (r *pgDiscussionRepository) Create(disc *model.Discussion) error {
	return r.db.Create(disc).Error
}

func (r *pgDiscussionRepository) GetByID(id uuid.UUID) (*model.Discussion, error) {
	var disc model.Discussion
	err := r.db.Preload("User.Profile").
		Preload("Comments", func(db *gorm.DB) *gorm.DB {
			return db.Order("comments.created_at asc")
		}).
		Preload("Comments.User.Profile").
		Where("id = ?", id).First(&disc).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &disc, nil
}

func (r *pgDiscussionRepository) List(search, category, role, sort string, page, limit int) ([]model.Discussion, int64, error) {
	var discussions []model.Discussion
	var total int64

	query := r.db.Model(&model.Discussion{}).Preload("User.Profile")

	if category != "" {
		query = query.Where("category = ?", category)
	}

	if role != "" {
		query = query.Joins("INNER JOIN users ON users.id = discussions.user_id").
			Where("users.role = ?", role)
	}

	if search != "" {
		searchTerm := "%" + search + "%"
		query = query.Where("title ILIKE ? OR content ILIKE ?", searchTerm, searchTerm)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply Sorting
	order := "created_at desc"
	if sort == "oldest" {
		order = "created_at asc"
	}

	// Fetch paginated results
	offset := (page - 1) * limit
	err := query.Limit(limit).Offset(offset).Order(order).Find(&discussions).Error
	if err != nil {
		return nil, 0, err
	}

	return discussions, total, nil
}

func (r *pgDiscussionRepository) Update(disc *model.Discussion) error {
	return r.db.Save(disc).Error
}

func (r *pgDiscussionRepository) Delete(id uuid.UUID) error {
	return r.db.Where("id = ?", id).Delete(&model.Discussion{}).Error
}

func (r *pgDiscussionRepository) CreateComment(comment *model.Comment) error {
	return r.db.Create(comment).Error
}

func (r *pgDiscussionRepository) GetCommentByID(id uuid.UUID) (*model.Comment, error) {
	var comment model.Comment
	err := r.db.Where("id = ?", id).First(&comment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &comment, nil
}

func (r *pgDiscussionRepository) DeleteComment(id uuid.UUID) error {
	return r.db.Where("id = ?", id).Delete(&model.Comment{}).Error
}
