package cache

import (
	"context"
	"github.com/rwpp/RzWeLook/reward/domain"
)

type RewardCache interface {
	GetCachedCodeURL(ctx context.Context, r domain.Reward) (domain.CodeURL, error)
	CachedCodeURL(ctx context.Context, cu domain.CodeURL, r domain.Reward) error
}
