package services

import (
	"context"
	"errors"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
	auth_models "github.com/alongkornn/Web-VRGame-Backend/internal/auth/models"
	checkpoint_models "github.com/alongkornn/Web-VRGame-Backend/internal/checkpoint/models"
	score_models "github.com/alongkornn/Web-VRGame-Backend/internal/score/models"
	"github.com/alongkornn/Web-VRGame-Backend/pkg/utils"
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

// func GetAllScoreByCheckpointId(checkpointId string, ctx context.Context) ([]*score_models.Score, int, error) {
// 	iter := config.DB.Collection("User").
// 		Where("is_deleted", "==", false).
// 		Where("role", "==", auth_models.Player).
// 		Where("status", "==", "approved").
// 		Documents(ctx)

// 	defer iter.Stop()

// 	var users_score []*score_models.Score
// 	for {
// 		userDoc, err := iter.Next()
// 		if err == iterator.Done {
// 			if len(users_score) == 0 {
// 				return nil, http.StatusNotFound, errors.New("user not found")
// 			}
// 			break
// 		}
// 		if err != nil {
// 			return nil, http.StatusInternalServerError, err
// 		}

// 		var user auth_models.User
// 		err = userDoc.DataTo(&user)
// 		if err != nil {
// 			return nil, http.StatusInternalServerError, err
// 		}

// 		hasCheckpoint := utils.GetCheckpointByID(user.CurrentCheckpoint)

// 		checkpointDoc, err := hasCheckpoint.Documents(ctx).Next()
// 		if err != nil {
// 			return nil, http.StatusNotFound, errors.New("checkpoint not found")
// 		}

// 		var currentCheckpoint checkpoint_models.Checkpoints
// 		if err := checkpointDoc.DataTo(&currentCheckpoint); err != nil {
// 			return nil, http.StatusInternalServerError, err
// 		}

// 		if user.CurrentCheckpoint == checkpointId {
// 			score := score_models.Score{
// 				CheckpointName: currentCheckpoint.Name,
// 				Category:       currentCheckpoint.Category,
// 				Name:           user.FirstName,
// 			}
// 			users_score = append(users_score, &score)
// 		} else {
// 			return nil, http.StatusBadRequest, errors.New("checkpoin id not found")
// 		}

// 	}
// 	return users_score, http.StatusOK, nil
// }

func SetProjectileScore(userId string, score int, playTime string, ctx context.Context) (int, error) {
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

	// ดึงข้อมูล checkpoint ปัจจุบัน
	hasCheckpoint := utils.GetCheckpointByID(user.ProjectileCurrentCheckpoint)
	checkpointDoc, err := hasCheckpoint.Documents(ctx).Next()
	if err != nil {
		return http.StatusNotFound, errors.New("checkpoint not found")
	}

	var currentCheckpoint checkpoint_models.Checkpoints
	if err := checkpointDoc.DataTo(&currentCheckpoint); err != nil {
		return http.StatusInternalServerError, err
	}

	// ตรวจสอบว่าด่านปัจจุบันเคยอยู่ใน CompletedCheckpoints หรือไม่
	updated := false
	for i, completed := range user.CompletedCheckpoints {
		if completed.CheckpointID == user.ProjectileCurrentCheckpoint {
			if completed.Score < score {
				// อัปเดตคะแนนด่านที่เคยเล่นผ่าน
				user.CompletedCheckpoints[i].Score = score

				// คำนวณคะแนนรวมใหม่
				newScore := user.Score + (score - completed.Score)

				// อัปเดตข้อมูล Firestore
				_, err = userDoc.Ref.Update(ctx, []firestore.Update{
					{
						Path:  "completed_checkpoint",
						Value: user.CompletedCheckpoints, // อัปเดตทั้ง array
					},
					{
						Path:  "score",
						Value: newScore,
					},
					{
						Path:  "updated_at",
						Value: time.Now(),
					},
				})
				if err != nil {
					return http.StatusInternalServerError, err
				}
			}
			updated = true
			break
		}
	}

	// ถ้า checkpoint นี้ยังไม่เคยอยู่ใน CompletedCheckpoints ให้เพิ่มเข้าไป
	if !updated && score >= currentCheckpoint.PassScore {
		completedCheckpoint := checkpoint_models.CompleteCheckpoint{
			CheckpointID: currentCheckpoint.ID,
			Name:         currentCheckpoint.Name,
			Category:     currentCheckpoint.Category,
			Score:        score,
			Time:         playTime,
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
	}

	return http.StatusOK, nil
}

func SetMomentumScore(userId string, score int, playTime string, ctx context.Context) (int, error) {
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

	// ดึงข้อมูล checkpoint ปัจจุบัน
	hasCheckpoint := utils.GetCheckpointByID(user.MomentumCurrentCheckpoint)
	checkpointDoc, err := hasCheckpoint.Documents(ctx).Next()
	if err != nil {
		return http.StatusNotFound, errors.New("checkpoint not found")
	}

	var currentCheckpoint checkpoint_models.Checkpoints
	if err := checkpointDoc.DataTo(&currentCheckpoint); err != nil {
		return http.StatusInternalServerError, err
	}

	// ตรวจสอบว่าด่านปัจจุบันเคยอยู่ใน CompletedCheckpoints หรือไม่
	updated := false
	for i, completed := range user.CompletedCheckpoints {
		if completed.CheckpointID == user.MomentumCurrentCheckpoint {
			if completed.Score < score {
				// อัปเดตคะแนนด่านที่เคยเล่นผ่าน
				user.CompletedCheckpoints[i].Score = score

				// คำนวณคะแนนรวมใหม่
				newScore := user.Score + (score - completed.Score)

				// อัปเดตข้อมูล Firestore
				_, err = userDoc.Ref.Update(ctx, []firestore.Update{
					{
						Path:  "completed_checkpoint",
						Value: user.CompletedCheckpoints,
					},
					{
						Path:  "score",
						Value: newScore,
					},
					{
						Path:  "updated_at",
						Value: time.Now(),
					},
				})
				if err != nil {
					return http.StatusInternalServerError, err
				}
			}
			updated = true
			break
		}
	}

	// ถ้า checkpoint นี้ยังไม่เคยอยู่ใน CompletedCheckpoints ให้เพิ่มเข้าไป
	if !updated && score >= currentCheckpoint.PassScore {
		completedCheckpoint := checkpoint_models.CompleteCheckpoint{
			CheckpointID: currentCheckpoint.ID,
			Name:         currentCheckpoint.Name,
			Category:     currentCheckpoint.Category,
			Score:        score,
			Time:         playTime,
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
	}

	return http.StatusOK, nil
}
func SetForceScore(userId string, score int, playTime string, ctx context.Context) (int, error) {
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

	// ดึงข้อมูล checkpoint ปัจจุบัน
	hasCheckpoint := utils.GetCheckpointByID(user.ForceCurrentCheckpoint)
	checkpointDoc, err := hasCheckpoint.Documents(ctx).Next()
	if err != nil {
		return http.StatusNotFound, errors.New("checkpoint not found")
	}

	var currentCheckpoint checkpoint_models.Checkpoints
	if err := checkpointDoc.DataTo(&currentCheckpoint); err != nil {
		return http.StatusInternalServerError, err
	}

	// ตรวจสอบว่าด่านปัจจุบันเคยอยู่ใน CompletedCheckpoints หรือไม่
	updated := false
	for i, completed := range user.CompletedCheckpoints {
		if completed.CheckpointID == user.ForceCurrentCheckpoint {
			if completed.Score < score {
				// อัปเดตคะแนนด่านที่เคยเล่นผ่าน
				user.CompletedCheckpoints[i].Score = score

				// คำนวณคะแนนรวมใหม่
				newScore := user.Score + (score - completed.Score)

				// อัปเดตข้อมูล Firestore
				_, err = userDoc.Ref.Update(ctx, []firestore.Update{
					{
						Path:  "completed_checkpoint",
						Value: user.CompletedCheckpoints,
					},
					{
						Path:  "score",
						Value: newScore,
					},
					{
						Path:  "updated_at",
						Value: time.Now(),
					},
				})
				if err != nil {
					return http.StatusInternalServerError, err
				}
			}
			updated = true
			break
		}
	}

	// ถ้า checkpoint นี้ยังไม่เคยอยู่ใน CompletedCheckpoints ให้เพิ่มเข้าไป
	if !updated && score >= currentCheckpoint.PassScore {
		completedCheckpoint := checkpoint_models.CompleteCheckpoint{
			CheckpointID: currentCheckpoint.ID,
			Name:         currentCheckpoint.Name,
			Category:     currentCheckpoint.Category,
			Score:        score,
			Time:         playTime,
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
	}

	return http.StatusOK, nil
}
