package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func RateLimiter(r rate.Limit, b int) gin.HandlerFunc {
	type clientLimiter struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	var mu sync.Mutex
	limiters := map[string]*clientLimiter{}

	getLimiter := func(ip string) *rate.Limiter {
		mu.Lock()
		defer mu.Unlock()

		if cl, ok := limiters[ip]; ok {
			cl.lastSeen = time.Now()
			return cl.limiter
		}

		l := rate.NewLimiter(r, b)
		limiters[ip] = &clientLimiter{limiter: l, lastSeen: time.Now()}
		return l
	}

	// lightweight cleanup goroutine to avoid map growing forever
	go func() {
		ticker := time.NewTicker(10 * time.Minute)
		defer ticker.Stop()
		for range ticker.C {
			cutoff := time.Now().Add(-30 * time.Minute)
			mu.Lock()
			for ip, cl := range limiters {
				if cl.lastSeen.Before(cutoff) {
					delete(limiters, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return func(c *gin.Context) {
		ip := c.ClientIP()
		limiter := getLimiter(ip)
		if !limiter.AllowN(time.Now(), 1) {
			c.Header("Retry-After", "60")
			c.String(http.StatusTooManyRequests, "Rate limit exceeded")
			c.Abort()
			return
		}
		c.Next()
	}
}
