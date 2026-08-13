package routes

import (
	"net/http"

	"nalakarsa/internal/config"
	"nalakarsa/internal/middleware"
	"nalakarsa/internal/utils"

	"github.com/gin-gonic/gin"
)

// NewRouter creates the application router and installs global middleware.
// Feature routes are registered by their own module packages from main.
func NewRouter(cfg *config.Config) *gin.Engine {
	r := gin.New()
	r.Use(middleware.LoggerMiddleware(cfg))
	r.Use(gin.Recovery())
	r.Use(middleware.CORSMiddleware(cfg))

	r.GET("/health", func(c *gin.Context) {
		utils.JSONResponse(c, http.StatusOK, "Server is healthy", gin.H{
			"status":  "healthy",
			"service": "Nalakarsa Backend API Server",
			"version": "1.0.0",
		}, nil)
	})

	return r
}
