package services

import (
	"context"
	"errors"
	"net/http"

	"github.com/alongkornn/Web-VRGame-Backend/config"
	auth_models "github.com/alongkornn/Web-VRGame-Backend/internal/auth/models"
	score_models "github.com/alongkornn/Web-VRGame-Backend/internal/score/models"
	"google.golang.org/api/iterator"
)

func GetScorebyID(id string, ctx context.Context) (*score_models.Score, int, error) {
	hasUser := config.DB.Collection("User").
		Where("is_deleted", "==", false).
		Where("id", "==", id).
		Limit(1)

	doc, err := hasUser.Documents(ctx).Next()
	if err != nil {
		return nil, http.StatusBadRequest, errors.New("user not found")
	}

	var user auth_models.User
	doc.DataTo(&user)

	score := score_models.Score{
		Name: user.FirstName,
		Score: user.Score,
	}

	

	return &score, http.StatusOK, nil
}

func GetAllScore(ctx context.Context) ([]*score_models.Score, int, error) {
	iter := config.DB.Collection("User").Where("is_deleted", "==", false).Documents(ctx)
	defer iter.Stop()

	var users_score []*score_models.Score
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}

		var user auth_models.User
		// var score score_models.Score
		err = doc.DataTo(&user)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
		
		score := score_models.Score{
			Name: user.FirstName,
			Score: user.Score,
		}

		users_score = append(users_score, &score)

	}
	return users_score, http.StatusOK, nil
} 
