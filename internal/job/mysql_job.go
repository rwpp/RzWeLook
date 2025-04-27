package job

import (
	"context"
	"fmt"
	"github.com/rwpp/RzWeLook/internal/domain"
	"github.com/rwpp/RzWeLook/pkg/logger"
	"golang.org/x/sync/semaphore"
	"time"

	"github.com/rwpp/RzWeLook/internal/service"
)

type Executor interface {
	Name() string
	Exec(ctx context.Context, j domain.Job) error
}

type LocalFuncExecutor struct {
	funcs map[string]func(ctx context.Context, j domain.Job) error
}

func NewLocalFuncExecutor() *LocalFuncExecutor {
	return &LocalFuncExecutor{
		funcs: make(map[string]func(ctx context.Context, j domain.Job) error)}
}

func (l *LocalFuncExecutor) Name() string {
	return "local"
}
func (l *LocalFuncExecutor) RegisterFunc(name string, fn func(ctx context.Context, j domain.Job) error) {
	l.funcs[name] = fn
}
func (l *LocalFuncExecutor) Exec(ctx context.Context, j domain.Job) error {
	fn, ok := l.funcs[j.Name]
	if !ok {
		return fmt.Errorf("未知任务: %s", j.Name)
	}
	return fn(ctx, j)
}

// Scheduler 调度器
type Scheduler struct {
	execs   map[string]Executor
	svc     service.JobService
	l       logger.LoggerV1
	limiter *semaphore.Weighted
}

func NewScheduler(svc service.JobService, l logger.LoggerV1) *Scheduler {
	return &Scheduler{
		limiter: semaphore.NewWeighted(200),
		execs:   make(map[string]Executor),
		svc:     svc,
		l:       l,
	}
}
func (s *Scheduler) RegisterExecutor(exec Executor) {
	s.execs[exec.Name()] = exec
}
func (s *Scheduler) Schedule(ctx context.Context) error {
	for {
		if ctx.Err() != nil {
			//退出调度循环
			return ctx.Err()
		}
		err := s.limiter.Acquire(ctx, 1)
		if err != nil {
			return err
		}
		//一次调度数据库所需的查询时间
		dbCtx, cancel := context.WithTimeout(ctx, time.Second)
		j, err := s.svc.Preempt(dbCtx)
		cancel()
		if err != nil {
			s.l.Error("抢占任务失败", logger.Error(err))
		}
		exec, ok := s.execs[j.Executor]
		if !ok {
			s.l.Error("未知执行器", logger.String("executor", j.Executor))
			continue
		}
		go func() {
			defer func() {
				s.limiter.Release(1)
				err1 := j.CancelFunc()
				if err1 != nil {
					s.l.Error("取消任务失败",
						logger.Error(err1),
						logger.Int64("jid", j.Id))
				}
			}()
			err1 := exec.Exec(ctx, j)
			if err1 != nil {
				s.l.Error("任务执行失败", logger.Error(err1))
			}
			ctx, cancel = context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			err1 = s.svc.ResetNextTime(ctx, j)
			if err1 != nil {
				s.l.Error("设置下一次执行时间失败", logger.Error(err1))
			}
		}()
	}
}
