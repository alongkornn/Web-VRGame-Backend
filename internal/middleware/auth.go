package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/alongkornn/Web-VRGame-Backend/config"
	"github.com/golang-jwt/jwt"
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

// validation jwt token
func JWTAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// ดึง Header Authorization
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		}

		// ดึง Token จาก Header
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// ตรวจสอบความถูกต้องของ Token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// ตรวจสอบว่า Algorithm เป็น HMAC
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			// คืนค่า Key ที่ใช้ในการตรวจสอบ Signature ของ Token
			return []byte(config.GetEnv("jwt.secret_key")), nil
		})

		// ถ้าเกิดข้อผิดพลาดหรือ Token ไม่ถูกต้อง
		if err != nil || !token.Valid {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
		}

		// ถ้าผ่านการตรวจสอบ ให้ทำงานต่อไป
		return next(c)
	}
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
