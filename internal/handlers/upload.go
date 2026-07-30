package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"repo-backend/pkg/apperr"
	"repo-backend/pkg/response"

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
		response.Fail(c, apperr.Validation("file is required"))
		return
	}
	if file.Size > 10<<20 {
		response.Fail(c, apperr.New("FILE_TOO_LARGE", "file exceeds 10 MB", http.StatusRequestEntityTooLarge))
		return
	}

	extension := strings.ToLower(filepath.Ext(file.Filename))
	if _, ok := allowedImageExtensions[extension]; !ok {
		response.Fail(c, apperr.Validation("unsupported image type"))
		return
	}

	userID := c.GetString("user_id")
	filename := fmt.Sprintf("%s-%d%s", userID, time.Now().UnixNano(), extension)
	destination := filepath.Join("uploads", filename)
	if err := os.MkdirAll("uploads", 0o750); err != nil {
		response.Fail(c, apperr.Internal(err))
		return
	}
	if err := c.SaveUploadedFile(file, destination); err != nil {
		response.Fail(c, apperr.Internal(err))
		return
	}

	publicBaseURL := strings.TrimRight(os.Getenv("PUBLIC_API_URL"), "/")
	if publicBaseURL == "" {
		publicBaseURL = strings.TrimRight("https://"+c.Request.Host, "/")
	}
	response.Created(c, gin.H{
		"url": publicBaseURL + "/uploads/" + filename,
	})
}
