package service

import (
	"fmt"
	"time"

	"nalakarsa/internal/dto"
	"nalakarsa/internal/repository"

	"github.com/google/uuid"
)

type DashboardService interface {
	GetDashboard(userID uuid.UUID) (*dto.DashboardResponse, error)
}

type dashboardService struct {
	projRepo  repository.ProjectRepository
	connRepo  repository.ConnectionRepository
	convRepo  repository.ConversationRepository
	notifRepo repository.NotificationRepository
}

func NewDashboardService(
	projRepo repository.ProjectRepository,
	connRepo repository.ConnectionRepository,
	convRepo repository.ConversationRepository,
	notifRepo repository.NotificationRepository,
) DashboardService {
	return &dashboardService{
		projRepo:  projRepo,
		connRepo:  connRepo,
		convRepo:  convRepo,
		notifRepo: notifRepo,
	}
}

func (s *dashboardService) GetDashboard(userID uuid.UUID) (*dto.DashboardResponse, error) {
	// Active projects count (where user is owner or member)
	activeCount, _ := s.projRepo.CountByOwner(userID, "active")

	// Unread messages
	unreadMessages, _ := s.convRepo.CountTotalUnread(userID)

	// Active connections
	activeConns, _ := s.connRepo.CountAccepted(userID)

	// My projects
	projects, _ := s.projRepo.ListByUser(userID)
	myProjects := make([]dto.DashboardProjectItem, 0, len(projects))
	for _, p := range projects {
		item := dto.DashboardProjectItem{
			ID:       p.ID.String(),
			Title:    p.Title,
			Role:     "owner",
			Progress: p.Progress,
			Status:   p.Status,
		}

		if p.OwnerID != userID {
			item.Role = "member"
		}

		// Next milestone
		nextMs, _ := s.projRepo.GetNextMilestone(p.ID)
		if nextMs != nil {
			item.NextMilestone = &nextMs.Title
		}

		if p.Deadline != nil {
			dl := p.Deadline.Format(time.RFC3339)
			item.Deadline = &dl
		}

		myProjects = append(myProjects, item)
	}

	// Recent activity (notifications)
	notifs, _, _ := s.notifRepo.ListByUser(userID, 1, 10)
	recentActivity := make([]dto.NotificationResponse, len(notifs))
	for i, n := range notifs {
		actorName := ""
		if n.Actor != nil {
			actorName = n.Actor.Profile.NamaLengkap
		}
		recentActivity[i] = dto.NotificationResponse{
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

	_ = fmt.Sprintf // suppress unused import if needed

	return &dto.DashboardResponse{
		ActiveProjectsCount: activeCount,
		UnreadMessagesCount: unreadMessages,
		ActiveConnections:   activeConns,
		MyProjects:          myProjects,
		RecentActivity:      recentActivity,
	}, nil
}
