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

var notificationCollection *mongo.Collection

func SetNotificationCollection(client *mongo.Client) {
	notificationCollection = client.Database(config.AppConfig.MongoDB_DB).Collection("notifications")

	// Create index for username field for faster queries
	_, err := notificationCollection.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys: bson.D{{Key: "username", Value: 1}},
		},
	)

	if err != nil {
		panic(err)
	}

	// Create compound index for username + postId + type for uniqueness
	_, err = notificationCollection.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys: bson.D{
				{Key: "username", Value: 1},
				{Key: "postId", Value: 1},
				{Key: "type", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
	)

	if err != nil {
		panic(err)
	}
}

// GetNotifications returns the authenticated user's notifications
func GetNotifications(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	username := c.GetString("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kullanıcı adı bulunamadı"})
		return
	}

	// Find all notifications for the user, sort by read status then by updatedAt
	findOptions := options.Find().
		SetSort(bson.D{
			{Key: "read", Value: 1},       // Unread (false) first
			{Key: "updatedAt", Value: -1}, // Newest first
		}).
		SetLimit(50).                        // Limit to 50 notifications
		SetProjection(bson.M{"username": 0}) // Exclude username field from results

	cursor, err := notificationCollection.Find(ctx, bson.M{"username": username}, findOptions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(ctx)

	var notifications []models.Notification
	if err = cursor.All(ctx, &notifications); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, notifications)
}

// GetUnreadCount returns the count of unread notifications for the authenticated user
func GetUnreadCount(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	username := c.GetString("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kullanıcı adı bulunamadı"})
		return
	}

	// Count unread notifications
	count, err := notificationCollection.CountDocuments(ctx, bson.M{
		"username": username,
		"read":     false,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"unreadCount": count})
}

// MarkAsRead marks a notification as read
func MarkAsRead(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	username := c.GetString("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kullanıcı adı bulunamadı"})
		return
	}

	notificationID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Geçersiz bildirim kimliği"})
		return
	}

	// Update notification to mark as read
	result, err := notificationCollection.UpdateOne(
		ctx,
		bson.M{
			"_id":      notificationID,
			"username": username, // Ensure the notification belongs to the user
		},
		bson.M{
			"$set": bson.M{
				"read": true,
			},
		},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.MatchedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Bildirim bulunamadı veya erişim izniniz yok"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Bildirim okundu olarak işaretlendi"})
}

// MarkAllAsRead marks all notifications as read for the authenticated user
func MarkAllAsRead(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	username := c.GetString("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kullanıcı adı bulunamadı"})
		return
	}

	// Update all notifications to mark them as read
	result, err := notificationCollection.UpdateMany(
		ctx,
		bson.M{
			"username": username,
			"read":     false,
		},
		bson.M{
			"$set": bson.M{
				"read": true,
			},
		},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tüm bildirimler okundu olarak işaretlendi", "modifiedCount": result.ModifiedCount})
}

// DeleteAllNotifications deletes all notifications for the authenticated user
func DeleteAllNotifications(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	username := c.GetString("username")
	if username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kullanıcı adı bulunamadı"})
		return
	}

	// Delete all notifications for the user (both read and unread)
	result, err := notificationCollection.DeleteMany(
		ctx,
		bson.M{
			"username": username,
		},
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Bildirimler silinirken bir hata oluştu: " + err.Error()})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusOK, gin.H{"message": "Silinecek bildirim bulunamadı"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tüm bildirimler silindi", "deletedCount": result.DeletedCount})
}

