package models

import "time"

type Role string

const (
	Player Role = "player"
	Admin  Role = "admin"
)

type User struct {
	ID         string    `json:"id" firestore:"id"`
	FirstName  string    `json:"firstname" firestore:"firstname" validate:"required,min=2,max=50"`
	LastName   string    `json:"lastname" firestore:"lastname" validate:"require, min=2, max=50"`
	Email      string    `json:"email" firestore:"email" validate:"required,email"`
	Class      string    `json:"class" firstore:"class" validate:"required"`
	Number     string    `json:"number" firstore:"number" validate:"required"`
	Role       Role      `json:"role" firestore:"role" validate:"required"`
	Score      int       `json:"score" firestore:"score"`
	Status     string    `json:"status" firestore:"status"`
	CreatedAt  time.Time `json:"created_at" firestore:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" firestore:"updated_at"`
	Is_Deleted time.Time `json:"is_deleted" firestore:"is_deleted"`
}
