package routes

import (
	"github.com/alongkornn/Web-VRGame-Backend/internal/admin/controllers"
	"github.com/labstack/echo/v4"
)

func AdminRoute(g *echo.Group) {
	adminGroup := g.Group("/admin")

	adminGroup.POST("/create/:id", controllers.CreateAdmin)
	adminGroup.POST("/get/:id", controllers.GetAdminByID)
	adminGroup.POST("/get", controllers.GetAllAdmin)
	adminGroup.PUT("/approve/:id", controllers.ApprovedRegister)
	adminGroup.DELETE("/delete/user/:id", controllers.RemoveUser)
	adminGroup.DELETE("/delete/admin/:id", controllers.RemoveAdmin)
	adminGroup.PUT("/update/admin/:id", controllers.UpdateDataAdmin)
	adminGroup.PUT("/updatepassword/admin/:id", controllers.UpdatePasswordAdmin)

	// protectedGroup := g.Group("")
	// adminGroup := protectedGroup.Group("/admin")
    // adminGroup.Use(middleware.RoleBasedMiddleware("admin"))
}