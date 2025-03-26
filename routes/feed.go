package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/sirridemirtas/anonsocial/controllers"
	"github.com/sirridemirtas/anonsocial/middleware"
)

func FeedRoutes(rg *gin.RouterGroup) {
	rg.GET("/feed", middleware.OptionalAuth(), controllers.GetFeedPosts)

	rg.GET("/posts/:id/replies", middleware.OptionalAuth(), controllers.GetFeedPostReplies)
	rg.GET("/users/:username/posts", middleware.OptionalAuth(), middleware.OptionalAuth(), controllers.GetFeedUserPosts)
	rg.GET("/posts/university/:universityId", controllers.GetFeedUniversityPosts)
}
