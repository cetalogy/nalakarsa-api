package middleware

import (
	"log/slog"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

var logger *slog.Logger

func init() {
	logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
}

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method

		if raw != "" {
			path = path + "?" + raw
		}

		if statusCode >= 400 {
			logger.Error("HTTP request failed",
				slog.Int("status", statusCode),
				slog.String("method", method),
				slog.String("path", path),
				slog.String("ip", clientIP),
				slog.Duration("latency", latency),
			)
		} else {
			logger.Info("HTTP request success",
				slog.Int("status", statusCode),
				slog.String("method", method),
				slog.String("path", path),
				slog.String("ip", clientIP),
				slog.Duration("latency", latency),
			)
		}
	}
}
