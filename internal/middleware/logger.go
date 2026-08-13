package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
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
		clientIP := c.ClientIP()

		if raw != "" {
			path = path + "?" + raw
		}

		transactionName := resolveTransactionName(method, path)
		latencyLabel := formatLatency(latency)
		level := resolveLogLevel(statusCode)
		methodLabel := colorMethod(method, isDevelopment)
		statusText := http.StatusText(statusCode)
		if statusText == "" {
			statusText = "Unknown Status"
		}
		statusLabel := colorStatus(statusCode, isDevelopment) + " " + statusText

		// Inject User information if authenticated
		userInfo := ""
		if email, exists := c.Get("email"); exists {
			userInfo = fmt.Sprintf(" | User: %v", email)
		} else if userID, exists := c.Get("user_id"); exists {
			userInfo = fmt.Sprintf(" | UserID: %v", userID)
		}

		message := fmt.Sprintf("%s %s %s -> %s (%s) | Client: %s%s",
			transactionName,
			methodLabel,
			path,
			statusLabel,
			latencyLabel,
			clientIP,
			userInfo,
		)
		logger.Log(c.Request.Context(), level, message)
	}
}

func resolveTransactionName(method, path string) string {
	method = strings.ToUpper(method)
	cleanPath := path
	if idx := strings.Index(cleanPath, "?"); idx != -1 {
		cleanPath = cleanPath[:idx]
	}

	switch {
	// Auth Module
	case method == "POST" && strings.HasPrefix(cleanPath, "/api/v1/auth/login"):
		return "🔐 [USER LOGIN]"
	case method == "POST" && strings.HasPrefix(cleanPath, "/api/v1/auth/register"):
		return "👤 [USER REGISTER]"
	case method == "POST" && strings.HasPrefix(cleanPath, "/api/v1/auth/refresh"):
		return "🔄 [REFRESH TOKEN]"
	case method == "GET" && strings.HasPrefix(cleanPath, "/api/v1/auth/me"):
		return "👤 [GET CURRENT USER]"

	// Discussion Module
	case method == "POST" && strings.Contains(cleanPath, "/collaboration"):
		return "🤝 [MARK COLLABORATION]"
	case method == "POST" && strings.Contains(cleanPath, "/votes"):
		return "👍 [UPVOTE DISCUSSION]"
	case method == "DELETE" && strings.Contains(cleanPath, "/votes"):
		return "👎 [UNVOTE DISCUSSION]"
	case method == "POST" && strings.Contains(cleanPath, "/replies"):
		return "💬 [ADD REPLY]"
	case method == "DELETE" && strings.Contains(cleanPath, "/replies"):
		return "🗑️ [DELETE REPLY]"
	case method == "GET" && strings.Contains(cleanPath, "/replies"):
		return "💬 [GET REPLIES]"
	case method == "POST" && strings.HasPrefix(cleanPath, "/api/v1/discussions"):
		return "📝 [CREATE DISCUSSION]"
	case method == "PATCH" && strings.HasPrefix(cleanPath, "/api/v1/discussions"):
		return "✏️ [UPDATE DISCUSSION]"
	case method == "DELETE" && strings.HasPrefix(cleanPath, "/api/v1/discussions"):
		return "🗑️ [DELETE DISCUSSION]"
	case method == "GET" && cleanPath == "/api/v1/discussions":
		return "💬 [LIST DISCUSSIONS]"
	case method == "GET" && strings.HasPrefix(cleanPath, "/api/v1/discussions/"):
		return "🔍 [GET DISCUSSION DETAIL]"

	// Project Module
	case method == "GET" && strings.HasPrefix(cleanPath, "/api/v1/projects"):
		return "🚀 [PROJECT SERVICE]"
	case method == "POST" && strings.HasPrefix(cleanPath, "/api/v1/projects"):
		return "🚀 [CREATE PROJECT]"

	// User / Connection / Conversation / Notification / Homepage
	case strings.HasPrefix(cleanPath, "/api/v1/users"):
		return "👤 [USER PROFILE]"
	case strings.HasPrefix(cleanPath, "/api/v1/connections"):
		return "👥 [CONNECTION SERVICE]"
	case strings.HasPrefix(cleanPath, "/api/v1/conversations"):
		return "💬 [CHAT CONVERSATION]"
	case strings.HasPrefix(cleanPath, "/api/v1/notifications"):
		return "🔔 [NOTIFICATION SERVICE]"
	case strings.HasPrefix(cleanPath, "/api/v1/homepage"):
		return "🏠 [HOMEPAGE LANDING]"

	// Health check
	case cleanPath == "/" || cleanPath == "/health":
		return "🟢 [HEALTH CHECK]"

	default:
		return fmt.Sprintf("🌐 [%s]", method)
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
