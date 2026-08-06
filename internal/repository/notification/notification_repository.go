package notificationrepository

import (
	"nalakarsa/internal/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NotificationRepository interface {
	Create(notif *model.Notification) error
	ListByUser(userID uuid.UUID, page, limit int) ([]model.Notification, int64, error)
	MarkRead(id uuid.UUID) error
	MarkAllRead(userID uuid.UUID) error
	CountUnread(userID uuid.UUID) (int64, error)
}

type pgNotificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) NotificationRepository {
	return &pgNotificationRepository{db: db}
}

func (r *pgNotificationRepository) Create(notif *model.Notification) error {
	return r.db.Create(notif).Error
}

func (r *pgNotificationRepository) ListByUser(userID uuid.UUID, page, limit int) ([]model.Notification, int64, error) {
	var notifs []model.Notification
	var total int64

	query := r.db.Model(&model.Notification{}).Where("user_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := r.db.Preload("Actor.Profile").
		Where("user_id = ?", userID).
		Order("created_at desc").
		Limit(limit).Offset(offset).Find(&notifs).Error
	return notifs, total, err
}

func (r *pgNotificationRepository) MarkRead(id uuid.UUID) error {
	return r.db.Exec("UPDATE notifications SET read_at = NOW() WHERE id = ? AND read_at IS NULL", id).Error
}

func (r *pgNotificationRepository) MarkAllRead(userID uuid.UUID) error {
	return r.db.Exec("UPDATE notifications SET read_at = NOW() WHERE user_id = ? AND read_at IS NULL", userID).Error
}

func (r *pgNotificationRepository) CountUnread(userID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&model.Notification{}).
		Where("user_id = ? AND read_at IS NULL", userID).
		Count(&count).Error
	return count, err
}
