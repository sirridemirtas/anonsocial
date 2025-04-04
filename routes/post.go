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

		posts.POST("", middleware.CustomRateLimit(1, 1), middleware.Auth(0), middleware.ActivityTracker(), controllers.CreatePost)
		posts.DELETE("/:id", middleware.Auth(0), controllers.DeletePost)

		posts.POST("/:id/like", middleware.CustomRateLimit(1, 3), middleware.Auth(0), controllers.LikePost)
		posts.POST("/:id/dislike", middleware.CustomRateLimit(1, 3), middleware.Auth(0), controllers.DislikePost)

		posts.DELETE("/:id/unlike", middleware.CustomRateLimit(1, 3), middleware.Auth(0), controllers.RemoveLikePost)
		posts.DELETE("/:id/undislike", middleware.CustomRateLimit(1, 3), middleware.Auth(0), controllers.RemoveDislikePost)
	}
}
