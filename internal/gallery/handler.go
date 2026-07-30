package gallery

import (
	"repo-backend/pkg/apperr"
	"repo-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

type Handler struct{ service *Service }

func NewHandler(service *Service) *Handler { return &Handler{service: service} }
func (h *Handler) List(c *gin.Context) {
	var query ListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Fail(c, apperr.Validation(err.Error()))
		return
	}
	query.Defaults(12)
	value, err := h.service.List(c, query)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, value)
}
func (h *Handler) Get(c *gin.Context) {
	value, err := h.service.Get(c, c.Param("id"))
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, value)
}
func (h *Handler) Create(c *gin.Context) {
	var input CreateDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Fail(c, apperr.Validation(err.Error()))
		return
	}
	value, err := h.service.Create(c, input, c.GetString("user_id"))
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.Created(c, value)
}
func (h *Handler) Update(c *gin.Context) {
	var input UpdateDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Fail(c, apperr.Validation(err.Error()))
		return
	}
	value, err := h.service.Update(c, c.Param("id"), input)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, value)
}
func (h *Handler) Delete(c *gin.Context) {
	if err := h.service.Delete(c, c.Param("id")); err != nil {
		response.Fail(c, err)
		return
	}
	response.NoContent(c)
}
