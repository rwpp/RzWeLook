package service

import (
	"context"
	"errors"
	"github.com/ecodeclub/ekit/queue"
	"github.com/ecodeclub/ekit/slice"
	intrv1 "github.com/rwpp/RzWeLook/api/proto/gen/intr/v1"
	"github.com/rwpp/RzWeLook/internal/repository"
	"math"
	"time"

	"github.com/rwpp/RzWeLook/internal/domain"
)

type RankingService interface {
	TopN(ctx context.Context) error
}

func NewRankingService(artSvc ArticleService,
	intr intrv1.InteractiveServiceClient,
	repo repository.RankingRepository) RankingService {
	return &BatchRankingService{
		artSvc:    artSvc,
		intr:      intr,
		batchSize: 100,
		n:         100,
		repo:      repo,
		scoreFunc: func(t time.Time, likeCnt int64) float64 {
			sec := time.Since(t).Seconds()
			return float64(likeCnt-1) / math.Pow(float64(sec+2), 1.5)
		},
	}
}

type BatchRankingService struct {
	artSvc    ArticleService
	intr      intrv1.InteractiveServiceClient
	repo      repository.RankingRepository
	batchSize int
	n         int
	scoreFunc func(t time.Time, likeCnt int64) float64
}

func (svc *BatchRankingService) TopN(ctx context.Context) error {
	arts, err := svc.topN(ctx)
	if err != nil {
		return err
	}
	return svc.repo.ReplaceTopN(ctx, arts)
}
func (svc *BatchRankingService) topN(ctx context.Context) ([]domain.Article, error) {
	now := time.Now()
	offset := 0
	type Score struct {
		art   domain.Article
		score float64
	}
	topN := queue.NewPriorityQueue[Score](svc.n,
		func(src Score, dst Score) int {
			if src.score > dst.score {
				return 1
			} else if src.score == dst.score {
				return 0
			}
			return -1
		})
	for {
		arts, err := svc.artSvc.ListPub(ctx, now, offset, svc.batchSize)
		if err != nil {
			return nil, err
		}
		ids := slice.Map[domain.Article, int64](arts,
			func(idx int, src domain.Article) int64 {
				return src.Id
			})
		intrs, err := svc.intr.GetByIds(ctx, &intrv1.GetByIdRequest{
			Biz:    "article",
			BizIds: ids,
		})
		if err != nil {
			return nil, err
		}
		if len(intrs.Intrs) == 0 {
			return nil, errors.New("没有数据")
		}
		for _, art := range arts {
			intr := intrs.Intrs[art.Id]
			score := svc.scoreFunc(art.Utime, intr.LikeCnt)
			err = topN.Enqueue(Score{
				art:   art,
				score: score,
			})
			if err == queue.ErrOutOfCapacity {
				val, _ := topN.Dequeue()
				if val.score < score {
					err = topN.Enqueue(Score{
						art:   art,
						score: score,
					})
				} else {
					_ = topN.Enqueue(val)
				}
			}
		}
		if len(arts) < svc.batchSize || now.Sub(arts[len(arts)-1].Utime).Hours() > 7*24 {
			break
		}
		offset = offset + len(arts)
	}
	res := make([]domain.Article, svc.n)
	for i := svc.n - 1; i >= 0; i-- {
		val, err := topN.Dequeue()
		if err != nil {
			break
		}
		res[i] = val.art
	}
	return res, nil
}
