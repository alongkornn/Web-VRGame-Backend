package dto

type CheckpointDTO struct {
	CheckpointID     string `json:"checkpointID" validate:"required"`
	UserID string `json:"userID" validate:"required"`
}
