package routes

import (
	"fmt"
	"time"

	"repo-backend/config"
	"repo-backend/internal/container"
	"repo-backend/internal/handlers"
	"repo-backend/internal/middleware"
	"repo-backend/pkg/apperr"
	"repo-backend/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB, dependencies *container.Container) (*gin.Engine, error) {
	supabaseURL, err := config.Required("SUPABASE_URL")
	if err != nil {
		return nil, err
	}
	supabaseAnonKey, err := config.Required("SUPABASE_ANON_KEY")
	if err != nil {
		return nil, err
	}

	accounts := dependencies.Accounts

	r := gin.New()
	if err := r.SetTrustedProxies([]string{"127.0.0.1", "172.16.0.0/12"}); err != nil {
		return nil, fmt.Errorf("configure trusted proxies: %w", err)
	}
	r.Use(middleware.RequestID(), middleware.SecurityHeaders(), gin.Logger(), middleware.Recovery())
	r.Use(middleware.CORS(config.Optional("ALLOWED_ORIGINS", "http://localhost:3000,https://bem-fteic.com,https://www.bem-fteic.com")))
	rateLimitStore := middleware.RateLimitStore(middleware.NewMemoryRateLimitStore())
	if redisURL := config.Optional("REDIS_URL", ""); redisURL != "" {
		rateLimitStore, err = middleware.NewRedisRateLimitStore(redisURL, "bem-fteic:rate-limit")
		if err != nil {
			return nil, err
		}
	}
	r.Use(middleware.RateLimitWithStore(rateLimitStore, 120, time.Minute))
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
	r.POST("/visitors", dependencies.Analytics.Track)
	r.GET("/visitors/count", dependencies.Analytics.Count)

	registerMediaRoutes(r, dependencies.Media, protected)
	registerPublicContentRoutes(r, dependencies)
	registerAdminRoutes(r, dependencies, accounts, protected)
	return r, nil
}

func registerMediaRoutes(r *gin.Engine, handler *handlers.Media, protected []gin.HandlerFunc) {
	r.POST("/uploads/images", append(protected, handler.UploadImage)...)
	r.DELETE("/uploads/images", append(protected, handler.DeleteImage)...)
	r.Static("/uploads", "./uploads")
}

func registerPublicContentRoutes(r *gin.Engine, dependencies *container.Container) {
	blogs := r.Group("/blogs")
	blogs.GET("/", dependencies.Blog.ListPublic)
	blogs.GET("/:id", dependencies.Blog.GetPublic)

	events := r.Group("/events")
	events.GET("/", dependencies.Event.ListPublic)
	events.GET("/:id", dependencies.Event.GetPublic)

	gallery := r.Group("/gallery")
	gallery.GET("/", dependencies.Gallery.List)
	gallery.GET("/:id", dependencies.Gallery.Get)
}

func registerAdminRoutes(r *gin.Engine, dependencies *container.Container, accounts *handlers.Account, protected []gin.HandlerFunc) {
	adminAPI := r.Group("/admin", protected...)
	adminAPI.GET("/blogs", dependencies.Blog.ListAdmin)
	adminAPI.GET("/blogs/:id", dependencies.Blog.GetAdmin)
	adminAPI.POST("/blogs", dependencies.Blog.Create)
	adminAPI.PUT("/blogs/:id", dependencies.Blog.Update)
	adminAPI.DELETE("/blogs/:id", dependencies.Blog.Delete)
	adminAPI.GET("/events", dependencies.Event.ListAdmin)
	adminAPI.GET("/events/:id", dependencies.Event.GetAdmin)
	adminAPI.POST("/events", dependencies.Event.Create)
	adminAPI.PUT("/events/:id", dependencies.Event.Update)
	adminAPI.DELETE("/events/:id", dependencies.Event.Delete)
	adminAPI.GET("/gallery", dependencies.Gallery.List)
	adminAPI.GET("/gallery/:id", dependencies.Gallery.Get)
	adminAPI.POST("/gallery", dependencies.Gallery.Create)
	adminAPI.PUT("/gallery/:id", dependencies.Gallery.Update)
	adminAPI.DELETE("/gallery/:id", dependencies.Gallery.Delete)
	adminAPI.GET("/whitelist", accounts.ListWhitelist)
	adminAPI.POST("/whitelist", accounts.AddWhitelist)
	adminAPI.DELETE("/whitelist/:id", accounts.DeleteWhitelist)
}
