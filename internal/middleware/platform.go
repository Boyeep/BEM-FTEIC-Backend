package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"repo-backend/pkg/apperr"
	"repo-backend/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader("X-Request-ID")
		if id == "" {
			b := make([]byte, 16)
			if _, err := rand.Read(b); err == nil {
				id = hex.EncodeToString(b)
			} else {
				id = time.Now().UTC().Format("20060102150405.000000000")
			}
		}
		c.Set("request_id", id)
		c.Header("X-Request-ID", id)
		c.Next()
	}
}

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if recovered := recover(); recovered != nil {
				log.Printf("panic request_id=%s: %v", c.GetString("request_id"), recovered)
				response.Fail(c, apperr.Internal(fmt.Errorf("panic: %v", recovered)))
				c.Abort()
			}
		}()
		c.Next()
	}
}

type visitor struct {
	count int
	reset time.Time
}

func RateLimit(limit int, window time.Duration) gin.HandlerFunc {
	var mu sync.Mutex
	visitors := map[string]*visitor{}
	return func(c *gin.Context) {
		now, key := time.Now(), c.ClientIP()
		mu.Lock()
		v := visitors[key]
		if v == nil || now.After(v.reset) {
			v = &visitor{reset: now.Add(window)}
			visitors[key] = v
		}
		v.count++
		blocked := v.count > limit
		retry := time.Until(v.reset).Seconds()
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

func RequireRole(db *gorm.DB, role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var current string
		err := db.WithContext(c).Table("profiles").Select("role").Where("id = ?", c.GetString("user_id")).Scan(&current).Error
		if err != nil {
			response.Fail(c, apperr.Internal(err))
			c.Abort()
			return
		}
		if current != role {
			response.Fail(c, apperr.New("FORBIDDEN", "admin role required", http.StatusForbidden))
			c.Abort()
			return
		}
		c.Set("user_role", current)
		c.Next()
	}
}
