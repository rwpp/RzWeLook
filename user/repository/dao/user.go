package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrUserDuplicate = errors.New("邮箱/手机号码冲突")
	ErrUserNotFound  = gorm.ErrRecordNotFound
)

type UserDAO interface {
	Insert(ctx context.Context, u User) error
	FindByEmail(ctx context.Context, email string) (User, error)
	FindById(ctx context.Context, id int64) (User, error)
	FindByPhone(ctx context.Context, phone string) (User, error)
	UpdateNonZeroFields(ctx context.Context, u User) error
	FindByWechat(ctx context.Context, openID string) (User, error)
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

func (dao *GormUserDAO) FindByWechat(ctx context.Context, openID string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("wechat_open_id = ?", openID).First(&u).Error
	return u, err
}

func (dao *GormUserDAO) UpdateNonZeroFields(ctx context.Context, u User) error {
	// 这种写法是很不清晰的，因为它依赖了 gorm 的两个默认语义
	// 会使用 ID 来作为 WHERE 条件
	// 会使用非零值来更新
	// 另外一种做法是显式指定只更新必要的字段，
	// 那么这意味着 DAO 和 service 中非敏感字段语义耦合了
	return dao.db.Updates(&u).Error
}

func (dao *GormUserDAO) FindByPhone(ctx context.Context, phone string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("phone = ?", phone).First(&u).Error
	return u, err
}

func (dao *GormUserDAO) FindById(ctx context.Context, id int64) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("id = ?", id).First(&u).Error
	return u, err
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
			return ErrUserDuplicate
		}
	}
	return err
}

type User struct {
	Id int64 `gorm:"primaryKey,autoIncrement"`
	// 设置为唯一索引
	Email    sql.NullString `gorm:"unique"`
	Password string

	//Phone *string
	Phone sql.NullString `gorm:"unique"`

	// 这三个字段表达为 sql.NullXXX 的意思，
	// 就是希望使用的人直到，这些字段在数据库中是可以为 NULL 的
	// 这种做法好处是看到这个定义就知道数据库中可以为 NULL，坏处就是用起来没那么方便
	// 大部分公司不推荐使用 NULL 的列
	// 所以你也可以直接使用 string, int64，那么对应的意思是零值就是每填写
	// 这种做法的好处是用起来好用，但是看代码的话要小心空字符串的问题
	// 生日。一样是毫秒数
	Birthday sql.NullInt64
	// 昵称
	Nickname sql.NullString
	// 自我介绍
	// 指定是 varchar 这个类型，并且长度是 1024
	// 因此你可以看到在 web 里面有这个校验
	AboutMe sql.NullString `gorm:"type:varchar(1024)"`

	// 微信有关数据。有些公司会尝试把这些数据分离出去做一个单独的表
	// 从而避免这个表有过多的列，但是暂时来说
	// 我们还没到这个地步
	WechatOpenId  sql.NullString `gorm:"type:varchar(256);unique"`
	WechatUnionId sql.NullString `gorm:"type:varchar(256)"`

	// 创建时间
	Ctime int64
	// 更新时间
	Utime int64
}
