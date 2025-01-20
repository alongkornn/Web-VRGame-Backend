package main

import (
	"fmt"
	"net/http"

	"github.com/alongkornn/Web-VRGame-Backend/config"
	adminRoute "github.com/alongkornn/Web-VRGame-Backend/internal/admin/routes"
	authRoute "github.com/alongkornn/Web-VRGame-Backend/internal/auth/routes"
	checkpointRoute "github.com/alongkornn/Web-VRGame-Backend/internal/checkpoint/routes"
	middlewares "github.com/alongkornn/Web-VRGame-Backend/internal/middleware"
	scoreRoute "github.com/alongkornn/Web-VRGame-Backend/internal/score/routes"
	userRoute "github.com/alongkornn/Web-VRGame-Backend/internal/user/routes"
	websocket_services "github.com/alongkornn/Web-VRGame-Backend/internal/websocket/services"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	config.InitConfig()
	// Firebase Config
	config.InitFirebase()
	// Redis Config
	config.InitRedis()

	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:3000"}, // Origin ที่อนุญาต
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowHeaders:     []string{echo.HeaderContentType, echo.HeaderAuthorization},
		AllowCredentials: true, // สำคัญ! เพื่ออนุญาตการส่งคุกกี้
	}))

	globalGroup := e.Group(config.GetEnv("app.prefix"))

	// เรียกใช้ middleware ในทุก api ที่เรียกโดย middleware

	globalGroup.POST("/", func(c echo.Context) error {
		c.JSON(http.StatusOK, map[string]string{"message": "Hello world."})
		return nil
	})

	// ตรวจสอบว่า token ถูกต้องหรือไม่
	protectRoute := e.Group("/api")
	protectRoute.Use(middlewares.JWTMiddlewareWithCookie((config.GetEnv("jwt.secret_key"))))
	protectRoute.GET("/protected", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "You are authorized",
		})
	})

	authRoute.AuthRoute(globalGroup)
	scoreRoute.ScoreRoute(globalGroup)
	adminRoute.AdminRoute(globalGroup)
	checkpointRoute.CheckpointRoute(globalGroup)
	userRoute.UserRoute(globalGroup)

	// เริ่มต้น WebSocket Server
	e.Any("/ws", websocket_services.HandleWebSocket)

	// เริ่มต้น Firestore Listener
	go config.ListenForUserScoreUpdates()

	port := config.GetEnv("app.port")
	e.Logger.Fatal(e.Start(":" + port))
	fmt.Println("Server started on port " + port)
}
