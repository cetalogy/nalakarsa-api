package middleware

import (
	"log/slog"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var logger *slog.Logger

func init() {
	logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
}

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Generate request/correlation ID
		requestID := uuid.New().String()
		c.Set("request_id", requestID)
		c.Writer.Header().Set("X-Request-ID", requestID)

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

		attrs := []slog.Attr{
			slog.String("request_id", requestID),
			slog.Int("status", statusCode),
			slog.String("method", method),
			slog.String("path", path),
			slog.String("ip", clientIP),
			slog.Duration("latency", latency),
		}

		if statusCode >= 400 {
			logger.LogAttrs(c.Request.Context(), slog.LevelError, "HTTP request failed", attrs...)
		} else {
			logger.LogAttrs(c.Request.Context(), slog.LevelInfo, "HTTP request success", attrs...)
		}
	}
}
