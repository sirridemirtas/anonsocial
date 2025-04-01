package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/sirridemirtas/anonsocial/controllers"
	"github.com/sirridemirtas/anonsocial/middleware"
)

func FeedRoutes(rg *gin.RouterGroup) {
	// Home feed, includes posts from all users
	rg.GET("/posts", middleware.OptionalAuth(), controllers.GetFeedPosts)

	// University feed, includes posts from specific university
	rg.GET("/posts/university/:universityId", controllers.GetFeedUniversityPosts)

	// User feed, includes posts from specific user
	rg.GET("/users/:username/posts", middleware.OptionalAuth(), middleware.OptionalAuth(), controllers.GetFeedUserPosts)
}
