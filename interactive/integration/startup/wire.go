//go:build wireinject

package startup

import (
	"github.com/google/wire"
	"github.com/rwpp/RzWeLook/interactive/grpc"
	"github.com/rwpp/RzWeLook/interactive/repository"
	"github.com/rwpp/RzWeLook/interactive/repository/cache"
	"github.com/rwpp/RzWeLook/interactive/repository/dao"
	"github.com/rwpp/RzWeLook/interactive/service"
)

var thirdProvider = wire.NewSet(
	InitRedis, InitTestDB,
	InitLog,
	InitKafka,
)

//var userSvcProvider = wire.NewSet(
//	dao.NewUserDAO,
//	cache.NewUserCache, cache.NewCodeCache,
//	repository.NewUserRepository,
//	service.NewUserService)

var interactiveSvcProvider = wire.NewSet(
	service.NewInteractiveService,
	repository.NewCachedInteractiveRepository,
	dao.NewGORMInteractiveDAO,
	cache.NewRedisInteractiveCache,
)

//func InitWebServer() *gin.Engine {
//	wire.Build(
//		thirdProvider,
//		userSvcProvider,
//		cache.NewCodeCache)
//}

func InitInteractiveService() service.InteractiveService {
	wire.Build(
		thirdProvider,
		interactiveSvcProvider,
	)
	return service.NewInteractiveService(nil, nil)
}

func InitGRPCServer() *grpc.InteractiveServiceServer {
	wire.Build(
		grpc.NewInteractiveServiceServer,
		thirdProvider,
		interactiveSvcProvider,
	)
	return new(grpc.InteractiveServiceServer)
}
