package routes

import (
	"github.com/alongkornn/Web-VRGame-Backend/internal/checkpoint/controllers"
	"github.com/labstack/echo/v4"
)

func CheckpointRoute(g *echo.Group) {
	checkpointGroup := g.Group("/checkpoint")

	checkpointGroup.POST("/set/:checkpoint/:user", controllers.GetCurrentCheckpointToUser)
	checkpointGroup.POST("/get", controllers.GetAllCheckpoint)
	checkpointGroup.POST("/create", controllers.CreateCheckpoint)
	checkpointGroup.POST("/get/category", controllers.GetCheckpointWithCategory)
}