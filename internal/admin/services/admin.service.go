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
	"github.com/alongkornn/Web-VRGame-Backend/internal/admin/dto"
	auth_models "github.com/alongkornn/Web-VRGame-Backend/internal/auth/models"
	checkpoint_models "github.com/alongkornn/Web-VRGame-Backend/internal/checkpoint/models"
	"github.com/alongkornn/Web-VRGame-Backend/pkg/utils"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/api/iterator"
)

// ผู้ดูแลระบบอนุมัติการลงทะเบียนของผู้เล่น
func AddminApprovedUserRegister(userId string, approved auth_models.Status, ctx context.Context) (int, error) {
	hasUser := config.DB.Collection("User").
		Where("is_deleted", "==", false).
		Where("status", "==", auth_models.Pending).
		Where("id", "==", userId).
		Limit(1)

	userDoc, err := hasUser.Documents(ctx).Next()
	if err != nil {
		return http.StatusNotFound, errors.New("user not found")
	}

	var user auth_models.User
	err = userDoc.DataTo(&user)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	_, err = userDoc.Ref.Update(ctx, []firestore.Update{
		{
			Path:  "status",
			Value: auth_models.Approved,
		},
		{
			Path:  "updated_at",
			Value: firestore.ServerTimestamp,
		},
	})
	if err != nil {
		return http.StatusBadRequest, errors.New("failed to approve")
	}
	return http.StatusOK, nil
}

// ผู้ดูแลระบบลบผู้เล่นออก
func AdminRemoveUser(userId string, ctx context.Context) (int, error) {
	hasUser := utils.HasUser(userId)

	userDoc, err := hasUser.Documents(ctx).Next()
	if err != nil {
		return http.StatusNotFound, errors.New("user not found")
	}

	var user auth_models.User
	err = userDoc.DataTo(&user)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	_, err = userDoc.Ref.Update(ctx, []firestore.Update{
		{
			Path:  "is_deleted",
			Value: true,
		},
		{
			Path:  "updated_at",
			Value: firestore.ServerTimestamp,
		},
		{
			Path:  "status",
			Value: auth_models.Deleted,
		},
	})
	if err != nil {
		return http.StatusBadRequest, err
	}
	return http.StatusOK, nil
}

// ลบผู้ดูแลระบบออก
func RemoveAdmin(adminId string, ctx context.Context) (int, error) {
	hasAdmin := utils.HasAdmin(adminId)

	adminDoc, err := hasAdmin.Documents(ctx).Next()
	if err != nil {
		return http.StatusNotFound, errors.New("admin not found")
	}

	var user auth_models.User
	if err = adminDoc.DataTo(&user); err != nil {
		return http.StatusInternalServerError, err
	}

	_, err = adminDoc.Ref.Update(ctx, []firestore.Update{
		{
			Path:  "is_deleted",
			Value: true,
		},
		{
			Path:  "role",
			Value: auth_models.Player,
		},
		{
			Path:  "updated_at",
			Value: firestore.ServerTimestamp,
		},
	})
	if err != nil {
		return http.StatusBadRequest, err
	}

	return http.StatusOK, nil
}

// แสดงผู้ดูแลระบบทั้งหมด
func GetAllAdmin(ctx context.Context) ([]*auth_models.User, int, error) {
	iter := config.DB.Collection("User").
		Where("role", "==", auth_models.Admin).
		Where("is_deleted", "==", false).
		Where("status", "==", auth_models.Approved).
		Documents(ctx)

	defer iter.Stop()

	var admins []*auth_models.User

	for {
		adminDoc, err := iter.Next()
		if err == iterator.Done {
			if len(admins) == 0 {
				return nil, http.StatusNotFound, errors.New("admin not found")
			}
			break
		}
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}

		var admin auth_models.User
		if err = adminDoc.DataTo(&admin); err != nil {
			return nil, http.StatusInternalServerError, err
		}

		admins = append(admins, &admin)
	}
	return admins, http.StatusOK, nil
}

