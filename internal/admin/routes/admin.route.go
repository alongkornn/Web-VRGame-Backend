package routes

import (
	"github.com/alongkornn/Web-VRGame-Backend/internal/admin/controllers"
	"github.com/labstack/echo/v4"
)

func AdminRoute(g *echo.Group) {
	adminGroup := g.Group("/admin")

	// เพิ่มผู้ดูแลระบบ
	adminGroup.POST("/create", controllers.CreateAdmin)
	// แสดงผู้ดูแลระบบแค่คนเดียว
	adminGroup.POST("/get/:adminId", controllers.GetAdminById)
	// แสดงผูู้ดูแลระบบทั้งหมด
	adminGroup.POST("/get", controllers.GetAllAdmin)
	// ผู้ดูแลระบบอนุมัติการลงทะเบียนของผู้เล่น
	adminGroup.PUT("/approve/:userId", controllers.AddminApprovedUserRegister)
	// ผู้ดูแลระบบลบผู้เล่น
	adminGroup.DELETE("/delete/user/:userId", controllers.AdminRemoveUser)
	// ลบผู้ดูแลระบบออก
	adminGroup.DELETE("/delete/admin/:adminId", controllers.RemoveAdmin)
	// แก้ไขข้อมูลผู้ดูแลระบบ
	adminGroup.PUT("/update/admin/:adminId", controllers.UpdateDataAdmin)
	// แก้ไขรหัสผ่านผู้ดูแลระบบ
	adminGroup.PUT("/updatepassword/admin/:Id", controllers.UpdatePasswordAdmin)

	// protectedGroup := g.Group("")
	// adminGroup := protectedGroup.Group("/admin")
	// adminGroup.Use(middleware.RoleBasedMiddleware("admin"))
}
