package routes

import (
	"fmt"
	"time"

	"repo-backend/config"
	"repo-backend/internal/handlers"
	"repo-backend/internal/media"
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
	mediaService := media.NewService(
		media.NewLocalStorage("uploads"),
		config.Optional("PUBLIC_API_URL", ""),
	)
	h := &handlers.Content{
		Blogs:   service.NewBlog(contents, mediaService),
		Events:  service.NewEvent(repository.NewEventRepository(db), mediaService),
		Gallery: service.NewGallery(repository.NewGalleryRepository(db), mediaService),
	}
	accounts := &handlers.Account{
		Service: service.NewAccount(repository.NewAccountRepository(db)),
	}
	analytics := &handlers.Analytics{
		Service: service.NewAnalytics(repository.NewAnalyticsRepository(db)),
	}
	mediaHandler := &handlers.Media{
		Service: mediaService,
	}

	r := gin.New()
	if err := r.SetTrustedProxies([]string{"127.0.0.1", "172.16.0.0/12"}); err != nil {
		return nil, fmt.Errorf("configure trusted proxies: %w", err)
	}
	r.Use(middleware.RequestID(), middleware.SecurityHeaders(), gin.Logger(), middleware.Recovery())
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
	r.GET("/me", auth, accounts.Me)
	r.PUT("/me", auth, accounts.UpdateMe)
	r.GET("/profiles", accounts.PublicProfiles)
	r.POST("/visitors", analytics.Track)
	r.GET("/visitors/count", analytics.Count)

	registerMediaRoutes(r, mediaHandler, protected)
	registerPublicContentRoutes(r, h)
	registerAdminRoutes(r, h, accounts, protected)
	return r, nil
}

func registerMediaRoutes(r *gin.Engine, handler *handlers.Media, protected []gin.HandlerFunc) {
	r.POST("/uploads/images", append(protected, handler.UploadImage)...)
	r.DELETE("/uploads/images", append(protected, handler.DeleteImage)...)
	r.Static("/uploads", "./uploads")
}

func registerPublicContentRoutes(r *gin.Engine, handler *handlers.Content) {
	blogs := r.Group("/blogs")
	blogs.GET("/", handler.ListBlogs)
	blogs.GET("/:id", handler.GetBlog)

	events := r.Group("/events")
	events.GET("/", handler.ListEvents)
	events.GET("/:id", handler.GetEvent)

	gallery := r.Group("/gallery")
	gallery.GET("/", handler.ListGallery)
	gallery.GET("/:id", handler.GetGallery)
}

func registerAdminRoutes(r *gin.Engine, h *handlers.Content, accounts *handlers.Account, protected []gin.HandlerFunc) {
	adminAPI := r.Group("/admin", protected...)
	adminAPI.GET("/blogs", h.ListAdminBlogs)
	adminAPI.GET("/blogs/:id", h.GetAdminBlog)
	adminAPI.POST("/blogs", h.CreateBlog)
	adminAPI.PUT("/blogs/:id", h.UpdateBlog)
	adminAPI.DELETE("/blogs/:id", h.DeleteBlog)
	adminAPI.GET("/events", h.ListAdminEvents)
	adminAPI.GET("/events/:id", h.GetAdminEvent)
	adminAPI.POST("/events", h.CreateEvent)
	adminAPI.PUT("/events/:id", h.UpdateEvent)
	adminAPI.DELETE("/events/:id", h.DeleteEvent)
	adminAPI.GET("/gallery", h.ListGallery)
	adminAPI.GET("/gallery/:id", h.GetGallery)
	adminAPI.POST("/gallery", h.CreateGallery)
	adminAPI.PUT("/gallery/:id", h.UpdateGallery)
	adminAPI.DELETE("/gallery/:id", h.DeleteGallery)
	adminAPI.GET("/whitelist", accounts.ListWhitelist)
	adminAPI.POST("/whitelist", accounts.AddWhitelist)
	adminAPI.DELETE("/whitelist/:id", accounts.DeleteWhitelist)
}
