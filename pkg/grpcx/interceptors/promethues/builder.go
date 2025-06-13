package promethues

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rwpp/RzWeLook/pkg/grpcx/interceptors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"strings"
	"time"
)

type InterceptorBuilder struct {
	Namespace string
	Subsystem string
	interceptors.Builder
}

func (b *InterceptorBuilder) buildUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	// ServerHandleHistogram ...
	summary := prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: b.Namespace,
			Subsystem: b.Subsystem,
			Name:      "server_handle_seconds",
			Objectives: map[float64]float64{
				0.5:   0.01,
				0.9:   0.01,
				0.95:  0.01,
				0.99:  0.001,
				0.999: 0.0001,
			},
		}, []string{"type", "service", "method", "peer", "code"})
	prometheus.MustRegister(summary)
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		start := time.Now()
		defer func() {
			serviceName, method := b.splitMethodName(info.FullMethod)
			st, _ := status.FromError(err)
			code := "OK"
			if st != nil {
				code = st.Code().String()
			}
			summary.WithLabelValues("unary", serviceName, method,
				b.PeerName(ctx), code).Observe(float64(time.Since(start).Milliseconds()))
		}()
		resp, err = handler(ctx, req)
		return
	}
}

func (b *InterceptorBuilder) splitMethodName(fullMethodName string) (string, string) {
	fullMethodName = strings.TrimPrefix(fullMethodName, "/") // remove leading slash
	if i := strings.Index(fullMethodName, "/"); i >= 0 {
		return fullMethodName[:i], fullMethodName[i+1:]
	}
	return "unknown", "unknown"
}
