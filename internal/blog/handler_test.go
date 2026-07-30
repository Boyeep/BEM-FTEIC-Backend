package blog

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	shared "repo-backend/internal/content"
	"repo-backend/internal/models"
	"repo-backend/internal/repository"

	"github.com/gin-gonic/gin"
)

type repositoryStub struct {
	items   []models.Blog
	item    *models.Blog
	listErr error
	findErr error
}

func (s *repositoryStub) List(context.Context, shared.ListOptions) ([]models.Blog, int64, error) {
	return s.items, int64(len(s.items)), s.listErr
}
func (s *repositoryStub) Find(context.Context, string, bool) (*models.Blog, error) {
	if s.findErr != nil {
		return nil, s.findErr
	}
	if s.item != nil {
		return s.item, nil
	}
	return &models.Blog{}, nil
}
func (s *repositoryStub) Create(context.Context, *models.Blog) error { return nil }
func (s *repositoryStub) Update(context.Context, *models.Blog) error { return nil }
func (s *repositoryStub) Delete(context.Context, string) error       { return nil }

func TestHandlerListContract(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewHandler(NewService(&repositoryStub{items: []models.Blog{{Title: "Blog"}}}, nil))
	router := gin.New()
	router.GET("/blogs/", handler.ListPublic)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/blogs/?page=1&page_size=6", nil))
	if recorder.Code != http.StatusOK || !bytes.Contains(recorder.Body.Bytes(), []byte(`"pagination"`)) {
		t.Fatalf("unexpected response %d: %s", recorder.Code, recorder.Body.String())
	}
}

func TestHandlerRejectsInvalidCreate(t *testing.T) {
	handler := NewHandler(NewService(&repositoryStub{}, nil))
	router := gin.New()
	router.POST("/admin/blogs", handler.Create)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodPost, "/admin/blogs", bytes.NewBufferString(`{"status":"INVALID"}`)))
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", recorder.Code)
	}
}

func TestHandlerCreateAndDeleteContracts(t *testing.T) {
	handler := NewHandler(NewService(&repositoryStub{}, nil))
	router := gin.New()
	router.POST("/admin/blogs", func(c *gin.Context) { c.Set("user_id", "11111111-1111-4111-8111-111111111111") }, handler.Create)
	router.DELETE("/admin/blogs/:id", handler.Delete)
	body := `{"title":"Blog","author":"Admin","category":"FTEIC","content":"Content","status":"DRAFT"}`
	request := httptest.NewRequest(http.MethodPost, "/admin/blogs", bytes.NewBufferString(body))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", recorder.Code, recorder.Body.String())
	}
	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodDelete, "/admin/blogs/id", nil))
	if recorder.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", recorder.Code)
	}
}

func TestHandlerGetAndUpdateContracts(t *testing.T) {
	item := &models.Blog{ID: "blog-id", Title: "Before", Author: "Admin", Category: "FTEIC", Content: "Content", Status: "DRAFT"}
	handler := NewHandler(NewService(&repositoryStub{item: item}, nil))
	router := gin.New()
	router.GET("/blogs/:id", handler.GetPublic)
	router.PUT("/admin/blogs/:id", handler.Update)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/blogs/blog-id", nil))
	if recorder.Code != http.StatusOK || !bytes.Contains(recorder.Body.Bytes(), []byte(`"title":"Before"`)) {
		t.Fatalf("unexpected get response %d: %s", recorder.Code, recorder.Body.String())
	}

	body := `{"title":"After","author":"Admin","category":"FTEIC","content":"Updated","status":"PUBLISHED","published_at":"2026-08-01T00:00:00Z"}`
	request := httptest.NewRequest(http.MethodPut, "/admin/blogs/blog-id", bytes.NewBufferString(body))
	request.Header.Set("Content-Type", "application/json")
	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK || !bytes.Contains(recorder.Body.Bytes(), []byte(`"title":"After"`)) {
		t.Fatalf("unexpected update response %d: %s", recorder.Code, recorder.Body.String())
	}
}

func TestHandlerValidationAndRepositoryErrors(t *testing.T) {
	router := gin.New()
	handler := NewHandler(NewService(&repositoryStub{findErr: repository.ErrNotFound}, nil))
	router.GET("/blogs/:id", handler.GetPublic)
	router.GET("/blogs/", handler.ListPublic)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/blogs/missing", nil))
	if recorder.Code != http.StatusNotFound || !bytes.Contains(recorder.Body.Bytes(), []byte(`"code":"NOT_FOUND"`)) {
		t.Fatalf("unexpected not-found response %d: %s", recorder.Code, recorder.Body.String())
	}

	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/blogs/?page=-1", nil))
	if recorder.Code != http.StatusBadRequest || !bytes.Contains(recorder.Body.Bytes(), []byte(`"code":"VALIDATION_ERROR"`)) {
		t.Fatalf("unexpected validation response %d: %s", recorder.Code, recorder.Body.String())
	}

	errorRouter := gin.New()
	errorRouter.GET("/blogs/", NewHandler(NewService(&repositoryStub{listErr: context.Canceled}, nil)).ListPublic)
	recorder = httptest.NewRecorder()
	errorRouter.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/blogs/", nil))
	if recorder.Code != http.StatusInternalServerError || !bytes.Contains(recorder.Body.Bytes(), []byte(`"code":"INTERNAL_ERROR"`)) {
		t.Fatalf("unexpected repository error response %d: %s", recorder.Code, recorder.Body.String())
	}
}
