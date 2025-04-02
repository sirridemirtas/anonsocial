package controllers

import (
	"encoding/xml"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/sirridemirtas/anonsocial/data"
	"github.com/sirridemirtas/anonsocial/models"
)

var sitemapPostCollection *mongo.Collection

// SetSitemapPostCollection sets the post collection for the sitemap controller
func SetSitemapPostCollection(client *mongo.Client) {
	sitemapPostCollection = client.Database("anonsocial").Collection("posts")
	models.SetPostCollectionForSitemap(client)
}

// GenerateSitemapXML handles the request for sitemap.xml
func GenerateSitemapXML(w http.ResponseWriter, r *http.Request) {
	baseURL := "https://dedimki.com"

	// Create sitemap structure
	sitemap := models.Sitemap{
		Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs:  []models.SitemapURL{},
	}

	defaultLastMod := "2025-04-02"

	homepageLastMod := models.GetLatestPostDate()
	sitemap.URLs = append(sitemap.URLs, models.SitemapURL{
		Loc:     baseURL,
		LastMod: homepageLastMod,
	})

	otherRoutes := []string{
		"/contact",
		"/login",
		"/register",
		"/university",
	}

	for _, route := range otherRoutes {
		sitemap.URLs = append(sitemap.URLs, models.SitemapURL{
			Loc:     baseURL + route,
			LastMod: defaultLastMod,
		})
	}

	for _, univ := range data.Universities {
		lastMod := models.GetLatestPostDateForUniversity(univ.ID)

		sitemap.URLs = append(sitemap.URLs, models.SitemapURL{
			Loc:     baseURL + "/university/" + univ.ID,
			LastMod: lastMod,
		})
	}

	w.Header().Set("Content-Type", "application/xml")
	w.Write([]byte(xml.Header))
	enc := xml.NewEncoder(w)
	enc.Indent("", "  ")
	if err := enc.Encode(sitemap); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
