package controller

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/url"
	"pixivImages/app/request/authorization_controller_requests"
	"pixivImages/app/service/ms_graph"
	"pixivImages/config"
	"pixivImages/utils"
	"pixivImages/utils/validator"
	"strings"
)

type OauthController struct {
}

func (o *OauthController) Receive(ctx *gin.Context) {
	requestData := authorization_controller_requests.Receive{}
	if err := ctx.ShouldBindQuery(&requestData); err != nil {
		utils.NewResponse(500, "").WithError(validator.FormatValidatorErrors(err).First()).
			ReturnJsonResponse(ctx, 200)
		return
	}
	result, err := ms_graph.NewAuthorization().GetToken(requestData.Code, false)
	if err != nil {
		if errors.Is(err, ms_graph.AuthorizationError) {
			utils.NewResponse(500, "").WithError("授权失败").
				WithMessage(gin.H{
					"error": result.Error,
					"desc":  result.ErrorDesc,
					"desc2": result.ErrorUrl,
				}).
				ReturnJsonResponse(ctx, 200)
			return
		}
		utils.NewResponse(500, "").WithError("授权失败").
			WithMessage(gin.H{
				"error": err.Error(),
			}).
			ReturnJsonResponse(ctx, 200)
		return
	}
	//service.SetToken(result.Token, result.RefreshToken)
	utils.NewResponse(20, "ok").
		ReturnJsonResponse(ctx, 200)
}

func (o *OauthController) AuthorizationUrl(ctx *gin.Context) {
	graphConfig := config.Get().MsGraph
	authUrl := fmt.Sprintf("https://login.microsoftonline.com/organizations/oauth2/v2.0/authorize?client_id=%s&response_type=code&redirect_uri=%s&response_mode=query&scope=%s",
		graphConfig.ClientId,
		url.PathEscape(graphConfig.ReceiveUrl),
		url.PathEscape(strings.Join(graphConfig.Scopes, " ")),
	)
	utils.NewResponse(200, "ok").WithData(authUrl).ReturnJsonResponse(ctx, 200)
}
