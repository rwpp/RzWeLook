package circuitbreaker

import (
	"context"
	"github.com/go-kratos/aegis/circuitbreaker"
	"google.golang.org/grpc"
)

type InterceptorBuilder struct {
	breaker circuitbreaker.CircuitBreaker
}

func (b *InterceptorBuilder) BuildServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any,
		info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		if b.breaker.Allow() == nil {
			resp, err = handler(ctx, req)
			//s, ok := status.FromError(err)
			//if s != nil && s.Code() == codes.Unavailable {
			//	b.breaker.MarkFailed()
			//}else{}
			//
			if err != nil {
				b.breaker.MarkFailed()
			} else {
				b.breaker.MarkSuccess()
			}
		}
		b.breaker.MarkFailed()
		//触发熔断器
		return nil, err
	}
}
