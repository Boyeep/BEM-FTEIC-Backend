package handlers

import (
	"repo-backend/internal/dto"
	"repo-backend/internal/service"
	"repo-backend/pkg/apperr"
	"repo-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

type Content struct {
	Blogs   *service.BlogService
	Events  *service.EventService
	Gallery *service.GalleryService
}

func bind(c *gin.Context, target any) bool {
	if err := c.ShouldBindJSON(target); err != nil {
		response.Fail(c, apperr.Validation(err.Error()))
		return false
	}
	return true
}

func bindListQuery(c *gin.Context, defaultPageSize int) (dto.ListQuery, bool) {
	var query dto.ListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Fail(c, apperr.Validation(err.Error()))
		return query, false
	}
	query.Defaults(defaultPageSize)
	return query, true
}

func (h *Content) ListBlogs(c *gin.Context) {
	query, ok := bindListQuery(c, 6)
	if !ok {
		return
	}
	v, e := h.Blogs.List(c, query, true)
	if e != nil {
		response.Fail(c, e)
		return
	}
	response.OK(c, v)
}
func (h *Content) GetBlog(c *gin.Context) {
	v, e := h.Blogs.GetPublic(c, c.Param("id"))
	if e != nil {
		response.Fail(c, e)
		return
	}
	response.OK(c, v)
}
func (h *Content) ListAdminBlogs(c *gin.Context) {
	query, ok := bindListQuery(c, 20)
	if !ok {
		return
	}
	v, e := h.Blogs.List(c, query, false)
	if e != nil {
		response.Fail(c, e)
		return
	}
	response.OK(c, v)
}
func (h *Content) GetAdminBlog(c *gin.Context) {
	v, e := h.Blogs.Get(c, c.Param("id"))
	if e != nil {
		response.Fail(c, e)
		return
	}
	response.OK(c, v)
}
func (h *Content) CreateBlog(c *gin.Context) {
	var in dto.CreateBlog
	if !bind(c, &in) {
		return
	}
	v, e := h.Blogs.Create(c, in, c.GetString("user_id"))
	if e != nil {
		response.Fail(c, e)
		return
	}
	response.Created(c, v)
}
func (h *Content) UpdateBlog(c *gin.Context) {
	var in dto.UpdateBlog
	if !bind(c, &in) {
		return
	}
	v, e := h.Blogs.Update(c, c.Param("id"), in)
	if e != nil {
		response.Fail(c, e)
		return
	}
	response.OK(c, v)
}
func (h *Content) DeleteBlog(c *gin.Context) {
	if e := h.Blogs.Delete(c, c.Param("id")); e != nil {
		response.Fail(c, e)
		return
	}
	response.NoContent(c)
}

func (h *Content) ListEvents(c *gin.Context) {
	query, ok := bindListQuery(c, 8)
	if !ok {
		return
	}
	v, e := h.Events.List(c, query, true)
	if e != nil {
		response.Fail(c, e)
		return
	}
	response.OK(c, v)
}
func (h *Content) GetEvent(c *gin.Context) {
	v, e := h.Events.GetPublic(c, c.Param("id"))
	if e != nil {
		response.Fail(c, e)
		return
	}
	response.OK(c, v)
}
func (h *Content) ListAdminEvents(c *gin.Context) {
	query, ok := bindListQuery(c, 20)
	if !ok {
		return
	}
	v, e := h.Events.List(c, query, false)
	if e != nil {
		response.Fail(c, e)
		return
	}
	response.OK(c, v)
}
func (h *Content) GetAdminEvent(c *gin.Context) {
	v, e := h.Events.Get(c, c.Param("id"))
	if e != nil {
		response.Fail(c, e)
		return
	}
	response.OK(c, v)
}
func (h *Content) CreateEvent(c *gin.Context) {
	var in dto.CreateEvent
	if !bind(c, &in) {
		return
	}
	v, e := h.Events.Create(c, in, c.GetString("user_id"))
	if e != nil {
		response.Fail(c, e)
		return
	}
	response.Created(c, v)
}
func (h *Content) UpdateEvent(c *gin.Context) {
	var in dto.UpdateEvent
	if !bind(c, &in) {
		return
	}
	v, e := h.Events.Update(c, c.Param("id"), in)
	if e != nil {
		response.Fail(c, e)
		return
	}
	response.OK(c, v)
}
func (h *Content) DeleteEvent(c *gin.Context) {
	if e := h.Events.Delete(c, c.Param("id")); e != nil {
		response.Fail(c, e)
		return
	}
	response.NoContent(c)
}

func (h *Content) ListGallery(c *gin.Context) {
	query, ok := bindListQuery(c, 12)
	if !ok {
		return
	}
	v, e := h.Gallery.List(c, query)
	if e != nil {
		response.Fail(c, e)
		return
	}
	response.OK(c, v)
}
func (h *Content) GetGallery(c *gin.Context) {
	v, e := h.Gallery.Get(c, c.Param("id"))
	if e != nil {
		response.Fail(c, e)
		return
	}
	response.OK(c, v)
}
func (h *Content) CreateGallery(c *gin.Context) {
	var in dto.CreateGallery
	if !bind(c, &in) {
		return
	}
	v, e := h.Gallery.Create(c, in, c.GetString("user_id"))
	if e != nil {
		response.Fail(c, e)
		return
	}
	response.Created(c, v)
}
func (h *Content) UpdateGallery(c *gin.Context) {
	var in dto.UpdateGallery
	if !bind(c, &in) {
		return
	}
	v, e := h.Gallery.Update(c, c.Param("id"), in)
	if e != nil {
		response.Fail(c, e)
		return
	}
	response.OK(c, v)
}
func (h *Content) DeleteGallery(c *gin.Context) {
	if e := h.Gallery.Delete(c, c.Param("id")); e != nil {
		response.Fail(c, e)
		return
	}
	response.NoContent(c)
}
