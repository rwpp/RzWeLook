//go:build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/rwpp/RzWeLook/user/grpc"
	"github.com/rwpp/RzWeLook/user/ioc"
	"github.com/rwpp/RzWeLook/user/repository"
	"github.com/rwpp/RzWeLook/user/repository/cache"
	"github.com/rwpp/RzWeLook/user/repository/dao"
	"github.com/rwpp/RzWeLook/user/service"
	"github.com/rwpp/RzWeLook/wego"
)

var thirdProvider = wire.NewSet(
	ioc.InitLogger,
	ioc.InitDB,
	ioc.InitRedis,
	ioc.InitEtcdClient,
)

func Init() *wego.App {
	wire.Build(
		thirdProvider,
		cache.NewUserCache,
		dao.NewUserDAO,
		repository.NewUserRepository,
		service.NewUserService,
		grpc.NewUserServiceServer,
		ioc.InitGRPCxServer,
		wire.Struct(new(wego.App), "GRPCServer"),
	)
	return new(wego.App)
}
