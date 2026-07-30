package handlers

import (
	"repo-backend/internal/dto"
	"repo-backend/internal/service"
	"repo-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

type Account struct{ Service *service.AccountService }

func (h *Account) PublicProfiles(c *gin.Context) {
	profiles, err := h.Service.PublicProfiles(c, c.Query("ids"))
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, profiles)
}

func (h *Account) Me(c *gin.Context) {
	profile, err := h.Service.Me(c, c.GetString("user_id"), c.GetString("user_email"))
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, profile)
}

func (h *Account) UpdateMe(c *gin.Context) {
	var in dto.UpdateProfile
	if !bind(c, &in) {
		return
	}
	profile, err := h.Service.UpdateProfile(c, c.GetString("user_id"), in.Username, in.AvatarURL)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, profile)
}

func (h *Account) ListWhitelist(c *gin.Context) {
	rows, err := h.Service.ListWhitelist(c)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, rows)
}

func (h *Account) AddWhitelist(c *gin.Context) {
	var in dto.AddWhitelist
	if !bind(c, &in) {
		return
	}
	row, err := h.Service.AddWhitelist(c, in.Email, c.GetString("user_id"))
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.Created(c, row)
}

func (h *Account) DeleteWhitelist(c *gin.Context) {
	if err := h.Service.DeleteWhitelist(c, c.Param("id")); err != nil {
		response.Fail(c, err)
		return
	}
	response.NoContent(c)
}
