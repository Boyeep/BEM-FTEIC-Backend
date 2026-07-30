package media

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestLocalStorageSaveAndDelete(t *testing.T) {
	root := t.TempDir()
	storage := NewLocalStorage(root)

	if err := storage.Save(context.Background(), "image.png", bytes.NewBufferString("content")); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(root, "image.png")); err != nil {
		t.Fatalf("saved image is missing: %v", err)
	}
	if err := storage.Delete(context.Background(), "image.png"); err != nil {
		t.Fatal(err)
	}
	if err := storage.Delete(context.Background(), "image.png"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestDeleteRejectsNonUploadURL(t *testing.T) {
	service := NewService(NewLocalStorage(t.TempDir()), "https://api.example.com")
	err := service.Delete(context.Background(), "https://example.com/avatar.png")
	if err == nil {
		t.Fatal("expected non-upload URL to be rejected")
	}
}
