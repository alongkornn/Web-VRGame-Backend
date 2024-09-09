package models

import (
	"time"
)

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
	FirstName  string    `json:"firstname" firestore:"firstname"`
	LastName   string    `json:"lastname" firestore:"lastname"`
	Email      string    `json:"email" firestore:"email"`
	Password   string    `json:"password" firestore:"password"`
	Class      string    `json:"class" firestore:"class"`
	Number     string    `json:"number" firestore:"number"`
	Score      int       `json:"score" firestore:"score"`
	Level      int    `json:"level" firestore:"level"`
	Role       Role      `json:"role" firestore:"role"`
	Status     Status    `json:"status" firestore:"status"`
	CreatedAt  time.Time `json:"created_at" firestore:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" firestore:"updated_at"`
	Is_Deleted bool      `json:"is_deleted" firestore:"is_deleted"`
}

