package models

import (
	"context"
	"encoding/xml"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SitemapURL represents a URL entry in the sitemap
type SitemapURL struct {
	XMLName xml.Name `xml:"url"`
	Loc     string   `xml:"loc"`
	LastMod string   `xml:"lastmod"`
}

// Sitemap represents the entire sitemap
type Sitemap struct {
	XMLName xml.Name     `xml:"urlset"`
	Xmlns   string       `xml:"xmlns,attr"`
	URLs    []SitemapURL `xml:"url"`
}

var postCollection *mongo.Collection

// SetPostCollectionForSitemap sets the post collection for the sitemap model
func SetPostCollectionForSitemap(client *mongo.Client) {
	postCollection = client.Database("anonsocial").Collection("posts")
}

// GetLatestPostDateForUniversity returns the formatted date of the latest post for a university
// that is not a reply to another post
func GetLatestPostDateForUniversity(universityID string) string {
	// Default date if no posts are found - simplified format YYYY-MM-DD
	defaultDate := "2025-04-02"

	// If collection isn't set, return default date
	if postCollection == nil {
		return defaultDate
	}

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Set up the query to find the latest post for this university that is not a reply
	filter := bson.M{
		"universityId": universityID,
		"replyTo":      bson.M{"$exists": false}, // Only posts that aren't replies
	}
	opts := options.FindOne().SetSort(bson.M{"createdAt": -1}) // Sort by createdAt in descending order

	// Execute the query
	var post struct {
		CreatedAt time.Time `bson:"createdAt"`
	}

	err := postCollection.FindOne(ctx, filter, opts).Decode(&post)
	if err != nil {
		return defaultDate
	}

	// Format the timestamp with just year-month-day
	return post.CreatedAt.Format("2006-01-02")
}

// GetLatestPostDate returns the formatted date of the latest post across all universities
// that is not a reply to another post
func GetLatestPostDate() string {
	// Default date if no posts are found - simplified format YYYY-MM-DD
	defaultDate := "2025-04-02"

	if postCollection == nil {
		return defaultDate
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Set up the query to find the latest post across all universities that is not a reply
	filter := bson.M{
		"replyTo": bson.M{"$exists": false}, // Only posts that aren't replies
	}
	opts := options.FindOne().SetSort(bson.M{"createdAt": -1}) // Sort by createdAt in descending order

	// Execute the query
	var post struct {
		CreatedAt time.Time `bson:"createdAt"`
	}

	err := postCollection.FindOne(ctx, filter, opts).Decode(&post)
	if err != nil {
		return defaultDate
	}

	// Format the timestamp with just year-month-day
	return post.CreatedAt.Format("2006-01-02")
}
