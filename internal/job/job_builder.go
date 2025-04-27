package job

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/robfig/cron/v3"
	"github.com/rwpp/RzWeLook/pkg/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"strconv"
	"time"
)

type CronJobBuilder struct {
	l      logger.LoggerV1
	p      *prometheus.SummaryVec
	tracer trace.Tracer
}

func NewCronJobBuilder(l logger.LoggerV1) *CronJobBuilder {
	p := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: "RzWeLook",
		Subsystem: "job",
		Help:      "统计定时任务的执行情况",
		Name:      "cron_job_duration",
	}, []string{"name", "success"})
	prometheus.MustRegister(p)
	return &CronJobBuilder{
		l:      l,
		p:      p,
		tracer: otel.GetTracerProvider().Tracer("internal/job/job_builder.go"),
	}
}
func (b *CronJobBuilder) Build(job Job) cron.Job {
	name := job.Name()
	b.l.Error("cron job %s failed: %v",
		logger.String("name", name))
	return cronJobFuncAdapter(func() error {
		_, span := b.tracer.Start(context.Background(), name)
		defer span.End()
		start := time.Now()
		var success bool
		defer func() {
			b.l.Error("任务结束",
				logger.String("name", name))
			duration := time.Since(start).Milliseconds()
			b.p.WithLabelValues(name,
				strconv.FormatBool(success)).Observe(float64(duration))
		}()
		err := job.Run()
		success = err == nil
		if err != nil {
			span.RecordError(err)
			b.l.Error("运行任务失败", logger.Error(err),
				logger.String("name", name))
		}
		return nil
	})
}

type cronJobFuncAdapter func() error

func (c cronJobFuncAdapter) Run() {
	_ = c()
}
