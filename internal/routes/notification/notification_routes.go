package notificationroutes

import (
	"nalakarsa/internal/handler/notification"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(protected *gin.RouterGroup, h *notificationhandler.NotificationHandler) {
	protected.GET("/notifications", h.List)
	protected.PATCH("/notifications/:id/read", h.MarkRead)
	protected.PATCH("/notifications/read-all", h.MarkAllRead)
}
