package cache

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"github.com/rwpp/RzWeLook/internal/domain"
	"time"
)

type RankingCache interface {
	Set(ctx context.Context, arts []domain.Article) error
	Get(ctx context.Context) ([]domain.Article, error)
}

func NewRankingCache(client redis.Cmdable) RankingCache {
	return &RankingRedisCached{
		client: client,
		key:    "ranking",
	}
}

type RankingRedisCached struct {
	client redis.Cmdable
	key    string
}

func (r *RankingRedisCached) Set(ctx context.Context, arts []domain.Article) error {
	for i := 0; i < len(arts); i++ {
		arts[i].Content = ""
	}
	val, err := json.Marshal(arts)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, r.key, string(val), time.Minute*10).Err()
}

func (r *RankingRedisCached) Get(ctx context.Context) ([]domain.Article, error) {
	data, err := r.client.Get(ctx, r.key).Bytes()
	if err != nil {
		return nil, err
	}
	var res []domain.Article
	err = json.Unmarshal(data, &res)
	return res, err
}
