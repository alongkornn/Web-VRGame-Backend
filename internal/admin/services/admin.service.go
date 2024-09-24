package services

import (
	"context"
	"errors"
	"net/http"

	"cloud.google.com/go/firestore"
	"github.com/alongkornn/Web-VRGame-Backend/config"
	"github.com/alongkornn/Web-VRGame-Backend/internal/admin/dto"
	auth_models "github.com/alongkornn/Web-VRGame-Backend/internal/auth/models"
	checkpoint_models "github.com/alongkornn/Web-VRGame-Backend/internal/checkpoint/models"
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
	hasAdmin := config.DB.Collection("User").
		Where("is_deleted", "==", false).
		Where("role", "==", auth_models.Admin).
		Where("id", "==", adminId).
		Limit(1)

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
	hasAdmin := config.DB.Collection("User").
		Where("is_deleted", "==", false).
		Where("role", "==", auth_models.Admin).
		Where("status", "==", auth_models.Approved).
		Where("id", "==", adminId).
		Limit(1)

	adminDoc, err := hasAdmin.Documents(ctx).Next()
	if err != nil {
		return nil, http.StatusNotFound, errors.New("admin not found")
	}

	var admin *auth_models.User
	if err := adminDoc.DataTo(&admin); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return admin, http.StatusOK, nil
}

// เพิ่มผู้ดูแลระบบ
func CreateAdmin(userId string, role auth_models.Role, ctx context.Context) (int, error) {
	hasUser := config.DB.Collection("User").
		Where("is_deleted", "==", false).
		Where("role", "==", auth_models.Player).
		Where("status", "==", auth_models.Approved).
		Where("id", "==", userId).
		Limit(1)

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
	hasAdmin := config.DB.Collection("User").
		Where("is_deleted", "==", false).
		Where("role", "==", auth_models.Admin).
		Where("status", "==", auth_models.Approved).
		Where("id", "==", adminId).
		Limit(1)

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
	hasAdmin := config.DB.Collection("User").
		Where("is_deleted", "==", false).
		Where("role", "==", auth_models.Admin).
		Where("id", "==", adminId).
		Limit(1)

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
	hasUser := config.DB.Collection("User").
		Where("is_deleted", "==", false).
		Where("role", "==", auth_models.Player).
		Where("status", "==", auth_models.Approved).
		Where("id", "==", userId).
		Limit(1)

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
				isStrength = append(isStrength, &checkpoint.Category)
			}
		}
	}

	return isStrength, http.StatusOK, nil
}

// แสดงจุดด้อยของผู้เล่น
func ShowScoreWiteWeaknesses(userId string, ctx context.Context) ([]*checkpoint_models.Category, int, error) {
	hasUser := config.DB.Collection("User").
		Where("is_deleted", "==", false).
		Where("role", "==", auth_models.Player).
		Where("status", "==", auth_models.Approved).
		Where("id", "==", userId).
		Limit(1)

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
				isWeaknesses = append(isWeaknesses, &checkpoint.Category)
			}
		}
	}

	return isWeaknesses, http.StatusOK, nil
}