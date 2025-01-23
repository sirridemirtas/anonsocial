package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/gin-gonic/gin"

	"github.com/sirridemirtas/anonsocial/config"
	"github.com/sirridemirtas/anonsocial/controllers"
	"github.com/sirridemirtas/anonsocial/middleware"
	"github.com/sirridemirtas/anonsocial/routes"
)

var client *mongo.Client

func connectDB() *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.AppConfig.MongoDBURI))
	if err != nil {
		log.Fatal(err)
	}

	// Bağlantıyı kontrol et
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("MongoDB bağlantı hatası:", err)
	}

	log.Println("MongoDB'ye başarıyla bağlanıldı!")
	return client
}

func main() {
	config.LoadConfig()

	client := connectDB()
	defer client.Disconnect(nil)

	controllers.SetUserCollection(client)
	controllers.SetPostCollection(client)

	router := gin.Default()
	router.Use(middleware.Cors()) // Add CORS middleware

	apiV1 := router.Group("/api/v1")
	routes.UserRoutes(apiV1)
	routes.AuthRoutes(apiV1)
	routes.PostRoutes(apiV1)

	routes.StaticRoutes(router)

	router.Run(":" + config.AppConfig.Port)
}
