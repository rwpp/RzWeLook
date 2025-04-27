package tencent

import (
	"context"
	"fmt"
	"github.com/ecodeclub/ekit"
	"github.com/ecodeclub/ekit/slice"
	"github.com/rwpp/RzWeLook/pkg/ratelimit"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
	"go.uber.org/zap"
)

type Service struct {
	client   *sms.Client
	appId    *string
	signName *string
	limiter  ratelimit.Limiter
}

func NewService(c *sms.Client, appId string,
	signName string, limiter ratelimit.Limiter) *Service {
	return &Service{
		client:   c,
		appId:    ekit.ToPtr[string](appId),
		signName: ekit.ToPtr[string](signName),
		limiter:  limiter,
	}
}

func (s *Service) Send(ctx context.Context, biz string,
	args []string, numbers ...string) error {
	req := sms.NewSendSmsRequest()
	req.PhoneNumberSet = toStringPtrSlice(numbers)
	req.SmsSdkAppId = s.appId
	// ctx 继续往下传
	req.SetContext(ctx)
	req.TemplateParamSet = toStringPtrSlice(args)
	req.TemplateId = ekit.ToPtr[string](biz)
	req.SignName = s.signName
	resp, err := s.client.SendSms(req)
	zap.L().Debug("调用腾讯短信服务",
		zap.Any("req", req),
		zap.Any("resp", resp))
	if err != nil {
		return err
	}
	for _, status := range resp.Response.SendStatusSet {
		if status.Code == nil || *(status.Code) != "Ok" {
			return fmt.Errorf("发送失败，code: %s, 原因：%s",
				*status.Code, *status.Message)
		}
	}
	return nil
}

func toStringPtrSlice(src []string) []*string {
	return slice.Map[string, *string](src, func(idx int, src string) *string {
		return &src
	})
}
