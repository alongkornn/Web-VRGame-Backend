package middleware

import (
	"net/http"

	"github.com/alongkornn/Web-VRGame-Backend/config"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// JWTConfig provides configuration for the JWT middleware
func JWTConfig() middleware.JWTConfig {
	secretKey := config.GetEnv("jwt.secret_key")
	return middleware.JWTConfig{
		Claims:     &jwtCustomClaims{},
		SigningKey: []byte(secretKey),
	}
}

// jwtCustomClaims are custom claims extending default ones
type jwtCustomClaims struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	jwt.StandardClaims
}

// JWTMiddleware returns a JWT middleware instance
func JWTMiddleware() echo.MiddlewareFunc {
	return middleware.JWTWithConfig(JWTConfig())
}

// RequireJWT is a middleware that requires a valid JWT token
func RequireJWT() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			user := c.Get("user").(*jwt.Token)
			claims := user.Claims.(*jwtCustomClaims)
			if claims == nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"status":  "error",
					"message": "unauthorized",
				})
			}
			return next(c)
		}
	}
}
