package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/sirridemirtas/anonsocial/controllers"
	"github.com/sirridemirtas/anonsocial/middleware"
)

func FeedRoutes(rg *gin.RouterGroup) {
	feeds := rg.Group("/feeds")
	rg.Use(middleware.OptionalAuth())

	// Home feed, includes posts from all users
	feeds.GET("/home", controllers.GetHomeFeed)

	// University feed, includes posts from specific university
	feeds.GET("/universities/:universityId", controllers.GetUniversityFeed)

	// User feed, includes posts from specific user
	feeds.GET("/users/:username", controllers.GetUserFeed)
}
