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
	base          *gorm.DB
	target        *gorm.DB
	l             logger.LoggerV1
	p             events.Producer
	direction     string
	batchSize     int
	utime         int64
	sleepInterval time.Duration
}

func NewValidator[T migrator.Entity](base *gorm.DB, target *gorm.DB, l logger.LoggerV1, p events.Producer, direction string, batchSize int) *Validator[T] {
	return &Validator[T]{
		base:      base,
		target:    target,
		l:         l,
		p:         p,
		direction: direction,
		batchSize: batchSize,
	}
}

// func NewIncrValidator[T migrator.Entity](base *gorm.DB,
//
//		target *gorm.DB, l logger.LoggerV1,
//		p events.Producer, direction string, batchSize int) {
//		v := newValidator[T](base, target, l, p, direction, batchSize)
//		v.order = "utime"
//	}
//
// func NewFullValidator[T migrator.Entity](base *gorm.DB,
//
//		target *gorm.DB, l logger.LoggerV1,
//		p events.Producer, direction string, batchSize int) {
//		v := newValidator[T](base, target, l, p, direction, batchSize)
//		v.order = "id"
//	}
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
	offset := 0
	for {
		dbCtx, cancel := context.WithTimeout(ctx, time.Second)

		var src T
		err := v.base.WithContext(dbCtx).
			Where("utime>?", v.utime).Offset(offset).
			Order("utime").First(&src).Error
		cancel()
		switch err {
		case nil:
			var dst T
			err = v.target.Where("id = ?", src.ID()).First(&dst).Error
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
			}

		case gorm.ErrRecordNotFound:
			if v.sleepInterval <= 0 {
				return
			}
			time.Sleep(v.sleepInterval)
			continue
		default:
			v.l.Error("校验数据，查询base失败", logger.Error(err))
			//数据库错误
		}
		offset++
	}
}
func (v *Validator[T]) validateTargetToBase(ctx context.Context) {
	offset := 0
	for {
		dbCtx, cancel := context.WithTimeout(ctx, time.Second)
		var dstTs []T
		err := v.target.WithContext(dbCtx).
			Where("utime>?", v.utime).
			Select("id").
			Offset(offset).Limit(v.batchSize).
			Order("utime").Find(&dstTs).Error
		cancel()
		if len(dstTs) == 0 {
			if v.sleepInterval <= 0 {
				return
			}
			time.Sleep(v.sleepInterval)
			return
		}
		switch err {
		case gorm.ErrRecordNotFound:
			if v.sleepInterval <= 0 {
				return
			}
			time.Sleep(v.sleepInterval)
			continue
		case nil:
			ids := slice.Map(dstTs, func(idx int, t T) int64 {
				return t.ID()
			})
			var srcTs []T
			err = v.base.Where("id IN ?", ids).Find(&srcTs).Error
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
			default:

			}
		default:
			v.l.Error("查询target失败", logger.Error(err))
		}
		offset += len(dstTs)
		if len(dstTs) < v.batchSize {
			if v.sleepInterval <= 0 {
				return
			}
			time.Sleep(v.sleepInterval)

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
