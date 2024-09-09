package middleware

import (
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// JWTMiddleware ตรวจสอบว่า JWT token ถูกต้องและ decode เพื่อใช้ข้อมูลข้างใน token
func JWTMiddleware(secretKey string) echo.MiddlewareFunc {
	return middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:  []byte(secretKey),
		TokenLookup: "header:Authorization",
		AuthScheme:  "Bearer",
		ContextKey:  "user",
	})
}

// AdminMiddleware ตรวจสอบว่า role ของผู้ใช้เป็น admin หรือไม่
func AdminMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)

		role := claims["role"].(string)
		if strings.ToLower(role) != "admin" {
			return c.JSON(http.StatusForbidden, map[string]string{"message": "Forbidden, admin only"})
		}

		return next(c)
	}
}
