package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/sirridemirtas/anonsocial/controllers"
	"github.com/sirridemirtas/anonsocial/middleware"
)

func PostRoutes(rg *gin.RouterGroup) {
	posts := rg.Group("/posts")
	{
		posts.GET("", controllers.GetPosts)
		posts.GET("/:id", controllers.GetPost)
		posts.GET("/university/:universityId", controllers.GetPostsByUniversity)
		posts.GET("/:id/replies", controllers.GetPostReplies)
		posts.POST("", middleware.Auth(0), controllers.CreatePost)
		posts.DELETE("/:id", middleware.Auth(0), controllers.DeletePost)
	}
	users := rg.Group("/users")
	{
		users.GET("/:username/posts", controllers.GetPostsByUser)
	}
}
