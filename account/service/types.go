package service

import (
	"context"
	"github.com/rwpp/RzWeLook/account/domain"
)

type AccountService interface {
	Credit(ctx context.Context, cr domain.Credit) error
}
