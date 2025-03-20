package models

import (
	"errors"
	"sort"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Maximum allowed messages per conversation
const MaxMessagesPerConversation = 100

// Maximum message content length
const MaxMessageLength = 500

// Message represents a single message in a conversation
type Message struct {
	Sender    string    `bson:"sender" json:"sender"`
	Content   string    `bson:"content" json:"content"`
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
}

// Conversation represents a messaging conversation between two users
type Conversation struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Participants   []string           `bson:"participants" json:"participants"`
	ParticipantKey string             `bson:"participantKey" json:"-"` // Used for uniqueness indexing
	CreatedAt      time.Time          `bson:"createdAt" json:"createdAt"`
	LastUpdated    time.Time          `bson:"lastUpdated" json:"lastUpdated"`
	Messages       []Message          `bson:"messages" json:"messages"`
	DeletedBy      []string           `bson:"deletedBy" json:"deletedBy,omitempty"`
	UnreadCounts   map[string]int     `bson:"unreadCounts" json:"unreadCounts"`
}

// CreateParticipantKey creates a unique key for the participants
func CreateParticipantKey(user1, user2 string) string {
	participants := []string{user1, user2}
	sort.Strings(participants)
	return strings.Join(participants, ":")
}

// NewConversation creates a new conversation between two users
// It ensures participants are always stored in a consistent order (alphabetically)
func NewConversation(user1, user2 string) *Conversation {
	now := time.Now()

	// Sort participants alphabetically to ensure consistent ordering
	participants := []string{user1, user2}
	sort.Strings(participants)

	// Create unique key for these participants
	participantKey := CreateParticipantKey(user1, user2)

	return &Conversation{
		Participants:   participants,
		ParticipantKey: participantKey,
		CreatedAt:      now,
		LastUpdated:    now,
		Messages:       []Message{},
		DeletedBy:      []string{},
		UnreadCounts: map[string]int{
			user1: 0,
			user2: 0,
		},
	}
}

// AddMessage adds a new message to the conversation
// It also ensures the conversation has at most MaxMessagesPerConversation messages
func (c *Conversation) AddMessage(sender, content string) error {
	if len(content) > MaxMessageLength {
		return errors.New("message content exceeds maximum length of 500 characters")
	}

	if sender != c.Participants[0] && sender != c.Participants[1] {
		return errors.New("sender is not a participant in this conversation")
	}

	// Find the receiver (the other participant)
	var receiver string
	if sender == c.Participants[0] {
		receiver = c.Participants[1]
	} else {
		receiver = c.Participants[0]
	}

	// Add message
	c.Messages = append(c.Messages, Message{
		Sender:    sender,
		Content:   content,
		CreatedAt: time.Now(),
	})

	// Update conversation metadata
	c.LastUpdated = time.Now()

	// Increment unread count for receiver
	c.UnreadCounts[receiver]++

	// If we exceed the maximum number of messages, remove oldest ones
	if len(c.Messages) > MaxMessagesPerConversation {
		c.Messages = c.Messages[len(c.Messages)-MaxMessagesPerConversation:]
	}

	return nil
}

// HasParticipant checks if a user is a participant in this conversation
func (c *Conversation) HasParticipant(username string) bool {
	for _, participant := range c.Participants {
		if participant == username {
			return true
		}
	}
	return false
}

// IsDeletedBy checks if the conversation was deleted by a user
func (c *Conversation) IsDeletedBy(username string) bool {
	for _, user := range c.DeletedBy {
		if user == username {
			return true
		}
	}
	return false
}

// MarkAsRead marks all messages as read for a specific user
func (c *Conversation) MarkAsRead(username string) error {
	if !c.HasParticipant(username) {
		return errors.New("user is not a participant in this conversation")
	}

	c.UnreadCounts[username] = 0
	return nil
}
