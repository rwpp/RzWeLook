//go:build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/rwpp/RzWeLook/interactive/events"
	"github.com/rwpp/RzWeLook/interactive/grpc"
	"github.com/rwpp/RzWeLook/interactive/ioc"
	"github.com/rwpp/RzWeLook/interactive/repository"
	"github.com/rwpp/RzWeLook/interactive/repository/cache"
	"github.com/rwpp/RzWeLook/interactive/repository/dao"
	"github.com/rwpp/RzWeLook/interactive/service"
)

var thirdPartySet = wire.NewSet(
	ioc.InitDB,
	ioc.InitRedis,
	ioc.InitKafka,
	ioc.InitLogger,
)
var interactiveSvcProvider = wire.NewSet(
	service.NewInteractiveService,
	repository.NewCachedInteractiveRepository,
	dao.NewGORMInteractiveDAO,
	cache.NewRedisInteractiveCache,
)

func InitApp() *App {
	wire.Build(interactiveSvcProvider,
		thirdPartySet,
		events.NewInteractiveReadEventConsumer,
		grpc.NewInteractiveServiceServer,
		ioc.NewConsumers,
		ioc.InitGRPCxServer,
		wire.Struct(new(App), "*"))
	return new(App)
}
