package ioc

import (
	"context"
	"github.com/rwpp/RzWeLook/internal/domain"
	"github.com/rwpp/RzWeLook/internal/job"
	"github.com/rwpp/RzWeLook/internal/service"
	"github.com/rwpp/RzWeLook/pkg/logger"
	"time"
)

func InitScheduler(svc service.JobService,
	local *job.LocalFuncExecutor,
	l logger.LoggerV1) *job.Scheduler {
	res := job.NewScheduler(svc, l)
	res.RegisterExecutor(local)
	return res
}
func InitLocalFuncExecutor(svc service.RankingService) *job.LocalFuncExecutor {
	res := job.NewLocalFuncExecutor()
	res.RegisterFunc("ranking", func(ctx context.Context, j domain.Job) error {
		ctx, cancel := context.WithTimeout(ctx, time.Second*30)
		defer cancel()
		return svc.TopN(ctx)
	})
	return res
}
