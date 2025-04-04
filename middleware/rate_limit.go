package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// client represents a rate limited client with identification
type client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// Store for rate limiters
type limiterStore struct {
	clients map[string]*client
	mu      sync.Mutex
	// Cleanup interval time
	cleanupInterval time.Duration
	// Rate limiter configuration
	rate  rate.Limit
	burst int
}

// Create a new limiter store with cleanup goroutine
func newLimiterStore(r rate.Limit, b int, cleanupInterval time.Duration) *limiterStore {
	store := &limiterStore{
		clients:         make(map[string]*client),
		cleanupInterval: cleanupInterval,
		rate:            r,
		burst:           b,
	}

	// Start the cleanup goroutine
	go store.cleanup()

	return store
}

// getClientIdentifier creates a more granular identifier than just IP
// For authenticated users, it combines user ID with IP
// For unauthenticated users, it falls back to IP + User-Agent
func getClientIdentifier(c *gin.Context) string {
	// Get client IP
	clientIP := c.ClientIP()

	// Check if user is authenticated (has a user ID)
	userID, exists := c.Get("userID")
	if exists && userID != nil {
		// For authenticated users, use IP + UserID
		return clientIP + ":" + userID.(string)
	}

	// For unauthenticated users, use IP + User-Agent (better than just IP)
	userAgent := c.GetHeader("User-Agent")
	if userAgent == "" {
		// Fallback to just IP if User-Agent is not available
		return clientIP
	}

	// Combine IP and User-Agent for better identification
	return clientIP + ":" + userAgent
}

// Get a rate limiter for a client
func (s *limiterStore) getClientLimiter(clientIdentifier string) *rate.Limiter {
	s.mu.Lock()
	defer s.mu.Unlock()

	// If the client doesn't exist, create a new one
	c, exists := s.clients[clientIdentifier]
	if !exists {
		limiter := rate.NewLimiter(s.rate, s.burst)
		s.clients[clientIdentifier] = &client{
			limiter:  limiter,
			lastSeen: time.Now(),
		}
		return limiter
	}

	// Update last seen time
	c.lastSeen = time.Now()
	return c.limiter
}

// Cleanup removes old clients
func (s *limiterStore) cleanup() {
	for {
		time.Sleep(s.cleanupInterval)

		s.mu.Lock()
		for clientIdentifier, client := range s.clients {
			if time.Since(client.lastSeen) > s.cleanupInterval {
				delete(s.clients, clientIdentifier)
			}
		}
		s.mu.Unlock()
	}
}

// Default limiter store with 10 requests per second and burst of 20
var defaultLimiterStore = newLimiterStore(10, 20, 5*time.Minute)

// RateLimit returns a middleware that limits the number of requests per second
func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get client identifier (more specific than just IP)
		clientIdentifier := getClientIdentifier(c)

		// Get rate limiter for this client
		limiter := defaultLimiterStore.getClientLimiter(clientIdentifier)

		// Check if rate limit is exceeded
		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"status":  "error",
				"message": "Çok fazla istek gönderdiniz, lütfen daha sonra tekrar deneyin",
			})
			return
		}

		// Continue to the next middleware/handler
		c.Next()
	}
}

// CustomRateLimit returns a middleware with custom rate and burst parameters
func CustomRateLimit(r rate.Limit, b int) gin.HandlerFunc {
	// Create a new store for this custom rate limiter
	store := newLimiterStore(r, b, 5*time.Minute)

	return func(c *gin.Context) {
		// Get client identifier (more specific than just IP)
		clientIdentifier := getClientIdentifier(c)

		// Get rate limiter for this client
		limiter := store.getClientLimiter(clientIdentifier)

		// Check if rate limit is exceeded
		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"status":  "error",
				"message": "Çok fazla istek gönderdiniz, lütfen daha sonra tekrar deneyin",
			})
			return
		}

		// Continue to the next middleware/handler
		c.Next()
	}
}
