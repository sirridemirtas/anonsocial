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
		posts.GET("", middleware.OptionalAuth(), controllers.GetPosts)
		posts.GET("/:id", middleware.OptionalAuth(), controllers.GetPost)
		//posts.GET("/:id/replies", middleware.OptionalAuth(), controllers.GetPostReplies)
		//posts.GET("/university/:universityId", middleware.OptionalAuth(), controllers.GetPostsByUniversity)

		// Authenticated routes
		posts.POST("", middleware.Auth(0), controllers.CreatePost)
		posts.DELETE("/:id", middleware.Auth(0), controllers.DeletePost)

		// reaction endpoints
		posts.POST("/:id/like", middleware.Auth(0), controllers.LikePost)
		posts.POST("/:id/dislike", middleware.Auth(0), controllers.DislikePost)

		// removing reactions with PUT method
		posts.POST("/:id/unlike", middleware.Auth(0), controllers.RemoveLikePost)
		posts.POST("/:id/undislike", middleware.Auth(0), controllers.RemoveDislikePost)
	}

	// Add a route for user posts with optional auth
	//rg.GET("/users/:username/posts", middleware.OptionalAuth(), controllers.GetPostsByUser)
}
