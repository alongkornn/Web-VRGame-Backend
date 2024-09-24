package controllers

import (
	"net/http"

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

	status, err := services.Register(ctx.Request().Context(), &registerDTO)
	if err != nil {
		return utils.SendError(ctx, status, "Failed to create User", err.Error())
	}
	return utils.SendSuccess(ctx, status, "Created User Successfully", nil)
}

// เข้าสู้ระบบ
func Login(ctx echo.Context) error {
	var loginDTO dto.LoginDTO
	if err := ctx.Bind(&loginDTO); err != nil {
		return utils.SendError(ctx, http.StatusBadRequest, "Invalid input", nil)
	}

	data, status, err := services.Login(loginDTO.Email, loginDTO.Password, ctx.Request().Context())
	if err != nil {
		return utils.SendError(ctx, status, err.Error(), nil)
	}
	return utils.SendSuccess(ctx, status, "Successfully to Login", data)
}

