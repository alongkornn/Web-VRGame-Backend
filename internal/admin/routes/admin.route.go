package routes

import (
	"github.com/alongkornn/Web-VRGame-Backend/internal/admin/controllers"
	"github.com/labstack/echo/v4"
)

func AdminRoute(e *echo.Group) {
	adminGroup := e.Group("/admin")

	adminGroup.PUT("/approve/:id", controllers.ApprovedRegister)
	adminGroup.DELETE("/delete/:id", controllers.RemoveUser)
}
