package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
)

func RequestID() gin.HandlerFunc {
	validRequestID := regexp.MustCompile(`^[A-Za-z0-9_-]{8,64}$`)
	return func(c *gin.Context) {
		id := c.GetHeader("X-Request-ID")
		if !validRequestID.MatchString(id) {
			bytes := make([]byte, 16)
			if _, err := rand.Read(bytes); err == nil {
				id = hex.EncodeToString(bytes)
			} else {
				id = time.Now().UTC().Format("20060102150405.000000000")
			}
		}
		c.Set("request_id", id)
		c.Header("X-Request-ID", id)
		c.Next()
	}
}
