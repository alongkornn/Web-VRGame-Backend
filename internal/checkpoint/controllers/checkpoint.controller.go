package controllers

import (
	"net/http"

	"github.com/alongkornn/Web-VRGame-Backend/internal/checkpoint/dto"
	"github.com/alongkornn/Web-VRGame-Backend/internal/checkpoint/services"
	"github.com/alongkornn/Web-VRGame-Backend/pkg/utils"
	"github.com/labstack/echo/v4"
)

// เพิ่มด่านปัจจุบันให้กับผู้เล่น
func GetCurrentCheckpointToUser(ctx echo.Context) error {
	checkpointID := ctx.Param("checkpoint")
	userID := ctx.Param("user")

	status, err := services.GetCurrentCheckpointToUser(checkpointID, userID, ctx.Request().Context())
	if err != nil {
		return utils.SendError(ctx, status, err.Error(), nil)
	}

	return utils.SendSuccess(ctx, status, "Add checkpoint success", nil)
}

func GetAllCheckpoint(ctx echo.Context) error {
	checkpoints, status, err := services.GetAllCheckpoint(ctx.Request().Context())
	if err != nil {
		return utils.SendError(ctx, status, err.Error(), nil)
	}

	return utils.SendSuccess(ctx, status, "Get checkpoint success", checkpoints)
}

func CreateCheckpoint(ctx echo.Context) error {
	var checkpointDTO dto.CreateCheckpointsDTO
	if err := ctx.Bind(&checkpointDTO); err != nil {
		return utils.SendError(ctx, http.StatusBadRequest, "Invalid inpout", nil)
	}

	status, err := services.CreateCheckpoint(checkpointDTO, ctx.Request().Context())
	if err != nil {
		return utils.SendError(ctx, status, err.Error(), nil)
	}

	return utils.SendSuccess(ctx, status, "Successfully to created", nil)
}

func GetCheckpointWithCategory(ctx echo.Context) error {
	var categoryDTO dto.GetCheckpointWithCategoryDTO
	if err := ctx.Bind(&categoryDTO); err != nil {
		return utils.SendError(ctx, http.StatusBadRequest, "Invalid input", nil)
	}

	checkpoints, status, err := services.GetCheckpointWithCategory(categoryDTO.Category, ctx.Request().Context())
	if err != nil {
		return utils.SendError(ctx, status, err.Error(), nil)
	}

	return utils.SendSuccess(ctx, status, "Successfully to get checkpoin with catery", checkpoints)
}