package ratelimit

import (
	"context"
	_ "embed"
	"github.com/redis/go-redis/v9"
	"time"
)

//go:embed slide_window.lua
var luaSlideWindow string

type RedisSlideWindowLimiter struct {
	cmd      redis.Cmdable
	interval time.Duration
	// 阈值
	rate int
}

func NewRedisSlideWindowLimiter(cmd redis.Cmdable,
	interval time.Duration, rate int) Limiter {
	return &RedisSlideWindowLimiter{
		cmd:      cmd,
		interval: interval,
		rate:     rate,
	}
}
func (r *RedisSlideWindowLimiter) Limit(ctx context.Context, key string) (bool, error) {
	return r.cmd.Eval(ctx, luaSlideWindow, []string{key},
		r.interval.Milliseconds(), r.rate, time.Now().UnixMilli()).Bool()
}
