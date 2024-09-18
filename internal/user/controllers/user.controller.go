package controllers

import (
	"github.com/alongkornn/Web-VRGame-Backend/internal/user/services"
	"github.com/alongkornn/Web-VRGame-Backend/pkg/utils"
	"github.com/labstack/echo/v4"
)

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