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
	protected.POST("/conversations/:id/attachments", h.UploadAttachment)
	protected.PATCH("/conversations/:id/read", h.MarkRead)

	// Aliases for /chat (FE requirement)
	protected.GET("/chat/contacts", h.ListConversations)
	protected.POST("/chat/start", h.StartChat)
	protected.GET("/chat/contacts/:id/messages", h.ListMessages)
	protected.POST("/chat/contacts/:id/messages", h.SendMessage)
	protected.POST("/chat/contacts/:id/attachments", h.UploadAttachment)

	// Group Chats (FE Contract Specification)
	protected.GET("/chats/groups", h.ListGroupChats)
	protected.GET("/chats/groups/:groupId/messages", h.ListGroupMessages)
	protected.POST("/chats/groups/:groupId/messages", h.SendGroupMessage)
	protected.POST("/chats/groups/:groupId/attachments", h.UploadGroupAttachment)

	// Delete Message (Direct & Group)
	protected.DELETE("/chats/messages/:id", h.DeleteMessage)
	protected.DELETE("/conversations/messages/:id", h.DeleteMessage)
	protected.DELETE("/chat/messages/:id", h.DeleteMessage)
}
