package event

import (
	"repo-backend/pkg/apperr"
	"repo-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

type Handler struct{ service *Service }

func NewHandler(service *Service) *Handler   { return &Handler{service: service} }
func (h *Handler) ListPublic(c *gin.Context) { h.list(c, 8, true) }
func (h *Handler) ListAdmin(c *gin.Context)  { h.list(c, 20, false) }
func (h *Handler) list(c *gin.Context, size int, published bool) {
	var query ListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Fail(c, apperr.Validation(err.Error()))
		return
	}
	query.Defaults(size)
	value, err := h.service.List(c, query, published)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, value)
}
func (h *Handler) GetPublic(c *gin.Context) { h.get(c, true) }
func (h *Handler) GetAdmin(c *gin.Context)  { h.get(c, false) }
func (h *Handler) get(c *gin.Context, published bool) {
	value, err := h.service.Get(c, c.Param("id"), published)
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
