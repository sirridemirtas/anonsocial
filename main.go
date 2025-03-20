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
	controllers.SetConversationCollection(database.GetClient()) // Initialize conversation collection

	router := gin.Default()
	router.Use(middleware.Cors())

	apiV1 := router.Group("/api/v1")
	routes.UserRoutes(apiV1)
	routes.AuthRoutes(apiV1)
	routes.PostRoutes(apiV1)
	routes.FeedRoutes(apiV1)
	routes.MessageRoutes(apiV1) // Add message routes
	routes.StaticRoutes(router)

	router.Run(":" + config.AppConfig.Port)
}
