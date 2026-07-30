package handlers

import (
	"net/http"
	"time"

	"repo-backend/internal/database"
	"repo-backend/internal/models"

	"github.com/gin-gonic/gin"
)

func CreateEvent(c *gin.Context) {
	var event models.Event
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := c.GetString("user_id")
	event.CreatedBy = &userID
	if err := database.DB.Create(&event).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, event)
}

func GetEvents(c *gin.Context) {
	var events []models.Event
	query := database.DB.Order("event_date DESC")
	if category := c.Query("category"); category != "" {
		query = query.Where("category = ?", category)
	}
	if c.Query("published") == "true" {
		query = query.Where("status = ?", "PUBLISHED")
	}
	if err := query.Find(&events).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, events)
}

func GetEvent(c *gin.Context) {
	var event models.Event
	if err := database.DB.First(&event, "id = ?", c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
		return
	}
	c.JSON(http.StatusOK, event)
}

func UpdateEvent(c *gin.Context) {
	var payload models.Event
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	payload.ID = ""
	payload.CreatedBy = nil
	payload.CreatedAt = time.Time{}
	payload.UpdatedAt = time.Now()
	result := database.DB.Model(&models.Event{}).
		Where("id = ?", c.Param("id")).
		Updates(payload)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
		return
	}
	c.Status(http.StatusNoContent)
}

func DeleteEvent(c *gin.Context) {
	result := database.DB.Delete(&models.Event{}, "id = ?", c.Param("id"))
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
