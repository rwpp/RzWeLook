//go:build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/rwpp/RzWeLook/oauth2/grpc"
	ioc2 "github.com/rwpp/RzWeLook/oauth2/ioc"
)

var thirdProvider = wire.NewSet(
	ioc2.InitLogger,
)

func Init() *App {
	wire.Build(
		thirdProvider,
		ioc2.InitPrometheus,
		grpc.NewOauth2ServiceServer,
		ioc2.InitGRPCxServer,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
