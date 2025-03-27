package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port         string
	MongoDBURI   string
	MongoDB_DB   string
	JWTSecret    string
	JWTExpiresIn string
}

var AppConfig Config

func LoadConfig() {
	env := os.Getenv("GO_ENV")
	if env == "" {
		env = "development"
	}

	err := godotenv.Load(".env." + env)
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	AppConfig = Config{
		Port:         os.Getenv("PORT"),
		MongoDBURI:   os.Getenv("MONGODB_URI"),
		MongoDB_DB:   os.Getenv("MONGODB_DB"),
		JWTSecret:    os.Getenv("JWT_SECRET"),
		JWTExpiresIn: os.Getenv("JWT_EXPIRES_IN"),
	}
}
