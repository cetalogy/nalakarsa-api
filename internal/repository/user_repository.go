package repository

import (
	"errors"

	"nalakarsa/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *model.User) error
	GetByEmail(email string) (*model.User, error)
	GetByID(id uuid.UUID) (*model.User, error)
	UpdateProfile(profile *model.Profile) error
	UpdateAvatar(userID uuid.UUID, avatarURL string) error
	ListUsers(search, role string, page, limit int) ([]model.User, int64, error)
	CreateRefreshToken(rt *model.RefreshToken) error
	GetRefreshToken(token string) (*model.RefreshToken, error)
	DeleteRefreshToken(token string) error
	DeleteRefreshTokensByUserID(userID uuid.UUID) error
}

type pgUserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &pgUserRepository{db: db}
}

func (r *pgUserRepository) Create(user *model.User) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}
		user.Profile.UserID = user.ID
		if err := tx.Create(&user.Profile).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *pgUserRepository) GetByEmail(email string) (*model.User, error) {
	var user model.User
	err := r.db.Preload("Profile").Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *pgUserRepository) GetByID(id uuid.UUID) (*model.User, error) {
	var user model.User
	err := r.db.Preload("Profile").Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *pgUserRepository) UpdateProfile(profile *model.Profile) error {
	return r.db.Model(&model.Profile{}).
		Where("user_id = ?", profile.UserID).
		Updates(profile).Error
}

func (r *pgUserRepository) UpdateAvatar(userID uuid.UUID, avatarURL string) error {
	return r.db.Model(&model.Profile{}).
		Where("user_id = ?", userID).
		Update("avatar_url", avatarURL).Error
}

func (r *pgUserRepository) ListUsers(search, role string, page, limit int) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	query := r.db.Model(&model.User{}).Preload("Profile")

	if role != "" {
		query = query.Where("role = ?", role)
	}

	if search != "" {
		searchTerm := "%" + search + "%"
		// Join profiles table to filter on profile fields
		query = query.Joins("INNER JOIN profiles ON profiles.user_id = users.id").
			Where("profiles.nama_lengkap ILIKE ? OR profiles.afiliasi ILIKE ? OR profiles.bidang_keahlian ILIKE ?", searchTerm, searchTerm, searchTerm)
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

func (r *pgUserRepository) CreateRefreshToken(rt *model.RefreshToken) error {
	return r.db.Create(rt).Error
}

func (r *pgUserRepository) GetRefreshToken(token string) (*model.RefreshToken, error) {
	var rt model.RefreshToken
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
	return r.db.Where("token = ?", token).Delete(&model.RefreshToken{}).Error
}

func (r *pgUserRepository) DeleteRefreshTokensByUserID(userID uuid.UUID) error {
	return r.db.Where("user_id = ?", userID).Delete(&model.RefreshToken{}).Error
}
