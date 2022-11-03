package cache

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"spider_douban/config"
)

func GetRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.RedisSetting.Host, config.RedisSetting.Port),
		Password: config.RedisSetting.Auth,
		DB:       0,
	})
}
