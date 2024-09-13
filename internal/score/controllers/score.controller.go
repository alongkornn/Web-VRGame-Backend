package controllers

import (
	"github.com/alongkornn/Web-VRGame-Backend/internal/score/services"
	"github.com/alongkornn/Web-VRGame-Backend/pkg/utils"
	"github.com/labstack/echo/v4"
)



func GetScorebyID(ctx echo.Context) error {
	id := ctx.Param("id")
	score, status, err := services.GetScorebyID(id, ctx.Request().Context())
	if err != nil {
		return utils.SendError(ctx, status, err.Error(), nil)
	}

	return utils.SendSuccess(ctx, status, "Successfully to get score", score)
} 

func GetAllScore(ctx echo.Context) error {
	users_score, status, err := services.GetAllScore(ctx.Request().Context())
	if err != nil {
		return utils.SendError(ctx, status, err.Error(), nil)
	}

	return utils.SendSuccess(ctx, status, "Successfully to get Score", users_score)
}