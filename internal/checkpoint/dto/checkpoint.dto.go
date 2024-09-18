package dto

import (
	"github.com/alongkornn/Web-VRGame-Backend/internal/checkpoint/models"
)
type CheckpointDTO struct {
	CheckpointID     string `json:"checkpointID" validate:"required"`
	UserID string `json:"userID" validate:"required"`
}


type CreateCheckpointsDTO struct {
	Name        string    `json:"name" firestore:"name"`
	MaxScore    int       `json:"max_score" firestore:"max_score"`
	PassScore   int       `json:"pass_score" firestore:"pass_score"`
	Category    models.Category  `json:"category" firestore:"category"`
}