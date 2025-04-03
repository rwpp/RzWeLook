package ratelimit

import (
	"context"
	"fmt"
	"github.com/rwpp/RzWeLook/internal/service/sms"
	"github.com/rwpp/RzWeLook/pkg/ratelimit"
)

var errLimitExceed = fmt.Errorf("限流触发")

type Service struct {
	svc     sms.Service
	limiter ratelimit.Limiter
}

func NewService(svc sms.Service, limiter ratelimit.Limiter) sms.Service {
	return &Service{
		limiter: limiter,
		svc:     svc,
	}
}
func (s *Service) Send(ctx context.Context, tplId string,
	args []string, numbers ...string) error {
	limited, err := s.limiter.Limit(ctx, "sms:tencent")
	if err != nil {
		return fmt.Errorf("限流失败: %w", err)
	}
	if limited {
		return errLimitExceed
	}
	err = s.svc.Send(ctx, tplId, args, numbers...)
	return err
}
