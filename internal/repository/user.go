package repository

import (
	"context"
	"errors"

	"github.com/rwpp/RzWeLook/internal/domain"
	"github.com/rwpp/RzWeLook/internal/repository/cache"
	"github.com/rwpp/RzWeLook/internal/repository/dao"
)

var (
	ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
	ErrorUserNotFound     = dao.ErrUserNotFound
)

type UserRepository interface {
	Create(ctx context.Context, u domain.User) error
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	FindById(ctx context.Context, id int64) (domain.User, error)
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

func (r *userRepository) FindById(ctx context.Context, id int64) (domain.User, error) {
	u, err := r.cache.Get(ctx, id)
	if err == nil {
		return u, nil
	}
	ue, err := r.dao.FindById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	u = domain.User{
		Id:       ue.Id,
		Email:    ue.Email,
		Password: ue.Password,
	}
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
	return domain.User{
		Id:       u.Id,
		Email:    u.Email,
		Password: u.Password,
	}, nil
}

func (r *userRepository) Create(ctx context.Context, u domain.User) error {
	if r.dao == nil {
		return errors.New("DAO is not initialized")
	}
	err := r.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})
	if err != nil {
		return err
	}
	return nil
}
