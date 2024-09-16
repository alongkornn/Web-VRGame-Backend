package services

import (
	"context"
	"errors"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/alongkornn/Web-VRGame-Backend/config"
	"github.com/alongkornn/Web-VRGame-Backend/internal/auth/models"
	"google.golang.org/api/iterator"
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
		"updated_at": firestore.ServerTimestamp,
	}, firestore.MergeAll)
	if err != nil {
		return http.StatusBadRequest, errors.New("failed to approve")
	}
	return http.StatusOK, nil
}

func RemoveUser(id string, ctx context.Context) (int, error) {
	hasUser := config.DB.Collection("User").
		Where("is_deleted", "==", false).
		Where("role", "==", models.Player).
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
		"updated_at": firestore.ServerTimestamp,
	}, firestore.MergeAll)
	if err != nil {
		return http.StatusBadRequest, err
	}
	return http.StatusOK, nil
}

func RemoveAdmin(id string, ctx context.Context) (int, error) {
	hasUser := config.DB.Collection("User").
		Where("is_deleted", "==", false).
		Where("role", "==", models.Admin).
		Where("id", "==", id).
		Limit(1)

	doc, err := hasUser.Documents(ctx).Next()
	if err != nil {
		return http.StatusNotFound, errors.New("admin not found")
	}

	var user models.User
	if err = doc.DataTo(&user); err != nil {
		return http.StatusInternalServerError, err
	}

	_, err = doc.Ref.Set(ctx, map[string]interface{}{
		"role": models.Player,
		"updated_at": firestore.ServerTimestamp,
	}, firestore.MergeAll)
	if err != nil {
		return http.StatusBadRequest, err
	}

	return http.StatusOK, nil
}

func GetAllAdmin(ctx context.Context) ([]*models.User, int, error) {
	iter := config.DB.Collection("User").
		Where("role", "==", models.Admin).
		Where("is_deleted", "==", false).
		Documents(ctx)
	
	defer iter.Stop()

	var users []*models.User

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			if len(users) == 0 {
				return nil, http.StatusNotFound, errors.New("admin not found")
			}
			break
		}
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}

		var user models.User
		if err = doc.DataTo(&user); err != nil {
			return nil, http.StatusInternalServerError, err
		}

		users = append(users, &user)
	}
	return users, http.StatusOK, nil
}

func GetAdminByID(id string, ctx context.Context) (*models.User, int, error) {
	hasUser := config.DB.Collection("User").
		Where("is_deleted", "==", false).
		Where("role", "==", models.Admin).
		Where("id", "==", id).
		Limit(1)

	doc, err := hasUser.Documents(ctx).Next()
	if err != nil {
		return nil, http.StatusNotFound, errors.New("admin not found")
	}

	var user *models.User
	if err = doc.DataTo(&user); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return user, http.StatusOK, nil
}

// admin
func CreateAdmin(id, role string, ctx context.Context) (int, error) {
	hasUser := config.DB.Collection("User").
		Where("is_deleted", "==", false).
		Where("role", "==", models.Player).
		Where("id", "==", id).
		Limit(1)
	doc, err := hasUser.Documents(ctx).Next()
	if err != nil {
		return http.StatusNotFound, errors.New("user not found")
	}

	var user models.User
	doc.DataTo(&user)

	// อัปเดตข้อมูลของ user ใน Firestore
	_, err = doc.Ref.Set(ctx, map[string]interface{}{
		"role": models.Admin,
		"updated_at": firestore.ServerTimestamp,
	}, firestore.MergeAll)
	if err != nil {
		return http.StatusInternalServerError, errors.New("failed to update user role")
	}

	return http.StatusCreated, nil
}

