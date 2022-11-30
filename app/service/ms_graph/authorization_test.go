package ms_graph

import (
	"fmt"
	"pixivImages/config"
	"pixivImages/database"
	"pixivImages/logger"
	"testing"
)

var auth *Authorization

func init() {
	config.LoadConfig()
	database.InitRedis()
	logger.InitLogger()
	auth = NewAuthorization()
}

func TestAuthorization_GetToken(t *testing.T) {
	token, err := auth.GetToken("123456", false)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("%+v", token)
}

func TestGetGraphRefreshToken(t *testing.T) {
	_, err := auth.GetToken(GetGraphRefreshToken(), true)
	if err != nil {
		logger.Logger.Error(err.Error())
	}
}
