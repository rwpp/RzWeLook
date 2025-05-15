//go:build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/rwpp/RzWeLook/cronjob/grpc"
	"github.com/rwpp/RzWeLook/cronjob/ioc"
	"github.com/rwpp/RzWeLook/cronjob/repository"
	"github.com/rwpp/RzWeLook/cronjob/repository/dao"
	"github.com/rwpp/RzWeLook/cronjob/service"
)

var serviceProviderSet = wire.NewSet(
	dao.NewGORMJobDAO,
	repository.NewPreemptCronJobRepository,
	service.NewCronJobService)

var thirdProvider = wire.NewSet(
	ioc.InitDB,
	ioc.InitEtcdClient,
	ioc.InitLogger,
)

func Init() *App {
	wire.Build(
		thirdProvider,
		serviceProviderSet,
		grpc.NewCronJobServiceServer,
		ioc.InitGRPCxServer,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
