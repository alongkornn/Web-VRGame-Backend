package main

import (
	"fmt"
	"net/http"

	"github.com/alongkornn/Web-VRGame-Backend/config"
	"github.com/alongkornn/Web-VRGame-Backend/internal/auth/routes"
	"github.com/labstack/echo/v4"
)

func main() {
	// Firebase Config
	config.InitFirebase()

	e := echo.New()
	globalGroup := e.Group("/api")

	globalGroup.POST("/", func(c echo.Context) error {
		c.JSON(http.StatusOK, map[string]string{"message": "Hello world."})
		return nil
	})

	routes.AuthRoute(globalGroup)

	e.Logger.Fatal(e.Start(":8000"))
	fmt.Println("Server started on port 8080")
}
