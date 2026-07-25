package service

import (
	"nalakarsa/internal/dto"
	"nalakarsa/internal/repository"

	"github.com/google/uuid"
)

type NotificationService interface {
	ListNotifications(userID uuid.UUID, page, limit int) ([]dto.NotificationResponse, int64, error)
	MarkRead(userID uuid.UUID, notifID uuid.UUID) error
	MarkAllRead(userID uuid.UUID) error
}

type notificationService struct {
	notifRepo repository.NotificationRepository
}

func NewNotificationService(notifRepo repository.NotificationRepository) NotificationService {
	return &notificationService{notifRepo: notifRepo}
}

func (s *notificationService) ListNotifications(userID uuid.UUID, page, limit int) ([]dto.NotificationResponse, int64, error) {
	notifs, total, err := s.notifRepo.ListByUser(userID, page, limit)
	if err != nil {
		return nil, 0, err
	}

	res := make([]dto.NotificationResponse, len(notifs))
	for i, n := range notifs {
		actorName := ""
		if n.Actor != nil {
			actorName = n.Actor.Profile.NamaLengkap
		}

		res[i] = dto.NotificationResponse{
			ID:           n.ID,
			Type:         n.Type,
			ActorID:      n.ActorID,
			ActorName:    actorName,
			ResourceType: n.ResourceType,
			ResourceID:   n.ResourceID,
			ReadAt:       n.ReadAt,
			CreatedAt:    n.CreatedAt,
		}
	}

	return res, total, nil
}

func (s *notificationService) MarkRead(userID uuid.UUID, notifID uuid.UUID) error {
	return s.notifRepo.MarkRead(notifID)
}

func (s *notificationService) MarkAllRead(userID uuid.UUID) error {
	return s.notifRepo.MarkAllRead(userID)
}
