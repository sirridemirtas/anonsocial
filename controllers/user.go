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

	"github.com/sirridemirtas/anonsocial/config"
	"github.com/sirridemirtas/anonsocial/middleware"
	"github.com/sirridemirtas/anonsocial/models"
	"github.com/sirridemirtas/anonsocial/utils"
)

var userCollection *mongo.Collection
var avatarCollection *mongo.Collection

func SetUserCollection(client *mongo.Client) {
	userCollection = client.Database(config.AppConfig.MongoDB_DB).Collection("users")
	avatarCollection = client.Database(config.AppConfig.MongoDB_DB).Collection("avatars")

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

	// Create unique index for avatar username
	_, err = avatarCollection.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys: bson.D{{Key: "username", Value: 1}},
			Options: options.Index().
				SetUnique(true).
				SetName("username_unique").
				SetCollation(&options.Collation{
					Locale:   "en",
					Strength: 2, // Case-insensitive comparison
				}),
		},
	)

	if err != nil {
		log.Fatal("Error creating unique index for avatar username:", err)
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kullanıcı adı parametresi zorunludur"}) // Username parameter is required
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
		c.JSON(http.StatusNotFound, gin.H{"error": "Kullanıcı bulunamadı"}) // User not found
		return
	}

	c.JSON(http.StatusOK, user)
}

func UpdateUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	id, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz ID"}) // Invalid ID
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
		c.JSON(http.StatusNotFound, gin.H{"error": "Kullanıcı bulunamadı"}) // User not found
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Kullanıcı güncellendi"}) // User updated
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
		c.JSON(http.StatusNotFound, gin.H{"error": "Kullanıcı bulunamadı"}) // User not found
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Kullanıcı silindi"}) // User deleted
}

func CheckUsernameAvailability(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	username := c.Param("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kullanıcı adı parametresi zorunludur"}) // Username parameter is required
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
			"message":   "Bu kullanıcı adı daha önce alınmış", // This username is already taken
		})
	} else {
		// Username is available
		c.JSON(http.StatusOK, gin.H{
			"available": true,
			"valid":     true,
			"message":   "Kullanıcı adı uygun", // Username is available
		})
	}
}

