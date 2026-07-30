package handlers

import (
	"net/http"
	"time"

	"repo-backend/internal/database"
	"repo-backend/internal/models"

	"github.com/gin-gonic/gin"
)

func CreateBlog(c *gin.Context) {
	var blog models.Blog
	if err := c.ShouldBindJSON(&blog); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := c.GetString("user_id")
	blog.CreatedBy = &userID
	if blog.PublishedAt.IsZero() {
		blog.PublishedAt = time.Now()
	}
	if err := database.DB.Create(&blog).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, blog)
}

func GetBlogs(c *gin.Context) {
	var blogs []models.Blog
	query := database.DB.Order("published_at DESC")
	if c.Query("published") == "true" {
		query = query.Where("status = ?", "PUBLISHED")
	}
	if err := query.Find(&blogs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, blogs)
}

func GetBlog(c *gin.Context) {
	var blog models.Blog
	if err := database.DB.First(&blog, "id = ?", c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "blog not found"})
		return
	}
	c.JSON(http.StatusOK, blog)
}

func UpdateBlog(c *gin.Context) {
	var payload models.Blog
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	payload.ID = ""
	payload.CreatedBy = nil
	payload.CreatedAt = time.Time{}
	payload.UpdatedAt = time.Now()
	result := database.DB.Model(&models.Blog{}).
		Where("id = ?", c.Param("id")).
		Updates(payload)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "blog not found"})
		return
	}
	c.Status(http.StatusNoContent)
}

func DeleteBlog(c *gin.Context) {
	result := database.DB.Delete(&models.Blog{}, "id = ?", c.Param("id"))
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func UploadBlogImage(c *gin.Context) {
	UploadImage(c)
}
