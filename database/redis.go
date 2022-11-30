package database

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"pixivImages/config"
)

var redisClient *redis.Client

var RedisBaseCtx = context.Background()

func InitRedis() {
	redisConfig := config.Get().Redis
	redisClient = redis.NewClient(&redis.Options{
		Password: redisConfig.Password,
		Addr:     fmt.Sprintf("%s:%d", redisConfig.Host, redisConfig.Port),
		DB:       redisConfig.Db,
		PoolSize: redisConfig.Pool,
	})
}

func GetRedisConn() *redis.Client {
	return redisClient
}
