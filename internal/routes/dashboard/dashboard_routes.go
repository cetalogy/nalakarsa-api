package dashboardroutes

import (
	"nalakarsa/internal/handler/dashboard"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(protected *gin.RouterGroup, h *dashboardhandler.DashboardHandler) {
	protected.GET("/dashboard", h.GetDashboard)
}
