package cache

import "github.com/go-redis/redis/v8"

func GetRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "120.78.67.238:6379",
		Password: "woaini1!",
		DB:       0,
	})
}
