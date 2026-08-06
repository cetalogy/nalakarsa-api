package dto

type DashboardResponse struct {
	ActiveProjectsCount int64                  `json:"active_projects_count" form:"active_projects_count"`
	UnreadMessagesCount int64                  `json:"unread_messages_count" form:"unread_messages_count"`
	ActiveConnections   int64                  `json:"active_connections" form:"active_connections"`
	MyProjects          []DashboardProjectItem `json:"my_projects" form:"my_projects"`
	RecentActivity      []NotificationResponse `json:"recent_activity" form:"recent_activity"`
}

type DashboardProjectItem struct {
	ID            string  `json:"id" form:"id"`
	Title         string  `json:"title" form:"title"`
	Role          string  `json:"role" form:"role"`
	Progress      int     `json:"progress" form:"progress"`
	Status        string  `json:"status" form:"status"`
	NextMilestone *string `json:"next_milestone,omitempty" form:"next_milestone"`
	Deadline      *string `json:"deadline,omitempty" form:"deadline"`
}
