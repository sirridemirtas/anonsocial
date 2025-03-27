package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirridemirtas/anonsocial/config"
	"github.com/sirridemirtas/anonsocial/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var conversationCollection *mongo.Collection

func SetConversationCollection(client *mongo.Client) {
	conversationCollection = client.Database(config.AppConfig.MongoDB_DB).Collection("conversations")

	// Try to drop all existing indexes to start fresh
	_, _ = conversationCollection.Indexes().DropAll(context.Background())

	// Create a unique index on the participantKey field
	_, err := conversationCollection.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys: bson.D{
				{Key: "participantKey", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
	)

	if err != nil {
		panic(err)
	}
}

// GetConversation retrieves a conversation between the current user and another user
func GetConversation(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get current authenticated user from context (set by Auth middleware)
	currentUser := c.GetString("username")
	if currentUser == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kullanıcı adı bulunamadı"})
		return
	}

	// Get target user from URL parameter
	targetUser := c.Param("username")

	// Validate that current user is not messaging themselves
	if currentUser == targetUser {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kendinize mesaj gönderemezsiniz"})
		return
	}

	// Create participant key for finding the conversation
	participantKey := models.CreateParticipantKey(currentUser, targetUser)

	// Find conversation using the participantKey instead of using $all on participants array
	// This ensures we consistently find the same conversation regardless of participant order
	var conversation models.Conversation
	err := conversationCollection.FindOne(ctx, bson.M{
		"participantKey": participantKey,
	}).Decode(&conversation)

	// If conversation doesn't exist, return empty conversation
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusOK, models.NewConversation(currentUser, targetUser))
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Check if current user has deleted this conversation
	if conversation.IsDeletedBy(currentUser) {
		c.JSON(http.StatusGone, gin.H{"error": "Bu görüşme silinmiş"})
		return
	}

	// REMOVED: No longer auto-marking messages as read when viewing a conversation
	// Let the explicit /messages/:username/read endpoint handle this

	c.JSON(http.StatusOK, conversation)
}

