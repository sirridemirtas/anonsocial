package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// NotificationType represents the different types of notifications
type NotificationType string

const (
	NotificationTypeReply        NotificationType = "reply"          // Someone replied to user's post
	NotificationTypeReplyToReply NotificationType = "reply_to_reply" // Someone replied to a post user also replied to
	NotificationTypeReaction     NotificationType = "reaction"       // Someone reacted to user's post
)

// Notification represents a user notification
type Notification struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Username     string             `bson:"username" json:"username"`                             // User who receives the notification
	PostID       primitive.ObjectID `bson:"postId" json:"postId"`                                 // Related post ID
	PostSnippet  string             `bson:"postSnippet" json:"postSnippet"`                       // First 50 chars of post content
	Type         NotificationType   `bson:"type" json:"type"`                                     // Type of notification
	LikeCount    int                `bson:"likeCount,omitempty" json:"likeCount,omitempty"`       // For reaction type
	DislikeCount int                `bson:"dislikeCount,omitempty" json:"dislikeCount,omitempty"` // For reaction type
	Read         bool               `bson:"read" json:"read"`                                     // Whether notification has been read
	CreatedAt    time.Time          `bson:"createdAt" json:"createdAt"`                           // When notification was created
	UpdatedAt    time.Time          `bson:"updatedAt" json:"updatedAt"`                           // When notification was last updated
}
