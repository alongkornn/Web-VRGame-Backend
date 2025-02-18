package middlewares

import (
	"net/http"
	"strings"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func JWTMiddleware(secretKey string) echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		SigningKey:  []byte(secretKey),
		TokenLookup: "header:Authorization,cookie:token",
		ContextKey:  "user",
		BeforeFunc: func(c echo.Context) {
			// ตรวจสอบ Authorization Header แล้วลบ "Bearer " ออก
			authHeader := c.Request().Header.Get("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") {
				token := strings.TrimPrefix(authHeader, "Bearer ")
				c.Request().Header.Set("Authorization", token)
			}
		},
	})
}

// RoleBasedMiddleware ตรวจสอบว่า role ของผู้ใช้เป็น admin หรือไม่
func RoleBasedMiddleware(role string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userRole := c.Get("role")
			if userRole != role {
				return c.JSON(http.StatusForbidden, map[string]string{"message": "Access forbidden"})
			}
			return next(c)
		}
	}
}
