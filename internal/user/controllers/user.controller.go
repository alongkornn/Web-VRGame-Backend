package controllers

import (
	"net/http"

	"github.com/alongkornn/Web-VRGame-Backend/internal/user/dto"
	"github.com/alongkornn/Web-VRGame-Backend/internal/user/services"
	"github.com/alongkornn/Web-VRGame-Backend/pkg/utils"
	"github.com/labstack/echo/v4"
)

// แสดงผู้เล่นโดยเข้าถึงผ่านไอดี
func GetUserByID(ctx echo.Context) error {
	id := ctx.Param("userId")

	user, status, err := services.GetUserByID(id, ctx.Request().Context())
	if err != nil {
		return utils.SendError(ctx, status, err.Error(), nil)
	}

	return utils.SendSuccess(ctx, status, "Successfully to get User", user)
}

// แสดงผู้เล่นทั้งหมด
func GetAllUser(ctx echo.Context) error {
	users, status, err := services.GetAllUser(ctx.Request().Context())
	if err != nil {
		return utils.SendError(ctx, status, err.Error(), nil)
	}

	return utils.SendSuccess(ctx, status, "Successfully to get User", users)
}

// แสดงผู้เล่นที่ยังไม่ได้อนุมัติ
func GetUserPending(ctx echo.Context) error {
	users, status, err := services.GetUserPending(ctx.Request().Context())
	if err != nil {
		return utils.SendError(ctx, status, err.Error(), nil)
	}

	return utils.SendSuccess(ctx, status, "Successfully to get User status is pending", users)
}

// แก้ไขข้อมูลผู้เล่น
func UpdateUser(ctx echo.Context) error {
	id := ctx.Param("userId")
	var updateUserDTO dto.UpdateUserDTO
	if err := ctx.Bind(&updateUserDTO); err != nil {
		return utils.SendError(ctx, http.StatusBadRequest, "Invalid input", nil)
	}

	status, err := services.UpdateUser(id, updateUserDTO, ctx.Request().Context())
	if err != nil {
		return utils.SendError(ctx, status, err.Error(), nil)
	}

	return utils.SendSuccess(ctx, status, "Successfully to update data", nil)
}

// แสดงคะแนนรวมของผู้เล่น
func GetSumScore(ctx echo.Context) error {
	userId := ctx.Param("userId")

	sumScore, status, err := services.GetSumScore(userId, ctx.Request().Context())
	if err != nil {
		return utils.SendError(ctx, status, err.Error(), nil)
	}

	return utils.SendSuccess(ctx, status, "Successfully to get score", sumScore)
}

// รวมคะแนนทั้งหมดของผู้เล่น
func SetSumSocore(ctx echo.Context) error {
	userId := ctx.Param("userId")
	status, err := services.SetSumScore(userId, ctx.Request().Context())
	if err != nil {
		return utils.SendError(ctx, status, err.Error(), nil)
	}

	return utils.SendSuccess(ctx, status, "Successfully to set score", nil)
}

func GetUserBySortScore(ctx echo.Context) error {
	users, status, err := services.GetUserBySortScore(ctx.Request().Context())

	if err != nil {
		return utils.SendError(ctx, status, err.Error(), nil)
	}

	return utils.SendSuccess(ctx, status, "Successfully to fetch User", users)
}
