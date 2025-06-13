// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package startup

import (
	"github.com/google/wire"
	"github.com/rwpp/RzWeLook/interactive/grpc"
	"github.com/rwpp/RzWeLook/interactive/repository"
	"github.com/rwpp/RzWeLook/interactive/repository/cache"
	"github.com/rwpp/RzWeLook/interactive/repository/dao"
	"github.com/rwpp/RzWeLook/interactive/service"
)

// Injectors from wire.go:

func InitInteractiveService() service.InteractiveService {
	gormDB := InitTestDB()
	interactiveDAO := dao.NewGORMInteractiveDAO(gormDB)
	cmdable := InitRedis()
	interactiveCache := cache.NewRedisInteractiveCache(cmdable)
	loggerV1 := InitLog()
	interactiveRepository := repository.NewCachedInteractiveRepository(interactiveDAO, interactiveCache, loggerV1)
	interactiveService := service.NewInteractiveService(interactiveRepository, loggerV1)
	return interactiveService
}

func InitGRPCServer() *grpc.InteractiveServiceServer {
	gormDB := InitTestDB()
	interactiveDAO := dao.NewGORMInteractiveDAO(gormDB)
	cmdable := InitRedis()
	interactiveCache := cache.NewRedisInteractiveCache(cmdable)
	loggerV1 := InitLog()
	interactiveRepository := repository.NewCachedInteractiveRepository(interactiveDAO, interactiveCache, loggerV1)
	interactiveService := service.NewInteractiveService(interactiveRepository, loggerV1)
	interactiveServiceServer := grpc.NewInteractiveServiceServer(interactiveService)
	return interactiveServiceServer
}

// wire.go:

var thirdProvider = wire.NewSet(
	InitRedis, InitTestDB,
	InitLog,
	InitKafka,
)

var interactiveSvcProvider = wire.NewSet(service.NewInteractiveService, repository.NewCachedInteractiveRepository, dao.NewGORMInteractiveDAO, cache.NewRedisInteractiveCache)
