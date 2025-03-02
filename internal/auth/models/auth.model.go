package models

import (
	"time"

	checkpoint_models "github.com/alongkornn/Web-VRGame-Backend/internal/checkpoint/models"
)

type Role string
type Status string

const (
	Player Role = "player"
	Admin  Role = "admin"
)

type User struct {
	ID                          string                                  `json:"id" firestore:"id"`
	FirstName                   string                                  `json:"firstname" firestore:"firstname"`
	LastName                    string                                  `json:"lastname" firestore:"lastname"`
	Email                       string                                  `json:"email" firestore:"email"`
	Password                    string                                  `json:"password" firestore:"password"`
	Score                       int                                     `json:"score,omitempty" firestore:"score"`
	Role                        Role                                    `json:"role" firestore:"role"`
	Status                      string                                  `json:"status" firestore:"status"`
	ProjectileCurrentCheckpoint string                                  `json:"projectilecurrent_checkpoint" firestore:"projectilecurrent_checkpoint"`
	MomentumCurrentCheckpoint   string                                  `json:"momentumcurrent_checkpoint" firestore:"momentumcurrent_checkpoint" `
	ForceCurrentCheckpoint      string                                  `json:"forcecurrent_checkpoint" firestore:"forcecurrent_checkpoint" `
	CompletedCheckpoints        []*checkpoint_models.CompleteCheckpoint `json:"completed_checkpoint,omitempty" firestore:"completed_checkpoint,omitempty"`
	Time                        string                                  `json:"time,omitempty" firestore:"time,omitempty"`
	CreatedAt                   time.Time                               `json:"created_at" firestore:"created_at"`
	UpdatedAt                   time.Time                               `json:"updated_at" firestore:"updated_at"`
	VerifyEmail                 bool                                    `json:"verify_email" firestore:"verify_email"`
	Is_Deleted                  bool                                    `json:"is_deleted" firestore:"is_deleted"`
}
