package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/gin-gonic/gin"

	"github.com/sirridemirtas/anonsocial/controllers"
	"github.com/sirridemirtas/anonsocial/routes"
)

var client *mongo.Client

func connectDB() *mongo.Client {
	uri := "mongodb://localhost:27017"

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
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
	client := connectDB()
	defer client.Disconnect(nil)

	controllers.SetUserCollection(client)

	router := gin.Default()

	apiV1 := router.Group("/api/v1")
	routes.UserRoutes(apiV1)

	router.Run(":8080")
}
