package services

import (
	"context"
	"errors"
	"net/http"

	"github.com/alongkornn/Web-VRGame-Backend/config"
	"github.com/alongkornn/Web-VRGame-Backend/internal/auth/models"
)

func GetUserByID(id string, ctx context.Context) (*models.User, int, error) {
	hasUser := config.DB.Collection("User").
		Where("is_deleted", "==", false).
		Where("status", "==", models.Approved).
		Where("id", "==", id).
		Limit(1)

	doc, err := hasUser.Documents(ctx).Next()
	if err != nil {
		return nil, http.StatusBadRequest, errors.New("user not found")
	}

	var user *models.User
	doc.DataTo(user)


	return user, http.StatusOK, nil
}
