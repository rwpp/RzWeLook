package repository

import (
	"context"
	"github.com/rwpp/RzWeLook/internal/repository/cache"
)

var (
	ErrCodeSendTooMany   = cache.ErrCodeSendTooMany
	ErrCodeVerifyTooMany = cache.ErrCodeVerifyTooMany
)

type CodeRepositoryInterface interface {
	Store(ctx context.Context, biz string, phone string, code string) error
	Verify(ctx context.Context, biz, phone string, inputCode string) (bool, error)
}

func NewCodeRepository(c cache.CodeCacheInterface) CodeRepositoryInterface {
	return &CodeRepository{
		cache: c,
	}
}

type CodeRepository struct {
	cache cache.CodeCacheInterface
}

func (repo *CodeRepository) Verify(ctx context.Context, biz, phone string, inputCode string) (bool, error) {
	return repo.cache.Verify(ctx, biz, phone, inputCode)
}

func (repo *CodeRepository) Store(ctx context.Context, biz string, phone string, code string) error {
	return repo.cache.Set(ctx, biz, phone, code)
}
