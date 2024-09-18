package routes

import (
	"github.com/alongkornn/Web-VRGame-Backend/internal/user/controllers"
	"github.com/labstack/echo/v4"
)

func UserRoute(e *echo.Group) {
	userGroup := e.Group("/user")

	userGroup.POST("/get/:id", controllers.GetUserByID)
	userGroup.POST("/add/checkpoint/:checkpointID/:userID", controllers.AddPlayerInCheckpoint)
	userGroup.POST("/get", controllers.GetAllUser)
	userGroup.POST("/get/pending", controllers.GetUserPending)
	userGroup.PUT("/update/:id", controllers.UpdateUser)

}