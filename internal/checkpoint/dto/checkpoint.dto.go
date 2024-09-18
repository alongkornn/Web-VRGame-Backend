package dto

import (
	"github.com/alongkornn/Web-VRGame-Backend/internal/checkpoint/models"
)

type CheckpointDTO struct {
	CheckpointID string `json:"checkpointID"`
	UserID       string `json:"userID"`
}

type CreateCheckpointsDTO struct {
	Name      string          `json:"name"`
	MaxScore  int             `json:"max_score"`
	PassScore int             `json:"pass_score"`
	Category  models.Category `json:"category"`
}

type GetCheckpointWithCategoryDTO struct {
	Category string `json:"category"`
}
