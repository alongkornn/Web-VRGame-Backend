package routes

import (
	"github.com/alongkornn/Web-VRGame-Backend/internal/checkpoint/controllers"
	"github.com/labstack/echo/v4"
)

func CheckpointRoute(g *echo.Group) {
	checkpointGroup := g.Group("/checkpoint")

	checkpointGroup.POST("/set/:checkpoint/:user", controllers.GetCurrentCheckpointFromUser)
	checkpointGroup.POST("/get", controllers.GetAllCheckpoint)
	checkpointGroup.POST("/create", controllers.CreateCheckpoint)
	checkpointGroup.POST("/save/checkpoint/complete/:userId", controllers.SaveCheckpointToComplete)
	checkpointGroup.POST("/get/checkpoint/complete/:userId", controllers.GetCompleteCheckpointByUserId)
}