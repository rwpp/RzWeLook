package ioc

import (
	wechat2 "github.com/rwpp/RzWeLook/oauth2/service/wechat"
	"github.com/rwpp/RzWeLook/pkg/logger"
	"github.com/spf13/viper"
)

func InitPrometheus(logv1 logger.LoggerV1) wechat2.Service {
	svc := InitService(logv1)
	type Config struct {
		NameSpace  string `yaml:"nameSpace"`
		Subsystem  string `yaml:"subsystem"`
		InstanceID string `yaml:"instanceId"`
		Name       string `yaml:"name"`
	}
	var cfg Config
	err := viper.UnmarshalKey("prometheus", &cfg)
	if err != nil {
		panic(err)
	}
	return wechat2.NewPrometheusDecorator(svc, cfg.NameSpace, cfg.Subsystem, cfg.InstanceID, cfg.Name)
}
