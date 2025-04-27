package ioc

import (
	rlock "github.com/gotomicro/redis-lock"
	"github.com/robfig/cron/v3"
	"github.com/rwpp/RzWeLook/internal/job"
	"github.com/rwpp/RzWeLook/internal/service"
	"github.com/rwpp/RzWeLook/pkg/logger"
	"time"
)

func InitRankingJob(svc service.RankingService,
	rlockClient *rlock.Client,
	l logger.LoggerV1) *job.RankingJob {
	return job.NewRankingJob(svc, rlockClient, l, time.Second*30)

}
func InitJobs(l logger.LoggerV1, rankingJob *job.RankingJob) *cron.Cron {
	res := cron.New(cron.WithSeconds())
	cbd := job.NewCronJobBuilder(l)
	_, err := res.AddJob("0 */3 * * * ?", cbd.Build(rankingJob))
	if err != nil {
		panic(err)
	}
	return res
}
