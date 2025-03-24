package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/sirridemirtas/anonsocial/controllers"
	"github.com/sirridemirtas/anonsocial/middleware"
)

func FeedRoutes(rg *gin.RouterGroup) {
	// Keep only the feed route to avoid conflicts with post routes
	rg.GET("/feed", middleware.OptionalAuth(), controllers.GetFeedPosts)

	// Remove these duplicated routes - they're already defined in PostRoutes
	// rg.GET("/posts/:id/replies", controllers.GetFeedPostReplies)         - DUPLICATE
	// rg.GET("/users/:username/posts", controllers.GetFeedUserPosts)       - DUPLICATE
	// rg.GET("/posts/university/:universityId", controllers.GetFeedUniversityPosts) - DUPLICATE

	// You can add other unique feed-related routes here if needed, for example:
	// rg.GET("/feed/trending", middleware.OptionalAuth(), controllers.GetTrendingPosts)
	// rg.GET("/feed/following", middleware.Auth(0), controllers.GetFollowingPosts)
}
