package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/sirridemirtas/anonsocial/controllers"
	"github.com/sirridemirtas/anonsocial/middleware"
)

func AuthRoutes(rg *gin.RouterGroup) {
	auth := rg.Group("/auth")
	{
		auth.POST("/register", controllers.Register)
		auth.POST("/login", controllers.Login)
		auth.POST("/logout", middleware.Auth(0), controllers.Logout)
		auth.GET("/token-info", middleware.Auth(0), controllers.TokenInfo)
		auth.POST("/refresh-token", middleware.Auth(0), controllers.RefreshToken)
	}
}
