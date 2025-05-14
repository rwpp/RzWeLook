package service

import (
	"context"
	"errors"
	"github.com/rwpp/RzWeLook/user/domain"
	"github.com/rwpp/RzWeLook/user/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserDuplicate         = repository.ErrUserDuplicate
	ErrInvalidUserOrPassword = errors.New("账号/邮箱或密码不对")
)

type UserService interface {
	FindOrCreate(ctx context.Context, phone string) (domain.User, error)
	SignUp(ctx context.Context, u domain.User) error
	Login(ctx context.Context, email, password string) (domain.User, error)
	Profile(ctx context.Context, id int64) (domain.User, error)
	UpdateNonSensitiveInfo(ctx context.Context, user domain.User) error
	FindOrCreateByWechat(ctx context.Context, wechatInfo domain.WechatInfo) (domain.User, error)
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

type userService struct {
	repo repository.UserRepository
}

func (svc *userService) FindOrCreateByWechat(ctx context.Context, info domain.WechatInfo) (domain.User, error) {
	u, err := svc.repo.FindByWechat(ctx, info.OpenID)
	if err != repository.ErrorUserNotFound {
		return u, err
	}
	u = domain.User{
		WechatInfo: info,
	}
	err = svc.repo.Create(ctx, u)
	if err != nil && err != repository.ErrUserDuplicate {
		return u, err
	}
	return svc.repo.FindByWechat(ctx, info.OpenID)
}

func (svc *userService) UpdateNonSensitiveInfo(ctx context.Context, user domain.User) error {
	// 写法1
	// 这种是简单的写法，依赖与 Web 层保证没有敏感数据被修改
	// 也就是说，你的基本假设是前端传过来的数据就是不会修改 Email，Phone 之类的信息的。
	//return svc.repo.Update(ctx, user)

	// 写法2
	// 这种是复杂写法，依赖于 repository 中更新会忽略 0 值
	// 这个转换的意义在于，你在 service 层面上维护住了什么是敏感字段这个语义
	user.Email = ""
	user.Phone = ""
	user.Password = ""
	return svc.repo.Update(ctx, user)
}
func (svc *userService) FindOrCreate(ctx context.Context, phone string) (domain.User, error) {
	u, err := svc.repo.FindByPhone(ctx, phone)
	if err != repository.ErrorUserNotFound {
		return u, err
	}
	u = domain.User{
		Phone: phone,
	}
	err = svc.repo.Create(ctx, u)
	if err != nil && err != repository.ErrUserDuplicate {
		return u, err
	}
	return svc.repo.FindByPhone(ctx, phone)
}

func (svc *userService) Profile(ctx context.Context, id int64) (domain.User, error) {
	u, err := svc.repo.FindById(ctx, id)
	return u, err
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
	if svc.repo == nil {
		return errors.New("repository is not initialized")
	}
	// 先检查邮箱是否已存在
	_, err := svc.repo.FindByEmail(ctx, u.Email)
	if err == nil {
		return ErrUserDuplicate
	}
	if !errors.Is(err, repository.ErrorUserNotFound) {
		return err
	}
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
