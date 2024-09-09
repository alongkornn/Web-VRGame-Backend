package models

import "time"

type Role string
type Status string

const (
	Player Role = "player"
	Admin  Role = "admin"
)

const (
	Pending Status = "pending"
	Done    Status = "done"
)

type User struct {
	ID         string    `json:"id" firestore:"id"`
	FirstName  string    `json:"firstname" firestore:"firstname" validate:"required,min=2,max=50"`
	LastName   string    `json:"lastname" firestore:"lastname" validate:"require, min=2, max=50"`
	Email      string    `json:"email" firestore:"email" validate:"required,email"`
	Password   string    `json:"password" firestore:"password" validate:"required"`
	Class      string    `json:"class" firestore:"class" validate:"required"`
	Number     string    `json:"number" firestore:"number" validate:"required"`
	Role       Role      `json:"role" firestore:"role"`
	Score      int       `json:"score" firestore:"score"`
	Status     Status    `json:"status" firestore:"status"`
	CreatedAt  time.Time `json:"created_at" firestore:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" firestore:"updated_at"`
	Is_Deleted bool      `json:"is_deleted" firestore:"is_deleted"`
}
