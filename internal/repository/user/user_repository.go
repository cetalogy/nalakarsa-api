package userrepository

import (
	"errors"
	"strings"
	"time"

	"nalakarsa/internal/model/discussion"
	"nalakarsa/internal/model/user"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(u *user.User) error
	GetByEmail(email string) (*user.User, error)
	GetByID(id uuid.UUID) (*user.User, error)
	UpdateProfile(u *user.User) error
	UpdateAvatar(userID uuid.UUID, avatarURL string) error
	IncrementViewCount(userID uuid.UUID) error
	ListUsers(search, role string, page, limit int) ([]user.User, int64, error)
	CreateRefreshToken(rt *user.RefreshToken) error
	GetRefreshToken(token string) (*user.RefreshToken, error)
	DeleteRefreshToken(token string) error
	DeleteRefreshTokensByUserID(userID uuid.UUID) error
	CountActiveRefreshTokens(userID uuid.UUID) (int64, error)
	DeleteOldestRefreshTokensByUser(userID uuid.UUID, keepLatest int) error
	CountDiscussions(userID uuid.UUID) (int64, error)
}

type pgUserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &pgUserRepository{db: db}
}

func (r *pgUserRepository) Create(u *user.User) error {
	// Normalize email to lowercase
	u.Email = strings.ToLower(strings.TrimSpace(u.Email))

	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(u).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *pgUserRepository) GetByEmail(email string) (*user.User, error) {
	var u user.User
	err := r.db.Where("email = ?", strings.ToLower(strings.TrimSpace(email))).First(&u).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *pgUserRepository) GetByID(id uuid.UUID) (*user.User, error) {
	var u user.User
	err := r.db.Where("id = ?", id).First(&u).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *pgUserRepository) UpdateProfile(u *user.User) error {
	return r.db.Model(&user.User{}).
		Where("id = ?", u.ID).
		Updates(u).Error
}

func (r *pgUserRepository) UpdateAvatar(userID uuid.UUID, avatarURL string) error {
	return r.db.Model(&user.User{}).
		Where("id = ?", userID).
		Update("avatar_url", avatarURL).Error
}

func (r *pgUserRepository) IncrementViewCount(userID uuid.UUID) error {
	return r.db.Model(&user.User{}).
		Where("id = ?", userID).
		UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error
}

func (r *pgUserRepository) ListUsers(search, role string, page, limit int) ([]user.User, int64, error) {
	var users []user.User
	var total int64

	query := r.db.Model(&user.User{})

	if role != "" {
		query = query.Where("role = ?", role)
	}

	if search != "" {
		searchTerm := "%" + search + "%"
		query = query.Where("full_name ILIKE ? OR affiliation ILIKE ? OR expertise ILIKE ?", searchTerm, searchTerm, searchTerm)
	}

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Fetch paginated results
	offset := (page - 1) * limit
	err := query.Limit(limit).Offset(offset).Order("users.created_at desc").Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *pgUserRepository) CreateRefreshToken(rt *user.RefreshToken) error {
	return r.db.Create(rt).Error
}

func (r *pgUserRepository) GetRefreshToken(token string) (*user.RefreshToken, error) {
	var rt user.RefreshToken
	err := r.db.Where("token = ?", token).First(&rt).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &rt, nil
}

func (r *pgUserRepository) DeleteRefreshToken(token string) error {
	return r.db.Where("token = ?", token).Delete(&user.RefreshToken{}).Error
}

func (r *pgUserRepository) DeleteRefreshTokensByUserID(userID uuid.UUID) error {
	return r.db.Where("user_id = ?", userID).Delete(&user.RefreshToken{}).Error
}

func (r *pgUserRepository) CountActiveRefreshTokens(userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&user.RefreshToken{}).
		Where("user_id = ?", userID).
		Where("expires_at > ?", time.Now()).
		Count(&count).Error
	return count, err
}

func (r *pgUserRepository) DeleteOldestRefreshTokensByUser(userID uuid.UUID, keepLatest int) error {
	if keepLatest < 0 {
		keepLatest = 0
	}

	// Cleanup expired refresh tokens first so active limit focuses on valid sessions.
	if err := r.db.Where("user_id = ? AND expires_at <= ?", userID, time.Now()).Delete(&user.RefreshToken{}).Error; err != nil {
		return err
	}

	var ids []uuid.UUID
	err := r.db.Model(&user.RefreshToken{}).
		Where("user_id = ?", userID).
		Where("expires_at > ?", time.Now()).
		Order("created_at desc").
		Offset(keepLatest).
		Pluck("id", &ids).Error
	if err != nil || len(ids) == 0 {
		return err
	}

	return r.db.Where("id IN ?", ids).Delete(&user.RefreshToken{}).Error
}

func (r *pgUserRepository) CountDiscussions(userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&discussion.Discussion{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}
