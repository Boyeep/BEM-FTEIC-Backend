package gallery

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
	items   []models.Gallery
	item    *models.Gallery
	listErr error
	findErr error
}

func (s *repositoryStub) List(context.Context, shared.ListOptions) ([]models.Gallery, int64, error) {
	return s.items, int64(len(s.items)), s.listErr
}
func (s *repositoryStub) Find(context.Context, string) (*models.Gallery, error) {
	if s.findErr != nil {
		return nil, s.findErr
	}
	if s.item != nil {
		return s.item, nil
	}
	return &models.Gallery{}, nil
}
func (s *repositoryStub) Create(context.Context, *models.Gallery) error { return nil }
func (s *repositoryStub) Update(context.Context, *models.Gallery) error { return nil }
func (s *repositoryStub) Delete(context.Context, string) error          { return nil }

func TestHandlerListContract(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewHandler(NewService(&repositoryStub{items: []models.Gallery{{Title: "Gallery"}}}, nil))
	router := gin.New()
	router.GET("/gallery/", handler.List)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/gallery/?page=1&page_size=12", nil))
	if recorder.Code != http.StatusOK || !bytes.Contains(recorder.Body.Bytes(), []byte(`"pagination"`)) {
		t.Fatalf("unexpected response %d: %s", recorder.Code, recorder.Body.String())
	}
}

func TestHandlerRejectsInvalidCreate(t *testing.T) {
	handler := NewHandler(NewService(&repositoryStub{}, nil))
	router := gin.New()
	router.POST("/admin/gallery", handler.Create)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodPost, "/admin/gallery", bytes.NewBufferString(`{"link":"bad"}`)))
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", recorder.Code)
	}
}

func TestHandlerCreateAndDeleteContracts(t *testing.T) {
	handler := NewHandler(NewService(&repositoryStub{}, nil))
	router := gin.New()
	router.POST("/admin/gallery", func(c *gin.Context) { c.Set("user_id", "11111111-1111-4111-8111-111111111111") }, handler.Create)
	router.DELETE("/admin/gallery/:id", handler.Delete)
	body := `{"title":"Gallery","link":"https://example.com/gallery","category":"all","taken_at":"2026-08-01"}`
	request := httptest.NewRequest(http.MethodPost, "/admin/gallery", bytes.NewBufferString(body))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", recorder.Code, recorder.Body.String())
	}
	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodDelete, "/admin/gallery/id", nil))
	if recorder.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", recorder.Code)
	}
}

func TestHandlerGetAndUpdateContracts(t *testing.T) {
	item := &models.Gallery{
		ID: "gallery-id", Title: "Before", Link: "https://example.com/before",
		Category: "all", TakenAt: "2026-08-01",
	}
	handler := NewHandler(NewService(&repositoryStub{item: item}, nil))
	router := gin.New()
	router.GET("/gallery/:id", handler.Get)
	router.PUT("/admin/gallery/:id", handler.Update)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/gallery/gallery-id", nil))
	if recorder.Code != http.StatusOK || !bytes.Contains(recorder.Body.Bytes(), []byte(`"title":"Before"`)) {
		t.Fatalf("unexpected get response %d: %s", recorder.Code, recorder.Body.String())
	}

	body := `{"title":"After","link":"https://example.com/after","category":"teknik_informatika","taken_at":"2026-08-02"}`
	request := httptest.NewRequest(http.MethodPut, "/admin/gallery/gallery-id", bytes.NewBufferString(body))
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
	router.GET("/gallery/:id", handler.Get)
	router.GET("/gallery/", handler.List)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/gallery/missing", nil))
	if recorder.Code != http.StatusNotFound || !bytes.Contains(recorder.Body.Bytes(), []byte(`"code":"NOT_FOUND"`)) {
		t.Fatalf("unexpected not-found response %d: %s", recorder.Code, recorder.Body.String())
	}

	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/gallery/?sort=invalid", nil))
	if recorder.Code != http.StatusBadRequest || !bytes.Contains(recorder.Body.Bytes(), []byte(`"code":"VALIDATION_ERROR"`)) {
		t.Fatalf("unexpected validation response %d: %s", recorder.Code, recorder.Body.String())
	}

	errorRouter := gin.New()
	errorRouter.GET("/gallery/", NewHandler(NewService(&repositoryStub{listErr: context.Canceled}, nil)).List)
	recorder = httptest.NewRecorder()
	errorRouter.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/gallery/", nil))
	if recorder.Code != http.StatusInternalServerError || !bytes.Contains(recorder.Body.Bytes(), []byte(`"code":"INTERNAL_ERROR"`)) {
		t.Fatalf("unexpected repository error response %d: %s", recorder.Code, recorder.Body.String())
	}
}
