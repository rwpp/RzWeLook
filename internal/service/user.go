package service

import (
	"context"
	"errors"
	"github.com/rwpp/RzWeLook/internal/domain"
	"github.com/rwpp/RzWeLook/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"log"
)

var (
	ErrUserDuplicateEmail    = repository.ErrUserDuplicateEmail
	ErrInvalidUserOrPassword = errors.New("账号/邮箱或密码不对")
)

type UserService interface {
	SignUp(ctx context.Context, u domain.User) error
	Login(ctx context.Context, email, password string) (domain.User, error)
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

type userService struct {
	repo repository.UserRepository
}

func (svc *userService) Login(ctx context.Context, email, password string) (domain.User, error) {
	u, err := svc.repo.FindByEmail(ctx, email)
	if err == repository.ErrorUserNotFound {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	if err != nil {
		return domain.User{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	return u, nil
}

func (svc *userService) SignUp(ctx context.Context, u domain.User) error {
	log.Printf("Service: Starting signup process for email: %s", u.Email)
	if svc.repo == nil {
		log.Println("Service: Repository is nil")
		return errors.New("repository is not initialized")
	}
	// 先检查邮箱是否已存在
	_, err := svc.repo.FindByEmail(ctx, u.Email)
	if err == nil {
		log.Printf("Service: Email %s already exists", u.Email)
		return ErrUserDuplicateEmail
	}
	if !errors.Is(err, repository.ErrorUserNotFound) {
		log.Printf("Service: Error checking existing email: %v", err)
		return err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Service: Failed to hash password: %v", err)
		return err
	}
	u.Password = string(hash)

	log.Println("Service: Calling repository.Create")
	err = svc.repo.Create(ctx, u)
	if err != nil {
		log.Printf("Service: Failed to create user: %v", err)
		return err
	}
	log.Println("Service: User created successfully")
	return nil
}
