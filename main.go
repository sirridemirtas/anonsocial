package main

import (
	"github.com/gin-gonic/gin"

	"github.com/sirridemirtas/anonsocial/config"
	"github.com/sirridemirtas/anonsocial/controllers"
	"github.com/sirridemirtas/anonsocial/database"
	"github.com/sirridemirtas/anonsocial/middleware"
	"github.com/sirridemirtas/anonsocial/routes"
)

func main() {
	config.LoadConfig()

	database.ConnectDB()
	defer database.DisconnectDB()

	controllers.SetUserCollection(database.GetClient())
	controllers.SetPostCollection(database.GetClient())
	controllers.SetConversationCollection(database.GetClient())
	controllers.SetNotificationCollection(database.GetClient()) // Add this line

	router := gin.Default()
	router.Use(middleware.Cors())

	apiV1 := router.Group("/api/v1")

	// Register routes in the correct order
	// First register basic routes
	routes.AuthRoutes(apiV1)
	routes.UserRoutes(apiV1)
	routes.PostRoutes(apiV1)

	// Then register composite routes (like feed which might depend on posts)
	routes.FeedRoutes(apiV1)
	routes.MessageRoutes(apiV1)
	routes.NotificationRoutes(apiV1) // Add this line
	routes.StaticRoutes(router)

	// Contact form route
	router.POST("/api/v1/contact", controllers.SubmitContactForm)

	router.Run(":" + config.AppConfig.Port)
}
