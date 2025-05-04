package validator

import (
	"context"
	"github.com/ecodeclub/ekit/slice"
	"github.com/rwpp/RzWeLook/migrator"
	"github.com/rwpp/RzWeLook/migrator/events"
	"github.com/rwpp/RzWeLook/pkg/logger"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
	"time"
)

type Validator[T migrator.Entity] struct {
	base      *gorm.DB
	target    *gorm.DB
	l         logger.LoggerV1
	p         events.Producer
	direction string
	batchSize int
}

func NewValidator[T migrator.Entity](base *gorm.DB, target *gorm.DB, l logger.LoggerV1, p events.Producer, direction string, batchSize int) *Validator[T] {
	return &Validator[T]{
		base:      base,
		target:    target,
		l:         l,
		p:         p,
		direction: direction,
		batchSize: batchSize}
}

func (v *Validator[T]) validate(ctx context.Context) error {
	var eg errgroup.Group
	eg.Go(func() error {
		v.validateBaseToTarget(ctx)
		return nil
	})
	eg.Go(func() error {
		v.validateTargetToBase(ctx)
		return nil
	})
	return eg.Wait()
}

// Validate调用者通过ctx来控制校验程序退出
func (v *Validator[T]) validateBaseToTarget(ctx context.Context) {
	offset := -1
	for {
		dbCtx, cancel := context.WithTimeout(ctx, time.Second)
		offset++
		var src T
		err := v.base.WithContext(dbCtx).Offset(offset).Order("id").First(&src).Error
		cancel()
		switch err {
		case nil:
			var dst T
			err := v.target.Where("id = ?", src.ID()).First(&dst).Error
			switch err {
			case nil:

				if !src.CompareTo(dst) {
					//不相等
					v.notify(ctx, src.ID(), events.InconsistentEventTypeNEQ)
				}
			case gorm.ErrRecordNotFound:
				v.notify(ctx, src.ID(), events.InconsistentEventTypeTargetMissing)
			default:
				v.l.Error("查询数据失败", logger.Error(err))
				continue
			}
		case gorm.ErrRecordNotFound:
			return
		default:
			v.l.Error("校验数据，查询base失败", logger.Error(err))
			continue
			//数据库错误
		}

	}
}
func (v *Validator[T]) validateTargetToBase(ctx context.Context) {
	offset := -v.batchSize
	for {
		offset = offset + v.batchSize
		dbCtx, cancel := context.WithTimeout(ctx, time.Second)
		var dstTs []T
		err := v.base.WithContext(dbCtx).
			Offset(offset).Limit(v.batchSize).
			Order("id").Find(&dstTs).Error
		cancel()
		if len(dstTs) == 0 {
			return
		}
		switch err {
		case gorm.ErrRecordNotFound:
			return
		case nil:
			ids := slice.Map(dstTs, func(idx int, t T) int64 {
				return t.ID()
			})
			var srcTs []T
			err = v.target.Where("id IN ?", ids).Find(&srcTs).Error
			switch err {
			case gorm.ErrRecordNotFound:
				v.notifyBaseMissing(ctx, ids)
			case nil:
				srcIds := slice.Map(srcTs, func(idx int, t T) int64 {
					return t.ID()
				})
				//计算差值
				diff := slice.DiffSet(ids, srcIds)
				v.notifyBaseMissing(ctx, diff)
			}
		default:
			continue
		}
		if len(dstTs) < v.batchSize {
			return
		}
	}
}
func (v *Validator[T]) notifyBaseMissing(ctx context.Context, ids []int64) {
	for _, id := range ids {
		v.notify(ctx, id, events.InconsistentEventTypeBaseMissing)
	}
}
func (v *Validator[T]) notify(ctx context.Context, id int64, typ string) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	err := v.p.ProduceInconsistentEvent(ctx,
		events.InconsistentEvent{
			ID:        id,
			Direction: v.direction,
			Type:      typ,
		})
	cancel()
	if err != nil {
		v.l.Error("校验数据，发送不一致事件失败", logger.Error(err))
	}
}
