package services

import (
	"context"
	"errors"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/alongkornn/Web-VRGame-Backend/config"
	auth_models "github.com/alongkornn/Web-VRGame-Backend/internal/auth/models"
	"github.com/alongkornn/Web-VRGame-Backend/internal/checkpoint/dto"
	checkpoint_models "github.com/alongkornn/Web-VRGame-Backend/internal/checkpoint/models"
	"github.com/google/uuid"
	"google.golang.org/api/iterator"
)

func GetCurrentCheckpointToUser(checkpointId, userId string, ctx context.Context) (int, error) {
	checkpointID := config.DB.Collection("Checkpoint").
		Where("is_deleted", "==", false).
		Where("id", "==", checkpointId).
		Limit(1)

	checkpointDoc, err := checkpointID.Documents(ctx).Next()
	if err != nil {
		return http.StatusNotFound, errors.New("checkpoint not found")
	}
	var checkpoint checkpoint_models.Checkpoints
	if err = checkpointDoc.DataTo(&checkpoint); err != nil {
		return http.StatusInternalServerError, err
	}

	hasUser := config.DB.Collection("User").
		Where("is_deleted", "==", false).
		Where("role", "==", auth_models.Player).
		Where("id", "==", userId).
		Limit(1)

	userDoc, err := hasUser.Documents(ctx).Next()
	if err != nil {
		return http.StatusNotFound, errors.New("checkpoint not found")
	}
	var user auth_models.User
	if err := userDoc.DataTo(&user); err != nil {
		return http.StatusInternalServerError, err
	}

	_, err = userDoc.Ref.Set(ctx, map[string]interface{}{
		"current_checkpoint": checkpoint,
		"updated_at":         firestore.ServerTimestamp,
	}, firestore.MergeAll)
	if err != nil {
		return http.StatusInternalServerError, err

	}

	return http.StatusOK, nil
}

func GetAllCheckpoint(ctx context.Context) ([]*checkpoint_models.Checkpoints, int, error) {
	iter := config.DB.Collection("Checkpoint").
		Where("is_deleted", "==", false).
		Documents(ctx)
	defer iter.Stop()

	var checkpoints []*checkpoint_models.Checkpoints
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			if len(checkpoints) == 0 {
				return nil, http.StatusNotFound, errors.New("checkpoint is empty")
			}
			break
		}
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}

		var checkpoint checkpoint_models.Checkpoints
		if err := doc.DataTo(&checkpoint); err != nil {
			return nil, http.StatusInternalServerError, err
		}

		checkpoints = append(checkpoints, &checkpoint)
	}
	return checkpoints, http.StatusOK, nil
}

func CreateCheckpoint(checkpointDTO dto.CreateCheckpointsDTO, ctx context.Context) (int, error) {
	checkpoint := checkpoint_models.Checkpoints{
		ID: uuid.New().String(),
		Name: checkpointDTO.Name,
		MaxScore: checkpointDTO.MaxScore,
		PassScore: checkpointDTO.PassScore,
		PlayerScore: nil,
		Category: checkpointDTO.Category,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Is_Deleted: false,
	}

	_, _, err := config.DB.Collection("Checkpoint").Add(ctx, checkpoint)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}