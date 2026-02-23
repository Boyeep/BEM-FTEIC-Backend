package routes

import (
	"repo-backend/internal/handlers"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	blogGroup := r.Group("/blogs")
	{
		blogGroup.POST("/", handlers.CreateBlog)
		blogGroup.POST("/upload-image", handlers.UploadBlogImage) // New endpoint for content images
		blogGroup.GET("/", handlers.GetBlogs)
		blogGroup.GET("/:id", handlers.GetBlog)
		blogGroup.PUT("/:id", handlers.UpdateBlog)
		blogGroup.DELETE("/:id", handlers.DeleteBlog)
	}

	adminGroup := r.Group("/admins")
	{
		adminGroup.POST("/", handlers.CreateAdmin)
		adminGroup.GET("/", handlers.GetAdmins)
		adminGroup.GET("/:id", handlers.GetAdmin)
		adminGroup.PUT("/:id", handlers.UpdateAdmin)
		adminGroup.DELETE("/:id", handlers.DeleteAdmin)
	}

	galleryGroup := r.Group("/gallery")
	{
		galleryGroup.POST("/", handlers.CreateGallery)
		galleryGroup.GET("/", handlers.GetGalleries)
		galleryGroup.GET("/:id", handlers.GetGalleryByID)
		galleryGroup.PUT("/:id", handlers.UpdateGallery)
		galleryGroup.DELETE("/:id", handlers.DeleteGallery)
	}

	eventGroup := r.Group("/events")
	{
		eventGroup.POST("/", handlers.CreateEvent)
		eventGroup.GET("/", handlers.GetEvents)
		eventGroup.GET("/:id", handlers.GetEvent)
		eventGroup.PUT("/:id", handlers.UpdateEvent)
		eventGroup.DELETE("/:id", handlers.DeleteEvent)
	}

	return r
}
