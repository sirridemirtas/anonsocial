package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/sirridemirtas/anonsocial/controllers"
)

func FeedRoutes(rg *gin.RouterGroup) {
	// All routes that return a feed of posts
	rg.GET("/posts", controllers.GetFeedPosts)
	rg.GET("/posts/:id/replies", controllers.GetFeedPostReplies)
	rg.GET("/users/:username/posts", controllers.GetFeedUserPosts)
	rg.GET("/posts/university/:universityId", controllers.GetFeedUniversityPosts)
}
