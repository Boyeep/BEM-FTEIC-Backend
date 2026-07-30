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

type RateLimitStore interface {
	Increment(key string, now time.Time, window time.Duration) (count int, reset time.Time, err error)
}

type MemoryRateLimitStore struct {
	mu       sync.Mutex
	visitors map[string]*rateLimitVisitor
}

func NewMemoryRateLimitStore() *MemoryRateLimitStore {
	return &MemoryRateLimitStore{visitors: make(map[string]*rateLimitVisitor)}
}

func (s *MemoryRateLimitStore) Increment(key string, now time.Time, window time.Duration) (int, time.Time, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.visitors) > 10000 {
		for ip, entry := range s.visitors {
			if now.After(entry.reset) {
				delete(s.visitors, ip)
			}
		}
	}
	visitor := s.visitors[key]
	if visitor == nil || now.After(visitor.reset) {
		visitor = &rateLimitVisitor{reset: now.Add(window)}
		s.visitors[key] = visitor
	}
	visitor.count++
	return visitor.count, visitor.reset, nil
}

func RateLimit(limit int, window time.Duration) gin.HandlerFunc {
	return RateLimitWithStore(NewMemoryRateLimitStore(), limit, window)
}

func RateLimitWithStore(store RateLimitStore, limit int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		now, key := time.Now(), c.ClientIP()
		count, reset, err := store.Increment(key, now, window)
		if err != nil {
			response.Fail(c, apperr.Internal(err))
			c.Abort()
			return
		}
		if count > limit {
			retry := time.Until(reset).Seconds()
			c.Header("Retry-After", strconv.Itoa(int(retry)+1))
			response.Fail(c, apperr.New("RATE_LIMITED", "too many requests", http.StatusTooManyRequests))
			c.Abort()
			return
		}
		c.Next()
	}
}
