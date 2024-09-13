package services

import (
	"context"
	"errors"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/alongkornn/Web-VRGame-Backend/config"
	"github.com/alongkornn/Web-VRGame-Backend/internal/auth/models"
)

func ApprovedRegister(id, approved string, ctx context.Context) (int, error) {
	hasUser := config.DB.Collection("User").
		Where("is_deleted", "==", false).
		Where("status", "==", models.Pending).
		Where("id", "==", id).
		Limit(1)

	doc, err := hasUser.Documents(ctx).Next()
	if err != nil {
		return http.StatusNotFound, errors.New("user not found")
	}

	var user models.User
	err = doc.DataTo(&user)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	_, err = doc.Ref.Set(ctx, map[string]interface{}{
		"status": approved,
	}, firestore.MergeAll)
	if err != nil {
		return http.StatusBadRequest, errors.New("failed to approve")
	}
	return http.StatusOK, nil
}

func RemoveUser(id string, ctx context.Context) (int, error) {
	hasUser := config.DB.Collection("User").
		Where("is_deleted", "==", false).
		Where("id", "==", id).
		Limit(1)

	doc, err := hasUser.Documents(ctx).Next()
	if err != nil {
		return http.StatusNotFound, errors.New("user not found")
	}

	var user models.User
	err = doc.DataTo(&user)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	_, err = doc.Ref.Set(ctx, map[string]interface{}{
		"is_deleted": true,
	}, firestore.MergeAll)
	if err != nil {
		return http.StatusBadRequest, err
	}
	return http.StatusOK, nil
}
