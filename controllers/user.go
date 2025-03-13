package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/sirridemirtas/anonsocial/models"
	"github.com/sirridemirtas/anonsocial/utils"
)

var userCollection *mongo.Collection

func SetUserCollection(client *mongo.Client) {
	userCollection = client.Database("anonsocial").Collection("users")

	// First, drop the existing index if it exists
	_, err := userCollection.Indexes().DropOne(
		context.Background(),
		"username_1",
	)

	// It's okay if the index doesn't exist, so we don't check the error

	// Create unique index for username with collation for case-insensitivity
	_, err = userCollection.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys: bson.D{{Key: "username", Value: 1}},
			Options: options.Index().
				SetUnique(true).
				SetName("username_case_insensitive").
				SetCollation(&options.Collation{
					Locale:   "en",
					Strength: 2, // Case-insensitive comparison
				}),
		},
	)

	if err != nil {
		log.Fatal("Error creating unique case-insensitive index for username:", err)
	}
}

func GetUsers(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	projection := bson.M{
		"_id":      0, // Exclude ID
		"password": 0,
		"salt":     0,
	}

	var users []models.User
	cursor, err := userCollection.Find(ctx, bson.M{}, options.Find().SetProjection(projection))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &users); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, users)
}

func GetUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	username := c.Param("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username parameter is required"})
		return
	}

	projection := bson.M{
		"_id":      0, // Exclude ID
		"password": 0,
		"salt":     0,
	}

	// Use collation option for case-insensitive search
	opts := options.FindOne().SetProjection(projection).SetCollation(&options.Collation{
		Locale:   "en",
		Strength: 2, // Case-insensitive comparison
	})

	var user models.User
	err := userCollection.FindOne(ctx, bson.M{"username": username}, opts).Decode(&user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func UpdateUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	update := bson.M{
		"$set": bson.M{
			"username":     user.Username,
			"isPrivate":    user.IsPrivate,
			"role":         user.Role,
			"universityId": user.UniversityID, // Changed from university to universityId
		},
	}

	result, err := userCollection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

func DeleteUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	result, err := userCollection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func CheckUsernameAvailability(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	username := c.Param("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username parameter is required"})
		return
	}

	// Validate username using the validator utility
	validationErrors := utils.ValidateUsername(username)

	// If there are validation errors, return them
	if len(validationErrors) > 0 {
		c.JSON(http.StatusOK, gin.H{
			"available": false,
			"valid":     false,
			"message":   validationErrors[0], // Return the first error message
		})
		return
	}

	// Use collation option for case-insensitive search
	opts := options.Count().SetCollation(&options.Collation{
		Locale:   "en",
		Strength: 2, // Case-insensitive comparison
	})

	// If format is valid, check if the username already exists in the database (case-insensitive)
	count, err := userCollection.CountDocuments(ctx, bson.M{"username": username}, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if count > 0 {
		// Username is already taken
		c.JSON(http.StatusOK, gin.H{
			"available": false,
			"valid":     true,
			"message":   "Bu kullanıcı adı daha önce adı alınmış",
		})
	} else {
		// Username is available
		c.JSON(http.StatusOK, gin.H{
			"available": true,
			"valid":     true,
			"message":   "Kullanıcı adı uygun",
		})
	}
}
