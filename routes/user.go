package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/sirridemirtas/anonsocial/controllers"
	"github.com/sirridemirtas/anonsocial/middleware"
)

func UserRoutes(rg *gin.RouterGroup) {
	userGroup := rg.Group("/users")
	{
		userGroup.GET("", middleware.Auth(2), controllers.GetUsers)
		userGroup.GET("/:username", controllers.GetUser)
		userGroup.GET("/check-username/:username", middleware.CustomRateLimit(1, 3), controllers.CheckUsernameAvailability)
		//userGroup.PUT("/:id", middleware.Auth(0), controllers.UpdateUser)
		userGroup.DELETE("/:id", middleware.Auth(1), controllers.DeleteUser)
		userGroup.PUT("/privacy", middleware.Auth(0), middleware.CustomRateLimit(2, 2), controllers.UpdateUserPrivacy)
		userGroup.PUT("/password/reset", middleware.Auth(0), controllers.ResetPassword)

		// Avatar endpoints
		userGroup.GET("/:username/avatar", middleware.OptionalAuth(), controllers.GetUserAvatar)
		userGroup.POST("/:username/avatar", middleware.Auth(0), controllers.UpdateUserAvatar)
	}
}
