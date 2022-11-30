package ms_graph

import (
	"github.com/imroc/req/v3"
	"github.com/pkg/errors"
	"pixivImages/config"
	"pixivImages/database"
	"pixivImages/utils"
	"strings"
	"time"
)

var (
	AuthorizationError = errors.New("授权失败")
)

const (
	TokenKey        = "graphToken"
	RefreshTokenKey = "graphRefreshToken"
)

type Token struct {
	AccessToken  string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

type AuthorizationErrorResp struct {
	Error     string `json:"error"`
	ErrorDesc string `json:"error_description"`
	ErrorUri  string `json:"error_uri"`
}

type Authorization struct {
	HttpClient *req.Client
	Config     config.MsGraph
}

type GetTokenResponse struct {
	ExpiresIn    int    `json:"expires_in"`
	Token        string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Error        string `json:"error"`
	ErrorDesc    string `json:"error_description"`
	ErrorUrl     string `json:"error_uri"`
}

func NewAuthorization() *Authorization {
	return &Authorization{
		HttpClient: req.C().SetTimeout(60 * time.Second).SetBaseURL("https://login.microsoftonline.com"),
		Config:     config.Get().MsGraph,
	}
}

func GetGraphToken() string {
	token, err := database.GetRedisConn().Get(database.RedisBaseCtx, TokenKey).Result()
	if err != nil {
		return ""
	}
	return token
}

func GetGraphRefreshToken() string {
	token, err := database.GetRedisConn().Get(database.RedisBaseCtx, RefreshTokenKey).Result()
	if err != nil {
		return ""
	}
	return token
}

func (auth *Authorization) GetToken(data string, refresh bool) (*GetTokenResponse, error) {
	errResp := AuthorizationErrorResp{}
	result := GetTokenResponse{}
	params := map[string]string{
		"client_id":     auth.Config.ClientId,
		"scope":         strings.Join(auth.Config.Scopes, " "),
		"redirect_uri":  auth.Config.ReceiveUrl,
		"grant_type":    utils.If[string](refresh, "refresh_token", "authorization_code"),
		"client_secret": auth.Config.SecretId,
	}
	if !refresh {
		params["code"] = data
	} else {
		params["refresh_token"] = data
	}

	resp, err := auth.HttpClient.R().SetFormData(params).SetError(&errResp).
		SetResult(&result).
		Post("/organizations/oauth2/v2.0/token")
	if err != nil {
		return nil, errors.WithMessage(AuthorizationError, err.Error())
	}

	if resp.StatusCode != 200 && errResp.Error != "" {
		return nil, errors.WithMessage(AuthorizationError, errResp.Error+":"+errResp.ErrorDesc)
	}

	if err := database.GetRedisConn().Set(database.RedisBaseCtx, TokenKey, result.Token, 2600*time.Second).Err(); err != nil {
		return nil, err
	}

	if err := database.GetRedisConn().Set(database.RedisBaseCtx, RefreshTokenKey, result.RefreshToken, 45*time.Hour).Err(); err != nil {
		return nil, err
	}

	return &result, nil
}
