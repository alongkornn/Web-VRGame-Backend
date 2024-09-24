package routes

import (
	"github.com/alongkornn/Web-VRGame-Backend/internal/user/controllers"
	"github.com/labstack/echo/v4"
)

func UserRoute(e *echo.Group) {
	userGroup := e.Group("/user")

	// แสดงผู้เล่นแค่คนเดียว
	userGroup.POST("/get/user/:userId", controllers.GetUserByID)
	// แสดงผู้เล่นทั้งหมด
	userGroup.POST("/get/user", controllers.GetAllUser)
	// แสดงผู้เล่นที่ยังไม่ได้รับการอนุมัติ
	userGroup.POST("/get/user/pending", controllers.GetUserPending)
	// แก้ไขข้อมูลผู้เล่น
	userGroup.PUT("/update/user/:userId", controllers.UpdateUser)
	// แสดงคะแนนรวมทั้งหมด
	userGroup.POST("/get/sumsocre/:userId", controllers.GetSumScore)
	// รวมคะแนนทั้งที่ผู้เล่นทำได้
	userGroup.POST("/set/sumscore/:userId", controllers.SetSumSocore)
}
