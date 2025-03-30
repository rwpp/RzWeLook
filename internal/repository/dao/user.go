package dao

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrUserDuplicateEmail = errors.New("邮箱冲突")
	ErrUserNotFound       = gorm.ErrRecordNotFound
)

type UserDAO interface {
	Insert(ctx context.Context, u User) error
	FindByEmail(ctx context.Context, email string) (User, error)
}

func NewUserDAO(db *gorm.DB) UserDAO {
	if db == nil {
		panic("database instance is nil")
	}
	return &GormUserDAO{
		db: db,
	}
}

type GormUserDAO struct {
	db *gorm.DB
}

func (dao *GormUserDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	return u, err
}

func (dao *GormUserDAO) Insert(ctx context.Context, u User) error {
	now := time.Now().UnixMilli()
	u.Utime = now
	u.Ctime = now
	//var u User
	err := dao.db.WithContext(ctx).Create(&u).Error
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		const uniqueConflictsErrNo uint16 = 1062
		if mysqlErr.Number == uniqueConflictsErrNo {
			//邮箱或手机号已存在
			return ErrUserDuplicateEmail
		}
	}
	return err
}

type User struct {
	Id    int64  `gorm:"primaryKey,autoIncrement"`
	Email string `gorm:"unique"`
	//Phone    string `gorm:"unique"`
	Password string
	//创建时间
	Ctime int64
	//更新时间
	Utime int64
}
