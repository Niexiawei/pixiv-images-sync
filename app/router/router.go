package router

import (
	"github.com/gin-gonic/gin"
	"pixivImages/app/controller"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	oauth2 := r.Group("/oauth2")
	{
		oauth2Controller := &controller.OauthController{}
		oauth2.GET("/authorization/url", oauth2Controller.AuthorizationUrl)
		oauth2.GET("/receive", oauth2Controller.Receive)
	}

	r.GET("/favicon.ico", func(context *gin.Context) {
		context.String(200, "")
	})

	r.NoRoute(func(context *gin.Context) {
		context.String(404, "404 Not Found")
	})

	r.NoMethod(func(context *gin.Context) {
		context.String(404, "404 Not Found Method")
	})

	return r
}
