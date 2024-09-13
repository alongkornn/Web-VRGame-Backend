package services

import (
	"context"
	"errors"
	"net/http"

	"github.com/alongkornn/Web-VRGame-Backend/config"
	"github.com/alongkornn/Web-VRGame-Backend/internal/auth/models"
	"google.golang.org/api/iterator"
)

// user approved
func GetUserByID(id string, ctx context.Context) (*models.User, int, error) {
	hasUser := config.DB.Collection("User").
		Where("is_deleted", "==", false).
		Where("status", "==", models.Approved).
		Where("id", "==", id).
		Limit(1)

	doc, err := hasUser.Documents(ctx).Next()
	if err != nil {
		if err == iterator.Done {
			return nil, http.StatusNotFound, errors.New("user not found")
		}
		return nil, http.StatusInternalServerError, err
	}
	var user *models.User
	err = doc.DataTo(&user)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return user, http.StatusOK, nil
}

func GetAllUser(ctx context.Context) ([]*models.User, int, error) {
	iter := config.DB.Collection("User").
		Where("is_deleted", "==", false).
		Where("status", "==", models.Approved).
		Documents(ctx)

	defer iter.Stop()

	var users []*models.User
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
		var user models.User
		err = doc.DataTo(&user)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}

		users = append(users, &user)
	}
	return users, http.StatusOK, nil
}

// user pending
func GetUserPending(ctx context.Context) ([]*models.User, int, error) {
	iter := config.DB.Collection("User").
		Where("is_deleted", "==", false).
		Where("status", "==", models.Pending).
		Documents(ctx)
	
	defer iter.Stop()

	var users []*models.User
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			if len(users) == 0 {
				return nil, http.StatusNotFound, errors.New("no pending users found")
			}
			break
		}
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}

		var user models.User
		err = doc.DataTo(&user)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}

		users = append(users, &user)
	}
	return users, http.StatusOK, nil
}
