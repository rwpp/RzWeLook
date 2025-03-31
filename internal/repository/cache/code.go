package cache

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
)

var (
	ErrCodeSendTooMany   = errors.New("发送验证码次数过多")
	ErrCodeVerifyTooMany = errors.New("验证码验证次数过多")
	ErrUnknownForCode    = errors.New("未知错误")
)

//go:embed lua/set_code.lua
var luaSetCode string

//go:embed lua/verify_code.lua
var luaVerifyCode string

type CodeCacheInterface interface {
	Set(ctx context.Context, biz, phone string, code string) error
	Verify(ctx context.Context, biz, phone string, code string) (bool, error)
}
type CodeCache struct {
	client redis.Cmdable
}

func NewCodeCache(client redis.Cmdable) CodeCacheInterface {
	return &CodeCache{
		client: client,
	}
}
func (c CodeCache) Set(ctx context.Context, biz, phone string, code string) error {
	res, err := c.client.Eval(ctx, luaSetCode, []string{c.key(biz, phone)}, code).Int()
	if err != nil {
		return err
	}
	switch res {
	case 0:
		return nil
	case -1:
		return ErrCodeSendTooMany
	default:
		return errors.New("未知错误")

	}
}

func (c *CodeCache) key(biz, phone string) string {
	return fmt.Sprintf("phone_code%s:%s", biz, phone)
}

func (c CodeCache) Verify(ctx context.Context, biz, phone string, inputCode string) (bool, error) {
	res, err := c.client.Eval(ctx, luaVerifyCode, []string{c.key(biz, phone)}, inputCode).Int()
	if err != nil {
		return false, err
	}
	switch res {
	case 0:
		return true, nil
	case -1:
		return false, ErrCodeVerifyTooMany
	case -2:
		return false, nil
	default:
		return false, ErrUnknownForCode
	}
}
