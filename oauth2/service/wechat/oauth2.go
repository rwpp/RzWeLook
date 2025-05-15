package wechat

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/rwpp/RzWeLook/oauth2/domain/wechat"
	"github.com/rwpp/RzWeLook/pkg/logger"
	"net/http"
	"net/url"
)

const authURLPattern = "https://open.weixin.qq.com/connect/qrconnect?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_login&state=%s#wechat_redire"

var redirectURL = url.PathEscape("https://meoying.com/oauth2/wechat/callback")

type service struct {
	appId     string
	appSecret string
	client    *http.Client
	logger    logger.LoggerV1
}

type Result struct {
	ErrCode int64  `json:"errcode"`
	ErrMsg  string `json:"errMsg"`

	Scope string `json:"scope"`

	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`

	OpenId  string `json:"openid"`
	UnionId string `json:"unionid"`
}

func NewService(appId, appSecret string,
	logger logger.LoggerV1) Service {
	return &service{
		appId:     appId,
		appSecret: appSecret,
		client:    http.DefaultClient,
		logger:    logger,
	}
}

func (s *service) VerifyCode(ctx context.Context, code string) (wechat.WechatInfo, error) {
	const baseURL = "https://api.weixin.qq.com/sns/oauth2/access_token"
	// 这是另外一种写法
	queryParams := url.Values{}
	queryParams.Set("appid", s.appId)
	queryParams.Set("secret", s.appSecret)
	queryParams.Set("code", code)
	queryParams.Set("grant_type", "authorization_code")
	accessTokenURL := baseURL + "?" + queryParams.Encode()
	req, err := http.NewRequest("GET", accessTokenURL, nil)
	if err != nil {
		return wechat.WechatInfo{}, err
	}
	req = req.WithContext(ctx)
	resp, err := s.client.Do(req)
	if err != nil {
		return wechat.WechatInfo{}, err
	}
	defer resp.Body.Close()
	var res Result

	err = json.NewDecoder(resp.Body).Decode(&res)
	if err != nil {
		return wechat.WechatInfo{}, err
	}
	if res.ErrCode != 0 {
		return wechat.WechatInfo{}, errors.New("换取 access_token 失败")
	}
	return wechat.WechatInfo{
		OpenId:  res.OpenId,
		UnionId: res.UnionId,
	}, nil
}

func (s *service) AuthURL(ctx context.Context, state string) (string, error) {
	return fmt.Sprintf(authURLPattern, s.appId, redirectURL, state), nil
}
