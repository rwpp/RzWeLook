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

var serviceProviderSet = wire.NewSet(
	dao.NewGORMInteractiveDAO,
	cache.NewRedisInteractiveCache,
	repository.NewCachedInteractiveRepository,
	service.NewInteractiveService)

var thirdProvider = wire.NewSet(
	ioc.InitSrc,
	ioc.InitDst,
	ioc.InitDoubleWritePool,
	ioc.InitBizDB,
	ioc.InitRedis,
	ioc.InitLogger,
	ioc.InitKafka,
	ioc.InitEtcdClient,
	ioc.InitSyncProducer,
)

var migratorSet = wire.NewSet(
	ioc.InitMigratorServer,
	ioc.InitFixDataConsumer,
	ioc.InitMigradatorProducer)

func Init() *App {
	wire.Build(
		thirdProvider,
		serviceProviderSet,
		migratorSet,
		grpc.NewInteractiveServiceServer,
		events.NewInteractiveReadEventConsumer,
		ioc.InitGRPCxServer,
		ioc.NewConsumers,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
