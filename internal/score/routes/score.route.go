package routes

import (
	"github.com/alongkornn/Web-VRGame-Backend/internal/score/controllers"
	"github.com/labstack/echo/v4"
)

func ScoreRoute(e *echo.Group) {
	scoreGroup := e.Group("/score")

	scoreGroup.POST("/get/score/:userId", controllers.GetScoreByUserId)
	scoreGroup.POST("/get/score/:checkpointId", controllers.GetAllScoreByCheckpointId)
	scoreGroup.POST("/set/score/:userId", controllers.SetScore)
} 