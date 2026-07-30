package routes

import (
	"fmt"
	"net/http"

	"repo-backend/config"
	"repo-backend/internal/handlers"
	"repo-backend/internal/middleware"

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

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	r.Use(middleware.CORS(config.Optional(
		"ALLOWED_ORIGINS",
		"http://localhost:3000,https://bem-fteic.com,https://www.bem-fteic.com",
	)))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	r.GET("/ready", func(c *gin.Context) {
		sqlDB, err := db.DB()
		if err != nil || sqlDB.PingContext(c.Request.Context()) != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not_ready"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ready"})
	})

	authRequired := middleware.NewSupabaseAuth(supabaseURL, supabaseAnonKey).Required()

	r.POST("/uploads/images", authRequired, handlers.UploadImage)

	blogGroup := r.Group("/blogs")
	{
		blogGroup.GET("/", handlers.GetBlogs)
		blogGroup.GET("/:id", handlers.GetBlog)
		blogGroup.POST("/", authRequired, handlers.CreateBlog)
		blogGroup.POST("/upload-image", authRequired, handlers.UploadBlogImage)
		blogGroup.PUT("/:id", authRequired, handlers.UpdateBlog)
		blogGroup.DELETE("/:id", authRequired, handlers.DeleteBlog)
	}

	adminGroup := r.Group("/admins", authRequired)
	{
		adminGroup.POST("/", handlers.CreateAdmin)
		adminGroup.GET("/", handlers.GetAdmins)
		adminGroup.GET("/:id", handlers.GetAdmin)
		adminGroup.PUT("/:id", handlers.UpdateAdmin)
		adminGroup.DELETE("/:id", handlers.DeleteAdmin)
	}

	galleryGroup := r.Group("/gallery")
	{
		galleryGroup.GET("/", handlers.GetGalleries)
		galleryGroup.GET("/:id", handlers.GetGalleryByID)
		galleryGroup.POST("/", authRequired, handlers.CreateGallery)
		galleryGroup.PUT("/:id", authRequired, handlers.UpdateGallery)
		galleryGroup.DELETE("/:id", authRequired, handlers.DeleteGallery)
	}

	eventGroup := r.Group("/events")
	{
		eventGroup.GET("/", handlers.GetEvents)
		eventGroup.GET("/:id", handlers.GetEvent)
		eventGroup.POST("/", authRequired, handlers.CreateEvent)
		eventGroup.PUT("/:id", authRequired, handlers.UpdateEvent)
		eventGroup.DELETE("/:id", authRequired, handlers.DeleteEvent)
	}

	r.Static("/uploads", "./uploads")

	if err := r.SetTrustedProxies(nil); err != nil {
		return nil, fmt.Errorf("disable trusted proxies: %w", err)
	}
	return r, nil
}
