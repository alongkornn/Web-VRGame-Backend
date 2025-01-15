package services

import (
	"context"
	"errors"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/alongkornn/Web-VRGame-Backend/config"
	auth_models "github.com/alongkornn/Web-VRGame-Backend/internal/auth/models"
	checkpoint_models "github.com/alongkornn/Web-VRGame-Backend/internal/checkpoint/models"
	score_models "github.com/alongkornn/Web-VRGame-Backend/internal/score/models"
	"github.com/alongkornn/Web-VRGame-Backend/pkg/utils"
	"google.golang.org/api/iterator"
)

func GetScoreByUserId(userId string, ctx context.Context) (*score_models.Score, int, error) {
	hasUser := utils.HasUser(userId)

	userDoc, err := hasUser.Documents(ctx).Next()
	if err != nil {
		return nil, http.StatusBadRequest, errors.New("user not found")
	}

	var user auth_models.User
	if err := userDoc.DataTo(&user); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	score := score_models.Score{
		Name:  user.FirstName,
		Score: user.Score,
	}

	return &score, http.StatusOK, nil
}

func GetAllScoreByCheckpointId(checkpointId string, ctx context.Context) ([]*score_models.Score, int, error) {
	iter := config.DB.Collection("User").
		Where("is_deleted", "==", false).
		Where("role", "==", auth_models.Player).
		Where("status", "==", auth_models.Approved).
		Documents(ctx)

	defer iter.Stop()

	var users_score []*score_models.Score
	for {
		userDoc, err := iter.Next()
		if err == iterator.Done {
			if len(users_score) == 0 {
				return nil, http.StatusNotFound, errors.New("user not found")
			}
			break
		}
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}

		var user auth_models.User
		err = userDoc.DataTo(&user)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}


		hasCheckpoint := utils.GetCheckpointByID(user.CurrentCheckpoint)

		checkpointDoc, err := hasCheckpoint.Documents(ctx).Next()
		if err != nil {
			return nil, http.StatusNotFound, errors.New("checkpoint not found")
		}

		var currentCheckpoint checkpoint_models.Checkpoints
		if err := checkpointDoc.DataTo(&currentCheckpoint); err != nil {
			return nil, http.StatusInternalServerError, err
		}

		if user.CurrentCheckpoint == checkpointId {
			score := score_models.Score{
				CheckpointName: currentCheckpoint.Name,
				Category:       currentCheckpoint.Category,
				Name:           user.FirstName,
			}
			users_score = append(users_score, &score)
		} else {
			return nil, http.StatusBadRequest, errors.New("checkpoin id not found")
		}

	}
	return users_score, http.StatusOK, nil
}

func SetScore(userId string, score int, ctx context.Context) (int, error) {
	hasUser := utils.HasUser(userId)

	userDoc, err := hasUser.Documents(ctx).Next()
	if err != nil {
		return http.StatusNotFound, errors.New("user not found")
	}

	var user auth_models.User
	if err := userDoc.DataTo(&user); err != nil {
		return http.StatusInternalServerError, err
	}

	_, err = userDoc.Ref.Update(ctx, []firestore.Update{
		{
			Path:  "user.score",
			Value: score,
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
