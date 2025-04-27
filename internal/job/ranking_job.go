package job

import (
	"context"
	rlock "github.com/gotomicro/redis-lock"
	"github.com/rwpp/RzWeLook/internal/service"
	"github.com/rwpp/RzWeLook/pkg/logger"
	"sync"
	"time"
)

type RankingJob struct {
	svc       service.RankingService
	timeout   time.Duration
	client    *rlock.Client
	key       string
	l         logger.LoggerV1
	lock      *rlock.Lock
	localLock *sync.Mutex
}

func NewRankingJob(svc service.RankingService,
	client *rlock.Client,
	l logger.LoggerV1,
	timeout time.Duration,
) *RankingJob {
	return &RankingJob{
		svc:       svc,
		timeout:   timeout,
		client:    client,
		l:         l,
		key:       "rlock:cron_job:ranking",
		localLock: &sync.Mutex{},
	}
}

func (r *RankingJob) Name() string {
	return "Ranking"
}

func (r *RankingJob) Run() error {
	r.localLock.Lock()
	defer r.localLock.Unlock()
	if r.lock == nil {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		lock, err := r.client.Lock(ctx, r.key, r.timeout, &rlock.FixIntervalRetry{
			Interval: time.Millisecond * 100,
			Max:      0,
		}, time.Second)
		if err != nil {
			return nil
		}
		r.lock = lock
		go func() {
			//自动续约
			err1 := lock.AutoRefresh(r.timeout/2, time.Second)
			if err1 != nil {
				r.l.Error("自动续约失败", logger.Error(err1))
			}
			r.localLock.Lock()
			r.lock = nil
			r.localLock.Unlock()
		}()
	}
	//defer func() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	//	err = lock.Unlock(ctx)
	//	if err != nil {
	//		r.l.Error("释放分布式锁失败", logger.Error(err))
	//	}
	//}()
	return r.svc.TopN(ctx)
}

func (r *RankingJob) Close() error {
	r.localLock.Lock()
	lock := r.lock
	r.lock = nil
	r.localLock.Unlock()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return lock.Unlock(ctx)
}
