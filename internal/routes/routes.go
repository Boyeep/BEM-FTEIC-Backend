package routes

import (
	"fmt"
	"time"

	"repo-backend/config"
	"repo-backend/internal/handlers"
	"repo-backend/internal/middleware"
	"repo-backend/internal/repository"
	"repo-backend/internal/service"
	"repo-backend/pkg/apperr"
	"repo-backend/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) (*gin.Engine, error) {
	supabaseURL, err := config.Required("SUPABASE_URL")
	if err != nil {
		return nil, err
	}
	supabaseAnonKey, err := config.Required("SUPABASE_ANON_KEY")
	if err != nil {
		return nil, err
	}

	contents := repository.NewContentRepository(db)
	h := &handlers.Content{
		Blogs:   service.NewBlog(contents),
		Events:  service.NewEvent(repository.NewEventRepository(db)),
		Gallery: service.NewGallery(repository.NewGalleryRepository(db)),
	}

	r := gin.New()
	if err := r.SetTrustedProxies([]string{"127.0.0.1", "172.16.0.0/12"}); err != nil {
		return nil, fmt.Errorf("configure trusted proxies: %w", err)
	}
	r.Use(middleware.RequestID(), gin.Logger(), middleware.Recovery())
	r.Use(middleware.CORS(config.Optional("ALLOWED_ORIGINS", "http://localhost:3000,https://bem-fteic.com,https://www.bem-fteic.com")))
	r.Use(middleware.RateLimit(120, time.Minute))
	r.NoRoute(func(c *gin.Context) {
		response.Fail(c, apperr.NotFound("route"))
	})

	r.GET("/health", func(c *gin.Context) { response.OK(c, gin.H{"status": "ok"}) })
	r.GET("/ready", func(c *gin.Context) {
		sqlDB, e := db.DB()
		if e != nil || sqlDB.PingContext(c) != nil {
			response.Fail(c, apperr.New("NOT_READY", "database unavailable", 503))
			return
		}
		response.OK(c, gin.H{"status": "ready"})
	})

	auth := middleware.NewSupabaseAuth(supabaseURL, supabaseAnonKey).Required()
	admin := middleware.RequireRole(db, "admin")
	protected := []gin.HandlerFunc{auth, admin}

	r.POST("/uploads/images", append(protected, handlers.UploadImage)...)
	r.Static("/uploads", "./uploads")

	blogs := r.Group("/blogs")
	blogs.GET("/", h.ListBlogs)
	blogs.GET("/:id", h.GetBlog)
	blogs.POST("/", append(protected, h.CreateBlog)...)
	blogs.POST("/upload-image", append(protected, handlers.UploadImage)...)
	blogs.PUT("/:id", append(protected, h.UpdateBlog)...)
	blogs.DELETE("/:id", append(protected, h.DeleteBlog)...)

	events := r.Group("/events")
	events.GET("/", h.ListEvents)
	events.GET("/:id", h.GetEvent)
	events.POST("/", append(protected, h.CreateEvent)...)
	events.PUT("/:id", append(protected, h.UpdateEvent)...)
	events.DELETE("/:id", append(protected, h.DeleteEvent)...)

	gallery := r.Group("/gallery")
	gallery.GET("/", h.ListGallery)
	gallery.GET("/:id", h.GetGallery)
	gallery.POST("/", append(protected, h.CreateGallery)...)
	gallery.PUT("/:id", append(protected, h.UpdateGallery)...)
	gallery.DELETE("/:id", append(protected, h.DeleteGallery)...)
	return r, nil
}
