//go:build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/rwpp/RzWeLook/code/grpc"
	"github.com/rwpp/RzWeLook/code/ioc"
	"github.com/rwpp/RzWeLook/code/repository"
	"github.com/rwpp/RzWeLook/code/repository/cache"
	"github.com/rwpp/RzWeLook/code/service"
)

var thirdProvider = wire.NewSet(
	ioc.InitRedis,
	ioc.InitEtcdClient,
	ioc.InitLogger,
)

func Init() *App {
	wire.Build(
		thirdProvider,
		ioc.InitSmsRpcClient,
		cache.NewCodeCache,
		repository.NewCodeRepository,
		service.NewCodeService,
		grpc.NewCodeServiceServer,
		ioc.InitGRPCxServer,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