// SendMessage sends a message from the current user to another user
func SendMessage(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get current authenticated user from context
	currentUser := c.GetString("username")
	if currentUser == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kullanıcı adı bulunamadı"})
		return
	}

	// Get target user from URL parameter
	targetUser := c.Param("username")

	// Validate that current user is not messaging themselves
	if currentUser == targetUser {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kendinize mesaj gönderemezsiniz"})
		return
	}

	// Check if target user exists
	var targetUserDoc models.User
	err := userCollection.FindOne(ctx, bson.M{"username": targetUser}).Decode(&targetUserDoc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Mesaj göndermek istediğiniz kullanıcı bulunamadı"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Kullanıcı kontrolü sırasında bir hata oluştu"})
		}
		return
	}

	// Parse request body
	var request struct {
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate message content length
	if len(request.Content) > models.MaxMessageLength {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Mesaj 500 karakterden uzun olamaz"})
		return
	}

	now := time.Now()

	// Create participant key for finding the conversation
	participantKey := models.CreateParticipantKey(currentUser, targetUser)

	// Try to find existing conversation using the participantKey
	var conversation models.Conversation
	err = conversationCollection.FindOne(ctx, bson.M{
		"participantKey": participantKey,
	}).Decode(&conversation)

	// If conversation doesn't exist, create a new one
	if err == mongo.ErrNoDocuments {
		conversation = *models.NewConversation(currentUser, targetUser)
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Check if current user has deleted this conversation
	if conversation.IsDeletedBy(currentUser) {
		// Remove current user from deletedBy list if they're sending a new message
		for i, user := range conversation.DeletedBy {
			if user == currentUser {
				conversation.DeletedBy = append(conversation.DeletedBy[:i], conversation.DeletedBy[i+1:]...)
				break
			}
		}
	}

	// Add message to conversation
	err = conversation.AddMessage(currentUser, request.Content)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Prepare update or insert operation
	if conversation.ID.IsZero() {
		// Insert new conversation
		result, err := conversationCollection.InsertOne(ctx, conversation)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		conversation.ID = result.InsertedID.(primitive.ObjectID)
	} else {
		// Update existing conversation
		update := bson.M{
			"$set": bson.M{
				"messages":                   conversation.Messages,
				"lastUpdated":                now,
				"unreadCounts." + targetUser: conversation.UnreadCounts[targetUser],
			},
		}

		// If user was previously in deletedBy, remove them
		if len(conversation.DeletedBy) > 0 {
			update["$pull"] = bson.M{"deletedBy": currentUser}
		}

		_, err = conversationCollection.UpdateOne(
			ctx,
			bson.M{"_id": conversation.ID},
			update,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, conversation)
}

// DeleteConversation marks a conversation as deleted for the current user
func DeleteConversation(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get current authenticated user from context
	currentUser := c.GetString("username")
	if currentUser == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kullanıcı adı bulunamadı"})
		return
	}

	// Get target user from URL parameter
	targetUser := c.Param("username")

	// Validate that current user is not messaging themselves
	if currentUser == targetUser {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kendinizle olan mesajları silemezsiniz"})
		return
	}

	// Create participant key for finding the conversation
	participantKey := models.CreateParticipantKey(currentUser, targetUser)

	// Find conversation using the participantKey
	var conversation models.Conversation
	err := conversationCollection.FindOne(ctx, bson.M{
		"participantKey": participantKey,
	}).Decode(&conversation)

	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusNotFound, gin.H{"error": "Görüşme bulunamadı"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Check if current user has already deleted this conversation
	if conversation.IsDeletedBy(currentUser) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bu görüşme zaten silinmiş"})
		return
	}

	// Add current user to deletedBy array
	_, err = conversationCollection.UpdateOne(
		ctx,
		bson.M{"_id": conversation.ID},
		bson.M{
			"$addToSet": bson.M{"deletedBy": currentUser},
		},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// If both users have deleted the conversation, actually delete it
	if len(conversation.DeletedBy) == 1 && conversation.DeletedBy[0] != currentUser {
		_, err = conversationCollection.DeleteOne(ctx, bson.M{"_id": conversation.ID})
		if err != nil {
			// Just log the error but return success to the user
			// c.Logger().Error(err)
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Görüşme silindi"})
}

// GetConversationList retrieves a list of all conversations for the current user
// Only returns the most recent message for each conversation
func GetConversationList(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get current authenticated user from context
	currentUser := c.GetString("username")
	if currentUser == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kullanıcı adı bulunamadı"})
		return
	}

	// Find all conversations where current user is a participant and hasn't deleted the conversation
	// Use projection to only include the last message for each conversation
	findOptions := options.Find().
		SetProjection(bson.M{
			"_id":            1,
			"participants":   1,
			"participantKey": 1,
			"createdAt":      1,
			"lastUpdated":    1,
			"deletedBy":      1,
			"unreadCounts":   1,
			// Use $slice to get only the last message (-1 means the last element)
			"messages": bson.M{"$slice": -1},
		}).
		SetSort(bson.D{{Key: "lastUpdated", Value: -1}}) // Sort by lastUpdated in descending order

	cursor, err := conversationCollection.Find(ctx, bson.M{
		"participants": currentUser,
		"deletedBy": bson.M{
			"$ne": currentUser,
		},
	}, findOptions)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(ctx)

	var conversations []models.Conversation
	if err = cursor.All(ctx, &conversations); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, conversations)
}

// MarkConversationAsRead marks all messages in a conversation as read for the authenticated user
func MarkConversationAsRead(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get current authenticated user from context
	currentUser := c.GetString("username")
	if currentUser == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kullanıcı adı bulunamadı"})
		return
	}

	// Get target user from URL parameter
	targetUser := c.Param("username")

	// Validate that current user is not messaging themselves
	if currentUser == targetUser {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kendinize mesaj gönderemezsiniz"})
		return
	}

	// Check if target user exists
	var targetUserDoc models.User
	err := userCollection.FindOne(ctx, bson.M{"username": targetUser}).Decode(&targetUserDoc)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Belirtilen kullanıcı bulunamadı"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Kullanıcı kontrolü sırasında bir hata oluştu"})
		}
		return
	}

	// Create participant key for finding the conversation
	participantKey := models.CreateParticipantKey(currentUser, targetUser)

	// Find conversation using the participantKey
	var conversation models.Conversation
	err = conversationCollection.FindOne(ctx, bson.M{
		"participantKey": participantKey,
	}).Decode(&conversation)

	// If conversation doesn't exist, return success (nothing to mark as read)
	if err == mongo.ErrNoDocuments {
		c.JSON(http.StatusOK, gin.H{"message": "Okunacak mesaj bulunamadı"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Check if current user has deleted this conversation
	if conversation.IsDeletedBy(currentUser) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bu görüşme silinmiş"})
		return
	}

	// Update the conversation in the database by setting unread count to 0 for current user
	_, err = conversationCollection.UpdateOne(
		ctx,
		bson.M{"_id": conversation.ID},
		bson.M{"$set": bson.M{"unreadCounts." + currentUser: 0}},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Mesajlar okundu olarak işaretlendi"})
}

// GetTotalUnreadCount gets the total number of unread messages for the authenticated user
func GetTotalUnreadCount(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get current authenticated user from context
	currentUser := c.GetString("username")
	if currentUser == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kullanıcı adı bulunamadı"})
		return
	}

	// Find all conversations where current user is a participant and hasn't deleted the conversation
	cursor, err := conversationCollection.Find(ctx, bson.M{
		"participants": currentUser,
		"deletedBy": bson.M{
			"$ne": currentUser,
		},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(ctx)

	var conversations []models.Conversation
	if err = cursor.All(ctx, &conversations); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Sum up unread counts for the current user
	totalUnread := 0
	for _, conversation := range conversations {
		totalUnread += conversation.UnreadCounts[currentUser]
	}

	c.JSON(http.StatusOK, gin.H{
		"unreadCount": totalUnread,
	})
}
