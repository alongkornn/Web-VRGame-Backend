package routes

import (
	"github.com/alongkornn/Web-VRGame-Backend/internal/checkpoint/controllers"
	"github.com/labstack/echo/v4"
)

func CheckpointRoute(g *echo.Group) {
	checkpointGroup := g.Group("/checkpoint")

	// แสดงด่านปัจจุบันโดยเข้าถึงผ่านไอดีผู้เล่น
	checkpointGroup.POST("/set/:checkpoint/:user", controllers.GetCurrentCheckpointFromUser)
	// แสดงทุกด่านทุกหมวดหมู่
	checkpointGroup.POST("/get", controllers.GetAllCheckpoint)
	// สร้างด่านใหม่
	checkpointGroup.POST("/create", controllers.CreateCheckpoint)
	// บันทึกด่านปัจจุบันลงในด่านที่เล่นผ่านแล้ว
	checkpointGroup.POST("/save/checkpoint/complete/:score/:userId", controllers.SaveCheckpointToComplete)
	// แสดงด่านที่เล่นผ่านโดยเข้าถึงผ่านไอดีผู้เล่น
	checkpointGroup.POST("/get/checkpoint/complete/:userId", controllers.GetCompleteCheckpointByUserId)
	// แสดงด่านทุกด่านตามหมวดหมู่
	checkpointGroup.POST("/get/category", controllers.GetCheckpointWithCategory)
}