// แสดงผู้ดูแลระบบโดยเข้าถึงผ่านไอดีผู้ดูแลระบบ
func GetAdminById(adminId string, ctx context.Context) (*auth_models.User, int, error) {
	// สร้าง key สำหรับ Redis
	adminCacheKey := fmt.Sprintf("admin:%s", adminId)

	// 1. ตรวจสอบใน Redis ก่อน
	cachedAdmin, err := config.RedisClient.Get(ctx, adminCacheKey).Result()
	if err == nil {
		var admin auth_models.User
		// หากมีข้อมูลใน Redis ก็จะนำมาใช้เลย
		if err := json.Unmarshal([]byte(cachedAdmin), &admin); err == nil {
			return &admin, http.StatusOK, nil
		}
	}

	// 2. ถ้าไม่มีใน Redis -> ไปดึงข้อมูลจาก Firestore
	hasAdmin := utils.HasAdmin(adminId)

	adminDoc, err := hasAdmin.Documents(ctx).Next()
	if err != nil {
		return nil, http.StatusNotFound, errors.New("admin not found")
	}

	var admin auth_models.User
	if err := adminDoc.DataTo(&admin); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	// 3. เก็บข้อมูลลง Redis
	adminData, _ := json.Marshal(admin)
	config.RedisClient.Set(ctx, adminCacheKey, adminData, 10*time.Minute) // ตั้งเวลาหมดอายุใน Redis เป็น 10 นาที

	return &admin, http.StatusOK, nil
}

// เพิ่มผู้ดูแลระบบ
func CreateAdmin(userId string, role auth_models.Role, ctx context.Context) (int, error) {
	hasUser := utils.HasUser(userId)

	userdoc, err := hasUser.Documents(ctx).Next()
	if err != nil {
		return http.StatusNotFound, errors.New("user not found")
	}

	var user auth_models.User
	if err := userdoc.DataTo(&user); err != nil {
		return http.StatusInternalServerError, err
	}

	// อัปเดตข้อมูลของ user ใน Firestore
	_, err = userdoc.Ref.Update(ctx, []firestore.Update{
		{
			Path:  "role",
			Value: auth_models.Admin,
		},
		{
			Path:  "updated_at",
			Value: firestore.ServerTimestamp,
		},
	})
	if err != nil {
		return http.StatusInternalServerError, errors.New("failed to update user role")
	}

	return http.StatusCreated, nil
}

