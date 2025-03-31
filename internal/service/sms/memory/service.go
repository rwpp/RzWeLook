package memory

import (
	"context"
	"fmt"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}
func (s *Service) Send(ctx context.Context, tpl string, args []string, number ...string) error {
	// 模拟发送短信
	// 实际上什么都不做
	fmt.Println(args)
	return nil
}
