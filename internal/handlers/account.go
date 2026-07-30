package handlers

import (
	"net/http"
	"strings"
	"time"

	"repo-backend/internal/dto"
	"repo-backend/pkg/apperr"
	"repo-backend/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Account struct{ DB *gorm.DB }

type profileView struct {
	ID        string  `json:"id"`
	Email     string  `json:"email"`
	Username  string  `json:"username"`
	AvatarURL *string `json:"avatar_url"`
	Role      string  `json:"role"`
}
type whitelistView struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	CreatedBy *string   `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
}

func (h *Account) Me(c *gin.Context) {
	email := c.GetString("user_email")
	username := strings.Split(email, "@")[0]
	if err := h.DB.WithContext(c).Exec(`INSERT INTO profiles (id,email,username,role)
		VALUES (?,?,?,'member') ON CONFLICT (id) DO NOTHING`, c.GetString("user_id"), email, username).Error; err != nil {
		response.Fail(c, apperr.Internal(err))
		return
	}
	var v profileView
	err := h.DB.WithContext(c).Table("profiles").Select("id,email,username,avatar_url,role").Where("id = ?", c.GetString("user_id")).Take(&v).Error
	if err != nil {
		response.Fail(c, mapDBError(err, "profile"))
		return
	}
	response.OK(c, v)
}
func (h *Account) UpdateMe(c *gin.Context) {
	var in dto.UpdateProfile
	if !bind(c, &in) {
		return
	}
	result := h.DB.WithContext(c).Table("profiles").Where("id = ?", c.GetString("user_id")).Updates(map[string]any{"username": strings.TrimSpace(in.Username), "avatar_url": in.AvatarURL, "updated_at": time.Now().UTC()})
	if result.Error != nil {
		response.Fail(c, apperr.Internal(result.Error))
		return
	}
	if result.RowsAffected == 0 {
		response.Fail(c, apperr.NotFound("profile"))
		return
	}
	h.Me(c)
}
func (h *Account) ListWhitelist(c *gin.Context) {
	var rows []whitelistView
	if err := h.DB.WithContext(c).Table("signup_whitelist").Order("created_at DESC").Find(&rows).Error; err != nil {
		response.Fail(c, apperr.Internal(err))
		return
	}
	response.OK(c, rows)
}
func (h *Account) AddWhitelist(c *gin.Context) {
	var in dto.AddWhitelist
	if !bind(c, &in) {
		return
	}
	row := whitelistView{Email: strings.ToLower(strings.TrimSpace(in.Email))}
	creator := c.GetString("user_id")
	row.CreatedBy = &creator
	err := h.DB.WithContext(c).Table("signup_whitelist").Create(&row).Error
	if err != nil {
		response.Fail(c, apperr.New("CONFLICT", "email is already whitelisted", http.StatusConflict))
		return
	}
	response.Created(c, row)
}
func (h *Account) DeleteWhitelist(c *gin.Context) {
	result := h.DB.WithContext(c).Exec("DELETE FROM signup_whitelist WHERE id = ?", c.Param("id"))
	if result.Error != nil {
		response.Fail(c, apperr.Internal(result.Error))
		return
	}
	if result.RowsAffected == 0 {
		response.Fail(c, apperr.NotFound("whitelist entry"))
		return
	}
	response.NoContent(c)
}
func mapDBError(err error, resource string) error {
	if err == gorm.ErrRecordNotFound {
		return apperr.NotFound(resource)
	}
	return apperr.Internal(err)
}