// แก้ไขข้อมูลผู้ดูแลระบบ
func UpdateDataAdmin(adminId string, updateDTO dto.UpdateDTO, ctx context.Context) (int, error) {
	hasAdmin := utils.HasAdmin(adminId)

	adminDoc, err := hasAdmin.Documents(ctx).Next()
	if err != nil {
		return http.StatusNotFound, errors.New("admin not found")
	}

	var admin auth_models.User
	err = adminDoc.DataTo(&admin)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	updateData := make(map[string]interface{})

	if updateDTO.FirstName != "" {
		updateData["firstname"] = updateDTO.FirstName
	}

	if updateDTO.LastName != "" {
		updateData["lastname"] = updateDTO.LastName
	}

	if updateDTO.Class != "" {
		updateData["class"] = updateDTO.Class
	}

	if updateDTO.Number != "" {
		updateData["number"] = updateDTO.Number
	}

	updateData["updated_at"] = firestore.ServerTimestamp

	if len(updateData) == 0 {
		return http.StatusBadRequest, errors.New("no data to update")
	}

	// อัปเดตเฉพาะข้อมูลที่มีการเปลี่ยนแปลง
	_, err = adminDoc.Ref.Set(ctx, updateData, firestore.MergeAll)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// แก้ไขรหัสผ่านของผู้ดูแลระบบ
func UpdatePasswordAdmin(adminId, password, newPassword string, ctx context.Context) (int, error) {
	hasAdmin := utils.HasAdmin(adminId)

	adminDoc, err := hasAdmin.Documents(ctx).Next()
	if err != nil {
		return http.StatusNotFound, errors.New("admin not found")
	}

	var admin auth_models.User
	err = adminDoc.DataTo(&admin)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password)); err != nil {
		return http.StatusBadRequest, errors.New("invalid password")
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return http.StatusBadRequest, errors.New("failed to hash")
	}

	_, err = adminDoc.Ref.Update(ctx, []firestore.Update{
		{
			Path:  "password",
			Value: hashPassword,
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

// แสดงจุดเด่นของผู้เล่น
func ShowScoreWiteStrength(userId string, ctx context.Context) ([]*checkpoint_models.Category, int, error) {
	hasUser := utils.HasUser(userId)

	userDoc, err := hasUser.Documents(ctx).GetAll()
	if err != nil || len(userDoc) == 0 {
		return nil, http.StatusBadRequest, errors.New("user not found")
	}

	var user auth_models.User
	if err := userDoc[0].DataTo(&user); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	isStrength := []*checkpoint_models.Category{}

	if user.CompletedCheckpoints != nil {
		for _, checkpoint := range user.CompletedCheckpoints {
			if checkpoint.Score >= 80 {
				hasCheckpoint := utils.GetCheckpointByID(checkpoint.CheckpointID)

				checkpointDoc, err := hasCheckpoint.Documents(ctx).Next()
				if err != nil {
					return nil, http.StatusNotFound, errors.New("checkpoint not found")
				}

				var currentCheckpoint checkpoint_models.Checkpoints
				if err := checkpointDoc.DataTo(&currentCheckpoint); err != nil {
					return nil, http.StatusInternalServerError, err
				}
				isStrength = append(isStrength, &currentCheckpoint.Category)
			}
		}
	}

	uniqueStrength := removeDuplicates(isStrength)

	return uniqueStrength, http.StatusOK, nil
}

// แสดงจุดด้อยของผู้เล่น
func ShowScoreWiteWeaknesses(userId string, ctx context.Context) ([]*checkpoint_models.Category, int, error) {
	hasUser := utils.HasUser(userId)

	userDoc, err := hasUser.Documents(ctx).GetAll()
	if err != nil || len(userDoc) == 0 {
		return nil, http.StatusBadRequest, errors.New("user not found")
	}

	var user auth_models.User
	if err := userDoc[0].DataTo(&user); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	isWeaknesses := []*checkpoint_models.Category{}

	if user.CompletedCheckpoints != nil {
		for _, checkpoint := range user.CompletedCheckpoints {
			if checkpoint.Score <= 50 {
				hasCheckpoint := utils.GetCheckpointByID(checkpoint.CheckpointID)

				checkpointDoc, err := hasCheckpoint.Documents(ctx).Next()
				if err != nil {
					return nil, http.StatusNotFound, errors.New("checkpoint not found")
				}

				var currentCheckpoint checkpoint_models.Checkpoints
				if err := checkpointDoc.DataTo(&currentCheckpoint); err != nil {
					return nil, http.StatusInternalServerError, err
				}

				isWeaknesses = append(isWeaknesses, &currentCheckpoint.Category)
			}
		}
	}

	uniqueWeaknesses := removeDuplicates(isWeaknesses)

	return uniqueWeaknesses, http.StatusOK, nil
}

func removeDuplicates(categories []*checkpoint_models.Category) []*checkpoint_models.Category {
	// ใช้ map เพื่อเก็บค่าที่ไม่ซ้ำกัน โดยใช้ pointer ของ Category เป็น key
	uniqueCategories := make(map[*checkpoint_models.Category]struct{})

	// วนลูปผ่าน categories แล้วเพิ่มลงใน map
	for _, category := range categories {
		uniqueCategories[category] = struct{}{} // struct{} เป็นค่าเปล่าที่ใช้เก็บข้อมูลใน map
	}

	// สร้าง slice ใหม่เพื่อเก็บค่าที่ไม่ซ้ำกัน
	var result []*checkpoint_models.Category
	for category := range uniqueCategories {
		result = append(result, category)
	}

	return result
}
