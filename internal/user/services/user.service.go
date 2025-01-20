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
	"github.com/alongkornn/Web-VRGame-Backend/internal/user/dto"
	"github.com/alongkornn/Web-VRGame-Backend/pkg/utils"
	"google.golang.org/api/iterator"
)

// แสดงผู้เล่นแค่คนเดียว
func GetUserByID(userId string, ctx context.Context) (*auth_models.User, int, error) {
	// สร้าง key สำหรับ Redis
	userCacheKey := fmt.Sprintf("user:%s", userId)

	// 1. ตรวจสอบใน Redis ก่อน
	cachedUser, err := config.RedisClient.Get(ctx, userCacheKey).Result()
	if err == nil {
		var user auth_models.User
		// หากมีข้อมูลใน Redis ก็จะนำมาใช้เลย
		if err := json.Unmarshal([]byte(cachedUser), &user); err == nil {
			return &user, http.StatusOK, nil
		}
	}

	// 2. ถ้าไม่มีใน Redis -> ไปดึงข้อมูลจาก Firestore
	hasUser := utils.HasUser(userId)

	userDoc, err := hasUser.Documents(ctx).Next()
	if err != nil {
		return nil, http.StatusBadRequest, errors.New("user not found")
	}

	var user auth_models.User
	if err := userDoc.DataTo(&user); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	// 3. เก็บข้อมูลลง Redis
	userData, _ := json.Marshal(user)
	config.RedisClient.Set(ctx, userCacheKey, userData, 10*time.Minute) // ตั้งเวลาหมดอายุใน Redis เป็น 10 นาที

	return &user, http.StatusOK, nil
}

// แสดงผู้เล่นทั้งหมด
func GetAllUser(ctx context.Context) ([]*auth_models.User, int, error) {
	// สร้าง key สำหรับ Redis
	cacheKey := "player:all"

	// 1. ตรวจสอบข้อมูลใน Redis ก่อน
	cachedData, err := config.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil && cachedData != "" {
		// ถ้ามีข้อมูลใน Redis แปลง JSON เป็น struct
		var users []*auth_models.User
		if err := json.Unmarshal([]byte(cachedData), &users); err == nil {
			return users, http.StatusOK, nil
		}
	}

	// 2. ถ้าไม่มีข้อมูลใน Redis -> ดึงข้อมูลจาก Firestore
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

	// 3. แคชข้อมูลลงใน Redis
	if len(users) > 0 {
		data, err := json.Marshal(users)
		if err == nil {
			// ตั้งค่าความหมดอายุ (เช่น 10 นาที)
			config.RedisClient.Set(ctx, cacheKey, data, 10*time.Minute)
		}
	}

	return users, http.StatusOK, nil
}

// แสดงผู้เล่นที่ยังไม่ได้รับการอนุมัติ
func GetUserPending(ctx context.Context) ([]*auth_models.User, int, error) {
	// สร้าง key สำหรับ Redis
	cacheKey := "player:pending"

	// 1. ตรวจสอบข้อมูลใน Redis ก่อน
	cachedData, err := config.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil && cachedData != "" {
		// ถ้ามีข้อมูลใน Redis แปลง JSON เป็น struct
		var users []*auth_models.User
		if err := json.Unmarshal([]byte(cachedData), &users); err == nil {
			return users, http.StatusOK, nil
		}
	}

	// 2. ถ้าไม่มีข้อมูลใน Redis -> ดึงข้อมูลจาก Firestore
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

	// 3. แคชข้อมูลลงใน Redis
	if len(users) > 0 {
		data, err := json.Marshal(users)
		if err == nil {
			// ตั้งค่าความหมดอายุ (เช่น 10 นาที)
			config.RedisClient.Set(ctx, cacheKey, data, 10*time.Minute)
		}
	}

	return users, http.StatusOK, nil
}

// แก้ไขข้อมูลผู้เล่น
func UpdateUser(userId string, updateUserDTO dto.UpdateUserDTO, ctx context.Context) (int, error) {
	hasUser := utils.HasUser(userId)

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

// แสดงคะแนนรวมทั้งหมด
func GetSumScore(userId string, ctx context.Context) (int, int, error) {
	hasUser := utils.HasUser(userId)

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

// รวมคะแนนทั้งหมดที่ผู้เล่นทำได้
func SetSumScore(userId string, ctx context.Context) (int, error) {
	hasUser := utils.HasUser(userId)

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
	var completeCheckpointScore int
	if user.CompletedCheckpoints != nil {
		for _, checkpoint := range user.CompletedCheckpoints {
			completeCheckpointScore += checkpoint.Score
		}
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

func GetUserBySortScore(ctx context.Context) ([]*auth_models.User, int, error) {
	iter := config.DB.Collection("User").Where("is_deleted", "==", false).
		Where("status", "==", auth_models.Approved).
		Where("role", "==", auth_models.Player).
		OrderBy("score", firestore.Desc).
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
