package controllers

import (
	"net/http"

	"github.com/alongkornn/Web-VRGame-Backend/config"
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
var store = sessions.NewCookieStore([]byte(config.GetEnv("super_secret_key")))
func Login(ctx echo.Context) error {
	var loginDTO dto.LoginDTO
	if err := ctx.Bind(&loginDTO); err != nil {
		return utils.SendError(ctx, http.StatusBadRequest, "Invalid input", nil)
	}

	token, status, err := services.Login(loginDTO.Email, loginDTO.Password, ctx.Request().Context())
	if err != nil {
		return utils.SendError(ctx, status, err.Error(), nil)
	}

	// บันทึก token ไว้ใน localstorage
	session, _ := store.Get(ctx.Request(), "authentication")

	session.Values["token"] = token
    session.Values["authenticated"] = true
    session.Save(ctx.Request(), ctx.Response())
	return utils.SendSuccess(ctx, status, "Successfully to Login", nil)











}

