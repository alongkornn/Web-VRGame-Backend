package controllers

import (
	"net/http"
	"strconv"

	"github.com/alongkornn/Web-VRGame-Backend/internal/checkpoint/dto"
	"github.com/alongkornn/Web-VRGame-Backend/internal/checkpoint/services"
	"github.com/alongkornn/Web-VRGame-Backend/pkg/utils"
	"github.com/labstack/echo/v4"
)

// แสดงด่านปัจจุบันของผู้เล่น(โดยจะเข้าถึงผ่านไอดีผู้เล่น)
func GetCurrentCheckpointFromUser(ctx echo.Context) error {
	userID := ctx.Param("user")

	checkpoint, status, err := services.GetCurrentCheckpointFromUserId(userID, ctx.Request().Context())
	if err != nil {
		return utils.SendError(ctx, status, err.Error(), nil)
	}

	return utils.SendSuccess(ctx, status, "Add checkpoint success", checkpoint)
}

// แสดงทุกด่านทุกหมวดหมู่
func GetAllCheckpoint(ctx echo.Context) error {
	checkpoints, status, err := services.GetAllCheckpoint(ctx.Request().Context())
	if err != nil {
		return utils.SendError(ctx, status, err.Error(), nil)
	}

	return utils.SendSuccess(ctx, status, "Get checkpoint success", checkpoints)
}

// สร้างด่านใหม่
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

// บันทึกด่านปัจจุบันลงในด่านที่เล่นผ่านแล้วโดยจะตรวจสอบว่าคะแนนผ่านเกณฑ์หรือยัง
func SaveCheckpointToComplete(ctx echo.Context) error {
	id := ctx.Param("userId")
	scoreStr := ctx.Param("score")

	// แปลงค่า scoreStr จาก string เป็น int
	score, err := strconv.Atoi(scoreStr)
	if err != nil {
		// ถ้าไม่สามารถแปลงได้ให้คืนค่า error
		return ctx.JSON(400, map[string]string{
			"error": "Invalid score parameter",
		})
	}

	status, err := services.SaveCheckpointToComplete(id, score, ctx.Request().Context())
	if err != nil {
		return utils.SendError(ctx, status, err.Error(), nil)
	}

	return utils.SendSuccess(ctx, status, "Successfully to save", nil)
}

// แสดงด่านที่เล่นผ่านแล้ว(โดยจะเข้าถึงผ่านไอดีของผู้เล่น)
func GetCompleteCheckpointByUserId(ctx echo.Context) error {
	id := ctx.Param("userId")
	completeCheckpoints, status, err := services.GetCheckpointDetails(id, ctx.Request().Context())
	if err != nil {
		return utils.SendError(ctx, status, err.Error(), nil)
	}
	return utils.SendSuccess(ctx, status, "Successfully to get checkpoinComplete", completeCheckpoints)
}

// แสดงทุกด่านโดยเข้าถึงผ่านหมวดหมู่
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

func SetTime(ctx echo.Context) error {
	userId := ctx.Param("userId")
	var time dto.SetTimeDTO
	if err := ctx.Bind(&time); err != nil {
		return utils.SendError(ctx, http.StatusBadRequest, "Invalid input", nil)
	}

	status, err := services.SetTime(userId, time.Time, ctx.Request().Context())
	if err != nil {
		return utils.SendError(ctx, status, err.Error(), nil)
	}

	return utils.SendSuccess(ctx, status, "Successfully to set time", nil)
}
