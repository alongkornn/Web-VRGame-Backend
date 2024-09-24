package controllers

import (
	"net/http"

	"github.com/alongkornn/Web-VRGame-Backend/internal/checkpoint/dto"
	"github.com/alongkornn/Web-VRGame-Backend/internal/checkpoint/services"
	"github.com/alongkornn/Web-VRGame-Backend/pkg/utils"
	"github.com/labstack/echo/v4"
)

func GetCurrentCheckpointFromUser(ctx echo.Context) error {
	userID := ctx.Param("user")

	checkpoint, status, err := services.GetCurrentCheckpointFromUserId(userID, ctx.Request().Context())
	if err != nil {
		return utils.SendError(ctx, status, err.Error(), nil)
	}

	return utils.SendSuccess(ctx, status, "Add checkpoint success", checkpoint)
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


func SaveCheckpointToComplete(ctx echo.Context) error {
	id := ctx.Param("userId")

	status, err := services.SaveCheckpointToComplete(id, ctx.Request().Context())
	if err != nil {
		return utils.SendError(ctx, status, err.Error(), nil)
	}

	return utils.SendSuccess(ctx, status, "Successfully to save", nil)
}

func GetCompleteCheckpointByUserId(ctx echo.Context) error {
	id := ctx.Param("userId")
	completeCheckpoints, status, err := services.GetCompleteCheckpointByUserId(id, ctx.Request().Context())
  if err != nil {
    return utils.SendError(ctx, status, err.Error(), nil)
  }
  return utils.SendSuccess(ctx, status, "Successfully to get checkpoinComplete", completeCheckpoints)
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
