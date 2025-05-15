package repository

import (
	"context"
	"github.com/rwpp/RzWeLook/account/domain"
)

type AccountRepository interface {
	AddCredit(ctx context.Context, c domain.Credit) error
}
