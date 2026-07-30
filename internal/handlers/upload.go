package handlers

import (
	"net/http"

	"repo-backend/internal/media"
	"repo-backend/pkg/apperr"
	"repo-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

type Media struct{ Service *media.Service }

func (h *Media) UploadImage(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, media.MaxImageSize+(1<<20))
	file, err := c.FormFile("file")
	if err != nil {
		response.Fail(c, apperr.Validation("file is required"))
		return
	}
	url, err := h.Service.Upload(c, c.GetString("user_id"), file, "https://"+c.Request.Host)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.Created(c, gin.H{"url": url})
}

func (h *Media) DeleteImage(c *gin.Context) {
	if err := h.Service.Delete(c, c.Query("url")); err != nil {
		response.Fail(c, err)
		return
	}
	response.NoContent(c)
}
