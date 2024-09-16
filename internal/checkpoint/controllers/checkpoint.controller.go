package controllers

import (
	"github.com/alongkornn/Web-VRGame-Backend/internal/checkpoint/services"
	"github.com/alongkornn/Web-VRGame-Backend/pkg/utils"
	"github.com/labstack/echo/v4"
)


func GetCurrentCheckpointToUser(ctx echo.Context) error {
	checkpointID := ctx.Param("checkpoint")
	userID := ctx.Param("user")

	status, err := services.GetCurrentCheckpointToUser(checkpointID, userID, ctx.Request().Context())
	if err != nil {
		return utils.SendError(ctx, status, err.Error(), nil)
	}

	return utils.SendSuccess(ctx, status, "Add checkpoint success", nil)
}