package metrics

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rwpp/RzWeLook/internal/service/sms"
	"time"
)

type PrometheusDecorator struct {
	svc    sms.Service
	vector *prometheus.SummaryVec
}

func NewPrometheusDecorator(svc sms.Service) *PrometheusDecorator {
	vector := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: "RzWeLook",
		Subsystem: "web",
		Name:      "sms_resp_time",
		Help:      "统计SMS执行时间",
	}, []string{"biz"})
	prometheus.MustRegister(vector)
	return &PrometheusDecorator{
		svc:    svc,
		vector: vector,
	}
}
func (p *PrometheusDecorator) Send(ctx context.Context,
	biz string, args []string, numbers ...string) error {
	startTime := time.Now()
	defer func() {
		duration := time.Since(startTime).Milliseconds()
		p.vector.WithLabelValues(biz).
			Observe(float64(duration))
	}()
	return p.svc.Send(ctx, biz, args, numbers...)
}
