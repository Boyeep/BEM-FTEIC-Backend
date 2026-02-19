package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"repo-backend/internal/database"
	"repo-backend/internal/models"

	"github.com/gin-gonic/gin"
)

// =========================
// CREATE GALLERY
// =========================
func CreateGallery(c *gin.Context) {

	title := c.PostForm("title")
	driveLink := c.PostForm("drive_link")

	if title == "" || driveLink == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Title and drive_link are required",
		})
		return
	}

	if !strings.Contains(driveLink, "drive.google.com") {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Google Drive link",
		})
		return
	}

	// Ambil file thumbnail
	file, err := c.FormFile("thumbnail")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Thumbnail is required",
		})
		return
	}

	// Buat nama file unik
	fileName := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)
	filePath := filepath.Join("uploads", fileName)

	// Simpan file
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save thumbnail",
		})
		return
	}

	gallery := models.Gallery{
		Title:     title,
		Thumbnail: filePath,
		DriveLink: driveLink,
	}

	if err := database.DB.Create(&gallery).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create gallery",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Gallery created successfully",
		"data":    gallery,
	})
}

// =========================
// GET ALL GALLERIES
// =========================
func GetGalleries(c *gin.Context) {
	var galleries []models.Gallery

	if err := database.DB.
		Order("created_at desc").
		Find(&galleries).Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch galleries",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": galleries,
	})
}

// =========================
// GET GALLERY BY ID
// =========================
func GetGalleryByID(c *gin.Context) {
	var gallery models.Gallery
	id := c.Param("id")

	if err := database.DB.First(&gallery, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Gallery not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gallery,
	})
}

// =========================
// UPDATE GALLERY
// (Bisa update dengan / tanpa ganti thumbnail)
// =========================
func UpdateGallery(c *gin.Context) {
	var gallery models.Gallery
	id := c.Param("id")

	if err := database.DB.First(&gallery, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Gallery not found",
		})
		return
	}

	title := c.PostForm("title")
	driveLink := c.PostForm("drive_link")

	if title != "" {
		gallery.Title = title
	}

	if driveLink != "" {
		if !strings.Contains(driveLink, "drive.google.com") {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid Google Drive link",
			})
			return
		}
		gallery.DriveLink = driveLink
	}

	// Cek apakah ada file baru
	file, err := c.FormFile("thumbnail")
	if err == nil {

		// Hapus file lama
		if gallery.Thumbnail != "" {
			_ = os.Remove(gallery.Thumbnail)
		}

		fileName := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)
		filePath := filepath.Join("uploads", fileName)

		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to update thumbnail",
			})
			return
		}

		gallery.Thumbnail = filePath
	}

	if err := database.DB.Save(&gallery).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update gallery",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Gallery updated successfully",
		"data":    gallery,
	})
}

// =========================
// DELETE GALLERY
// (Sekaligus hapus thumbnail)
// =========================
func DeleteGallery(c *gin.Context) {
	var gallery models.Gallery
	id := c.Param("id")

	if err := database.DB.First(&gallery, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Gallery not found",
		})
		return
	}

	// Hapus file thumbnail jika ada
	if gallery.Thumbnail != "" {
		_ = os.Remove(gallery.Thumbnail)
	}

	if err := database.DB.Delete(&gallery).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete gallery",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Gallery deleted successfully",
	})
}
