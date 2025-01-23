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
			c.String(http.StatusMethodNotAllowed, "405 method not allowed")
			return
		}

		path := c.Request.URL.Path
		if strings.HasPrefix(path, "/api") {
			c.String(http.StatusNotFound, "404 API endpoint not found")
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

		c.String(http.StatusNotFound, "404 page not found")
	})
}
