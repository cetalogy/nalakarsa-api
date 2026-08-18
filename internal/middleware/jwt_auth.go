package middleware

import (
	"log"
	"net/http"
	"strings"

	"nalakarsa/internal/config"
	"nalakarsa/internal/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Printf("[AUTH FAILED] %s %s from %s - Reason: Missing Authorization header", c.Request.Method, c.Request.URL.Path, c.ClientIP())
			utils.ErrorJSONResponseWithMessage(c, http.StatusUnauthorized, "Authorization header is required")
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			log.Printf("[AUTH FAILED] %s %s from %s - Reason: Invalid Bearer format (Header was: '%s')", c.Request.Method, c.Request.URL.Path, c.ClientIP(), authHeader)
			utils.ErrorJSONResponseWithMessage(c, http.StatusUnauthorized, "Authorization header format must be Bearer <token>")
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := utils.ValidateToken(tokenString, cfg.JWTSecret)
		if err != nil {
			log.Printf("[AUTH FAILED] %s %s from %s - Reason: Token validation failed: %v", c.Request.Method, c.Request.URL.Path, c.ClientIP(), err)
			utils.ErrorJSONResponseWithMessage(c, http.StatusUnauthorized, "Invalid or expired access token")
			c.Abort()
			return
		}

		// Inject user info into context
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)

		c.Next()
	})
}

