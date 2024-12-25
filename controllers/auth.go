package controllers

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"

	"github.com/sirridemirtas/anonsocial/config"
	"github.com/sirridemirtas/anonsocial/middleware"
	"github.com/sirridemirtas/anonsocial/models"
	"github.com/sirridemirtas/anonsocial/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func Register(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get form data
	username := c.PostForm("username")
	password := c.PostForm("password")
	universityId := c.PostForm("universityId")

	// Validate required fields
	if username == "" || password == "" || universityId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username, password and universityId are required"})
		return
	}

	// Create user object with only allowed fields
	user := models.User{
		Username:     username,
		UniversityID: universityId,
		CreatedAt:    time.Now(),
		Role:         0, // Default role
		IsPrivate:    false,
	}

	// Generate salt and hash password
	user.Salt = models.GenerateSalt()
	user.Password = user.HashPassword(password)

	// Validate user data
	if errors := utils.ValidateUser(&user); len(errors) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"errors": errors})
		return
	}

	result, err := userCollection.InsertOne(ctx, user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) || strings.Contains(err.Error(), "duplicate key error") {
			c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user.ID = result.InsertedID.(primitive.ObjectID)

	// Generate JWT token after successful registration
	claims := &middleware.Claims{
		UserID: user.ID.Hex(),
		Role:   user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.AppConfig.JWTSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	c.SetCookie("token", tokenString, 3600*24, "/", "", false, true)
	c.JSON(http.StatusCreated, gin.H{
		"message": "Registration successful",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"role":     user.Role,
		},
	})
}

func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	if username == "" || password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username and password are required"})
		return
	}

	var user models.User
	err := userCollection.FindOne(context.Background(), bson.M{"username": username}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if !user.ValidatePassword(password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	claims := &middleware.Claims{
		UserID: user.ID.Hex(),
		Role:   user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.AppConfig.JWTSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	c.SetCookie("token", tokenString, 3600*24, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

func Logout(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}
