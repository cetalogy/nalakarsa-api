package dto

type DashboardResponse struct {
	ActiveProjectsCount int64                   `json:"active_projects_count"`
	UnreadMessagesCount int64                   `json:"unread_messages_count"`
	ActiveConnections   int64                   `json:"active_connections"`
	MyProjects          []DashboardProjectItem  `json:"my_projects"`
	RecentActivity      []NotificationResponse  `json:"recent_activity"`
}

type DashboardProjectItem struct {
	ID            string  `json:"id"`
	Title         string  `json:"title"`
	Role          string  `json:"role"`
	Progress      int     `json:"progress"`
	Status        string  `json:"status"`
	NextMilestone *string `json:"next_milestone,omitempty"`
	Deadline      *string `json:"deadline,omitempty"`
}
