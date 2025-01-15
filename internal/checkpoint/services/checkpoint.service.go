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
	"github.com/alongkornn/Web-VRGame-Backend/pkg/utils"
	"github.com/google/uuid"
	"google.golang.org/api/iterator"
)

// แสดงด่านปัจจุบันของผู้เล่น(โดยจะเข้าถึงผ่านไอดีของผู้เล่น)
func GetCurrentCheckpointFromUserId(userId string, ctx context.Context) (*checkpoint_models.Checkpoints, int, error) {
	hasUser := utils.HasUser(userId)

	userDoc, err := hasUser.Documents(ctx).Next()
	if err != nil {
		return nil, http.StatusNotFound, errors.New("user not found")
	}

	var user auth_models.User
	if err := userDoc.DataTo(&user); err != nil {
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

	return &currentCheckpoint, http.StatusOK, nil
}

// แสดงด่านทั้งหมดทุกหมวดหมู่
func GetAllCheckpoint(ctx context.Context) ([]*checkpoint_models.Checkpoints, int, error) {
	iter := config.DB.Collection("Checkpoint").
		Where("is_deleted", "==", false).
		Documents(ctx)

	defer iter.Stop()

	var checkpoints []*checkpoint_models.Checkpoints
	for {
		checkpointDoc, err := iter.Next()
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
		if err := checkpointDoc.DataTo(&checkpoint); err != nil {
			return nil, http.StatusInternalServerError, err
		}

		checkpoints = append(checkpoints, &checkpoint)
	}
	return checkpoints, http.StatusOK, nil
}

// สร้างด่านใหม่
func CreateCheckpoint(checkpointDTO dto.CreateCheckpointsDTO, ctx context.Context) (int, error) {
	checkpoint := checkpoint_models.Checkpoints{
		ID:         uuid.New().String(),
		Name:       checkpointDTO.Name,
		Category:   checkpointDTO.Category,
		MaxScore:   100,
		PassScore:  50,
		TimeLimit:  "5 นาที",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		Is_Deleted: false,
	}

	_, _, err := config.DB.Collection("Checkpoint").Add(ctx, checkpoint)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// บันทึกด่านปัจจุบันลงในด่านที่เล่นผ่านแล้วโดยจะตรวจสอบว่าคะแนนผ่านเกณฑ์หรือยัง
func SaveCheckpointToComplete(userID string, score int, ctx context.Context) (int, error) {
	hasUser := utils.HasUser(userID)

	userDoc, err := hasUser.Documents(ctx).Next()
	if err != nil {
		return http.StatusNotFound, errors.New("user not found")
	}

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

	completeCheckpoint := checkpoint_models.CompleteCheckpoint{
		CheckpointID: currentCheckpoint.ID,
		Score:        score,
	}

	if score >= currentCheckpoint.PassScore && score <= currentCheckpoint.MaxScore {
		_, err := userDoc.Ref.Update(ctx, []firestore.Update{
			{
				Path:  "completed_checkpoints",
				Value: completeCheckpoint,
			},
			{
				Path:  "updated_at",
				Value: firestore.ServerTimestamp,
			},
			{
				Path:  "current_checkpoint",
				Value: nil,
			},
		})
		if err != nil {
			return http.StatusInternalServerError, err
		}
	}
	return http.StatusOK, nil
}

// แสดงด่านที่ผู้เล่นเล่นผ่าน(โดยจะเข้าถึงผ่านไอดีของผู้เล่น)
func GetCheckpointDetails(userID string, ctx context.Context) ([]checkpoint_models.CheckpointDetail, int, error) {
	var checkpointDetails []checkpoint_models.CheckpointDetail

	hasUser := utils.HasUser(userID)

	userDoc, err := hasUser.Documents(ctx).Next()
	if err != nil {
		return nil, http.StatusNotFound, errors.New("user not found")
	}

	var user auth_models.User
	if err := userDoc.DataTo(&user); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	// วนลูปผ่าน CompletedCheckpoints ของผู้ใช้
	for _, completedCheckpoint := range user.CompletedCheckpoints {
		checkpointID := completedCheckpoint.CheckpointID
		score := completedCheckpoint.Score

		// ดึงข้อมูลรายละเอียดของ checkpoint จาก Firestore
		checkpointRef := config.DB.Collection("Checkpoints").Doc(checkpointID)
		doc, err := checkpointRef.Get(context.Background())
		if err != nil {
			return nil, http.StatusInternalServerError, err // ถ้าไม่พบข้อมูลให้หยุดและแสดง error
		}

		// ดึงข้อมูลจากเอกสารที่ได้รับมา
		var checkpointData map[string]interface{}
		doc.DataTo(&checkpointData)

		// ดึงชื่อด่านและหมวดหมู่
		name := checkpointData["name"].(string)
		category := checkpointData["category"].(string)

		// รวมข้อมูลที่ได้ใน CheckpointDetail
		checkpointDetails = append(checkpointDetails, checkpoint_models.CheckpointDetail{
			CheckpointID: checkpointID,
			Name:         name,
			Category:     category,
			Score:        score,
		})
	}

	return checkpointDetails, http.StatusOK, nil
}

// แสดงทุกด่านตามหมวดหมู่
func GetCheckpointWithCategory(category string, ctx context.Context) ([]*checkpoint_models.Checkpoints, int, error) {
	iter := config.DB.Collection("Checkpoint").
		Where("is_deleted", "==", false).
		Where("category", "==", category).
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

// เพิ่มเวลาในด่านปัจจุบัน
func SetTime(userId string, time time.Duration, ctx context.Context) (int, error) {
	hasUser := utils.HasUser(userId)

	userDoc, err := hasUser.Documents(ctx).GetAll()
	if err != nil || len(userDoc) == 0 {
		return http.StatusBadRequest, errors.New("user not found")
	}

	var user auth_models.User
	if err := userDoc[0].DataTo(&user); err != nil {
		return http.StatusInternalServerError, err
	}

	_, err = userDoc[0].Ref.Update(ctx, []firestore.Update{
		{
			Path:  "user.current_checkpoint.time",
			Value: time,
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
