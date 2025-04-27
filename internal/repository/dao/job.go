package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type JubDAO interface {
	Preempt(ctx context.Context) (Job, error)
	Release(ctx context.Context, id int64) error
	UpdateUtime(ctx context.Context, id int64) error
	UpdateNextTime(ctx context.Context, id int64, next time.Time) error
	Stop(ctx context.Context, id int64) error
}

func NewGormJobDAO(db *gorm.DB) JubDAO {
	return &GORMJobDAO{
		db: db,
	}
}

type GORMJobDAO struct {
	db *gorm.DB
}

func (g *GORMJobDAO) Stop(ctx context.Context, id int64) error {
	return g.db.WithContext(ctx).
		Where("id=?", id).Updates(map[string]any{
		"status": jobStatusParse,
		"utime":  time.Now().UnixMilli(),
	}).Error
}

func (g *GORMJobDAO) UpdateUtime(ctx context.Context, id int64) error {
	return g.db.WithContext(ctx).Model(&Job{}).
		Where("id = ?", id).Updates(map[string]any{
		"utime": time.Now().UnixMilli(),
	}).Error
}

func (g *GORMJobDAO) UpdateNextTime(ctx context.Context, id int64, next time.Time) error {
	return g.db.WithContext(ctx).Model(&Job{}).
		Where("id = ?", id).Updates(map[string]any{
		"next_time": next.UnixMilli(),
	}).Error
}

func (g *GORMJobDAO) Release(ctx context.Context, id int64) error {
	return g.db.WithContext(ctx).Model(&Job{}).Where("id = ?", id).
		Updates(map[string]any{
			"status": jobStatusWaiting,
			"utime":  time.Now().UnixMilli(),
		}).Error
}

func (g *GORMJobDAO) Preempt(ctx context.Context) (Job, error) {
	db := g.db.WithContext(ctx)
	for {
		now := time.Now()
		var j Job
		err := db.WithContext(ctx).Where("status = ? AND next_time <= ?", jobStatusWaiting, now).
			First(&j).Error
		if err != nil {
			return Job{}, err
		}
		res := db.Where("id = ?AND version = ?", j.Id,
			j.Version).Model(&Job{}).Updates(map[string]any{
			"status":  jobStatusRunning,
			"utime":   now,
			"version": j.Version + 1,
		})
		if res.Error != nil {
			return Job{}, res.Error
		}
		if res.RowsAffected == 0 {
			continue
		}
		return j, nil
	}
}

type Job struct {
	Id       int64 `gorm:"primaryKey,autoIncrement"`
	Cfg      string
	Executor string
	Name     string `gorm:"unique"`
	Status   int
	NextTime int64 `gorm:"index"`
	Cron     string
	Version  int
	// 创建时间
	Ctime int64
	// 更新时间
	Utime int64
}

const (
	jobStatusWaiting = iota
	jobStatusRunning
	jobStatusParse
)
