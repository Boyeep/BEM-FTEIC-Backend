package handlers

import (
	"net/http"
	"repo-backend/internal/database"
	"repo-backend/internal/models"

	"github.com/gin-gonic/gin"
)

// CreateBlog creates a new blog post with image upload
func CreateBlog(c *gin.Context) {
	// Parse multipart form
	authorName := c.PostForm("author_name")
	description := c.PostForm("description")
	authorPhoto := c.PostForm("author_photo") // Assuming this is also a string/url for now, or another upload

	// Handle Image Upload
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image is required"})
		return
	}

	// Save the file to specific destination
	filePath := "uploads/" + file.Filename
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
		return
	}

	blog := models.Blog{
		AuthorName:  authorName,
		Description: description,
		AuthorPhoto: authorPhoto,
		ImagePath:   filePath,
	}

	result := database.DB.Create(&blog)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, blog)
}

// GetBlogs retrieves all blog posts
func GetBlogs(c *gin.Context) {
	var blogs []models.Blog
	result := database.DB.Find(&blogs)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, blogs)
}

// GetBlog retrieves a single blog post by ID
func GetBlog(c *gin.Context) {
	id := c.Param("id")
	var blog models.Blog
	result := database.DB.First(&blog, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Blog not found"})
		return
	}

	c.JSON(http.StatusOK, blog)
}

// UpdateBlog updates a blog post by ID
func UpdateBlog(c *gin.Context) {
	id := c.Param("id")
	var blog models.Blog
	if err := database.DB.First(&blog, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Blog not found"})
		return
	}

	var updateData models.Blog
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	database.DB.Model(&blog).Updates(updateData)
	c.JSON(http.StatusOK, blog)
}

// DeleteBlog deletes a blog post by ID
func DeleteBlog(c *gin.Context) {
	id := c.Param("id")
	result := database.DB.Delete(&models.Blog{}, id)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Blog deleted successfully"})
}
