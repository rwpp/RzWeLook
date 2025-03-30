package repository

import (
	"context"
	"errors"
	"github.com/rwpp/RzWeLook/internal/domain"
	"github.com/rwpp/RzWeLook/internal/repository/dao"
	"log"
)

var (
	ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
	ErrorUserNotFound     = dao.ErrUserNotFound
)

type UserRepository interface {
	Create(ctx context.Context, u domain.User) error
	FindByEmail(ctx context.Context, email string) (domain.User, error)
}

func NewUserRepository(dao dao.UserDAO) UserRepository {
	return &userRepository{
		dao: dao,
	}
}

type userRepository struct {
	dao dao.UserDAO
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
	log.Printf("Repository: Creating user with email: %s", u.Email)
	if r.dao == nil {
		log.Println("Repository: DAO is nil")
		return errors.New("DAO is not initialized")
	}

	err := r.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})
	if err != nil {
		log.Printf("Repository: Failed to create user: %v", err)
		return err
	}
	log.Println("Repository: User created successfully")
	return nil
}
