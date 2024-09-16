package services

import (
	"context"
	"errors"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/alongkornn/Web-VRGame-Backend/config"
	auth_models "github.com/alongkornn/Web-VRGame-Backend/internal/auth/models"
	checkpoint_models "github.com/alongkornn/Web-VRGame-Backend/internal/checkpoint/models"
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
		"updated_at": firestore.ServerTimestamp,
	}, firestore.MergeAll)
	if err != nil {
		return http.StatusInternalServerError, err
 
	}

	return http.StatusOK, nil
}
