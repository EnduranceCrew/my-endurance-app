package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type ipBucket struct {
	tokens     float64
	lastRefill time.Time
	mu         sync.Mutex
}

var buckets sync.Map

func getOrCreateBucket(ip string) *ipBucket {
	v, _ := buckets.LoadOrStore(ip, &ipBucket{tokens: 30, lastRefill: time.Now()})
	return v.(*ipBucket)
}

// RateLimit applies per-IP token bucket rate limiting.
// rps: requests per second allowed; burstSize: max burst.
func RateLimit(rps, burstSize float64) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		b := getOrCreateBucket(ip)
		b.mu.Lock()
		now := time.Now()
		elapsed := now.Sub(b.lastRefill).Seconds()
		b.tokens += elapsed * rps
		if b.tokens > burstSize {
			b.tokens = burstSize
		}
		b.lastRefill = now
		if b.tokens < 1 {
			b.mu.Unlock()
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"error":   "muitas requisições, tente novamente em breve",
			})
			return
		}
		b.tokens--
		b.mu.Unlock()
		c.Next()
	}
}
