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
		blogGroup.GET("/", handlers.GetBlogs)
		blogGroup.GET("/:id", handlers.GetBlog)
		blogGroup.PUT("/:id", handlers.UpdateBlog)
		blogGroup.DELETE("/:id", handlers.DeleteBlog)
	}

	return r
}
