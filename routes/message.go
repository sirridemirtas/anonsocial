package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/sirridemirtas/anonsocial/controllers"
	"github.com/sirridemirtas/anonsocial/middleware"
)

func MessageRoutes(rg *gin.RouterGroup) {
	messages := rg.Group("/messages")
	messages.Use(middleware.Auth(0)) // All message routes require authentication

	// Get conversation list for current user
	messages.GET("", controllers.GetConversationList)

	// Get total unread message count
	messages.GET("/unread_count", controllers.GetTotalUnreadCount)

	// Get conversation with specific user
	messages.GET("/:username", controllers.GetConversation)

	// Send message to specific user
	messages.POST("/:username", controllers.SendMessage)

	// Mark messages as read
	messages.POST("/:username/read", controllers.MarkConversationAsRead)

	// Delete conversation with specific user
	messages.DELETE("/:username", controllers.DeleteConversation)
}
