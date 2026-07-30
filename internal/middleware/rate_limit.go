package middleware

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"repo-backend/pkg/apperr"
	"repo-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

type rateLimitVisitor struct {
	count int
	reset time.Time
}

func RateLimit(limit int, window time.Duration) gin.HandlerFunc {
	var mu sync.Mutex
	visitors := map[string]*rateLimitVisitor{}
	return func(c *gin.Context) {
		now, key := time.Now(), c.ClientIP()
		mu.Lock()
		if len(visitors) > 10000 {
			for ip, entry := range visitors {
				if now.After(entry.reset) {
					delete(visitors, ip)
				}
			}
		}
		visitor := visitors[key]
		if visitor == nil || now.After(visitor.reset) {
			visitor = &rateLimitVisitor{reset: now.Add(window)}
			visitors[key] = visitor
		}
		visitor.count++
		blocked := visitor.count > limit
		retry := time.Until(visitor.reset).Seconds()
		mu.Unlock()
		if blocked {
			c.Header("Retry-After", strconv.Itoa(int(retry)+1))
			response.Fail(c, apperr.New("RATE_LIMITED", "too many requests", http.StatusTooManyRequests))
			c.Abort()
			return
		}
		c.Next()
	}
}
