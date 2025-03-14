package routes

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func StaticRoutes(router *gin.Engine) {
	router.NoRoute(func(c *gin.Context) {
		if c.Request.Method != http.MethodGet {
			c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "405 Method Not Allowed (Kabul Edilmeyen Metot: " + c.Request.Method + ")"})
			return
		}

		path := c.Request.URL.Path
		if strings.HasPrefix(path, "/api") {
			c.JSON(http.StatusNotFound, gin.H{"error": "404 Not Found (API Erişim Noktası Bulunamadı)"})
			return
		}

		// Redirect /university/:id to static/university.html
		if strings.HasPrefix(path, "/university/") {
			c.File("./static/university.html")
			return
		}

		// Redirect /@:username to static/user.html
		if strings.HasPrefix(path, "/@") {
			c.File("./static/user.html")
			return
		}

		// Redirect /post/:id to static/post.html
		if strings.HasPrefix(path, "/post/") {
			c.File("./static/post.html")
			return
		}

		staticPath := "./static" + path
		if _, err := os.Stat(staticPath); err == nil {
			c.File(staticPath)
			return
		}

		htmlPath := staticPath + ".html"
		if _, err := os.Stat(htmlPath); err == nil {
			c.File(htmlPath)
			return
		}

		error404 := "./static/404.html"
		if _, err := os.Stat(error404); err == nil {
			c.File(error404)
			return
		}

		c.JSON(http.StatusNotFound, gin.H{"error": "404 Not Found (Sayfa Bulunamadı)"})
	})
}
