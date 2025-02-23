package services

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

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
		Where("status", "==", "approved").
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
	// ตรวจสอบว่าผู้ใช้มีอยู่ในระบบหรือไม่
	hasUser := utils.HasUser(userId)
	userDoc, err := hasUser.Documents(ctx).Next()
	if err != nil {
		return http.StatusNotFound, errors.New("user not found")
	}

	// แปลงข้อมูลผู้ใช้
	var user auth_models.User
	if err := userDoc.DataTo(&user); err != nil {
		return http.StatusInternalServerError, err
	}

	hasCheckpoint := utils.GetCheckpointByID(user.CurrentCheckpoint)
	checkpointDoc, err := hasCheckpoint.Documents(ctx).Next()
	if err != nil {
		return http.StatusNotFound, errors.New("checkpoint not found")
	}

	var currentCheckpoint checkpoint_models.Checkpoints
	if err := checkpointDoc.DataTo(&currentCheckpoint); err != nil {
		return http.StatusInternalServerError, err
	}

	// ตรวจสอบว่าคะแนนมากกว่าหรือเท่ากับคะแนนที่ต้องผ่าน
	if score >= currentCheckpoint.PassScore {
		completedCheckpoint := checkpoint_models.CompleteCheckpoint{
			CheckpointID: currentCheckpoint.ID,
			Name:         currentCheckpoint.Name,
			Category:     currentCheckpoint.Category,
			Score:        score,
		}

		sumScore := user.Score + score

		_, err = userDoc.Ref.Update(ctx, []firestore.Update{
			{
				Path:  "completed_checkpoint",
				Value: firestore.ArrayUnion(completedCheckpoint),
			},
			{
				Path:  "current_checkpoint",
				Value: currentCheckpoint.NextCheckpoint,
			},
			{
				Path:  "score",
				Value: sumScore,
			},
			{
				Path:  "updated_at",
				Value: time.Now(),
			},
		})
		if err != nil {
			return http.StatusInternalServerError, err
		}

		err = config.UpdateCurrentCheckpointInRealtimeDB(userId, currentCheckpoint.NextCheckpoint, sumScore)
		if err != nil {
			log.Printf("Failed to update level in Realtime Database: %v\n", err)
			return http.StatusInternalServerError, errors.New("failed to update current level in Realtime Database")
		}
	}

	return http.StatusOK, nil
}
