package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func authorizationDatabase(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	t.Helper()
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = sqlDB.Close() })
	db, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	return db, mock
}

func TestRequireRoleAllowsAdmin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock := authorizationDatabase(t)
	mock.ExpectQuery(`SELECT role FROM "profiles"`).
		WithArgs("user-1").WillReturnRows(sqlmock.NewRows([]string{"role"}).AddRow("admin"))
	router := gin.New()
	router.GET("/", func(c *gin.Context) { c.Set("user_id", "user-1") }, RequireRole(db, "admin"), func(c *gin.Context) {
		if c.GetString("user_role") != "admin" {
			t.Fatal("expected role to be available to downstream handlers")
		}
		c.Status(http.StatusNoContent)
	})
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/", nil))
	if recorder.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d: %s", recorder.Code, recorder.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestRequireRoleRejectsMember(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock := authorizationDatabase(t)
	mock.ExpectQuery(`SELECT role FROM "profiles"`).
		WithArgs("user-1").WillReturnRows(sqlmock.NewRows([]string{"role"}).AddRow("member"))
	router := gin.New()
	router.GET("/", func(c *gin.Context) { c.Set("user_id", "user-1") }, RequireRole(db, "admin"), func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/", nil))
	if recorder.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d: %s", recorder.Code, recorder.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestRequireRoleRejectsMissingIdentity(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock := authorizationDatabase(t)
	mock.ExpectQuery(`SELECT role FROM "profiles"`).
		WithArgs("").WillReturnRows(sqlmock.NewRows([]string{"role"}))
	router := gin.New()
	router.GET("/", RequireRole(db, "admin"), func(c *gin.Context) {
		t.Fatal("protected handler must not run")
	})
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/", nil))
	if recorder.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d: %s", recorder.Code, recorder.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestRequireRoleMapsDatabaseFailure(t *testing.T) {
	gin.SetMode(gin.TestMode)
	db, mock := authorizationDatabase(t)
	mock.ExpectQuery(`SELECT role FROM "profiles"`).
		WithArgs("user-1").WillReturnError(errors.New("database unavailable"))
	router := gin.New()
	router.GET("/", func(c *gin.Context) { c.Set("user_id", "user-1") }, RequireRole(db, "admin"), func(c *gin.Context) {
		t.Fatal("protected handler must not run")
	})
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/", nil))
	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d: %s", recorder.Code, recorder.Body.String())
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
