package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/sirridemirtas/anonsocial/controllers"
	"github.com/sirridemirtas/anonsocial/middleware"
)

func NotificationRoutes(rg *gin.RouterGroup) {
	notifications := rg.Group("/notifications")
	notifications.Use(middleware.Auth(0)) // All notification routes require authentication

	// Get all notifications for the authenticated user
	notifications.GET("", controllers.GetNotifications)

	// Get the count of unread notifications
	notifications.GET("/unread-count", controllers.GetUnreadCount)

	// Mark a notification as read
	notifications.PUT("/:id", controllers.MarkAsRead)

	// Mark all notifications as read
	notifications.PUT("/mark-all-read", controllers.MarkAllAsRead)

	// Delete all notifications
	notifications.DELETE("/delete-all", controllers.DeleteAllNotifications)
}
