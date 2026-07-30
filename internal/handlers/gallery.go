package handlers

import (
	"net/http"
	"time"

	"repo-backend/internal/database"
	"repo-backend/internal/models"

	"github.com/gin-gonic/gin"
)

func CreateGallery(c *gin.Context) {
	var gallery models.Gallery
	if err := c.ShouldBindJSON(&gallery); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := c.GetString("user_id")
	gallery.CreatedBy = &userID
	if err := database.DB.Create(&gallery).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gallery)
}

func GetGalleries(c *gin.Context) {
	var galleries []models.Gallery
	if err := database.DB.Order("taken_at DESC").Find(&galleries).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, galleries)
}

func GetGalleryByID(c *gin.Context) {
	var gallery models.Gallery
	if err := database.DB.First(&gallery, "id = ?", c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "gallery item not found"})
		return
	}
	c.JSON(http.StatusOK, gallery)
}

func UpdateGallery(c *gin.Context) {
	var payload models.Gallery
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	payload.ID = ""
	payload.CreatedBy = nil
	payload.CreatedAt = time.Time{}
	payload.UpdatedAt = time.Now()
	result := database.DB.Model(&models.Gallery{}).
		Where("id = ?", c.Param("id")).
		Updates(payload)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "gallery item not found"})
		return
	}
	c.Status(http.StatusNoContent)
}

func DeleteGallery(c *gin.Context) {
	result := database.DB.Delete(&models.Gallery{}, "id = ?", c.Param("id"))
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
