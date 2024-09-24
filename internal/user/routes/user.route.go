package routes

import (
	"github.com/alongkornn/Web-VRGame-Backend/internal/user/controllers"
	"github.com/labstack/echo/v4"
)

func UserRoute(e *echo.Group) {
	userGroup := e.Group("/user")

	userGroup.POST("/get/user/:userId", controllers.GetUserByID)
	userGroup.POST("/get/user", controllers.GetAllUser)
	userGroup.POST("/get/user/pending", controllers.GetUserPending)
	userGroup.PUT("/update/user/:userId", controllers.UpdateUser)
}
