package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/sirridemirtas/anonsocial/controllers"
	"github.com/sirridemirtas/anonsocial/middleware"
)

func FeedRoutes(rg *gin.RouterGroup) {
	rg.GET("/posts", middleware.OptionalAuth(), controllers.GetFeedPosts)
	rg.GET("/posts/:id/replies", middleware.OptionalAuth(), controllers.GetFeedPostReplies)
	rg.GET("/posts/university/:universityId", controllers.GetFeedUniversityPosts)
	rg.GET("/users/:username/posts", middleware.OptionalAuth(), middleware.OptionalAuth(), controllers.GetFeedUserPosts)
}
