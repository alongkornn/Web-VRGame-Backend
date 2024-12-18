package middlewares

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// JWTMiddleware ตรวจสอบว่า JWT token ถูกต้องและ decode เพื่อใช้ข้อมูลข้างใน token
func JWTMiddlewareWithCookie(secretKey string) echo.MiddlewareFunc {
	return middleware.JWTWithConfig(middleware.JWTConfig{
		SigningKey:  []byte(secretKey),
		TokenLookup: "cookie:token", // ดึง token จาก cookie ชื่อ "token"
		ContextKey:  "user",         // บันทึกข้อมูลผู้ใช้ใน Context Key ชื่อ "user"
	})
}

// AdminMiddleware ตรวจสอบว่า role ของผู้ใช้เป็น admin หรือไม่
func RoleBasedMiddleware(role string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userRole := c.Get("role") // assuming the role is stored in context after JWT validation
			if userRole != role {
				return c.JSON(http.StatusForbidden, map[string]string{"message": "Access forbidden"})
			}
			return next(c)
		}
	}
}
