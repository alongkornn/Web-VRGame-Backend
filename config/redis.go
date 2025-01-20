package config

import (
	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client

func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // ชี้ไปยัง Redis ที่รันผ่าน Docker
		Password: "",               // ใช้ password ถ้าคุณตั้งไว้
		DB:       0,                // ใช้ default database
	})
}
