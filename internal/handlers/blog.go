package handlers

import (
	"net/http"
	"repo-backend/internal/database"
	"repo-backend/internal/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// CreateBlog creates a new blog post with image upload
func CreateBlog(c *gin.Context) {
	// Parse form
	title := c.PostForm("title")
	description := c.PostForm("description")

	adminIDStr := c.PostForm("admin_id")
	var adminID uint
	if adminIDStr != "" {
		if id, err := strconv.ParseUint(adminIDStr, 10, 32); err == nil {
			adminID = uint(id)
		}
	}

	var authorName string
	var authorPhoto string

	// Ambil info Admin (Wajib ada AdminID untuk mengisi penulis)
	if adminID != 0 {
		var admin models.Admin
		if err := database.DB.First(&admin, adminID).Error; err == nil {
			authorName = admin.Nama

			// TODO: Uncomment jika Admin sudah punya ProfilePhoto
			/*
				authorPhoto = admin.ProfilePhoto
			*/
		}
	}

	// Handle upload thumbnail
	var thumbnailPath string
	file, err := c.FormFile("thumbnail")
	if err == nil {
		filename := strconv.FormatInt(time.Now().UnixNano(), 10) + "_" + file.Filename
		filePath := "uploads/" + filename
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan thumbnail"})
			return
		}
		thumbnailPath = filePath
	}

	blog := models.Blog{
		Title:       title,
		Thumbnail:   thumbnailPath,
		AuthorName:  authorName,
		Description: description,
		AuthorPhoto: authorPhoto,
		AdminID:     adminID,
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

	// Update field teks jika ada
	if title := c.PostForm("title"); title != "" {
		blog.Title = title
	}
	if description := c.PostForm("description"); description != "" {
		blog.Description = description
	}

	// Update thumbnail jika ada file baru
	file, err := c.FormFile("thumbnail")
	if err == nil {
		filename := strconv.FormatInt(time.Now().UnixNano(), 10) + "_" + file.Filename
		filePath := "uploads/" + filename
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan thumbnail"})
			return
		}
		blog.Thumbnail = filePath
	}

	if result := database.DB.Save(&blog); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, blog)
}

// UploadBlogImage handles image uploads for blog content (WYSIWYG editor)
func UploadBlogImage(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image is required"})
		return
	}

	// Create unique filename to prevent overwrite
	filename := strconv.FormatInt(time.Now().UnixNano(), 10) + "_" + file.Filename
	filePath := "uploads/" + filename

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
		return
	}

	// Return URL that can be used in <img> tag
	// Adjust base URL as needed (e.g., http://localhost:8080/)
	c.JSON(http.StatusOK, gin.H{
		"url": filePath,
	})
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
