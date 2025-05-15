//go:build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/rwpp/RzWeLook/ranking/grpc"
	"github.com/rwpp/RzWeLook/ranking/ioc"
	"github.com/rwpp/RzWeLook/ranking/repository"
	"github.com/rwpp/RzWeLook/ranking/repository/cache"
	"github.com/rwpp/RzWeLook/ranking/service"
)

var serviceProviderSet = wire.NewSet(
	cache.NewRankingLocalCache,
	cache.NewRankingCache,
	repository.NewRankingRepository,
	service.NewBatchRankingService,
)

var thirdProvider = wire.NewSet(
	ioc.InitRedis,
	ioc.InitInterActiveRpcClient,
	ioc.InitArticleRpcClient,
	ioc.InitEtcdClient,
	ioc.InitLogger,
)

func Init() *App {
	wire.Build(
		thirdProvider,
		serviceProviderSet,
		grpc.NewRankingServiceServer,
		ioc.InitGRPCxServer,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
