package ioc

import (
	"github.com/rwpp/RzWeLook/internal/service/oauth/wechat"
	"github.com/rwpp/RzWeLook/internal/web"
	"os"
)

func InitOAuthWechatService() wechat.Service {
	appId, ok := os.LookupEnv("WECHAT_APP_ID")
	if !ok {
		panic("找不到环境变量WECHAT_APP_ID")
	}
	appKey, ok := os.LookupEnv("WECHAT_APP_SECRET")
	if !ok {
		panic("找不到环境变量WECHAT_APP_KEY")
	}
	return wechat.NewService(appId, appKey)
}

func NewWechatHandler() web.WechatHandlerConfig {
	return web.WechatHandlerConfig{
		Secure: false,
	}
}
