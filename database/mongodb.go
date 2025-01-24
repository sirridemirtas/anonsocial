package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/sirridemirtas/anonsocial/config"
)

var Client *mongo.Client

func ConnectDB() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.AppConfig.MongoDBURI))
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("MongoDB bağlantı hatası:", err)
	}

	Client = client
	log.Println("\033[32m", "MongoDB'ye başarıyla bağlanıldı!", "\033[0m")
}

func GetClient() *mongo.Client {
	return Client
}

func DisconnectDB() {
	if err := Client.Disconnect(context.Background()); err != nil {
		log.Printf("MongoDB bağlantısı kapatılırken hata oluştu: %v", err)
	}
} 