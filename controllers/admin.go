package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirridemirtas/anonsocial/middleware"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var activityCollection *mongo.Collection

// SetActivityCollection sets the activity collection for the admin controller
func SetActivityCollection(client *mongo.Client, dbName string) {
	activityCollection = client.Database(dbName).Collection("user_activities")
}

// UpdateUserRole allows administrators (role=2) to update other users' roles
func UpdateUserRole(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get the admin's role from the context (set by Auth middleware)
	adminRole := c.GetInt("userRole")
	if adminRole != 2 {
		c.JSON(http.StatusForbidden, gin.H{"error": "Bu işlem için yönetici yetkileri gerekiyor"}) // Admin privileges required
		return
	}

	// Get username of the user to update
	username := c.Param("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kullanıcı adı parametresi zorunludur"}) // Username parameter is required
		return
	}

	// Parse the requested role from request body
	// Using pointer to int (*int) to distinguish between "not provided" and "value is 0"
	var input struct {
		Role *int `json:"role" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if Role was provided
	if input.Role == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Yetki seviyesi zorunludur"}) // Role level is required
		return
	}

	// Validate the role value (must be 0 or 1)
	if *input.Role != 0 && *input.Role != 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Yetki seviyesi 0 veya 1 olmalıdır"}) // Role must be 0 or 1
		return
	}

	// Use collation option for case-insensitive search
	opts := options.Update().SetCollation(&options.Collation{
		Locale:   "en",
		Strength: 2, // Case-insensitive comparison
	})

	// Update the user's role
	update := bson.M{
		"$set": bson.M{
			"role": *input.Role,
		},
	}

	result, err := userCollection.UpdateOne(ctx, bson.M{"username": username}, update, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kullanıcı bulunamadı"}) // User not found
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Kullanıcı yetkisi güncellendi"}) // User role updated
}

// GetUserActivities retrieves all activity records for a specific user
func GetUserActivities(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get username from URL parameter
	username := c.Param("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kullanıcı adı parametresi zorunludur"})
		return
	}

	// Check that the user exists before proceeding
	err := userCollection.FindOne(ctx, bson.M{"username": username}).Err()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kullanıcı bulunamadı"})
		return
	}

	// Set up options to sort by updatedAt in descending order
	findOptions := options.FindOne().SetSort(bson.M{"updatedAt": -1})

	// Find the user's activity record
	var userActivity middleware.UserActivity
	err = activityCollection.FindOne(ctx, bson.M{"username": username}, findOptions).Decode(&userActivity)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			c.JSON(http.StatusOK, gin.H{"message": "Kullanıcı için aktivite kaydı bulunamadı"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Aktivite kayıtları alınırken bir hata oluştu"})
		return
	}

	c.JSON(http.StatusOK, userActivity)
}
