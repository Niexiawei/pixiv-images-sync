package cmd

import (
	"github.com/urfave/cli/v2"
	"os"
	"pixivImages/cmd/images"
	"pixivImages/cmd/server"
)

var App = &cli.App{
	Name:  "pixivImages",
	Usage: "pixiv图片管理系统",
	Commands: []*cli.Command{
		server.RunCmd,
		images.ImagesCmd,
	},
	Action: func(context *cli.Context) error {
		return cli.ShowAppHelp(context)
	},
}

func RunApp() {
	if err := App.Run(os.Args); err != nil {
		panic(err)
	}
}
