package routes

import (
	"github.com/alongkornn/Web-VRGame-Backend/internal/checkpoint/controllers"
	"github.com/labstack/echo/v4"
)

func CheckpointRoute(g *echo.Group) {
	checkpointGroup := g.Group("/checkpoint")

	// แสดงด่านปัจจุบันโดยเข้าถึงผ่านไอดีผู้เล่น
	checkpointGroup.GET("/current/:userId", controllers.GetCurrentCheckpointFromUser)
	// แสดงทุกด่านทุกหมวดหมู่
	checkpointGroup.GET("/", controllers.GetAllCheckpoint)
	// สร้างด่านใหม่
	checkpointGroup.POST("/create", controllers.CreateCheckpoint)
	// บันทึกด่านปัจจุบันลงในด่านที่เล่นผ่านแล้ว
	checkpointGroup.POST("/save/complete/:score/:userId", controllers.SaveCheckpointToComplete)
	// แสดงด่านที่เล่นผ่านโดยเข้าถึงผ่านไอดีผู้เล่น
	checkpointGroup.GET("/complete/:userId", controllers.GetCompleteCheckpointByUserId)
	// แสดงด่านทุกด่านตามหมวดหมู่
	checkpointGroup.GET("/category", controllers.GetCheckpointWithCategory)
}