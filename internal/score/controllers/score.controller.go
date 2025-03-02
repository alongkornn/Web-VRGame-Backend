package controllers

import (
	"net/http"

	"github.com/alongkornn/Web-VRGame-Backend/internal/score/dto"
	"github.com/alongkornn/Web-VRGame-Backend/internal/score/services"
	"github.com/alongkornn/Web-VRGame-Backend/pkg/utils"
	"github.com/labstack/echo/v4"
)

func GetScoreByUserId(ctx echo.Context) error {
	userId := ctx.Param("userId")
	score, status, err := services.GetScoreByUserId(userId, ctx.Request().Context())
	if err != nil {
		return utils.SendError(ctx, status, err.Error(), nil)
	}

	return utils.SendSuccess(ctx, status, "Successfully to get score", score)
}

// func GetAllScoreByCheckpointId(ctx echo.Context) error {
// 	checkpointId := ctx.Param("checkpointId")
// 	users_score, status, err := services.GetAllScoreByCheckpointId(checkpointId, ctx.Request().Context())
// 	if err != nil {
// 		return utils.SendError(ctx, status, err.Error(), nil)
// 	}

// 	return utils.SendSuccess(ctx, status, "Successfully to get Score", users_score)
// }

func SetProjectileScore(ctx echo.Context) error {
	id := ctx.Param("userId")
	var scoreDTO dto.SetScoreDTO
	if err := ctx.Bind(&scoreDTO); err != nil {
		return utils.SendError(ctx, http.StatusBadRequest, "Invalid input", nil)
	}

	status, err := services.SetProjectileScore(id, scoreDTO.Score, scoreDTO.Time, ctx.Request().Context())
	if err != nil {
		return utils.SendError(ctx, status, err.Error(), nil)
	}

	return utils.SendSuccess(ctx, status, "Successfully to save score", nil)
}
func SetMomentumScore(ctx echo.Context) error {
	id := ctx.Param("userId")
	var scoreDTO dto.SetScoreDTO
	if err := ctx.Bind(&scoreDTO); err != nil {
		return utils.SendError(ctx, http.StatusBadRequest, "Invalid input", nil)
	}

	status, err := services.SetMomentumScore(id, scoreDTO.Score, scoreDTO.Time, ctx.Request().Context())
	if err != nil {
		return utils.SendError(ctx, status, err.Error(), nil)
	}

	return utils.SendSuccess(ctx, status, "Successfully to save score", nil)
}
func SetForceScore(ctx echo.Context) error {
	id := ctx.Param("userId")
	var scoreDTO dto.SetScoreDTO
	if err := ctx.Bind(&scoreDTO); err != nil {
		return utils.SendError(ctx, http.StatusBadRequest, "Invalid input", nil)
	}

	status, err := services.SetForceScore(id, scoreDTO.Score, scoreDTO.Time, ctx.Request().Context())
	if err != nil {
		return utils.SendError(ctx, status, err.Error(), nil)
	}

	return utils.SendSuccess(ctx, status, "Successfully to save score", nil)
}
