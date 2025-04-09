package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/rwpp/RzWeLook/internal/domain"
	"github.com/rwpp/RzWeLook/internal/repository/cache"
	"github.com/rwpp/RzWeLook/internal/repository/dao"
)

var (
	ErrUserDuplicate  = dao.ErrUserDuplicate
	ErrorUserNotFound = dao.ErrUserNotFound
)

type UserRepository interface {
	Create(ctx context.Context, u domain.User) error
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	FindById(ctx context.Context, id int64) (domain.User, error)
	FindByPhone(ctx context.Context, phone string) (domain.User, error)
	Update(ctx context.Context, u domain.User) error
	FindByWechat(ctx context.Context, openID string) (domain.User, error)
}

func NewUserRepository(dao dao.UserDAO, cache cache.UserCacheInterface) UserRepository {
	return &userRepository{
		dao:   dao,
		cache: cache,
	}
}

type userRepository struct {
	dao   dao.UserDAO
	cache cache.UserCacheInterface
}

func (r *userRepository) FindByWechat(ctx context.Context, openID string) (domain.User, error) {
	u, err := r.dao.FindByWechat(ctx, openID)
	if err != nil {
		return domain.User{}, err
	}
	return r.entityToDomain(u), nil
}

func (r *userRepository) Update(ctx context.Context, u domain.User) error {
	err := r.dao.UpdateNonZeroFields(ctx, r.domainToEntity(u))
	if err != nil {
		return err
	}
	return r.cache.Delete(ctx, u.Id)
}

func (r *userRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	u, err := r.dao.FindByPhone(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}
	return r.entityToDomain(u), nil
}

func (r *userRepository) FindById(ctx context.Context, id int64) (domain.User, error) {
	u, err := r.cache.Get(ctx, id)
	if err == nil {
		return u, nil
	}
	ue, err := r.dao.FindById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	u = r.entityToDomain(ue)
	go func() {
		err = r.cache.Set(ctx, u)
		if err != nil {
			//return domain.User{}, err
		}
	}()
	return u, err
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := r.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return r.entityToDomain(u), nil
}

func (r *userRepository) Create(ctx context.Context, u domain.User) error {
	if r.dao == nil {
		return errors.New("DAO is not initialized")
	}
	err := r.dao.Insert(ctx, r.domainToEntity(u))
	if err != nil {
		return err
	}
	return nil
}

func (r *userRepository) domainToEntity(u domain.User) dao.User {
	return dao.User{
		Id: u.Id,
		Email: sql.NullString{
			String: u.Email,
			Valid:  u.Email != ""},
		Phone: sql.NullString{
			String: u.Phone,
			Valid:  u.Phone != ""},
		Password: u.Password,
		WechatOpenId: sql.NullString{
			String: u.WechatInfo.OpenID,
			Valid:  u.WechatInfo.OpenID != "",
		},
		WechatUnionId: sql.NullString{
			String: u.WechatInfo.UnionID,
			Valid:  u.WechatInfo.UnionID != "",
		},
		Birthday: sql.NullInt64{
			Int64: u.Birthday.UnixMilli(),
			Valid: !u.Birthday.IsZero(),
		},
		Nickname: sql.NullString{
			String: u.Nickname,
			Valid:  u.Nickname != "",
		},
		AboutMe: sql.NullString{
			String: u.AboutMe,
			Valid:  u.AboutMe != "",
		},
		Ctime: u.Ctime.UnixMilli(),
	}
}
func (r *userRepository) entityToDomain(u dao.User) domain.User {
	var birthday time.Time
	if u.Birthday.Valid {
		birthday = time.UnixMilli(u.Birthday.Int64)
	}
	return domain.User{
		Id:       u.Id,
		Email:    u.Email.String,
		Password: u.Password,
		Phone:    u.Phone.String,
		Nickname: u.Nickname.String,
		AboutMe:  u.AboutMe.String,
		Birthday: birthday,
		WechatInfo: domain.WechatInfo{
			UnionID: u.WechatUnionId.String,
			OpenID:  u.WechatOpenId.String,
		},
		Ctime: time.UnixMilli(u.Ctime),
		Utime: time.UnixMilli(u.Utime),
	}
}
