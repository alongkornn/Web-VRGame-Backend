package controllers

import (
	"net/http"

	"github.com/alongkornn/Web-VRGame-Backend/internal/admin/dto"
	"github.com/alongkornn/Web-VRGame-Backend/internal/admin/services"
	"github.com/alongkornn/Web-VRGame-Backend/pkg/utils"
	"github.com/labstack/echo/v4"
)

// ผู้ดูแลระบบอนุมัติการลงทะเบียนของผู้เล่น
func AddminApprovedUserRegister(ctx echo.Context) error {
	userId := ctx.Param("userId")
	var approveDTO dto.Approved
	if err := ctx.Bind(&approveDTO); err != nil {
		return utils.SendError(ctx, http.StatusBadRequest, "Invalid input", nil)
	}

	status, err := services.AddminApprovedUserRegister(userId, approveDTO.Status, ctx.Request().Context())
	if err != nil {
		return utils.SendError(ctx, status, err.Error(), nil)
	}

	return utils.SendSuccess(ctx, status, "Successfully to approved", nil)
}

// ลบผู้เล่นออก
func AdminRemoveUser(ctx echo.Context) error {
	userId := ctx.Param("userId")

	status, err := services.AdminRemoveUser(userId, ctx.Request().Context())
	if err != nil {
		return utils.SendError(ctx, status, err.Error(), nil)
	}

	return utils.SendSuccess(ctx, status, "Successfully to Delete", nil)
}

// ลบผู้ดูแลระบบออก
func RemoveAdmin(ctx echo.Context) error {
	adminId := ctx.Param("adminId")

	status, err := services.RemoveAdmin(adminId, ctx.Request().Context())
	if err != nil {
		return utils.SendError(ctx, status, err.Error(), nil)
	}

	return utils.SendSuccess(ctx, status, "Successfully to Delete", nil)
}

// แสดงแอดมินทั้งหมด
func GetAllAdmin(ctx echo.Context) error {
	users, status, err := services.GetAllAdmin(ctx.Request().Context())
	if err != nil {
		return utils.SendError(ctx, status, err.Error(), nil)
	}

	return utils.SendSuccess(ctx, status, "Successfuly to get admin", users)
}

// แสดงผู้ดูแลระบบโดยเข้าถึงผ่านไอดีของผู้ดูแลระบบ
func GetAdminById(ctx echo.Context) error {
	adminId := ctx.Param("adminId")

	user, status, err := services.GetAdminById(adminId, ctx.Request().Context())
	if err != nil {
		return utils.SendError(ctx, status, err.Error(), nil)
	}

	return utils.SendSuccess(ctx, status, "Successfully to Get user", user)
}

// สร้างผู้ดูแลระบบ
func CreateAdmin(ctx echo.Context) error {
	var roleDTO dto.RoleDTO
	if err := ctx.Bind(&roleDTO); err != nil {
		return utils.SendError(ctx, http.StatusBadRequest, "Invalid", nil)
	}
	status, err := services.CreateAdmin(roleDTO.UserId, roleDTO.Role, ctx.Request().Context())
	if err != nil {
		return utils.SendError(ctx, status, err.Error(), nil)
	}

	return utils.SendSuccess(ctx, status, "Successfully to Created", nil)
}

// แก้ไขข้อมูลผู้ดูแลระบบ
func UpdateDataAdmin(ctx echo.Context) error {
	adminId := ctx.Param("adminId")
	var updateDTO dto.UpdateDTO
	if err := ctx.Bind(&updateDTO); err != nil {
		return utils.SendError(ctx, http.StatusBadRequest, "Invalid input", nil)
	}

	status, err := services.UpdateDataAdmin(adminId, updateDTO, ctx.Request().Context())
	if err != nil {
		return utils.SendError(ctx, status, err.Error(), nil)
	}

	return utils.SendSuccess(ctx, status, "Successfully to update data", nil)
}

// แก้ไขรหัสผ่านของผู้ดูแลระบบ
func UpdatePasswordAdmin(ctx echo.Context) error {
	adminId := ctx.Param("adminId")

	var updatePasswordDTO dto.UpdatePasswordDTO
	if err := ctx.Bind(&updatePasswordDTO); err != nil {
		return utils.SendError(ctx, http.StatusBadRequest, "Invalid input", nil)
	}

	status, err := services.UpdatePasswordAdmin(adminId, updatePasswordDTO.Password, updatePasswordDTO.NewPassword, ctx.Request().Context())
	if err != nil {
		return utils.SendError(ctx, status, err.Error(), nil)
	}

	return utils.SendSuccess(ctx, status, "Successfully to change password", nil)
}

func ShowScoreWiteStrength(ctx echo.Context) error {
	userId := ctx.Param("userId")

	isStrong, status, err := services.ShowScoreWiteStrength(userId, ctx.Request().Context())
	if err != nil {
		return utils.SendError(ctx, status, err.Error(), nil)
	}

	return utils.SendSuccess(ctx, status, "Successfully to fetch score", isStrong)
}
