package wechat

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rwpp/RzWeLook/internal/domain"
	"net/http"
	"net/url"
)

var redirectURI = url.PathEscape("https://www.example.com/oauth/wechat/callback")

type Service interface {
	AuthURL(ctx context.Context, state string) (string, error)
	VerifyCode(ctx context.Context, code string) (domain.WechatInfo, error)
}

func NewService(appid string, appSecret string) Service {
	return &service{
		appId:     appid,
		appSecret: appSecret,
		client:    http.DefaultClient,
	}
}

type service struct {
	appId     string
	appSecret string
	client    *http.Client
}

func (s *service) VerifyCode(ctx context.Context, code string) (domain.WechatInfo, error) {
	const targetPattern = "https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code"
	target := fmt.Sprintf(targetPattern, s.appId, s.appSecret, code)
	req, err := http.NewRequest(http.MethodGet, target, nil)
	//resp, err := http.Get(target)
	if err != nil {
		return domain.WechatInfo{}, err
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return domain.WechatInfo{}, err
	}

	decoder := json.NewDecoder(resp.Body)
	var res Result
	err = decoder.Decode(&res)
	//body, err := io.ReadAll(resp.Body)
	if err != nil {
		return domain.WechatInfo{}, err
	}
	if res.ErrCode != 0 {
		return domain.WechatInfo{}, errors.New("微信返回错误信息")
	}
	return domain.WechatInfo{
		OpenID:  res.Openid,
		UnionID: res.Unionid,
	}, nil
	//req := http.NewRequest()
	//http.Client.Do()
}

func (s *service) AuthURL(ctx context.Context, state string) (string, error) {
	const urlPattern = "https://open.weixin.qq.com/connect/oauth2/authorize?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_userinfo&state=%s#wechat_redirect"
	return fmt.Sprintf(urlPattern, s.appId, redirectURI, state), nil
}

type Result struct {
	ErrCode      int64  `json:"errcode"`
	ErrMsg       string `json:"errmsg"`
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Unionid      string `json:"unionid"`
	Openid       string `json:"openid"`
	Scope        string `json:"scope"`
}
