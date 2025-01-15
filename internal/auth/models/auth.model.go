package models

import (
	"time"

	"github.com/alongkornn/Web-VRGame-Backend/internal/checkpoint/models"
)

type Role string
type Status string

const (
	Player Role = "player"
	Admin  Role = "admin"
)

const (
	Pending  Status = "pending"
	Approved Status = "approved"
	Deleted  Status = "deleted"
)

type User struct {
	ID                   string                `json:"id" firestore:"id"`
	FirstName            string                `json:"firstname" firestore:"firstname"`
	LastName             string                `json:"lastname" firestore:"lastname"`
	Email                string                `json:"email" firestore:"email"`
	Password             string                `json:"password" firestore:"password"`
	Score                int                   `json:"score,omitempty" firestore:"score"`
	Role                 Role                  `json:"role" firestore:"role"`
	Status               Status                `json:"status" firestore:"status"`
	CurrentCheckpoint    *models.Checkpoints   `json:"current_checkpoint,omitempty" firestore:"current_checkpoint,omitempty"`
	CompletedCheckpoints []*models.Checkpoints `json:"completed_checkpoints,omitempty" firestore:"completed_checkpoints,omitempty"`
	Time                 string                `json:"time,omitempty" firestore:"time,omitempty"`
	CreatedAt            time.Time             `json:"created_at" firestore:"created_at"`
	UpdatedAt            time.Time             `json:"updated_at" firestore:"updated_at"`
	VerifyEmail          bool                  `json:"verify_email" firestore:"verify_email"`
	Is_Deleted           bool                  `json:"is_deleted" firestore:"is_deleted"`
}
