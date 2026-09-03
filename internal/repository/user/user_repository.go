package userrepository

import (
	"errors"
	"strings"
	"time"

	"nalakarsa/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(u *model.User) error
	GetByEmail(email string) (*model.User, error)
	GetByID(id uuid.UUID) (*model.User, error)
	UpdateProfile(u *model.User) error
	UpdateAvatar(userID uuid.UUID, avatarURL string) error
	IncrementViewCount(userID uuid.UUID) error
	ListUsers(search, role string, page, limit int) ([]model.User, int64, error)
	CreateRefreshToken(rt *model.RefreshToken) error
	GetRefreshToken(token string) (*model.RefreshToken, error)
	DeleteRefreshToken(token string) error
	DeleteRefreshTokensByUserID(userID uuid.UUID) error
	CountActiveRefreshTokens(userID uuid.UUID) (int64, error)
	DeleteOldestRefreshTokensByUser(userID uuid.UUID, keepLatest int) error
	CountDiscussions(userID uuid.UUID) (int64, error)
	ToggleFollow(followerID, followingID uuid.UUID) (bool, error)
	IsFollowing(followerID, followingID uuid.UUID) (bool, error)
	GetFollowers(userID uuid.UUID, page, limit int) ([]model.User, int64, error)
	GetFollowing(userID uuid.UUID, page, limit int) ([]model.User, int64, error)
	CountFollowers(userID uuid.UUID) (int64, error)
	CountFollowing(userID uuid.UUID) (int64, error)
	GetByIDOrIdentifier(identifier string) (*model.User, error)
}

type pgUserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &pgUserRepository{db: db}
}

func (r *pgUserRepository) Create(u *model.User) error {
	u.Email = strings.ToLower(strings.TrimSpace(u.Email))

	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(u).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *pgUserRepository) GetByEmail(email string) (*model.User, error) {
	var u model.User
	err := r.db.Where("email = ?", strings.ToLower(strings.TrimSpace(email))).First(&u).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *pgUserRepository) GetByID(id uuid.UUID) (*model.User, error) {
	var u model.User
	err := r.db.Where("id = ?", id).First(&u).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *pgUserRepository) GetByIDOrIdentifier(identifier string) (*model.User, error) {
	trimmed := strings.TrimSpace(identifier)
	if trimmed == "" {
		return nil, nil
	}

	trimmedClean := strings.TrimPrefix(trimmed, "user_")
	if uid, err := uuid.Parse(trimmedClean); err == nil {
		return r.GetByID(uid)
	}
	if strings.Contains(trimmed, "@") {
		return r.GetByEmail(trimmed)
	}
	var u model.User
	err := r.db.Where("LOWER(TRIM(full_name)) = LOWER(TRIM(?))", trimmedClean).First(&u).Error
	if err == nil {
		return &u, nil
	}

	err = r.db.Where("LOWER(full_name) LIKE LOWER(?)", "%"+trimmedClean+"%").First(&u).Error
	if err == nil {
		return &u, nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return nil, err
}

func (r *pgUserRepository) UpdateProfile(u *model.User) error {
	updates := map[string]interface{}{
		"first_name": u.FirstName, "middle_name": u.MiddleName, "last_name": u.LastName,
		"full_name": u.FullName, "prefix_title": u.PrefixTitle, "suffix_title": u.SuffixTitle,
		"affiliation": u.Affiliation, "location": u.Location, "expertise": u.Expertise,
		"industry": u.Industry, "bio": u.Bio, "avatar_url": u.AvatarURL,
	}
	if u.PasswordHash != "" {
		updates["password_hash"] = u.PasswordHash
	}
	return r.db.Model(&model.User{}).
		Where("id = ?", u.ID).
		Updates(updates).Error
}

func (r *pgUserRepository) UpdateAvatar(userID uuid.UUID, avatarURL string) error {
	return r.db.Model(&model.User{}).
		Where("id = ?", userID).
		Update("avatar_url", avatarURL).Error
}

func (r *pgUserRepository) IncrementViewCount(userID uuid.UUID) error {
	return r.db.Model(&model.User{}).
		Where("id = ?", userID).
		UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error
}

func (r *pgUserRepository) ListUsers(search, role string, page, limit int) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	query := r.db.Model(&model.User{})

	if role != "" {
		query = query.Where("role = ?", role)
	}

	if search != "" {
		searchTerm := "%" + search + "%"
		query = query.Where("full_name ILIKE ? OR affiliation ILIKE ? OR expertise ILIKE ?", searchTerm, searchTerm, searchTerm)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * limit
	err := query.Limit(limit).Offset(offset).Order("users.created_at desc").Find(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *pgUserRepository) CreateRefreshToken(rt *model.RefreshToken) error {
	if err := r.db.Where("expires_at <= ?", time.Now()).Delete(&model.RefreshToken{}).Error; err != nil {
		return err
	}
	return r.db.Create(rt).Error
}

func (r *pgUserRepository) GetRefreshToken(token string) (*model.RefreshToken, error) {
	if err := r.db.Where("expires_at <= ?", time.Now()).Delete(&model.RefreshToken{}).Error; err != nil {
		return nil, err
	}
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

func (r *pgUserRepository) CountActiveRefreshTokens(userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&model.RefreshToken{}).
		Where("user_id = ?", userID).
		Where("expires_at > ?", time.Now()).
		Count(&count).Error
	return count, err
}

func (r *pgUserRepository) DeleteOldestRefreshTokensByUser(userID uuid.UUID, keepLatest int) error {
	if keepLatest < 0 {
		keepLatest = 0
	}
	if err := r.db.Where("user_id = ? AND expires_at <= ?", userID, time.Now()).Delete(&model.RefreshToken{}).Error; err != nil {
		return err
	}

	var ids []uuid.UUID
	err := r.db.Model(&model.RefreshToken{}).
		Where("user_id = ?", userID).
		Where("expires_at > ?", time.Now()).
		Order("created_at desc").
		Offset(keepLatest).
		Pluck("id", &ids).Error
	if err != nil || len(ids) == 0 {
		return err
	}

	return r.db.Where("id IN ?", ids).Delete(&model.RefreshToken{}).Error
}

func (r *pgUserRepository) CountDiscussions(userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&model.Discussion{}).Where("user_id = ?", userID).Count(&count).Error
	return count, err
}

func (r *pgUserRepository) ToggleFollow(followerID, followingID uuid.UUID) (bool, error) {
	var existing model.UserFollower
	err := r.db.Where("follower_id = ? AND following_id = ?", followerID, followingID).First(&existing).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newFollow := model.UserFollower{
				FollowerID:  followerID,
				FollowingID: followingID,
			}
			if createErr := r.db.Create(&newFollow).Error; createErr != nil {
				return false, createErr
			}
			return true, nil
		}
		return false, err
	}
	if deleteErr := r.db.Where("follower_id = ? AND following_id = ?", followerID, followingID).Delete(&model.UserFollower{}).Error; deleteErr != nil {
		return false, deleteErr
	}
	return false, nil
}

