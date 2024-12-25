package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/sirridemirtas/anonsocial/controllers"
	"github.com/sirridemirtas/anonsocial/middleware"
)

func AuthRoutes(rg *gin.RouterGroup) {
	auth := rg.Group("/auth")
	{
		auth.POST("/register", middleware.FormData(), controllers.Register)
		auth.POST("/login", middleware.FormData(), controllers.Login)
		auth.POST("/logout", middleware.Auth(0), controllers.Logout)
	}
}
