package discussionrepository

import (
	"errors"

	"nalakarsa/internal/model/discussion"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DiscussionRepository interface {
	Create(disc *discussion.Discussion) error
	GetByID(id uuid.UUID) (*discussion.Discussion, error)
	List(search, category, role, sort string, page, limit int) ([]discussion.Discussion, int64, error)
	Update(disc *discussion.Discussion) error
	Delete(id uuid.UUID) error

	// Replies (was Comments)
	CreateReply(reply *discussion.DiscussionReply) error
	GetReplyByID(id uuid.UUID) (*discussion.DiscussionReply, error)
	DeleteReply(id uuid.UUID) error
	ListReplies(discussionID uuid.UUID, page, limit int) ([]discussion.DiscussionReply, int64, error)
	CountReplies(discussionID uuid.UUID) (int64, error)

	// Votes
	CreateVote(vote *discussion.DiscussionVote) error
	DeleteVote(userID, discussionID uuid.UUID) error
	HasVoted(userID, discussionID uuid.UUID) (bool, error)
	CountVotes(discussionID uuid.UUID) (int64, error)
}

type pgDiscussionRepository struct {
	db *gorm.DB
}

func NewDiscussionRepository(db *gorm.DB) DiscussionRepository {
	return &pgDiscussionRepository{db: db}
}

func (r *pgDiscussionRepository) Create(disc *discussion.Discussion) error {
	return r.db.Create(disc).Error
}

func (r *pgDiscussionRepository) GetByID(id uuid.UUID) (*discussion.Discussion, error) {
	var disc discussion.Discussion
	err := r.db.Preload("User").
		Preload("Replies", func(db *gorm.DB) *gorm.DB {
			return db.Order("discussion_replies.created_at asc")
		}).
		Preload("Replies.User").
		Where("id = ?", id).First(&disc).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &disc, nil
}

func (r *pgDiscussionRepository) List(search, category, role, sort string, page, limit int) ([]discussion.Discussion, int64, error) {
	var discussions []discussion.Discussion
	var total int64

	query := r.db.Model(&discussion.Discussion{}).Preload("User")

	if category != "" {
		query = query.Where("discussions.category = ?", category)
	}

	if role != "" {
		query = query.Joins("INNER JOIN users ON users.id = discussions.user_id").
			Where("users.role = ?", role)
	}

	if search != "" {
		searchTerm := "%" + search + "%"
		query = query.Where("discussions.title ILIKE ? OR discussions.description ILIKE ?", searchTerm, searchTerm)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply Sorting
	order := "discussions.created_at desc"
	if sort == "oldest" {
		order = "discussions.created_at asc"
	}

	// Fetch paginated results
	offset := (page - 1) * limit
	err := query.Limit(limit).Offset(offset).Order(order).Find(&discussions).Error
	if err != nil {
		return nil, 0, err
	}

	return discussions, total, nil
}

func (r *pgDiscussionRepository) Update(disc *discussion.Discussion) error {
	return r.db.Save(disc).Error
}

func (r *pgDiscussionRepository) Delete(id uuid.UUID) error {
	return r.db.Where("id = ?", id).Delete(&discussion.Discussion{}).Error
}

// --- Replies ---

func (r *pgDiscussionRepository) CreateReply(reply *discussion.DiscussionReply) error {
	return r.db.Create(reply).Error
}

func (r *pgDiscussionRepository) GetReplyByID(id uuid.UUID) (*discussion.DiscussionReply, error) {
	var reply discussion.DiscussionReply
	err := r.db.Where("id = ?", id).First(&reply).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &reply, nil
}

func (r *pgDiscussionRepository) ListReplies(discussionID uuid.UUID, page, limit int) ([]discussion.DiscussionReply, int64, error) {
	var replies []discussion.DiscussionReply
	var total int64

	query := r.db.Model(&discussion.DiscussionReply{}).Where("discussion_id = ?", discussionID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := query.
		Preload("User").
		Order("discussion_replies.created_at asc").
		Limit(limit).Offset(offset).
		Find(&replies).Error
	return replies, total, err
}

func (r *pgDiscussionRepository) DeleteReply(id uuid.UUID) error {
	return r.db.Where("id = ?", id).Delete(&discussion.DiscussionReply{}).Error
}

func (r *pgDiscussionRepository) CountReplies(discussionID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&discussion.DiscussionReply{}).Where("discussion_id = ?", discussionID).Count(&count).Error
	return count, err
}

// --- Votes ---

func (r *pgDiscussionRepository) CreateVote(vote *discussion.DiscussionVote) error {
	return r.db.Create(vote).Error
}

func (r *pgDiscussionRepository) DeleteVote(userID, discussionID uuid.UUID) error {
	return r.db.Where("user_id = ? AND discussion_id = ?", userID, discussionID).Delete(&discussion.DiscussionVote{}).Error
}

func (r *pgDiscussionRepository) HasVoted(userID, discussionID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Model(&discussion.DiscussionVote{}).Where("user_id = ? AND discussion_id = ?", userID, discussionID).Count(&count).Error
	return count > 0, err
}

func (r *pgDiscussionRepository) CountVotes(discussionID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&discussion.DiscussionVote{}).Where("discussion_id = ?", discussionID).Count(&count).Error
	return count, err
}
