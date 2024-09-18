package main

import (
	"fmt"
	"net/http"

	"github.com/alongkornn/Web-VRGame-Backend/config"
	adminRoute "github.com/alongkornn/Web-VRGame-Backend/internal/admin/routes"
	authRoute "github.com/alongkornn/Web-VRGame-Backend/internal/auth/routes"
	checkpointRoute "github.com/alongkornn/Web-VRGame-Backend/internal/checkpoint/routes"
	scoreRoute "github.com/alongkornn/Web-VRGame-Backend/internal/score/routes"
	userRoute "github.com/alongkornn/Web-VRGame-Backend/internal/user/routes"
	"github.com/labstack/echo/v4"
)

func main() {
	config.InitConfig()
	// Firebase Config
	config.InitFirebase()

	e := echo.New()
	globalGroup := e.Group(config.GetEnv("app.prefix"))

	globalGroup.POST("/", func(c echo.Context) error {
		c.JSON(http.StatusOK, map[string]string{"message": "Hello world."})
		return nil
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
