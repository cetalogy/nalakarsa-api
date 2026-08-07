package dashboardroutes

import (
	"nalakarsa/internal/handler/dashboard"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(protected *gin.RouterGroup, h *dashboardhandler.DashboardHandler) {
	protected.GET("/dashboard", h.GetDashboard)
	protected.GET("/dashboard/summary", h.GetDashboardSummary)
	protected.GET("/dashboard/ongoing-projects", h.GetDashboardOngoingProjects)
	protected.GET("/dashboard/notifications", h.GetDashboardNotifications)
}
