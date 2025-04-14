package ioc

import (
	"github.com/rwpp/RzWeLook/internal/service/sms"
	"github.com/rwpp/RzWeLook/internal/service/sms/memory"
)

func InitSMSService() sms.Service {
	// 这里可以初始化短信服务
	// 比如使用阿里云短信服务
	// 实际上什么都不做
	//return metrics.NewPrometheusDecorator(memory.NewService())
	return memory.NewService()
}
