package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/sirridemirtas/anonsocial/controllers"
	"github.com/sirridemirtas/anonsocial/middleware"
)

func AdminRoutes(rg *gin.RouterGroup) {
	admin := rg.Group("/admin")
	admin.Use(middleware.Auth(2)) // All admin routes require administrator role (2)

	// Update user role (only admins can update roles, and only to 0 or 1)
	admin.PUT("/users/:username/role", controllers.UpdateUserRole)
}
