package controllers

import (
	"net/http"

	"github.com/alongkornn/Web-VRGame-Backend/internal/user/dto"
	"github.com/alongkornn/Web-VRGame-Backend/internal/user/services"
	"github.com/alongkornn/Web-VRGame-Backend/pkg/utils"
	"github.com/labstack/echo/v4"
)

// user approved
func GetUserByID(ctx echo.Context) error {
	id := ctx.Param("id")

	user, status, err := services.GetUserByID(id, ctx.Request().Context())
	if err != nil {
		return utils.SendError(ctx, status, err.Error(), nil)
	}


	return utils.SendSuccess(ctx, status, "Successfully to get User", user)
}


// เวลาผู้เล่นเข้าเล่นด่านไหนให้เพิ่มผู้เล่นไปในด่านด้วย
func AddPlayerInCheckpoint(ctx echo.Context) error {
	checkpointID := ctx.Param("checkpointID")
	userID := ctx.Param("userID")

	status, err := services.AddPlayerInCheckpoint(checkpointID, userID, ctx.Request().Context())
	if err != nil {
		return utils.SendError(ctx, status, err.Error(), nil)
	}

	return utils.SendSuccess(ctx, status, "Add player success", nil)
} 

func GetAllUser(ctx echo.Context) error {
	users, status, err := services.GetAllUser(ctx.Request().Context())
	if err != nil {
		return utils.SendError(ctx, status, err.Error(), nil)
	}


	return utils.SendSuccess(ctx, status, "Successfully to get User", users)
}

// user pending

func GetUserPending(ctx echo.Context) error {
	users, status, err := services.GetUserPending(ctx.Request().Context())
	if err != nil {
		return utils.SendError(ctx, status, err.Error(), nil)
	}

	return utils.SendSuccess(ctx, status, "Successfully to get User status is pending", users)
}

func UpdateUser(ctx echo.Context) error {
	id := ctx.Param("id")
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