package ioc

import (
	wechat2 "github.com/rwpp/RzWeLook/oauth2/service/wechat"
	"github.com/rwpp/RzWeLook/pkg/logger"
	"github.com/spf13/viper"
)

func InitService(logv1 logger.LoggerV1) wechat2.Service {
	type Config struct {
		AppID     string `yaml:"appId"`
		AppSecret string `yaml:"appSecret"`
	}
	var cfg Config
	err := viper.UnmarshalKey("weChatConf", &cfg)
	if err != nil {
		panic(err)
	}
	return wechat2.NewService(cfg.AppID, cfg.AppSecret, logv1)
}
