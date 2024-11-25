package controllers

import (
	"fmt"
	"net/http"

	"github.com/alongkornn/Web-VRGame-Backend/internal/auth/dto"
	"github.com/alongkornn/Web-VRGame-Backend/internal/auth/services"
	"github.com/alongkornn/Web-VRGame-Backend/pkg/utils"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

// ลงทะเบียน
func Register(ctx echo.Context) error {
	var registerDTO dto.RegisterDTO
	if err := ctx.Bind(&registerDTO); err != nil {
		return utils.SendError(ctx, http.StatusBadRequest, "Invalid input", nil)
	}

	status, err := services.Register(ctx.Request().Context(), &registerDTO)
	if err != nil {
		return utils.SendError(ctx, status, "Failed to create User", err.Error())
	}
	return utils.SendSuccess(ctx, status, "Created User Successfully", nil)
}

// เข้าสู้ระบบ
// สร้าง key ที่ปลอดภัย - ควรมีความยาวอย่างน้อย 32 bytes
var (
	store *sessions.CookieStore
	key   = "a@lkDKP%1!skeLOkd#" // ต้องแน่ใจว่าค่านี้ไม่เป็นค่าว่าง
)

func init() {
	// ตรวจสอบว่ามี key หรือไม่
	if len(key) == 0 {
		panic("SESSION_KEY is not set in environment variables")

	}

	// สร้าง store ด้วย key ที่กำหนด
	store = sessions.NewCookieStore([]byte(key))

	// กำหนดค่าเริ่มต้นสำหรับ store
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		Domain:   "", // ระบุ domain ถ้าจำเป็น เช่น "localhost"
	}
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

	tokenString := fmt.Sprintf("%v", token)

	// ใช้ echo.Context wrapper สำหรับ standard http
	req := ctx.Request()
	res := ctx.Response()

	// สร้าง session ใหม่
	session, err := store.Get(req, "authentication")
	if err != nil {
		// log error เพื่อ debug
		println("Session Error:", err.Error())
		return utils.SendError(ctx, http.StatusInternalServerError, "Session error", nil)
	}

	// เก็บข้อมูลใน session
	session.Values["token"] = tokenString
	session.Values["authenticated"] = true

	// บันทึก session ด้วย http.ResponseWriter
	err = session.Save(req, res.Writer)
	if err != nil {
		// log error เพื่อ debug
		println("Save Session Error:", err.Error())
		return utils.SendError(ctx, http.StatusInternalServerError, "Cannot save session", err.Error())
	}

	// สร้าง cookie แบบ manual เพิ่มเติม
	cookie := new(http.Cookie)
	cookie.Name = "authentication"
	cookie.Value = tokenString
	cookie.Path = "/"
	cookie.MaxAge = 3600
	cookie.HttpOnly = true
	cookie.Secure = false
	cookie.SameSite = http.SameSiteLaxMode
	ctx.SetCookie(cookie)

	return utils.SendSuccess(ctx, status, "Successfully to Login", token)
}
