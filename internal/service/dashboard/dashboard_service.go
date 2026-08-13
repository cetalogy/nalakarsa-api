package dashboardservice

import (
	"time"

	"nalakarsa/internal/dto"
	"nalakarsa/internal/model/notification"
	"nalakarsa/internal/model/project"
	connectionrepository "nalakarsa/internal/repository/connection"
	conversationrepository "nalakarsa/internal/repository/conversation"
	notificationrepository "nalakarsa/internal/repository/notification"
	projectrepository "nalakarsa/internal/repository/project"

	"github.com/google/uuid"
)

type DashboardService interface {
	GetDashboard(userID uuid.UUID) (*dto.DashboardResponse, error)
	GetDashboardSummary(userID uuid.UUID) (*dto.DashboardSummaryResponse, error)
	GetDashboardOngoingProjects(userID uuid.UUID) (*dto.DashboardOngoingProjectsResponse, error)
	GetDashboardNotifications(userID uuid.UUID, page, limit int) ([]dto.NotificationResponse, int64, error)
}

type dashboardService struct {
	projRepo  projectrepository.ProjectRepository
	connRepo  connectionrepository.ConnectionRepository
	convRepo  conversationrepository.ConversationRepository
	notifRepo notificationrepository.NotificationRepository
}

func NewDashboardService(
	projRepo projectrepository.ProjectRepository,
	connRepo connectionrepository.ConnectionRepository,
	convRepo conversationrepository.ConversationRepository,
	notifRepo notificationrepository.NotificationRepository,
) DashboardService {
	return &dashboardService{
		projRepo:  projRepo,
		connRepo:  connRepo,
		convRepo:  convRepo,
		notifRepo: notifRepo,
	}
}

func (s *dashboardService) GetDashboard(userID uuid.UUID) (*dto.DashboardResponse, error) {
	summary, err := s.GetDashboardSummary(userID)
	if err != nil {
		return nil, err
	}

	ongoing, err := s.GetDashboardOngoingProjects(userID)
	if err != nil {
		return nil, err
	}

	recentActivity, _, err := s.GetDashboardNotifications(userID, 1, 10)
	if err != nil {
		return nil, err
	}

	return &dto.DashboardResponse{
		ActiveProjectsCount: summary.ActiveProjectsCount,
		UnreadMessagesCount: summary.UnreadMessagesCount,
		ActiveConnections:   summary.ActiveConnections,
		MyProjects:          ongoing.Projects,
		RecentActivity:      recentActivity,
	}, nil
}

func (s *dashboardService) GetDashboardSummary(userID uuid.UUID) (*dto.DashboardSummaryResponse, error) {
	activeProjectsCount, err := s.projRepo.CountByOwner(userID, "active")
	if err != nil {
		return nil, err
	}

	unreadMessagesCount, err := s.convRepo.CountTotalUnread(userID)
	if err != nil {
		return nil, err
	}

	activeConnections, err := s.connRepo.CountAccepted(userID)
	if err != nil {
		return nil, err
	}

	return &dto.DashboardSummaryResponse{
		ActiveProjectsCount: activeProjectsCount,
		UnreadMessagesCount: unreadMessagesCount,
		ActiveConnections:   activeConnections,
	}, nil
}

func (s *dashboardService) GetDashboardOngoingProjects(userID uuid.UUID) (*dto.DashboardOngoingProjectsResponse, error) {
	projects, err := s.projRepo.ListByUser(userID)
	if err != nil {
		return nil, err
	}

	items := make([]dto.DashboardProjectItem, 0, len(projects))
	for _, p := range projects {
		items = append(items, s.toDashboardProjectItem(p, userID))
	}

	return &dto.DashboardOngoingProjectsResponse{
		Projects: items,
		Total:    int64(len(items)),
	}, nil
}

func (s *dashboardService) GetDashboardNotifications(userID uuid.UUID, page, limit int) ([]dto.NotificationResponse, int64, error) {
	notifs, total, err := s.notifRepo.ListByUser(userID, page, limit)
	if err != nil {
		return nil, 0, err
	}

	recentActivity := make([]dto.NotificationResponse, len(notifs))
	for i, n := range notifs {
		recentActivity[i] = s.toNotificationResponse(n)
	}

	return recentActivity, total, nil
}

func (s *dashboardService) toDashboardProjectItem(p project.Project, userID uuid.UUID) dto.DashboardProjectItem {
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

	nextMs, err := s.projRepo.GetNextMilestone(p.ID)
	if err == nil && nextMs != nil {
		item.NextMilestone = &nextMs.Title
	}

	if p.Deadline != nil {
		dl := p.Deadline.Format(time.RFC3339)
		item.Deadline = &dl
	}

	return item
}

func (s *dashboardService) toNotificationResponse(n notification.Notification) dto.NotificationResponse {
	actorName := ""
	if n.Actor != nil {
		actorName = n.Actor.FullName
	}

	return dto.NotificationResponse{
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
