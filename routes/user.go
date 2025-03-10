package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/sirridemirtas/anonsocial/controllers"
	"github.com/sirridemirtas/anonsocial/middleware"
)

func UserRoutes(rg *gin.RouterGroup) {
	userGroup := rg.Group("/users")
	{
		userGroup.GET("/", controllers.GetUsers)
		userGroup.GET("/:id", controllers.GetUser)
		userGroup.GET("/check-username/:username", controllers.CheckUsernameAvailability)
		userGroup.PUT("/:id", middleware.Auth(0), controllers.UpdateUser)    // Normal user
		userGroup.DELETE("/:id", middleware.Auth(1), controllers.DeleteUser) // Admin only
	}
}
