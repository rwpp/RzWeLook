//go:build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/rwpp/RzWeLook/reward/grpc"
	"github.com/rwpp/RzWeLook/reward/ioc"
	"github.com/rwpp/RzWeLook/reward/repository"
	"github.com/rwpp/RzWeLook/reward/repository/cache"
	"github.com/rwpp/RzWeLook/reward/repository/dao"
	"github.com/rwpp/RzWeLook/reward/service"
	"github.com/rwpp/RzWeLook/wego"
)

var thirdPartySet = wire.NewSet(
	ioc.InitDB,
	ioc.InitLogger,
	ioc.InitEtcdClient,
	ioc.InitRedis)

func Init() *wego.App {
	wire.Build(thirdPartySet,
		service.NewWechatNativeRewardService,
		ioc.InitAccountClient,
		ioc.InitGRPCxServer,
		ioc.InitPaymentClient,
		repository.NewRewardRepository,
		cache.NewRewardRedisCache,
		dao.NewRewardGORMDAO,
		grpc.NewRewardServiceServer,
		wire.Struct(new(wego.App), "GRPCServer"),
	)
	return new(wego.App)
}
