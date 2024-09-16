package controllers

import (
	"net/http"

	"github.com/alongkornn/Web-VRGame-Backend/internal/admin/dto"
	"github.com/alongkornn/Web-VRGame-Backend/internal/admin/services"
	"github.com/alongkornn/Web-VRGame-Backend/pkg/utils"
	"github.com/labstack/echo/v4"
)

// approve user register
func ApprovedRegister(ctx echo.Context) error {
	id := ctx.Param("id")
	var approveDTO dto.Approved
	if err := ctx.Bind(&approveDTO); err != nil {
		return utils.SendError(ctx, http.StatusBadRequest, "Invalid input", nil)
	}

	status, err := services.ApprovedRegister(id, approveDTO.Status,  ctx.Request().Context())
	if err != nil {
		return utils.SendError(ctx, status, err.Error(), nil)
	}

	return utils.SendSuccess(ctx, status, "Successfully to approved", nil)
}

// remove user
func RemoveUser(ctx echo.Context) error {
	id := ctx.Param("id")

	status, err := services.RemoveUser(id, ctx.Request().Context())
	if err != nil {
		return utils.SendError(ctx, status, err.Error(), nil)
	}

	return utils.SendSuccess(ctx, status, "Successfully to Delete", nil)
}

// remove admin
func RemoveAdmin(ctx echo.Context) error {
	id := ctx.Param("id")

	status, err := services.RemoveAdmin(id, ctx.Request().Context())
	if err != nil {
		return utils.SendError(ctx, status, err.Error(), nil)
	}

	return utils.SendSuccess(ctx, status, "Successfully to Delete", nil)
}

func GetAllAdmin(ctx echo.Context) error {
	users, status, err := services.GetAllAdmin(ctx.Request().Context())
	if err != nil {
		return utils.SendError(ctx, status, err.Error(), nil)
	}

	return utils.SendSuccess(ctx, status, "Successfuly to get admin", users)
}

func GetAdminByID(ctx echo.Context) error {
	id := ctx.Param("id")

	user, status, err := services.GetAdminByID(id, ctx.Request().Context())
	if err != nil {
		return utils.SendError(ctx, status, err.Error(), nil)
	}

	return utils.SendSuccess(ctx, status, "Successfully to Get user", user)
}

// admin
func CreateAdmin(ctx echo.Context) error {
	id := ctx.Param("id")
	var roleDTO dto.RoleDTO
	if err := ctx.Bind(&roleDTO); err != nil {
		return utils.SendError(ctx, http.StatusBadRequest, "Invalid", nil)
	} 
	status, err := services.CreateAdmin(id, roleDTO.Role, ctx.Request().Context())
	if err != nil {
		return utils.SendError(ctx, status, err.Error(), nil)
	}

	return utils.SendSuccess(ctx, status, "Successfully to Created", nil)
}

func UpdateDataAdmin(ctx echo.Context) error {
	id := ctx.Param("id")
	var updateDTO dto.UpdateDTO
	if err := ctx.Bind(&updateDTO); err != nil {
		return utils.SendError(ctx, http.StatusBadRequest, "Invalid input", nil)
	}

	status, err := services.UpdateDataAdmin(id, updateDTO, ctx.Request().Context())
	if err != nil {
		return utils.SendError(ctx, status, err.Error(), nil)
	}

	return utils.SendSuccess(ctx, status, "Successfully to update data", nil)
}

func UpdatePasswordAdmin(ctx echo.Context) error {
	id := ctx.Param("id")

	var updatePasswordDTO dto.UpdatePasswordDTO
	if err := ctx.Bind(&updatePasswordDTO); err != nil {
		return utils.SendError(ctx, http.StatusBadRequest, "Invalid input", nil)
	}

	status, err := services.UpdatePasswordAdmin(id, updatePasswordDTO.Password, updatePasswordDTO.NewPassword, ctx.Request().Context())
	if err != nil {
		return utils.SendError(ctx, status, err.Error(), nil)
	}

	return utils.SendSuccess(ctx, status, "Successfully to change password", nil)
}