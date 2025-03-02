package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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
func GetCurrentCheckpointFromUserId(userId string, ctx context.Context) ([]checkpoint_models.Checkpoints, int, error) {
	// ตรวจสอบว่าผู้ใช้มีอยู่จริงหรือไม่
	userQuery := utils.HasUser(userId)
	userSnapshot, err := userQuery.Documents(ctx).GetAll()
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error querying Firestore: %v", err)
	}

	if len(userSnapshot) == 0 {
		return nil, http.StatusNotFound, errors.New("user not found")
	}

	// ดึงข้อมูลผู้ใช้
	var user auth_models.User
	if err := userSnapshot[0].DataTo(&user); err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("error unmarshaling user data: %v", err)
	}

	// ตรวจสอบว่าผู้ใช้มี current checkpoint หรือไม่
	checkpointIDs := []string{
		user.ProjectileCurrentCheckpoint,
		user.MomentumCurrentCheckpoint,
		user.ForceCurrentCheckpoint,
	}

	// ตรวจสอบว่ามี checkpoint ไหนที่ยังไม่ได้กำหนด
	for _, checkpointID := range checkpointIDs {
		if checkpointID == "" {
			return nil, http.StatusBadRequest, errors.New("user does not have all required checkpoints")
		}
	}

	// สร้าง slice สำหรับเก็บ Checkpoints
	var currentCheckpoints []checkpoint_models.Checkpoints

	// ดึงข้อมูลจาก Firestore
	for _, checkpointID := range checkpointIDs {
		checkpointQuery := config.DB.Collection("Checkpoint").Where("id", "==", checkpointID)
		checkpointSnapshot, err := checkpointQuery.Documents(ctx).GetAll()
		if err != nil {
			return nil, http.StatusInternalServerError, fmt.Errorf("error getting checkpoint: %v", err)
		}
		if len(checkpointSnapshot) == 0 {
			return nil, http.StatusNotFound, fmt.Errorf("checkpoint not found: %s", checkpointID)
		}

		// แปลงข้อมูลจาก Firestore เป็น struct
		var checkpoint checkpoint_models.Checkpoints
		if err := checkpointSnapshot[0].DataTo(&checkpoint); err != nil {
			return nil, http.StatusInternalServerError, fmt.Errorf("error unmarshaling checkpoint data: %v", err)
		}

		// เพิ่ม Checkpoint ลงใน Slice
		currentCheckpoints = append(currentCheckpoints, checkpoint)
	}

	return currentCheckpoints, http.StatusOK, nil
}

// แสดงด่านทั้งหมดทุกหมวดหมู่
func GetAllCheckpoint(ctx context.Context) ([]*checkpoint_models.Checkpoints, int, error) {
	// สร้าง key สำหรับ Redis
	cacheKey := "checkpoints:all"

	// 1. ตรวจสอบข้อมูลใน Redis ก่อน
	cachedData, err := config.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil && cachedData != "" {
		// ถ้ามีข้อมูลใน Redis แปลง JSON เป็น struct
		var checkpoints []*checkpoint_models.Checkpoints
		if err := json.Unmarshal([]byte(cachedData), &checkpoints); err == nil {
			return checkpoints, http.StatusOK, nil
		}
	}

	// 2. ถ้าไม่มีข้อมูลใน Redis -> ดึงข้อมูลจาก Firestore
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

	// 3. เก็บข้อมูลลง Redis
	if len(checkpoints) > 0 {
		data, err := json.Marshal(checkpoints)
		if err == nil {
			// ตั้งค่าความหมดอายุ (เช่น 10 นาที)
			config.RedisClient.Set(ctx, cacheKey, data, 10*time.Minute)
		}
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
func SaveCheckpointToComplete(userID string, score int, time string, ctx context.Context) (int, error) {
	// ตรวจสอบว่ามีผู้ใช้หรือไม่
	hasUser := utils.HasUser(userID)
	userDoc, err := hasUser.Documents(ctx).Next()
	if err != nil {
		return http.StatusNotFound, errors.New("user not found")
	}

	// แปลงข้อมูลผู้ใช้จาก Firestore
	var user auth_models.User
	if err := userDoc.DataTo(&user); err != nil {
		return http.StatusInternalServerError, err
	}

	// ตรวจสอบว่ามี Checkpoint ปัจจุบันหรือไม่
	hasCheckpoint := utils.GetCheckpointByID(user.ProjectileCurrentCheckpoint)
	checkpointDoc, err := hasCheckpoint.Documents(ctx).Next()
	if err != nil {
		return http.StatusNotFound, errors.New("checkpoint not found")
	}

	// แปลงข้อมูล Checkpoint
	var currentCheckpoint checkpoint_models.Checkpoints
	if err := checkpointDoc.DataTo(&currentCheckpoint); err != nil {
		return http.StatusInternalServerError, err
	}

	completeCheckpoint := &checkpoint_models.CompleteCheckpoint{
		CheckpointID: currentCheckpoint.ID,
		Name:         currentCheckpoint.Name,
		Category:     currentCheckpoint.Category,
		Score:        score,
		Time:         time,
	}

	// ตรวจสอบว่าผู้ใช้ผ่าน checkpoint หรือไม่
	if score >= currentCheckpoint.PassScore && score <= currentCheckpoint.MaxScore {
		var completedCheckpoints []*checkpoint_models.CompleteCheckpoint

		// ถ้ามี completedCheckpoints อยู่แล้ว ให้นำมาใส่ใน slice
		if user.CompletedCheckpoints != nil {
			completedCheckpoints = append(completedCheckpoints, user.CompletedCheckpoints...)
		}

		completedCheckpoints = append(completedCheckpoints, completeCheckpoint)

		var firestoreCheckpoints []map[string]interface{}
		for _, c := range completedCheckpoints {
			firestoreCheckpoints = append(firestoreCheckpoints, map[string]interface{}{
				"checkpoint_id": c.CheckpointID,
				"name":          c.Name,
				"category":      c.Category,
				"score":         c.Score,
			})
		}

		// อัปเดต Firestore
		_, err := userDoc.Ref.Update(ctx, []firestore.Update{
			{
				Path:  "completed_checkpoints",
				Value: firestoreCheckpoints,
			},
			{
				Path:  "updated_at",
				Value: firestore.ServerTimestamp,
			},
		})
		if err != nil {
			return http.StatusInternalServerError, err
		}

		user.CompletedCheckpoints = completedCheckpoints

		completedCheckpointsJSON, err := json.Marshal(completedCheckpoints)
		if err != nil {
			return http.StatusInternalServerError, errors.New("failed to encode data for Redis")
		}

		redisKey := fmt.Sprintf("user:%s", userID) // ใช้ key แบบ user:id
		err = config.RedisClient.Set(ctx, redisKey, completedCheckpointsJSON, 0).Err()
		if err != nil {
			return http.StatusInternalServerError, errors.New("failed to update status in Redis")
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