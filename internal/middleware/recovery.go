package middleware

import (
	"fmt"
	"log"

	"repo-backend/pkg/apperr"
	"repo-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

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
