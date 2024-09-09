package routes

import (
	"github.com/alongkornn/Web-VRGame-Backend/internal/auth/controllers"
	"github.com/labstack/echo/v4"
)



func AuthRoute(g *echo.Group) {
	authGroup := g.Group("/auth")
	authGroup.POST("/register", controllers.Register)
}