func (r *pgUserRepository) IsFollowing(followerID, followingID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Model(&model.UserFollower{}).
		Where("follower_id = ? AND following_id = ?", followerID, followingID).
		Count(&count).Error
	return count > 0, err
}

func (r *pgUserRepository) GetFollowers(userID uuid.UUID, page, limit int) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	subQuery := r.db.Model(&model.UserFollower{}).Select("follower_id").Where("following_id = ?", userID)

	if err := r.db.Model(&model.User{}).Where("id IN (?)", subQuery).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := r.db.Model(&model.User{}).
		Joins("JOIN user_followers ON user_followers.follower_id = users.id").
		Where("user_followers.following_id = ?", userID).
		Order("user_followers.created_at desc").
		Limit(limit).Offset(offset).
		Find(&users).Error

	return users, total, err
}

func (r *pgUserRepository) GetFollowing(userID uuid.UUID, page, limit int) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	subQuery := r.db.Model(&model.UserFollower{}).Select("following_id").Where("follower_id = ?", userID)

	if err := r.db.Model(&model.User{}).Where("id IN (?)", subQuery).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := r.db.Model(&model.User{}).
		Joins("JOIN user_followers ON user_followers.following_id = users.id").
		Where("user_followers.follower_id = ?", userID).
		Order("user_followers.created_at desc").
		Limit(limit).Offset(offset).
		Find(&users).Error

	return users, total, err
}

func (r *pgUserRepository) CountFollowers(userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&model.UserFollower{}).Where("following_id = ?", userID).Count(&count).Error
	return count, err
}

func (r *pgUserRepository) CountFollowing(userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&model.UserFollower{}).Where("follower_id = ?", userID).Count(&count).Error
	return count, err
}
