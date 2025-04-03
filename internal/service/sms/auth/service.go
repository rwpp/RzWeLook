package auth

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rwpp/RzWeLook/internal/service/sms"
)

type SMSService struct {
	svc sms.Service
	key string
}

func (s *SMSService) Send(ctx context.Context, biz string,
	args []string, numbers ...string) error {
	var tc Claims
	_, err := jwt.ParseWithClaims(biz, &tc, func(token *jwt.Token) (interface{}, error) {
		return s.key, nil
	})
	if err != nil {
		return err
	}
	return s.svc.Send(ctx, tc.Tpl, args, numbers...)
}

type Claims struct {
	jwt.RegisteredClaims
	Tpl string `json:"tpl"`
}
