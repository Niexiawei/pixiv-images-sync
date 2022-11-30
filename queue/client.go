package queue

import (
	"fmt"
	"github.com/hibiken/asynq"
	"pixivImages/config"
)

var client *asynq.Client

var inspector *asynq.Inspector

func InitClient() {
	redisConfig := config.Get().Redis
	redisOpt := asynq.RedisClientOpt{
		Addr:     fmt.Sprintf("%s:%d", redisConfig.Host, redisConfig.Port),
		Password: redisConfig.Password,
		DB:       redisConfig.Db,
	}
	client = asynq.NewClient(redisOpt)
	inspector = asynq.NewInspector(redisOpt)
}

func GetClient() *asynq.Client {
	return client
}

func GetInspector() *asynq.Inspector {
	return inspector
}
