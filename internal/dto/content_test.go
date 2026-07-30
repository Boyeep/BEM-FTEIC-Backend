package dto

import (
	"testing"

	"github.com/go-playground/validator/v10"
)

func TestContentStatusContracts(t *testing.T) {
	validate := validator.New()
	validate.SetTagName("binding")
	valid := []any{
		CreateBlog{Title: "Blog", Author: "Admin", Category: "FTEIC", Content: "Content", Status: "ARCHIVED"},
		CreateEvent{Title: "Event", Description: "Description", Author: "Admin", Category: "FTEIC", EventDate: "2026-08-01", Status: "UPCOMING", PublicationStatus: "DRAFT"},
		CreateGallery{Title: "Gallery", Link: "https://example.com", Category: "teknik_informatika", TakenAt: "2026-08-01"},
	}
	for _, input := range valid {
		if err := validate.Struct(input); err != nil {
			t.Fatalf("expected valid input %#v: %v", input, err)
		}
	}

	invalid := CreateEvent{Title: "Event", Description: "Description", Author: "Admin", Category: "FTEIC", EventDate: "2026-08-01", Status: "PUBLISHED", PublicationStatus: "ONGOING"}
	if err := validate.Struct(invalid); err == nil {
		t.Fatal("expected lifecycle/publication validation failure")
	}
}
