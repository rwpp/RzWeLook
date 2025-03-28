package service

import (
	"context"
	"github.com/rwpp/RzWeLook/internal/domain"
	"github.com/rwpp/RzWeLook/internal/repository"
)

type UserService interface {
	SignUp(ctx context.Context, u domain.User) error
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

type userService struct {
	repo repository.UserRepository
}

func (svc userService) SignUp(ctx context.Context, u domain.User) error {
	err := svc.repo.Create(ctx, u)
	if err != nil {
		return err
	}
	return nil
}
