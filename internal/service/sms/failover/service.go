package failover

import (
	"context"
	"errors"
	"github.com/rwpp/RzWeLook/internal/service/sms"
	"log"
)

type FailoverSMSService struct {
	svcs []sms.Service
}

func NewFailoverSMSService(svcs []sms.Service) sms.Service {
	return &FailoverSMSService{
		svcs: svcs,
	}
}
func (f *FailoverSMSService) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	for _, svc := range f.svcs {
		err := svc.Send(ctx, tpl, args, numbers...)
		//发送成功
		if err == nil {
			return err
		}

		log.Println(err)
	}
	return errors.New("全部服务商都失败")
}
