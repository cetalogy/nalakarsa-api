package connectionroutes

import (
	"nalakarsa/internal/handler/connection"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(protected *gin.RouterGroup, h *connectionhandler.ConnectionHandler) {
	protected.GET("/connections", h.ListConnections)
	protected.GET("/partners", h.ListConnections)
	protected.GET("/users/me/partners", h.ListConnections)
	protected.GET("/connections/requests", h.ListRequests)
	protected.POST("/connections/requests", h.SendRequest)
	protected.PATCH("/connections/requests/:id/accept", h.AcceptRequest)
	protected.PATCH("/connections/requests/:id/reject", h.RejectRequest)
	protected.DELETE("/connections/requests/:id", h.CancelRequest)
	protected.DELETE("/connections/:userId", h.RemoveConnection)
}
