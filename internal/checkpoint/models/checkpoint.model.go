package models

import "time"

type Category string

const (
	Projectile     Category = "โพรเจกไทล์"
	Momentum       Category = "โมเมนตัมและการชน"
	ForceAndMotion Category = "แรงและกฎการเคลื่อนที่"
)

type Checkpoints struct {
	ID         string    `json:"id" firestore:"id"`
	Name       string    `json:"name" firestore:"name"`
	Category   Category  `json:"category" firestore:"category"`
	CreatedAt  time.Time `json:"created_at" firestore:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" firestore:"updated_at"`
	Is_Deleted bool      `json:"is_deleted" firestore:"is_deleted"`
}
