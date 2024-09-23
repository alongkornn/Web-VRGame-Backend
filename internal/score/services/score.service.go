package services

import (
	"context"
	"errors"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/alongkornn/Web-VRGame-Backend/config"
	auth_models "github.com/alongkornn/Web-VRGame-Backend/internal/auth/models"
	score_models "github.com/alongkornn/Web-VRGame-Backend/internal/score/models"
	"google.golang.org/api/iterator"
)

func GetScoreByUserId(userId string, ctx context.Context) (*score_models.Score, int, error) {
	hasUser := config.DB.Collection("User").
		Where("is_deleted", "==", false).
		Where("role", "==", auth_models.Player).
		Where("status", "==", auth_models.Approved).
		Where("id", "==", userId).
		Limit(1)

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

func GetAllScore(ctx context.Context) ([]*score_models.Score, int, error) {
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

		score := score_models.Score{
			Name:  user.FirstName,
			Score: user.Score,
		}

		users_score = append(users_score, &score)

	}
	return users_score, http.StatusOK, nil
}

func SetScore(userId string, score int, ctx context.Context) (int, error) {
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

	var user auth_models.User
	if err := userDoc.DataTo(&user); err != nil {
		return http.StatusInternalServerError, err
	}

	sumScore := user.Score + score
	maxScore := user.CurrentCheckpoint.MaxScore

	if sumScore <= maxScore {
		_, err = userDoc.Ref.Update(ctx, []firestore.Update{
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
	} else {
		return http.StatusBadRequest, errors.New("score exceeds maximum score")
	}

	return http.StatusOK, nil
}
