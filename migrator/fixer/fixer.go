package fixer

import (
	"context"
	"errors"
	"github.com/rwpp/RzWeLook/migrator"
	"github.com/rwpp/RzWeLook/migrator/events"
	"gorm.io/gorm"
)

type Fixer[T migrator.Entity] struct {
	base   *gorm.DB
	target *gorm.DB
}

func (f *Fixer[T]) Fix(ctx context.Context, evt events.InconsistentEvent) error {
	switch evt.Type {
	case events.InconsistentEventTypeTargetMissing:
		//执行插入
		var t T
		err := f.base.WithContext(ctx).
			Where("id =?", evt.ID).First(&t).Error
		switch err {
		case nil:
			return f.target.Create(&t).Error
		case gorm.ErrRecordNotFound:
			//刚被删
			return nil
		default:
			return err
		}
	case events.InconsistentEventTypeNEQ:
		//执行更新
		var t T
		err := f.base.WithContext(ctx).
			Where("id =?", evt.ID).First(&t).Error
		switch err {
		case nil:
			return f.target.Updates(&t).Error
		case gorm.ErrRecordNotFound:
			//刚被删
			return f.target.WithContext(ctx).
				Where("id=?", evt.ID).First(&t).Error
		default:
			return err
		}
	case events.InconsistentEventTypeBaseMissing:
		return f.target.WithContext(ctx).
			Where("id=?", evt.ID).Delete(new(T)).Error
	default:
		return errors.New("未知错误")
	}
}
