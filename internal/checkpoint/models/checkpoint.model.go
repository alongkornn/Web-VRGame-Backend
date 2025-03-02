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
	ID             string    `json:"id" firestore:"id"`
	Name           string    `json:"name" firestore:"name"`
	Category       Category  `json:"category" firestore:"category"`
	MaxScore       int       `json:"max_score" firestore:"max_score"`
	PassScore      int       `json:"pass_score" firestore:"pass_score"`
	TimeLimit      string    `json:"time_limit" firestore:"time_limit"`
	NextCheckpoint string    `json:"next_checkpoint" firestore:"next_checkpoint"`
	CreatedAt      time.Time `json:"created_at" firestore:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" firestore:"updated_at"`
	Is_Deleted     bool      `json:"is_deleted" firestore:"is_deleted"`
}

type CompleteCheckpoint struct {
	CheckpointID string   `json:"checkpoint_id" firestore:"checkpoint_id"`
	Name         string   `json:"name" firestore:"name"`
	Category     Category `json:"category" firestore:"category"`
	Score        int      `json:"score" firestore:"score"`
	Time         string   `json:"time" firestore:"time"`
}

type CheckpointDetail struct {
	CheckpointID string `json:"checkpoint_id" firestore:"checkpoint_id"`
	Name         string `json:"name" firestore:"name"`
	Category     string `json:"category" firestore:"category"`
	Score        int    `json:"score" firestore:"score"`
}
