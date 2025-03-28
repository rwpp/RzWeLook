package dao

import (
	"context"
	"gorm.io/gorm"
)

type UserDAO interface {
	Insert(ctx context.Context, u User) error
}

func NewUserDAO(db *gorm.DB) UserDAO {
	return &GormUserDAO{
		db: db,
	}
}

type GormUserDAO struct {
	db *gorm.DB
}

func (dao *GormUserDAO) Insert(ctx context.Context, u User) error {
	//TODO implement me
	panic("implement me")
}

type User struct {
	Id       int64  `gorm:"primaryKey"`
	Email    string `gorm:"unique"`
	Phone    string `gorm:"unique"`
	Password string
	//创建时间
	Ctime int64
	//更新时间
	Utime int64
}
