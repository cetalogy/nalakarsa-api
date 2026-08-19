package middleware

import (
	"net/http"
	"strings"

	"nalakarsa/internal/config"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		allowedOrigin := ""

		if cfg.CORSAllowedOrigins == "*" {
			if origin != "" {
				allowedOrigin = origin
			} else {
				allowedOrigin = "*"
			}
		} else {
			allowed := strings.Split(cfg.CORSAllowedOrigins, ",")
			for _, o := range allowed {
				if strings.TrimSpace(o) == origin {
					allowedOrigin = origin
					break
				}
			}

			if allowedOrigin == "" && strings.HasSuffix(origin, ".vercel.app") {
				allowedOrigin = origin
			}
		}

		if allowedOrigin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		}

		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
