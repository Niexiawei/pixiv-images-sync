package crontab

import (
	"github.com/go-co-op/gocron"
	"pixivImages/app/service/ms_graph"
	"pixivImages/logger"
	"time"
)

var Cron = gocron.NewScheduler(time.Local)

func msGraphTokenRefresh() {
	if token := ms_graph.GetGraphToken(); token == "" {
		refresh := ms_graph.GetGraphRefreshToken()
		_, err := ms_graph.NewAuthorization().GetToken(refresh, true)
		if err != nil {
			logger.Logger.Error(err.Error())
		}
	}

	_, err := Cron.Every(15).Minutes().Do(func() {
		refresh := ms_graph.GetGraphRefreshToken()
		_, err := ms_graph.NewAuthorization().GetToken(refresh, true)
		if err != nil {
			logger.Logger.Error(err.Error())
		}
	})

	if err != nil {
		panic(err)
	}
}

func Run() {
	msGraphTokenRefresh()
	Cron.StartAsync()
}
