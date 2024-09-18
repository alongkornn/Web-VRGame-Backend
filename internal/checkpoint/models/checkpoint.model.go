package models

import (
	"time"
)

type Category string

const (
	Projectile     Category = "โพรเจกไทล์"
	Momentum       Category = "โมเมนตัมและการชน"
	ForceAndMotion Category = "แรงและกฎการเคลื่อนที่"
)

type Checkpoints struct {
	ID          string    `json:"id" firestore:"id"`
	Name        string    `json:"name" firestore:"name"`
	MaxScore    int       `json:"max_score" firestore:"max_score"`
	PassScore   int       `json:"pass_score" firestore:"pass_score"`
	PlayerScore []Player  `json:"player_score" firestore:"player_score"`
	Category    Category  `json:"category" firestore:"category"`
	CreatedAt   time.Time `json:"created_at" firestore:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" firestore:"updated_at"`
	Is_Deleted  bool      `json:"is_deleted" firestore:"is_deleted"`
}

type Player struct {
	ID    string `json:"id" firestore:"id"`
	Name  string `json:"name" firestore:"name"`
	Score int    `json:"score" firestore:"score"`
}
