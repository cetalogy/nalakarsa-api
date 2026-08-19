package userroutes

import (
	userhandler "nalakarsa/internal/handler/user"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(public *gin.RouterGroup, protected *gin.RouterGroup, h *userhandler.UserHandler, suggestions gin.HandlerFunc) {
	users := public.Group("/users")
	users.GET("", h.ListUsers)
	users.GET("/search", h.ListUsers)
	users.GET("/:id", h.GetPublicProfile)
	users.GET("/:id/followers", h.GetFollowers)
	users.GET("/:id/following", h.GetFollowing)

	protected.GET("/users/me", h.GetProfile)
	protected.GET("/auth/me", h.GetProfile)
	protected.PATCH("/users/me", h.UpdateProfile)
	protected.PATCH("/auth/me", h.UpdateProfile)
	protected.GET("/users/me/projects", h.GetMyProjects)
	protected.GET("/users/me/stats", h.GetMyStats)
	protected.GET("/users/me/followers", h.GetFollowers)
	protected.GET("/users/me/following", h.GetFollowing)
	protected.POST("/users/me/avatar", h.UploadAvatar)
	protected.POST("/users/:id/follow", h.ToggleFollow)
	protected.GET("/users/suggestions", suggestions)
}
