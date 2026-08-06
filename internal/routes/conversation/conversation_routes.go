package conversationroutes

import (
	"nalakarsa/internal/handler/conversation"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(protected *gin.RouterGroup, h *conversationhandler.ConversationHandler) {
	protected.GET("/conversations", h.ListConversations)
	protected.POST("/conversations/direct", h.GetOrCreateDirect)
	protected.GET("/conversations/:id/messages", h.ListMessages)
	protected.POST("/conversations/:id/messages", h.SendMessage)
	protected.PATCH("/conversations/:id/read", h.MarkRead)
}
