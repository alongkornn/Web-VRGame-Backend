package routes

import (
	// "github.com/alongkornn/Web-VRGame-Backend/config"

	"github.com/alongkornn/Web-VRGame-Backend/internal/auth/controllers"
	// "github.com/alongkornn/Web-VRGame-Backend/internal/middleware"
	"github.com/labstack/echo/v4"
)

func AuthRoute(g *echo.Group) {
	authGroup := g.Group("/auth")
	authGroup.POST("/register", controllers.Register)
	authGroup.POST("/login", controllers.Login)
	// Apply JWT middleware only to the routes that require authentication
	// protectedGroup := g.Group("")
	// protectedGroup.Use(middleware.JWTMiddleware(config.GetEnv("jwt.secret_key")))

	// authGroup.GET("/verify-email", controllers.VerifyEmail)

}
