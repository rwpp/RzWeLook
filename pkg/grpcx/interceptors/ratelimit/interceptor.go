package ratelimit

import (
	"context"
	"github.com/rwpp/RzWeLook/pkg/logger"
	"github.com/rwpp/RzWeLook/pkg/ratelimit"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
)

type InterceptorBuilder struct {
	limiter ratelimit.Limiter
	key     string
	l       logger.LoggerV1
}

func NewInterceptorBuilder(limiter ratelimit.Limiter, key string, l logger.LoggerV1) *InterceptorBuilder {
	return &InterceptorBuilder{limiter: limiter, key: key, l: l}
}

func (b *InterceptorBuilder) BuildServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any,
		info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		limited, err := b.limiter.Limit(ctx, b.key)
		if err != nil {
			b.l.Error("限流失败", logger.Error(err))
			return nil, status.Errorf(codes.ResourceExhausted, "触发限流")
		}
		if limited {
			b.l.Error("触发限流", logger.String("key", b.key))
			return nil, status.Errorf(codes.ResourceExhausted, "触发限流")
		}
		return handler(ctx, req)
	}
}
func (b *InterceptorBuilder) BuildClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (err error) {
		limited, err := b.limiter.Limit(ctx, b.key)
		if err != nil {
			b.l.Error("限流失败", logger.Error(err))
			return status.Errorf(codes.ResourceExhausted, "触发限流")
		}
		if limited {
			return status.Errorf(codes.ResourceExhausted, "触发限流")
		}
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
func (b *InterceptorBuilder) BuildServerInterceptorService() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any,
		info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		if strings.HasPrefix(info.FullMethod, "/UserService") {
			limited, err := b.limiter.Limit(ctx, "limiter:service:user:UserService")
			if err != nil {
				b.l.Error("限流失败", logger.Error(err))
				return nil, status.Errorf(codes.ResourceExhausted, "触发限流")
			}
			if limited {
				return nil, status.Errorf(codes.ResourceExhausted, "触发限流")
			}
		}
		return handler(ctx, req)
	}
}
