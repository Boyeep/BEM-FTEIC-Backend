package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var allowedImageExtensions = map[string]struct{}{
	".gif":  {},
	".jpeg": {},
	".jpg":  {},
	".png":  {},
	".webp": {},
}

func UploadImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}
	if file.Size > 10<<20 {
		c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "file exceeds 10 MB"})
		return
	}

	extension := strings.ToLower(filepath.Ext(file.Filename))
	if _, ok := allowedImageExtensions[extension]; !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported image type"})
		return
	}

	userID := c.GetString("user_id")
	filename := fmt.Sprintf("%s-%d%s", userID, time.Now().UnixNano(), extension)
	destination := filepath.Join("uploads", filename)
	if err := os.MkdirAll("uploads", 0o750); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "prepare upload directory"})
		return
	}
	if err := c.SaveUploadedFile(file, destination); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "save image"})
		return
	}

	publicBaseURL := strings.TrimRight(os.Getenv("PUBLIC_API_URL"), "/")
	if publicBaseURL == "" {
		publicBaseURL = strings.TrimRight("https://"+c.Request.Host, "/")
	}
	c.JSON(http.StatusCreated, gin.H{
		"url": publicBaseURL + "/uploads/" + filename,
	})
}