// UpdateUserPrivacy updates the isPrivate field for the authenticated user
func UpdateUserPrivacy(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get username from token that was set by the Auth middleware
	username := c.GetString("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token içinde kullanıcı adı bulunamadı"}) // Username not found in token
		return
	}

	var input struct {
		IsPrivate bool `json:"isPrivate"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	update := bson.M{
		"$set": bson.M{
			"isPrivate": input.IsPrivate,
		},
	}

	// Use collation option for case-insensitive search
	opts := options.Update().SetCollation(&options.Collation{
		Locale:   "en",
		Strength: 2, // Case-insensitive comparison
	})

	result, err := userCollection.UpdateOne(ctx, bson.M{"username": username}, update, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kullanıcı bulunamadı"}) // User not found
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Kullanıcı gizlilik ayarı güncellendi", // User privacy setting updated
		"isPrivate": input.IsPrivate,
	})
}

// GetUserAvatar retrieves a user's avatar
func GetUserAvatar(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get the username from the URL parameter
	targetUsername := c.Param("username")
	if targetUsername == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kullanıcı adı parametresi zorunludur"}) // Username parameter is required
		return
	}

	// Get the requesting user's username (if authenticated)
	requestingUsername := getUsernameFromUserRequest(c)

	// First check if the target user exists
	var user models.User
	err := userCollection.FindOne(ctx, bson.M{"username": targetUsername}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kullanıcı bulunamadı"}) // User not found
		return
	}

	// If the user is private and the requester is not the same user, return 403 Forbidden
	if user.IsPrivate && targetUsername != requestingUsername {
		c.JSON(http.StatusForbidden, gin.H{"error": "Bu kullanıcının profili gizlidir"}) // This user's profile is private
		return
	}

	// Use collation option for case-insensitive search
	opts := options.FindOne().SetCollation(&options.Collation{
		Locale:   "en",
		Strength: 2, // Case-insensitive comparison
	})

	// Retrieve the avatar
	var avatar models.Avatar
	err = avatarCollection.FindOne(ctx, bson.M{"username": targetUsername}, opts).Decode(&avatar)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Avatar bulunamadı"}) // Avatar not found
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create response without username field
	response := gin.H{
		"faceColor":    avatar.FaceColor,
		"earSize":      avatar.EarSize,
		"hairStyle":    avatar.HairStyle,
		"hairColor":    avatar.HairColor,
		"hatStyle":     avatar.HatStyle,
		"hatColor":     avatar.HatColor,
		"eyeStyle":     avatar.EyeStyle,
		"glassesStyle": avatar.GlassesStyle,
		"noseStyle":    avatar.NoseStyle,
		"mouthStyle":   avatar.MouthStyle,
		"shirtStyle":   avatar.ShirtStyle,
		"shirtColor":   avatar.ShirtColor,
		"bgColor":      avatar.BgColor,
	}

	c.JSON(http.StatusOK, response)
}

// UpdateUserAvatar creates or updates a user's avatar
func UpdateUserAvatar(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get the target username from the URL parameter
	targetUsername := c.Param("username")
	if targetUsername == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kullanıcı adı parametresi zorunludur"}) // Username parameter is required
		return
	}

	// Get the authenticated user's username
	authUsername := c.GetString("username")
	if authUsername == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Bu işlem için giriş yapmalısınız"}) // You must be logged in for this operation
		return
	}

	// Only allow users to update their own avatar
	if targetUsername != authUsername {
		c.JSON(http.StatusForbidden, gin.H{"error": "Sadece kendi avatarınızı güncelleyebilirsiniz"}) // You can only update your own avatar
		return
	}

	// Parse and validate the avatar data
	var avatar models.Avatar
	if err := c.ShouldBindJSON(&avatar); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate the avatar data
	validationErrors := utils.ValidateAvatar(avatar)
	if len(validationErrors) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"errors": validationErrors})
		return
	}

	// Set the username to ensure it matches the authenticated user
	avatar.Username = authUsername

	// Use collation option for case-insensitive search
	opts := options.Update().SetCollation(&options.Collation{
		Locale:   "en",
		Strength: 2, // Case-insensitive comparison
	}).SetUpsert(true) // Create if not exists

	// Update or create the avatar
	update := bson.M{
		"$set": bson.M{
			"username":     avatar.Username,
			"faceColor":    avatar.FaceColor,
			"earSize":      avatar.EarSize,
			"hairStyle":    avatar.HairStyle,
			"hairColor":    avatar.HairColor,
			"hatStyle":     avatar.HatStyle,
			"hatColor":     avatar.HatColor,
			"eyeStyle":     avatar.EyeStyle,
			"glassesStyle": avatar.GlassesStyle,
			"noseStyle":    avatar.NoseStyle,
			"mouthStyle":   avatar.MouthStyle,
			"shirtStyle":   avatar.ShirtStyle,
			"shirtColor":   avatar.ShirtColor,
			"bgColor":      avatar.BgColor,
		},
	}

	result, err := avatarCollection.UpdateOne(ctx, bson.M{"username": authUsername}, update, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.MatchedCount > 0 {
		c.JSON(http.StatusOK, gin.H{"message": "Avatar güncellendi"}) // Avatar updated
	} else {
		c.JSON(http.StatusCreated, gin.H{"message": "Avatar oluşturuldu"}) // Avatar created
	}
}

// ResetPassword resets the password for the authenticated user
func ResetPassword(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get username from token that was set by the Auth middleware
	username := c.GetString("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Token içinde kullanıcı adı bulunamadı"}) // Username not found in token
		return
	}

	var input struct {
		CurrentPassword string `json:"currentPassword" binding:"required"`
		NewPassword     string `json:"newPassword" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find the user
	var user models.User
	err := userCollection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kullanıcı bulunamadı"}) // User not found
		return
	}

	// Verify current password
	if !user.ValidatePassword(input.CurrentPassword) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Mevcut şifre yanlış"}) // Current password is incorrect
		return
	}

	// Generate new salt and hash password
	user.Salt = models.GenerateSalt()
	user.Password = user.HashPassword(input.NewPassword)

	// Update the user in the database
	update := bson.M{
		"$set": bson.M{
			"password": user.Password,
			"salt":     user.Salt,
		},
	}

	// Use collation option for case-insensitive search
	opts := options.Update().SetCollation(&options.Collation{
		Locale:   "en",
		Strength: 2, // Case-insensitive comparison
	})

	result, err := userCollection.UpdateOne(ctx, bson.M{"username": username}, update, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kullanıcı bulunamadı"}) // User not found
		return
	}

	// Clear the auth cookie to log the user out
	c.SetCookie("token", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "Sıfırlama işlemi başarılı, yeni şifrenizle giriş yapabilirsiniz"})
}

// Helper function to get the username from the request
func getUsernameFromUserRequest(c *gin.Context) string {
	// First try to get from context (set by Auth middleware)
	username := c.GetString("username")
	if username != "" {
		return username
	}

	// If not in context, try to get from token
	cookie, err := c.Cookie("token")
	if err == nil {
		username, valid := middleware.GetUsernameFromToken(cookie)
		if valid {
			return username
		}
	}

	return ""
}
