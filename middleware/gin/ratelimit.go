package ginmiddleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// ipLimiterStore holds one rate.Limiter per IP address.
type ipLimiterStore struct {
	mu      sync.Map
	rateVal rate.Limit
	burst   int
}

func newStore(max int, window time.Duration) *ipLimiterStore {
	return &ipLimiterStore{
		rateVal: rate.Every(window / time.Duration(max)),
		burst:   max,
	}
}

func (s *ipLimiterStore) get(ip string) *rate.Limiter {
	v, _ := s.mu.LoadOrStore(ip, rate.NewLimiter(s.rateVal, s.burst))
	return v.(*rate.Limiter)
}

func rateLimitMiddleware(store *ipLimiterStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !store.get(c.ClientIP()).Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code":    "RATE_LIMIT_EXCEEDED",
				"message": "Too many requests, please try again later",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// AuthRateLimit applies a strict limit on auth endpoints (login, create)
// to protect against brute-force attacks: 10 requests per IP per minute.
func AuthRateLimit() gin.HandlerFunc {
	return rateLimitMiddleware(newStore(10, time.Minute))
}

// APIRateLimit applies a general limit on authenticated API endpoints:
// 100 requests per IP per minute.
func APIRateLimit() gin.HandlerFunc {
	return rateLimitMiddleware(newStore(100, time.Minute))
}
