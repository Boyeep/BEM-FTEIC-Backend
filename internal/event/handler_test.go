package event

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
	items   []models.Event
	item    *models.Event
	listErr error
	findErr error
}

func (s *repositoryStub) List(context.Context, shared.ListOptions) ([]models.Event, int64, error) {
	return s.items, int64(len(s.items)), s.listErr
}
func (s *repositoryStub) Find(context.Context, string, bool) (*models.Event, error) {
	if s.findErr != nil {
		return nil, s.findErr
	}
	if s.item != nil {
		return s.item, nil
	}
	return &models.Event{}, nil
}
func (s *repositoryStub) Create(context.Context, *models.Event) error { return nil }
func (s *repositoryStub) Update(context.Context, *models.Event) error { return nil }
func (s *repositoryStub) Delete(context.Context, string) error        { return nil }

func TestHandlerListContract(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := NewHandler(NewService(&repositoryStub{items: []models.Event{{Title: "Event"}}}, nil))
	router := gin.New()
	router.GET("/events/", handler.ListPublic)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/events/?page=1&page_size=8", nil))
	if recorder.Code != http.StatusOK || !bytes.Contains(recorder.Body.Bytes(), []byte(`"pagination"`)) {
		t.Fatalf("unexpected response %d: %s", recorder.Code, recorder.Body.String())
	}
}

func TestHandlerRejectsInvalidCreate(t *testing.T) {
	handler := NewHandler(NewService(&repositoryStub{}, nil))
	router := gin.New()
	router.POST("/admin/events", handler.Create)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodPost, "/admin/events", bytes.NewBufferString(`{"status":"PUBLISHED"}`)))
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", recorder.Code)
	}
}

func TestHandlerCreateAndDeleteContracts(t *testing.T) {
	handler := NewHandler(NewService(&repositoryStub{}, nil))
	router := gin.New()
	router.POST("/admin/events", func(c *gin.Context) { c.Set("user_id", "11111111-1111-4111-8111-111111111111") }, handler.Create)
	router.DELETE("/admin/events/:id", handler.Delete)
	body := `{"title":"Event","description":"Description","author":"Admin","category":"FTEIC","event_date":"2026-08-01","status":"UPCOMING","publication_status":"DRAFT"}`
	request := httptest.NewRequest(http.MethodPost, "/admin/events", bytes.NewBufferString(body))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", recorder.Code, recorder.Body.String())
	}
	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodDelete, "/admin/events/id", nil))
	if recorder.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", recorder.Code)
	}
}

func TestHandlerGetAndUpdateContracts(t *testing.T) {
	item := &models.Event{
		ID: "event-id", Title: "Before", Description: "Description", Author: "Admin",
		Category: "FTEIC", EventDate: "2026-08-01", Status: "UPCOMING", PublicationStatus: "DRAFT",
	}
	handler := NewHandler(NewService(&repositoryStub{item: item}, nil))
	router := gin.New()
	router.GET("/events/:id", handler.GetPublic)
	router.PUT("/admin/events/:id", handler.Update)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/events/event-id", nil))
	if recorder.Code != http.StatusOK || !bytes.Contains(recorder.Body.Bytes(), []byte(`"title":"Before"`)) {
		t.Fatalf("unexpected get response %d: %s", recorder.Code, recorder.Body.String())
	}

	body := `{"title":"After","description":"Updated","author":"Admin","category":"FTEIC","event_date":"2026-08-02","status":"ONGOING","publication_status":"PUBLISHED"}`
	request := httptest.NewRequest(http.MethodPut, "/admin/events/event-id", bytes.NewBufferString(body))
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
	router.GET("/events/:id", handler.GetPublic)
	router.GET("/events/", handler.ListPublic)

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/events/missing", nil))
	if recorder.Code != http.StatusNotFound || !bytes.Contains(recorder.Body.Bytes(), []byte(`"code":"NOT_FOUND"`)) {
		t.Fatalf("unexpected not-found response %d: %s", recorder.Code, recorder.Body.String())
	}

	recorder = httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/events/?start_date=invalid", nil))
	if recorder.Code != http.StatusBadRequest || !bytes.Contains(recorder.Body.Bytes(), []byte(`"code":"VALIDATION_ERROR"`)) {
		t.Fatalf("unexpected validation response %d: %s", recorder.Code, recorder.Body.String())
	}

	errorRouter := gin.New()
	errorRouter.GET("/events/", NewHandler(NewService(&repositoryStub{listErr: context.Canceled}, nil)).ListPublic)
	recorder = httptest.NewRecorder()
	errorRouter.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/events/", nil))
	if recorder.Code != http.StatusInternalServerError || !bytes.Contains(recorder.Body.Bytes(), []byte(`"code":"INTERNAL_ERROR"`)) {
		t.Fatalf("unexpected repository error response %d: %s", recorder.Code, recorder.Body.String())
	}
}
