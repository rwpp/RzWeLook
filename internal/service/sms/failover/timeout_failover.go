package failover

import (
	"context"
	"github.com/rwpp/RzWeLook/internal/service/sms"
	"sync/atomic"
)

type TimeoutFailoverSMSService struct {
	svcs []sms.Service
	//连续超时个数
	cnt int32
	idx int32
	//阈值，连续超过阈值就要切换服务商
	threshold int32
}

func NewTimeoutFailoverSMSService(svcs []sms.Service) sms.Service {
	return &TimeoutFailoverSMSService{
		svcs: svcs,
	}
}

func (t TimeoutFailoverSMSService) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	idx := atomic.LoadInt32(&t.idx)
	cnt := atomic.LoadInt32(&t.cnt)
	if cnt > t.threshold {
		newIdx := (idx + 1) % int32(len(t.svcs))
		if atomic.CompareAndSwapInt32(&t.idx, idx, newIdx) {
			atomic.StoreInt32(&t.cnt, 0)
		}
		idx = atomic.LoadInt32(&t.idx)
	}
	svc := t.svcs[idx]
	err := svc.Send(ctx, tpl, args, numbers...)
	switch err {
	case context.DeadlineExceeded:
		atomic.AddInt32(&t.cnt, 1)
		return err
	case nil:
		atomic.StoreInt32(&t.cnt, 0)
		return nil
	default:
		//未知错误
		return err
	}
}
