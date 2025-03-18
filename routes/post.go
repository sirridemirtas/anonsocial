package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/sirridemirtas/anonsocial/controllers"
	"github.com/sirridemirtas/anonsocial/middleware"
)

func PostRoutes(rg *gin.RouterGroup) {
	posts := rg.Group("/posts")
	{
		posts.GET("/:id", controllers.GetPost)
		posts.POST("", middleware.Auth(0), controllers.CreatePost)
		posts.DELETE("/:id", middleware.Auth(0), controllers.DeletePost)

		// reaction endpoints
		posts.POST("/:id/like", middleware.Auth(0), controllers.LikePost)
		posts.POST("/:id/dislike", middleware.Auth(0), controllers.DislikePost)

		// removing reactions
		posts.POST("/:id/unlike", middleware.Auth(0), controllers.RemoveLikePost)
		posts.POST("/:id/undislike", middleware.Auth(0), controllers.RemoveDislikePost)
	}
}
