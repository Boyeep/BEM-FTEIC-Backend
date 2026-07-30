package handlers

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"repo-backend/pkg/apperr"
	"repo-backend/pkg/response"

	"github.com/gabriel-vasile/mimetype"
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
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 10<<20)
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
	opened, err := file.Open()
	if err != nil {
		response.Fail(c, apperr.Validation("invalid image"))
		return
	}
	detected, err := mimetype.DetectReader(io.LimitReader(opened, 512*1024))
	opened.Close()
	if err != nil || !strings.HasPrefix(detected.String(), "image/") ||
		!imageMIMEMatchesExtension(detected.String(), extension) {
		response.Fail(c, apperr.Validation("file content does not match a supported image type"))
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

func imageMIMEMatchesExtension(mime, extension string) bool {
	switch extension {
	case ".jpg", ".jpeg":
		return mime == "image/jpeg"
	case ".png":
		return mime == "image/png"
	case ".gif":
		return mime == "image/gif"
	case ".webp":
		return mime == "image/webp"
	default:
		return false
	}
}

func DeleteImage(c *gin.Context) {
	rawURL := c.Query("url")
	parsed, err := url.Parse(rawURL)
	if err != nil || parsed.Path == "" {
		response.Fail(c, apperr.Validation("valid upload URL is required"))
		return
	}
	filename := filepath.Base(parsed.Path)
	if filename == "." || filename == "/" {
		response.Fail(c, apperr.Validation("invalid upload URL"))
		return
	}
	extension := strings.ToLower(filepath.Ext(filename))
	if _, ok := allowedImageExtensions[extension]; !ok {
		response.Fail(c, apperr.Validation("invalid upload file"))
		return
	}
	err = os.Remove(filepath.Join("uploads", filename))
	if err != nil && !os.IsNotExist(err) {
		response.Fail(c, apperr.Internal(err))
		return
	}
	response.NoContent(c)
}
