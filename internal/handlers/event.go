package handlers

import (
	"net/http"
	"repo-backend/internal/database"
	"repo-backend/internal/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// CreateEvent membuat event baru dengan upload foto
func CreateEvent(c *gin.Context) {
	title := c.PostForm("title")
	description := c.PostForm("description")
	organizer := c.PostForm("organizer")
	location := c.PostForm("location")
	isPublishedStr := c.PostForm("is_published")
	startDateStr := c.PostForm("start_date")
	endDateStr := c.PostForm("end_date")

	if title == "" || organizer == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title dan organizer wajib diisi"})
		return
	}

	// Parse tanggal (format: 2006-01-02)
	var startDate, endDate time.Time
	if startDateStr != "" {
		if t, err := time.Parse("2006-01-02", startDateStr); err == nil {
			startDate = t
		}
	}
	if endDateStr != "" {
		if t, err := time.Parse("2006-01-02", endDateStr); err == nil {
			endDate = t
		}
	}

	isPublished := isPublishedStr == "true" || isPublishedStr == "1"

	// Handle upload foto
	var photoPath string
	file, err := c.FormFile("photo")
	if err == nil {
		filename := strconv.FormatInt(time.Now().UnixNano(), 10) + "_" + file.Filename
		filePath := "uploads/" + filename
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan foto"})
			return
		}
		photoPath = filePath
	}

	event := models.Event{
		Title:       title,
		Description: description,
		Photo:       photoPath,
		Organizer:   models.Organizer(organizer),
		StartDate:   startDate,
		EndDate:     endDate,
		Location:    location,
		IsPublished: isPublished,
	}

	if result := database.DB.Create(&event); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, event)
}

// GetEvents mengambil semua event, bisa difilter by organizer
func GetEvents(c *gin.Context) {
	var events []models.Event

	query := database.DB.Model(&models.Event{})

	// Filter opsional berdasarkan organizer
	if organizer := c.Query("organizer"); organizer != "" {
		query = query.Where("organizer = ?", organizer)
	}

	// Filter opsional hanya yang sudah dipublish
	if published := c.Query("published"); published == "true" {
		query = query.Where("is_published = ?", true)
	}

	query = query.Order("start_date ASC")

	if result := query.Find(&events); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, events)
}

// GetEvent mengambil satu event berdasarkan ID
func GetEvent(c *gin.Context) {
	id := c.Param("id")
	var event models.Event

	if result := database.DB.First(&event, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, event)
}

// UpdateEvent mengupdate event berdasarkan ID
func UpdateEvent(c *gin.Context) {
	id := c.Param("id")
	var event models.Event

	if result := database.DB.First(&event, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event tidak ditemukan"})
		return
	}

	// Update field teks jika ada
	if title := c.PostForm("title"); title != "" {
		event.Title = title
	}
	if description := c.PostForm("description"); description != "" {
		event.Description = description
	}
	if organizer := c.PostForm("organizer"); organizer != "" {
		event.Organizer = models.Organizer(organizer)
	}
	if location := c.PostForm("location"); location != "" {
		event.Location = location
	}
	if startDateStr := c.PostForm("start_date"); startDateStr != "" {
		if t, err := time.Parse("2006-01-02", startDateStr); err == nil {
			event.StartDate = t
		}
	}
	if endDateStr := c.PostForm("end_date"); endDateStr != "" {
		if t, err := time.Parse("2006-01-02", endDateStr); err == nil {
			event.EndDate = t
		}
	}
	if isPublishedStr := c.PostForm("is_published"); isPublishedStr != "" {
		event.IsPublished = isPublishedStr == "true" || isPublishedStr == "1"
	}

	// Update foto jika ada file baru
	file, err := c.FormFile("photo")
	if err == nil {
		filename := strconv.FormatInt(time.Now().UnixNano(), 10) + "_" + file.Filename
		filePath := "uploads/" + filename
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan foto"})
			return
		}
		event.Photo = filePath
	}

	if result := database.DB.Save(&event); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, event)
}

// DeleteEvent menghapus event berdasarkan ID
func DeleteEvent(c *gin.Context) {
	id := c.Param("id")
	var event models.Event

	if result := database.DB.First(&event, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event tidak ditemukan"})
		return
	}

	if result := database.DB.Delete(&event); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Event berhasil dihapus"})
}
