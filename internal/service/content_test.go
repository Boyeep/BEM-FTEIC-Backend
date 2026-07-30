package service

import (
	"testing"
)

func TestNewPage(t *testing.T) {
	page := newPage([]string{"a", "b"}, 2, 2, 5)
	if page.Pagination.TotalPages != 3 {
		t.Fatalf("expected 3 pages, got %d", page.Pagination.TotalPages)
	}
	if !page.Pagination.HasNextPage || !page.Pagination.HasPreviousPage {
		t.Fatal("second page should have previous and next pages")
	}
	if len(page.Items) != 2 {
		t.Fatalf("expected 2 items, got %d", len(page.Items))
	}
}

func TestNewPageEmptyResult(t *testing.T) {
	page := newPage([]string{}, 1, 20, 0)
	if page.Pagination.TotalPages != 1 {
		t.Fatalf("empty result should still expose one page, got %d", page.Pagination.TotalPages)
	}
	if page.Pagination.HasNextPage || page.Pagination.HasPreviousPage {
		t.Fatal("empty first page must not expose navigation")
	}
}
