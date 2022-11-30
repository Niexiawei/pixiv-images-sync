package service_images

import (
	"pixivImages/config"
	"pixivImages/database"
	"testing"
)

func Test_RefreshPixivDownloadUrl(t *testing.T) {
	config.LoadConfig()
	database.InitRedis()
	database.InitMysql()
	RefreshPixivDownloadUrl()
}
