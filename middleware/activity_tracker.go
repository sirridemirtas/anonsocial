package middleware

import (
	"context"
	"log"
	"net"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var activityCollection *mongo.Collection

// UserActivity represents the MongoDB document structure for activity tracking
type UserActivity struct {
	Username  string    `bson:"username" json:"username"`
	IPEntries []IPEntry `bson:"ipEntries" json:"ipEntries"`
	UpdatedAt time.Time `bson:"updatedAt" json:"updatedAt"`
}

// IPEntry represents an IP address and port combination with activity history
type IPEntry struct {
	IP        string        `bson:"ip" json:"ip"`
	Port      string        `bson:"port" json:"port"`
	Actions   []ActionEntry `bson:"actions" json:"actions"`
	UpdatedAt time.Time     `bson:"updatedAt" json:"updatedAt"`
}

// ActionEntry represents a single user action
type ActionEntry struct {
	Endpoint  string    `bson:"endpoint" json:"endpoint"`
	Method    string    `bson:"method" json:"method"`
	Timestamp time.Time `bson:"timestamp" json:"timestamp"`
}

// SetActivityCollection sets the MongoDB collection for activity tracking
func SetActivityCollection(client *mongo.Client, dbName string) {
	activityCollection = client.Database(dbName).Collection("user_activities")

	// Create index on username for faster lookups
	indexModel := mongo.IndexModel{
		Keys: bson.M{"username": 1},
	}
	_, err := activityCollection.Indexes().CreateOne(context.Background(), indexModel)
	if err != nil {
		log.Printf("Error creating index on user_activities: %v", err)
	}
}

// ActivityTracker middleware records user activity on specific endpoints
func ActivityTracker() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Execute the endpoint handler first
		c.Next()

		// Get username from context or token
		username := extractUsername(c)
		if username == "" {
			return
		}

		// Get client's IP address and port
		ip := getRealIP(c)
		port := getPort(c)

		// Get endpoint and HTTP method
		endpoint := c.FullPath()
		method := c.Request.Method

		// Get current time in ISO 8601 format (UTC)
		now := time.Now().UTC()

		// Skip if collection not initialized
		if activityCollection == nil {
			log.Printf("Activity collection not initialized")
			return
		}

		// Try to update existing IP entry first
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		filter := bson.M{
			"username": username,
			"ipEntries": bson.M{
				"$elemMatch": bson.M{
					"ip":   ip,
					"port": port,
				},
			},
		}

		update := bson.M{
			"$push": bson.M{
				"ipEntries.$.actions": bson.M{
					"endpoint":  endpoint,
					"method":    method,
					"timestamp": now,
				},
			},
			"$set": bson.M{
				"ipEntries.$.updatedAt": now,
				"updatedAt":             now,
			},
		}

		result, err := activityCollection.UpdateOne(ctx, filter, update)

		if err != nil {
			log.Printf("Error updating user activity: %v", err)
			return
		}

		if result.MatchedCount == 0 {
			// If no document was updated, try to add a new IP entry
			newIPEntry := bson.M{
				"ip":   ip,
				"port": port,
				"actions": []bson.M{
					{
						"endpoint":  endpoint,
						"method":    method,
						"timestamp": now,
					},
				},
				"updatedAt": now,
			}

			filter = bson.M{"username": username}
			update = bson.M{
				"$push":        bson.M{"ipEntries": newIPEntry},
				"$set":         bson.M{"updatedAt": now},
				"$setOnInsert": bson.M{"username": username},
			}

			_, err := activityCollection.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))

			if err != nil {
				log.Printf("Error adding new IP entry: %v", err)
			}
		}
	}
}

// Extract username from context or JWT token
func extractUsername(c *gin.Context) string {
	// First try to get from context (set by Auth middleware)
	username := c.GetString("username")
	if username != "" {
		return username
	}

	// If not in context, try to get from token
	cookie, err := c.Cookie("token")
	if err == nil {
		username, valid := GetUsernameFromToken(cookie)
		if valid {
			return username
		}
	}

	return ""
}

// Get the real IP address, considering X-Forwarded-For
func getRealIP(c *gin.Context) string {
	// Check X-Forwarded-For header first
	forwardedFor := c.Request.Header.Get("X-Forwarded-For")
	if forwardedFor != "" {
		// X-Forwarded-For can contain multiple IPs; take the first one
		ips := strings.Split(forwardedFor, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// Check X-Real-IP header next
	xRealIP := c.Request.Header.Get("X-Real-IP")
	if xRealIP != "" {
		return xRealIP
	}

	// If no headers, use RemoteAddr
	remoteAddr := c.Request.RemoteAddr
	ip, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		// If can't split, just return the whole RemoteAddr
		return remoteAddr
	}
	return ip
}

// Get the port from RemoteAddr
func getPort(c *gin.Context) string {
	remoteAddr := c.Request.RemoteAddr
	_, port, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		return ""
	}
	return port
}
