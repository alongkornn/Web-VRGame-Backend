package services

import (
	"context"
	"errors"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/alongkornn/Web-VRGame-Backend/config"
	auth_models "github.com/alongkornn/Web-VRGame-Backend/internal/auth/models"
	"github.com/alongkornn/Web-VRGame-Backend/internal/user/dto"
	"google.golang.org/api/iterator"
)

func GetUserByID(userId string, ctx context.Context) (*auth_models.User, int, error) {
	hasUser := config.DB.Collection("User").
		Where("is_deleted", "==", false).
		Where("status", "==", auth_models.Approved).
		Where("role", "==", auth_models.Player).
		Where("id", "==", userId).
		Limit(1)

	userDoc, err := hasUser.Documents(ctx).Next()
	if err != nil {
		return nil, http.StatusBadRequest, errors.New("user not found")
	}

	var user *auth_models.User
	if err := userDoc.DataTo(user); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return user, http.StatusOK, nil
}

func GetAllUser(ctx context.Context) ([]*auth_models.User, int, error) {
	iter := config.DB.Collection("User").
		Where("is_deleted", "==", false).
		Where("status", "==", auth_models.Approved).
		Where("role", "==", auth_models.Player).
		Documents(ctx)

	defer iter.Stop()

	var users []*auth_models.User
	for {
		userDoc, err := iter.Next()
		if err == iterator.Done {
			if len(users) == 0 {
				return nil, http.StatusNotFound, errors.New("user not found")
			}
			break
		}
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
		var user auth_models.User
		if err := userDoc.DataTo(&user); err != nil {
			return nil, http.StatusInternalServerError, err
		}
		users = append(users, &user)
	}
	return users, http.StatusOK, nil
}

func GetUserPending(ctx context.Context) ([]*auth_models.User, int, error) {
	iter := config.DB.Collection("User").
		Where("is_deleted", "==", false).
		Where("status", "==", auth_models.Pending).
		Documents(ctx)

	defer iter.Stop()

	var users []*auth_models.User
	for {
		userDoc, err := iter.Next()
		if err == iterator.Done {
			if len(users) == 0 {
				return nil, http.StatusNotFound, errors.New("no pending users found")
			}
			break
		}
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}

		var user auth_models.User
		if err := userDoc.DataTo(&user); err != nil {
			return nil, http.StatusInternalServerError, err
		}
		users = append(users, &user)
	}
	return users, http.StatusOK, nil
}

func UpdateUser(userId string, updateUserDTO dto.UpdateUserDTO, ctx context.Context) (int, error) {
	hasUser := config.DB.Collection("User").
		Where("is_deleted", "==", false).
		Where("role", "==", auth_models.Player).
		Where("status", "==", auth_models.Approved).
		Where("id", "==", userId).
		Limit(1)

	userDoc, err := hasUser.Documents(ctx).Next()
	if err != nil {
		return http.StatusNotFound, errors.New("user not found")
	}

	updateUser := make(map[string]interface{})

	if updateUserDTO.FirstName != "" {
		updateUser["firstname"] = updateUserDTO.FirstName
	}

	if updateUserDTO.LastName != "" {
		updateUser["lastname"] = updateUserDTO.LastName
	}

	if updateUserDTO.Class != "" {
		updateUser["class"] = updateUserDTO.Class
	}

	if updateUserDTO.Number != "" {
		updateUser["number"] = updateUserDTO.Number
	}

	updateUser["updated_at"] = firestore.ServerTimestamp

	var user auth_models.User
	if err := userDoc.DataTo(&user); err != nil {
		return http.StatusInternalServerError, err
	}

	_, err = userDoc.Ref.Set(ctx, updateUser, firestore.MergeAll)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func GetSumScore(userId string, ctx context.Context) (int, int, error) {
	hasUser := config.DB.Collection("User").
		Where("is_deleted", "==", false).
		Where("role", "==", auth_models.Player).
		Where("status", "==", auth_models.Approved).
		Limit(1)

	userDoc, err := hasUser.Documents(ctx).GetAll()
	if err != nil || len(userDoc) == 0 {
		return 0, http.StatusBadRequest, errors.New("user not found")
	}

	var user auth_models.User
	if err := userDoc[0].DataTo(&user); err != nil {
		return 0, http.StatusInternalServerError, err
	}

	sumScore := user.Score

	return sumScore, http.StatusOK, nil
}

func SetSumScore(userId string, ctx context.Context) (int, error) {
	hasUser := config.DB.Collection("User").
		Where("is_deleted", "==", false).
		Where("role", "==", auth_models.Player).
		Where("status", "==", auth_models.Approved).
		Where("id", "==", userId).
		Limit(1)

	userDocs, err := hasUser.Documents(ctx).GetAll()
	if err != nil || len(userDocs) == 0 {
		return http.StatusBadRequest, errors.New("user not found")
	}

	var user auth_models.User
	if err := userDocs[0].DataTo(&user); err != nil {
		return http.StatusInternalServerError, err
	}
	
	// ตรวจสอบว่า CompletedCheckpoints มีข้อมูล
	var sumScore int
	if user.CompletedCheckpoints != nil {
		for _, checkpoint := range user.CompletedCheckpoints {
			sumScore += checkpoint.Score
		}
	}

	// เพิ่มคะแนนจาก current_checkpoint หากมีค่า
	if user.CurrentCheckpoint != nil {
		sumScore += user.CurrentCheckpoint.Score
	}

	_, err = userDocs[0].Ref.Update(ctx, []firestore.Update{
		{
			Path:  "score",
			Value: sumScore,
		},
		{
			Path:  "updated_at",
			Value: firestore.ServerTimestamp,
		},
	})
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
