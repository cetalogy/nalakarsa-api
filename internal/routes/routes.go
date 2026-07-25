package routes

import (
	"net/http"

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
	projH *handler.ProjectHandler,
	connH *handler.ConnectionHandler,
	convH *handler.ConversationHandler,
	notifH *handler.NotificationHandler,
	dashH *handler.DashboardHandler,
) *gin.Engine {
	r := gin.New()

	// Global Middlewares
	r.Use(middleware.LoggerMiddleware())
	r.Use(gin.Recovery())
	r.Use(middleware.CORSMiddleware(cfg))

	// Health Check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"service": "Nalakarsa Backend API Server",
			"version": "1.0.0",
		})
	})

	// Rate limiter for auth endpoints (5 requests/second, burst 10)
	authRateLimiter := middleware.RateLimiter(5, 10)

	// API Group
	v1 := r.Group("/api/v1")
	{
		// ============================================================
		// Public Auth
		// ============================================================
		auth := v1.Group("/auth")
		auth.Use(authRateLimiter)
		{
			auth.POST("/register", authH.Register)
			auth.POST("/login", authH.Login)
			auth.POST("/refresh", authH.RefreshToken)
			auth.POST("/forgot-password", authH.ForgotPassword)
			auth.POST("/reset-password", authH.ResetPassword)
		}

		// ============================================================
		// Public Users Directory
		// ============================================================
		users := v1.Group("/users")
		{
			users.GET("", userH.ListUsers)
			users.GET("/:id", userH.GetPublicProfile)
		}

		// ============================================================
		// Public Discussions
		// ============================================================
		discussions := v1.Group("/discussions")
		{
			discussions.GET("", discH.List)
			discussions.GET("/:id", discH.GetByID)
			discussions.GET("/:id/replies", discH.GetByID) // replies are included in detail
		}

		// ============================================================
		// Public Projects
		// ============================================================
		projects := v1.Group("/projects")
		{
			projects.GET("", projH.List)
			projects.GET("/:id", projH.GetByID)
			projects.GET("/:id/members", projH.ListMembers)
		}

		// ============================================================
		// Protected Routes Group
		// ============================================================
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(cfg))
		{
			// --- Auth Protected ---
			protected.POST("/auth/logout", authH.Logout)

			// --- User Protected ---
			protected.GET("/users/me", userH.GetProfile)
			protected.PATCH("/users/me", userH.UpdateProfile)
			protected.POST("/users/me/avatar", userH.UploadAvatar)
			protected.GET("/users/suggestions", connH.GetSuggestions)

			// --- Discussion Protected ---
			protected.POST("/discussions", discH.Create)
			protected.PATCH("/discussions/:id", discH.Update)
			protected.DELETE("/discussions/:id", discH.Delete)

			// Discussion Replies
			protected.POST("/discussions/:id/replies", discH.AddReply)
			protected.DELETE("/discussions/:id/replies/:reply_id", discH.DeleteReply)

			// Discussion Votes
			protected.POST("/discussions/:id/votes", discH.Vote)
			protected.DELETE("/discussions/:id/votes", discH.Unvote)

			// --- Connections ---
			protected.GET("/connections", connH.ListConnections)
			protected.GET("/connections/requests", connH.ListRequests)
			protected.POST("/connections/requests", connH.SendRequest)
			protected.PATCH("/connections/requests/:id/accept", connH.AcceptRequest)
			protected.PATCH("/connections/requests/:id/reject", connH.RejectRequest)
			protected.DELETE("/connections/requests/:id", connH.CancelRequest)
			protected.DELETE("/connections/:userId", connH.RemoveConnection)

			// --- Projects Protected ---
			protected.POST("/projects", projH.Create)
			protected.PATCH("/projects/:id", projH.Update)
			protected.DELETE("/projects/:id", projH.Delete)

			// Project Applications
			protected.POST("/projects/:id/applications", projH.Apply)
			protected.GET("/projects/:id/applications", projH.ListApplications)
			protected.PATCH("/projects/:id/applications/:applicationId", projH.UpdateApplicationStatus)

			// Project Milestones
			protected.POST("/projects/:id/milestones", projH.CreateMilestone)
			protected.PATCH("/projects/:id/milestones/:milestoneId", projH.UpdateMilestone)

			// --- Conversations ---
			protected.GET("/conversations", convH.ListConversations)
			protected.POST("/conversations/direct", convH.GetOrCreateDirect)
			protected.GET("/conversations/:id/messages", convH.ListMessages)
			protected.POST("/conversations/:id/messages", convH.SendMessage)
			protected.PATCH("/conversations/:id/read", convH.MarkRead)

			// --- Notifications ---
			protected.GET("/notifications", notifH.List)
			protected.PATCH("/notifications/:id/read", notifH.MarkRead)
			protected.PATCH("/notifications/read-all", notifH.MarkAllRead)

			// --- Dashboard ---
			protected.GET("/dashboard", dashH.GetDashboard)
		}
	}

	return r
}
