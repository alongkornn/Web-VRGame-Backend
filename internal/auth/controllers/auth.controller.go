package controllers

import (
	"net/http"

	"github.com/alongkornn/Web-VRGame-Backend/internal/auth/dto"
	"github.com/alongkornn/Web-VRGame-Backend/internal/auth/services"
	"github.com/alongkornn/Web-VRGame-Backend/pkg/utils"
	"github.com/labstack/echo/v4"
)


func Register(c echo.Context) error {
	var userRegister dto.RegisterDTO
	if err := c.Bind(&userRegister); err != nil {
		return utils.SendError(c, http.StatusBadRequest, "Invalid input", nil)
	}

	status, err := services.Register(userRegister)
	if (status != http.StatusOK && err != nil) {
		return utils.SendError(c, status, "Failed to create User", nil)
	}
	return utils.SendSuccess(c, status, "Created User Successfully", nil)
}