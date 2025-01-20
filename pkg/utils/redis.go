package utils

import "fmt"

func GetRedisKeys(userId string) (string, string) {
	userCacheKey := fmt.Sprintf("user:%s", userId)
	checkpointCacheKey := fmt.Sprintf("checkpoint:%s", userId)
	return userCacheKey, checkpointCacheKey
}
