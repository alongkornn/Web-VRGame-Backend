package controllers

import (
	"net/http"

	"github.com/alongkornn/Web-VRGame-Backend/internal/admin/dto"
	"github.com/alongkornn/Web-VRGame-Backend/internal/admin/services"
	"github.com/alongkornn/Web-VRGame-Backend/pkg/utils"
	"github.com/labstack/echo/v4"
)

// approve user register
func ApprovedRegister(ctx echo.Context) error {
	id := ctx.Param("id")
	var approveDTO dto.Approved
	if err := ctx.Bind(&approveDTO); err != nil {
		return utils.SendError(ctx, http.StatusBadRequest, "Invalid input", nil)
	}

	status, err := services.ApprovedRegister(id, approveDTO.Status,  ctx.Request().Context())
	if err != nil {
		return utils.SendError(ctx, status, err.Error(), nil)
	}

	return utils.SendSuccess(ctx, status, "Successfully to approved", nil)
}