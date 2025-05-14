//go:build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/rwpp/RzWeLook/sms/grpc"
	"github.com/rwpp/RzWeLook/sms/ioc"
	"github.com/rwpp/RzWeLook/wego"
)

func Init() *wego.App {
	wire.Build(
		ioc.InitLogger,
		ioc.InitEtcdClient,
		ioc.InitSmsTencentService,
		grpc.NewSmsServiceServer,
		ioc.InitGRPCxServer,
		wire.Struct(new(wego.App), "GRPCServer"),
	)
	return new(wego.App)
}
