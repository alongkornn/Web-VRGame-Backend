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
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	config.InitConfig()
	// Firebase Config
	config.InitFirebase()

	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowHeaders: []string{"Content-Type", "Authorization"},
	}))

	// เรียกใช้ middleware ในทุก api ที่เรียกโดย e

	globalGroup := e.Group(config.GetEnv("app.prefix"))

	middleware := e.Group("/test")
	middleware.Use(middlewares.JWTMiddlewareWithCookie((config.GetEnv("jwt.secret_key"))))

	globalGroup.POST("/", func(c echo.Context) error {
		c.JSON(http.StatusOK, map[string]string{"message": "Hello world."})
		return nil
	})

	middleware.GET("/protected", func(c echo.Context) error {
		// หากตรวจสอบผ่าน, จะมีข้อมูลของผู้ใช้ใน context
		user := c.Get("user")
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "Welcome, authorized user!",
			"user":    user,
		})
	})

	authRoute.AuthRoute(globalGroup)
	scoreRoute.ScoreRoute(globalGroup)
	adminRoute.AdminRoute(globalGroup)
	checkpointRoute.CheckpointRoute(globalGroup)
	userRoute.UserRoute(globalGroup)

	port := config.GetEnv("app.port")
	e.Logger.Fatal(e.Start(":" + port))
	fmt.Println("Server started on port " + port)
}
