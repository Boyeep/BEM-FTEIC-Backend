package media

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"repo-backend/pkg/apperr"

	"github.com/gabriel-vasile/mimetype"
)

const MaxImageSize = 10 << 20

var allowedExtensions = map[string]string{
	".gif":  "image/gif",
	".jpeg": "image/jpeg",
	".jpg":  "image/jpeg",
	".png":  "image/png",
	".webp": "image/webp",
}

type Service struct {
	storage       Storage
	publicBaseURL string
}

func NewService(storage Storage, publicBaseURL string) *Service {
	return &Service{storage: storage, publicBaseURL: strings.TrimRight(publicBaseURL, "/")}
}

func (s *Service) Upload(ctx context.Context, userID string, file *multipart.FileHeader, fallbackBaseURL string) (string, error) {
	if file.Size > MaxImageSize {
		return "", apperr.New("FILE_TOO_LARGE", "file exceeds 10 MB", http.StatusRequestEntityTooLarge)
	}
	extension := strings.ToLower(filepath.Ext(file.Filename))
	expectedMIME, ok := allowedExtensions[extension]
	if !ok {
		return "", apperr.Validation("unsupported image type")
	}
	opened, err := file.Open()
	if err != nil {
		return "", apperr.Validation("invalid image")
	}
	detected, err := mimetype.DetectReader(io.LimitReader(opened, 512*1024))
	opened.Close()
	if err != nil || detected.String() != expectedMIME {
		return "", apperr.Validation("file content does not match a supported image type")
	}
	source, err := file.Open()
	if err != nil {
		return "", apperr.Validation("invalid image")
	}
	defer source.Close()
	key := fmt.Sprintf("%s-%d%s", userID, time.Now().UnixNano(), extension)
	if err := s.storage.Save(ctx, key, source); err != nil {
		return "", apperr.Internal(err)
	}
	baseURL := s.publicBaseURL
	if baseURL == "" {
		baseURL = strings.TrimRight(fallbackBaseURL, "/")
	}
	return baseURL + "/uploads/" + key, nil
}

func (s *Service) Delete(ctx context.Context, rawURL string) error {
	parsed, err := url.Parse(rawURL)
	if err != nil || !strings.HasPrefix(parsed.Path, "/uploads/") {
		return apperr.Validation("valid upload URL is required")
	}
	key := filepath.Base(parsed.Path)
	extension := strings.ToLower(filepath.Ext(key))
	if key == "." || key == "/" || allowedExtensions[extension] == "" {
		return apperr.Validation("invalid upload file")
	}
	err = s.storage.Delete(ctx, key)
	if errors.Is(err, ErrNotFound) {
		return nil
	}
	if err != nil {
		return apperr.Internal(err)
	}
	return nil
}
