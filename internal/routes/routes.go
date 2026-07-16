package routes

import (
	"nalakarsa/internal/config"
	"nalakarsa/internal/handler"
	"nalakarsa/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(
	cfg *config.Config,
	authH *handler.AuthHandler,
	userH *handler.UserHandler,
	discH *handler.DiscussionHandler,
	collabH *handler.CollaborationHandler,
) *gin.Engine {
	r := gin.New()

	// Global Middlewares
	r.Use(middleware.LoggerMiddleware())
	r.Use(gin.Recovery())
	r.Use(middleware.CORSMiddleware())

	// API Group
	v1 := r.Group("/api/v1")
	{
		// Public Auth
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authH.Register)
			auth.POST("/login", authH.Login)
			auth.POST("/refresh", authH.RefreshToken)
		}

		// Public Users Directory
		users := v1.Group("/users")
		{
			users.GET("", userH.ListUsers)
			users.GET("/:id", userH.GetPublicProfile)
		}

		// Public Discussions
		discussions := v1.Group("/discussions")
		{
			discussions.GET("", discH.List)
			discussions.GET("/:id", discH.GetByID)
		}

		// Public Collaborations
		collaborations := v1.Group("/collaborations")
		{
			collaborations.GET("", collabH.List)
			collaborations.GET("/:id", collabH.GetByID)
		}

		// Protected Routes Group
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(cfg))
		{
			// Auth Protected
			protected.POST("/auth/logout", authH.Logout)

			// User Protected
			protected.GET("/users/me", userH.GetProfile)
			protected.PUT("/users/me", userH.UpdateProfile)
			protected.POST("/users/me/avatar", userH.UploadAvatar)

			// Discussion Protected
			protected.POST("/discussions", discH.Create)
			protected.PUT("/discussions/:id", discH.Update)
			protected.DELETE("/discussions/:id", discH.Delete)

			// Comments Protected
			protected.POST("/discussions/:id/comments", discH.AddComment)
			protected.DELETE("/discussions/:id/comments/:comment_id", discH.DeleteComment)

			// Collaboration Protected
			protected.POST("/collaborations", collabH.Create)
			protected.PUT("/collaborations/:id", collabH.Update)
			protected.DELETE("/collaborations/:id", collabH.Delete)
			protected.POST("/collaborations/:id/apply", collabH.Apply)
			protected.GET("/collaborations/:id/applications", collabH.ListApplications)
			protected.PUT("/collaborations/:id/applications/:app_id", collabH.UpdateApplicationStatus)
		}
	}

	return r
}
