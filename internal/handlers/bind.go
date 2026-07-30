package handlers

import (
	"repo-backend/pkg/apperr"
	"repo-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

func bind(c *gin.Context, target any) bool {
	if err := c.ShouldBindJSON(target); err != nil {
		response.Fail(c, apperr.Validation(err.Error()))
		return false
	}
	return true
}
