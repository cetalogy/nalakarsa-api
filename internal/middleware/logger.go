package middleware

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"nalakarsa/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/lmittmann/tint"
)

var logger *slog.Logger

func LoggerMiddleware(cfg *config.Config) gin.HandlerFunc {
	logger = buildHTTPLogger(cfg)
	isDevelopment := cfg != nil && cfg.Env == "development"

	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		method := c.Request.Method
		statusCode := c.Writer.Status()

		if raw != "" {
			path = path + "?" + raw
		}

		latencyLabel := formatLatency(latency)
		level := resolveLogLevel(statusCode)
		methodLabel := colorMethod(method, isDevelopment)
		statusLabel := colorStatus(statusCode, isDevelopment)

		message := fmt.Sprintf("%s %s %s (%s)", methodLabel, path, statusLabel, latencyLabel)
		logger.Log(c.Request.Context(), level, message)
	}
}

func buildHTTPLogger(cfg *config.Config) *slog.Logger {
	if cfg != nil && cfg.Env == "production" {
		return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
	}

	return slog.New(tint.NewTextHandler(os.Stdout, &tint.Options{
		TimeFormat: "15:04:05",
		NoColor:    false,
	}))
}

func formatLatency(latency time.Duration) string {
	if latency >= time.Millisecond {
		return fmt.Sprintf("%.2fms", float64(latency)/float64(time.Millisecond))
	}

	if latency >= time.Microsecond {
		return fmt.Sprintf("%.2fµs", float64(latency)/float64(time.Microsecond))
	}

	return fmt.Sprintf("%dns", latency.Nanoseconds())
}

func resolveLogLevel(statusCode int) slog.Level {
	if statusCode >= 500 {
		return slog.LevelError
	}
	if statusCode >= 400 {
		return slog.LevelWarn
	}
	return slog.LevelInfo
}

func colorMethod(method string, isDevelopment bool) string {
	method = strings.ToUpper(method)
	if !isDevelopment {
		return method
	}

	switch method {
	case "GET":
		return colorText("GET", "32")
	case "POST":
		return colorText("POST", "33")
	case "PUT":
		return colorText("PUT", "34")
	case "DELETE":
		return colorText("DELETE", "31")
	case "PATCH":
		return colorText("PATCH", "35")
	default:
		return method
	}
}

func colorStatus(statusCode int, isDevelopment bool) string {
	if !isDevelopment {
		return fmt.Sprintf("%d", statusCode)
	}

	switch {
	case statusCode >= 200 && statusCode < 300:
		return colorText(fmt.Sprintf("%d", statusCode), "32")
	case statusCode >= 300 && statusCode < 400:
		return colorText(fmt.Sprintf("%d", statusCode), "36")
	case statusCode >= 400 && statusCode < 500:
		return colorText(fmt.Sprintf("%d", statusCode), "33")
	case statusCode >= 500:
		return colorText(fmt.Sprintf("%d", statusCode), "1;31")
	default:
		return fmt.Sprintf("%d", statusCode)
	}
}

func colorText(value, ansiColor string) string {
	return "\x1b[" + ansiColor + "m" + value + "\x1b[0m"
}
