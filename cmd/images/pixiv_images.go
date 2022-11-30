package images

import (
	"github.com/urfave/cli/v2"
	service_images "pixivImages/app/service/images"
	"pixivImages/config"
	"pixivImages/database"
)

type PixivImages struct {
}

var ImagesCmd = &cli.Command{
	Usage: "images 操作",
	Name:  "images",
	Subcommands: []*cli.Command{
		RefreshImagesDownloadUrlCmd,
	},
	Action: func(context *cli.Context) error {
		return cli.ShowSubcommandHelp(context)
	},
}

func bootstrap() {
	config.LoadConfig()
	database.InitRedis()
	database.InitMysql()
}

var RefreshImagesDownloadUrlCmd = &cli.Command{
	Usage: "刷新images下载url",
	Name:  "refreshDownloadUrl",
	Action: func(context *cli.Context) error {
		bootstrap()
		service_images.RefreshPixivDownloadUrl()
		return nil
	},
}