// CreateOrUpdateReactionNotification handles notifications for reactions (likes/dislikes)
func CreateOrUpdateReactionNotification(postID primitive.ObjectID, postOwner string, postContent string, isLike bool) {
	// Don't notify if the content owner reacts to their own content
	if postOwner == "" {
		return
	}

	// Create snippet
	snippet := createSnippet(postContent)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	now := time.Now()

	// Find if notification already exists
	var notification models.Notification
	err := notificationCollection.FindOne(
		ctx,
		bson.M{
			"username": postOwner,
			"postId":   postID,
			"type":     models.NotificationTypeReaction,
		},
	).Decode(&notification)

	// Create or update notification
	if err == mongo.ErrNoDocuments {
		// Create new notification
		notification = models.Notification{
			Username:     postOwner,
			PostID:       postID,
			PostSnippet:  snippet,
			Type:         models.NotificationTypeReaction,
			LikeCount:    0,
			DislikeCount: 0,
			Read:         false,
			CreatedAt:    now,
			UpdatedAt:    now,
		}

		// Set initial counts
		if isLike {
			notification.LikeCount = 1
		} else {
			notification.DislikeCount = 1
		}

		_, err := notificationCollection.InsertOne(ctx, notification)
		if err != nil {
			// Just log error, don't fail the main operation
			return
		}
	} else if err == nil {
		// Update existing notification
		update := bson.M{
			"$set": bson.M{
				"updatedAt": now,
				"read":      false, // Mark as unread when updated
			},
		}

		// Increment appropriate counter
		if isLike {
			update["$inc"] = bson.M{"likeCount": 1}
		} else {
			update["$inc"] = bson.M{"dislikeCount": 1}
		}

		_, err := notificationCollection.UpdateOne(
			ctx,
			bson.M{
				"_id": notification.ID,
			},
			update,
		)

		if err != nil {
			// Just log error, don't fail the main operation
			return
		}
	}

	// Cleanup old notifications
	go CleanupOldNotifications(postOwner)
}

// CreateOrUpdateReplyNotification handles notifications for replies
func CreateOrUpdateReplyNotification(postID primitive.ObjectID, postOwner string, postContent string, replyOwner string, isReplyToReply bool) {
	// Don't notify if user replies to their own content
	if postOwner == replyOwner {
		return
	}

	// Create snippet
	snippet := createSnippet(postContent)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	now := time.Now()

	// Determine notification type
	notificationType := models.NotificationTypeReply
	if isReplyToReply {
		notificationType = models.NotificationTypeReplyToReply
	}

	// Find if notification already exists
	var notification models.Notification
	err := notificationCollection.FindOne(
		ctx,
		bson.M{
			"username": postOwner,
			"postId":   postID,
			"type":     notificationType,
		},
	).Decode(&notification)

	// Create or update notification
	if err == mongo.ErrNoDocuments {
		// Create new notification
		notification = models.Notification{
			Username:    postOwner,
			PostID:      postID,
			PostSnippet: snippet,
			Type:        notificationType,
			Read:        false,
			CreatedAt:   now,
			UpdatedAt:   now,
		}

		_, err := notificationCollection.InsertOne(ctx, notification)
		if err != nil {
			// Just log error, don't fail the main operation
			return
		}
	} else if err == nil {
		// If notification was read and the same user replied again, mark as unread
		// Otherwise, just update the timestamp
		update := bson.M{
			"$set": bson.M{
				"updatedAt": now,
				"read":      false, // Mark as unread when new reply comes in
			},
		}

		_, err := notificationCollection.UpdateOne(
			ctx,
			bson.M{
				"_id": notification.ID,
			},
			update,
		)

		if err != nil {
			// Just log error, don't fail the main operation
			return
		}
	}

	// Cleanup old notifications
	go CleanupOldNotifications(postOwner)
}

// Helper function to create a snippet of text (first 50 chars)
func createSnippet(content string) string {
	maxLength := 50
	if len(content) <= maxLength {
		return content
	}
	return content[:maxLength] + "..."
}

// CleanupOldNotifications removes notifications beyond the limit for a user
// This should be called periodically or after creating new notifications
func CleanupOldNotifications(username string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Find all notifications for the user, sorted by updatedAt
	findOptions := options.Find().
		SetSort(bson.D{{Key: "updatedAt", Value: -1}}).
		SetSkip(50).                    // Skip the first 50 (newest)
		SetProjection(bson.M{"_id": 1}) // Only get IDs

	cursor, err := notificationCollection.Find(ctx, bson.M{"username": username}, findOptions)
	if err != nil {
		return
	}
	defer cursor.Close(ctx)

	var toDelete []primitive.ObjectID
	var temp struct {
		ID primitive.ObjectID `bson:"_id"`
	}

	for cursor.Next(ctx) {
		if err := cursor.Decode(&temp); err == nil {
			toDelete = append(toDelete, temp.ID)
		}
	}

	if len(toDelete) > 0 {
		_, _ = notificationCollection.DeleteMany(ctx, bson.M{
			"_id": bson.M{"$in": toDelete},
		})
	}
}
