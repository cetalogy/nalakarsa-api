package discussionroutes

import (
	"nalakarsa/internal/handler/discussion"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(public *gin.RouterGroup, protected *gin.RouterGroup, h *discussionhandler.DiscussionHandler) {
	discussions := public.Group("/discussions")
	discussions.GET("", h.List)
	discussions.GET("/:id", h.GetByID)
	discussions.GET("/:id/replies", h.GetReplies)

	protected.POST("/discussions", h.Create)
	protected.PATCH("/discussions/:id", h.Update)
	protected.DELETE("/discussions/:id", h.Delete)
	protected.POST("/discussions/:id/collaboration", h.MarkCollaboration)
	protected.POST("/discussions/:id/replies", h.AddReply)
	protected.DELETE("/discussions/:id/replies/:reply_id", h.DeleteReply)
	protected.POST("/discussions/:id/votes", h.Vote)
	protected.DELETE("/discussions/:id/votes", h.Unvote)
}
