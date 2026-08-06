package authroutes

import (
	"nalakarsa/internal/handler/auth"
	"nalakarsa/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(public *gin.RouterGroup, protected *gin.RouterGroup, h *authhandler.AuthHandler) {
	auth := public.Group("/auth")
	auth.Use(middleware.RateLimiter(5, 10))
	auth.POST("/register", h.Register)
	auth.POST("/login", h.Login)
	auth.POST("/refresh", h.RefreshToken)
	auth.POST("/forgot-password", h.ForgotPassword)
	auth.POST("/reset-password", h.ResetPassword)

	protected.POST("/auth/logout", h.Logout)
}
