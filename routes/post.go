package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/sirridemirtas/anonsocial/controllers"
	"github.com/sirridemirtas/anonsocial/middleware"
)

func PostRoutes(rg *gin.RouterGroup) {
	posts := rg.Group("/posts")
	{
		// Public routes with optional auth to handle privacy
		posts.GET("/:id", middleware.OptionalAuth(), controllers.GetPost)
		posts.GET("/:id/replies", middleware.OptionalAuth(), controllers.GetPostReplies)

		posts.POST("", middleware.Auth(0), controllers.CreatePost)
		posts.DELETE("/:id", middleware.Auth(0), controllers.DeletePost)

		posts.POST("/:id/like", middleware.Auth(0), controllers.LikePost)
		posts.POST("/:id/dislike", middleware.Auth(0), controllers.DislikePost)

		posts.DELETE("/:id/unlike", middleware.Auth(0), controllers.RemoveLikePost)
		posts.DELETE("/:id/undislike", middleware.Auth(0), controllers.RemoveDislikePost)
	}
}
