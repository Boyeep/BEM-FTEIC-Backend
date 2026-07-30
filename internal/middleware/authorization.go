package middleware

import (
	"net/http"

	"repo-backend/pkg/apperr"
	"repo-backend/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RequireRole(db *gorm.DB, role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var current string
		err := db.WithContext(c).Table("profiles").Select("role").
			Where("id = ?", c.GetString("user_id")).Scan(&current).Error
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
