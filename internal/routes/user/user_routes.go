package userroutes

import (
	"nalakarsa/internal/handler/user"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(public *gin.RouterGroup, protected *gin.RouterGroup, h *userhandler.UserHandler, suggestions gin.HandlerFunc) {
	users := public.Group("/users")
	users.GET("", h.ListUsers)
	users.GET("/:id", h.GetPublicProfile)

	protected.GET("/users/me", h.GetProfile)
	protected.PATCH("/users/me", h.UpdateProfile)
	protected.GET("/users/me/projects", h.GetMyProjects)
	protected.GET("/users/me/stats", h.GetMyStats)
	protected.GET("/auth/me", h.GetProfile)
	protected.PATCH("/auth/me", h.UpdateProfile)
	protected.POST("/users/me/avatar", h.UploadAvatar)
	protected.GET("/users/suggestions", suggestions)
}
