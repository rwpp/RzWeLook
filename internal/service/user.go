package service

import (
	"context"
	"errors"
	"github.com/rwpp/RzWeLook/internal/domain"
	"github.com/rwpp/RzWeLook/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserDuplicateEmail    = repository.ErrUserDuplicateEmail
	ErrInvalidUserOrPassword = errors.New("账号/邮箱或密码不对")
)

type UserService interface {
	SignUp(ctx context.Context, u domain.User) error
	Login(ctx context.Context, email, password string) error
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

type userService struct {
	repo repository.UserRepository
}

func (svc *userService) Login(ctx context.Context, email, password string) error {
	u, err := svc.repo.FindByEmail(ctx, email)
	if err == repository.ErrorUserNotFound {
		return ErrInvalidUserOrPassword
	}
	if err != nil {
		return err
	}
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return ErrInvalidUserOrPassword
	}
	return nil
}

func (svc *userService) SignUp(ctx context.Context, u domain.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	err = svc.repo.Create(ctx, u)
	if err != nil {
		return err
	}
	return nil
}
