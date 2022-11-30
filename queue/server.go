package queue

import (
	"fmt"
	"github.com/hibiken/asynq"
	"pixivImages/config"
	"pixivImages/queue/tasks/pixiv_images_save"
)

var server *asynq.Server

func InitServer() {
	redisConfig := config.Get().Redis
	server = asynq.NewServer(
		asynq.RedisClientOpt{
			Addr:     fmt.Sprintf("%s:%d", redisConfig.Host, redisConfig.Port),
			Password: redisConfig.Password,
			DB:       redisConfig.Db,
		},
		asynq.Config{
			Concurrency:  10,
			ErrorHandler: &TaskErrorHandler{},
		},
	)
	mux := asynq.NewServeMux()
	mux.HandleFunc(pixiv_images_save.PixivImagesSave, pixiv_images_save.ProcessTask)
	go func() {
		if err := server.Run(mux); err != nil {
			fmt.Println(err)
		}
	}()
}
