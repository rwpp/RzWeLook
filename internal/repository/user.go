package repository

import (
	"context"
	"github.com/rwpp/RzWeLook/internal/domain"
	"github.com/rwpp/RzWeLook/internal/repository/dao"
)

type UserRepository interface {
	Create(ctx context.Context, user domain.User) error
}

func NewUserRepository(dao dao.UserDAO) UserRepository {
	return &userRepository{
		dao: dao,
	}
}

type userRepository struct {
	dao dao.UserDAO
}

func (u userRepository) Create(ctx context.Context, user domain.User) error {
	//TODO implement me
	panic("implement me")
}
