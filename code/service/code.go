package service

import (
	"context"
	"fmt"
	smsv1 "github.com/rwpp/RzWeLook/api/proto/gen/sms/v1"
	"github.com/rwpp/RzWeLook/code/repository"
	"math/rand"
)

const codeTpId = "1877556"

var (
	ErrCodeVerifyTooMany = repository.ErrCodeVerifyTooMany
	ErrCodeSendTooMany   = repository.ErrCodeSendTooMany
)

type CodeServiceInterface interface {
	Send(ctx context.Context, biz string, phone string) error
	Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error)
}

func NewCodeService(repo repository.CodeRepositoryInterface, smsSvc smsv1.SmsServiceClient) CodeServiceInterface {
	return &CodeService{
		repo:   repo,
		smsSvc: smsSvc,
	}
}

type CodeService struct {
	repo   repository.CodeRepositoryInterface
	smsSvc smsv1.SmsServiceClient
}

func (svc *CodeService) Send(ctx context.Context, biz string, phone string) error {
	code := svc.generateCode()
	err := svc.repo.Store(ctx, biz, phone, code)
	if err != nil {
		return err
	}
	_, err = svc.smsSvc.Send(ctx, &smsv1.SmsSendRequest{
		TplId:   codeTpId,
		Args:    []string{code},
		Numbers: []string{phone},
	})
	return err
}

func (svc *CodeService) Verify(ctx context.Context,
	biz string, phone string, inputCode string) (bool, error) {
	return svc.repo.Verify(ctx, biz, phone, inputCode)
}

func (svc *CodeService) generateCode() string {
	num := rand.Intn(1000000)
	return fmt.Sprintf("%06d", num)
}
