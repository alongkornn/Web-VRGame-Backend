package controllers

import (
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