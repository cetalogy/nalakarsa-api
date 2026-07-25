package middleware

import (
	"net/http"
	"sync"
	"time"

	"nalakarsa/internal/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// ipLimiter stores a rate limiter per IP address
type ipLimiter struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

var (
	visitors = make(map[string]*ipLimiter)
	mu       sync.Mutex
)

func init() {
	// Cleanup stale entries every 3 minutes
	go func() {
		for {
			time.Sleep(3 * time.Minute)
			mu.Lock()
			for ip, v := range visitors {
				if time.Since(v.lastSeen) > 5*time.Minute {
					delete(visitors, ip)
				}
			}
			mu.Unlock()
		}
	}()
}

func getVisitor(ip string, r rate.Limit, b int) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	v, exists := visitors[ip]
	if !exists {
		limiter := rate.NewLimiter(r, b)
		visitors[ip] = &ipLimiter{limiter: limiter, lastSeen: time.Now()}
		return limiter
	}

	v.lastSeen = time.Now()
	return v.limiter
}

// RateLimiter creates a rate limiting middleware.
// rps: requests per second allowed
// burst: maximum burst size
func RateLimiter(rps float64, burst int) gin.HandlerFunc {
	return func(c *gin.Context) {
		limiter := getVisitor(c.ClientIP(), rate.Limit(rps), burst)
		if !limiter.Allow() {
			utils.ErrorJSONResponseWithMessage(c, http.StatusTooManyRequests, "Too many requests. Please try again later.")
			c.Abort()
			return
		}
		c.Next()
	}
}
