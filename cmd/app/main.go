package main

import (
	"fmt"
	"net/http"

	"github.com/alongkornn/Web-VRGame-Backend/config"
	"github.com/labstack/echo/v4"
)

func main() {
    // Firebase Config
    config.InitFirebase()
	e := echo.New()

	globalGroup := e.Group("/api/v1")

	globalGroup.POST("/", func(c echo.Context) error {
		c.JSON(http.StatusOK, map[string]string{"message": "Hello World."})
		return nil
	})




    port := config.GetEnv("app.port")
	e.Logger.Fatal(e.Start(":" + port))
	fmt.Println("Server started on port " + port)
}
