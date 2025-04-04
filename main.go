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
	controllers.SetNotificationCollection(database.GetClient())
	controllers.SetSitemapPostCollection(database.GetClient())

	middleware.SetActivityCollection(database.GetClient(), config.AppConfig.MongoDB_DB)
	controllers.SetActivityCollection(database.GetClient(), config.AppConfig.MongoDB_DB)

	router := gin.Default()
	router.Use(middleware.Cors())
	router.Use(middleware.RateLimit())

	apiV1 := router.Group("/api/v1")

	routes.AuthRoutes(apiV1)
	routes.UserRoutes(apiV1)
	routes.PostRoutes(apiV1)
	routes.FeedRoutes(apiV1)
	routes.MessageRoutes(apiV1)
	routes.NotificationRoutes(apiV1)
	routes.AdminRoutes(apiV1)

	routes.StaticRoutes(router)

	apiV1.POST("/contact", controllers.SubmitContactForm)

	apiV1.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "Server is running",
		})
	})

	// Serve the sitemap.xml file for the Frontend App
	apiV1.GET("/sitemap.xml", func(c *gin.Context) {
		controllers.GenerateSitemapXML(c.Writer, c.Request)
	})

	router.Run(":" + config.AppConfig.Port)
}
