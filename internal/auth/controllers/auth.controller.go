package controllers

import (
	"net/http"
	"time"

	"github.com/alongkornn/Web-VRGame-Backend/internal/auth/dto"
	"github.com/alongkornn/Web-VRGame-Backend/internal/auth/services"
	"github.com/alongkornn/Web-VRGame-Backend/pkg/utils"
	"github.com/labstack/echo/v4"
)

// ลงทะเบียน
func Register(ctx echo.Context) error {
	var registerDTO dto.RegisterDTO

	if err := ctx.Bind(&registerDTO); err != nil {
		return utils.SendError(ctx, http.StatusBadRequest, "Invalid input", nil)
	}

	status, err := services.Register(ctx.Request().Context(), registerDTO)
	if err != nil {
		return utils.SendError(ctx, status, err.Error(), nil)
	}

	return utils.SendSuccess(ctx, status, "Created User Successfully", nil)
}

func Login(ctx echo.Context) error {
	var loginDTO dto.LoginDTO
	if err := ctx.Bind(&loginDTO); err != nil {
		return utils.SendError(ctx, http.StatusBadRequest, "Invalid input", nil)
	}

	// เรียกใช้ service เพื่อ login
	token, status, err := services.Login(loginDTO.Email, loginDTO.Password, ctx.Request().Context())
	if err != nil {
		return utils.SendError(ctx, status, err.Error(), nil)
	}

	// เก็บ JWT ใน Cookies
	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = token
	cookie.HttpOnly = false
	cookie.Secure = false                  // ใช้ true สำหรับ HTTPS เท่านั้น ในกรณีนี้เราต้องการทดสอบทเท่านั้น
	cookie.SameSite = http.SameSiteLaxMode // สำหรับการ Cross-origin
	cookie.Path = "/"
	cookie.Expires = time.Now().Add(24 * time.Hour) // คุกกี้หมดอายุใน 1 วัน
	ctx.SetCookie(cookie)

	return utils.SendSuccess(ctx, status, "Successfully to Login", token)
}

// func VerifyEmail(ctx echo.Context) error {
// 	// รับโทเค็นจาก URL
// 	token := ctx.QueryParam("token")

// 	// เรียกใช้ฟังก์ชัน VerifyEmail เพื่อตรวจสอบโทเค็นและอัปเดตสถานะการยืนยัน
// 	status, err := services.VerifyEmail(ctx.Request().Context(), token)
// 	if err != nil {
// 		return utils.SendError(ctx, status, "Failed to verify email", err.Error())
// 	}

// 	// ส่งข้อความแจ้งผู้ใช้ว่าอีเมลยืนยันสำเร็จ
// 	return utils.SendSuccess(ctx, http.StatusOK, "Email verified successfully", nil)
// }
