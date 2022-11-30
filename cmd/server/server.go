package server

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/urfave/cli/v2"
	"net/http"
	"os"
	"os/signal"
	"pixivImages/app/router"
	"pixivImages/config"
	"pixivImages/crontab"
	"pixivImages/database"
	"pixivImages/logger"
	"pixivImages/queue"
	"pixivImages/utils"
	"pixivImages/utils/validator"
	"syscall"
	"time"
)

type Server struct {
}

var RunCmd = &cli.Command{
	Usage: "启动Api服务",
	Name:  "server",
	Action: func(context *cli.Context) error {
		s := &Server{}
		s.Bootstrap()
		s.Run()
		return nil
	},
}

func (s *Server) Bootstrap() {
	config.LoadConfig()
	logger.InitLogger()
	database.InitRedis()
	crontab.Run()
	queue.InitServer()
	database.InitMysql()
	validator.InitValidatorTrans()
	utils.InitSnowflake()
}

func (s *Server) Run() {
	ginEngine := router.InitRouter()
	stopHttpServer := startHttpServer(ginEngine)
	logger.Logger.Info("server start successful ...")
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	stopHttpServer()
}

func startHttpServer(r *gin.Engine) func() {
	addr := fmt.Sprintf(":%d", config.Get().HttpServer.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	logger.Logger.Info("http server address " + addr)

	go func() {
		// 服务连接
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Logger.Error(fmt.Sprintf("listen: %s\n", err))
		}
	}()

	shutdown := func() {
		logger.Logger.Info("Shutdown Server ...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			logger.Logger.Error("Server Shutdown:", err)
		}
		logger.Logger.Info("Server exiting ...")
	}
	return shutdown
}
