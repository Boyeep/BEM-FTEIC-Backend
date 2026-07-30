package handlers

import (
	"repo-backend/internal/dto"
	"repo-backend/internal/service"
	"repo-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

type Analytics struct{ Service *service.AnalyticsService }

func (h *Analytics) Track(c *gin.Context) {
	var input dto.TrackVisitor
	if !bind(c, &input) {
		return
	}
	if err := h.Service.Track(c, input); err != nil {
		response.Fail(c, err)
		return
	}
	response.NoContent(c)
}

func (h *Analytics) Count(c *gin.Context) {
	count, err := h.Service.Count(c)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, gin.H{"count": count})
}
