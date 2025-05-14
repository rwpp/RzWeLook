//go:build wireinject

package main

import (
	"github.com/google/wire"
	repository2 "github.com/rwpp/RzWeLook/interactive/repository"
	cache2 "github.com/rwpp/RzWeLook/interactive/repository/cache"
	dao2 "github.com/rwpp/RzWeLook/interactive/repository/dao"
	service2 "github.com/rwpp/RzWeLook/interactive/service"
	"github.com/rwpp/RzWeLook/internal/events/article"
	"github.com/rwpp/RzWeLook/internal/ioc"
	"github.com/rwpp/RzWeLook/internal/repository"
	"github.com/rwpp/RzWeLook/internal/repository/cache"
	"github.com/rwpp/RzWeLook/internal/repository/dao"
	"github.com/rwpp/RzWeLook/internal/service"
	"github.com/rwpp/RzWeLook/internal/web"
	ijwt "github.com/rwpp/RzWeLook/internal/web/jwt"
)

var interactiveSvcProvider = wire.NewSet(
	service2.NewInteractiveService,
	repository2.NewCachedInteractiveRepository,
	dao2.NewGORMInteractiveDAO,
	cache2.NewRedisInteractiveCache,
)
var rankingServiceSet = wire.NewSet(
	cache.NewRankingCache,
	cache.NewRankingLocalCache,
	repository.NewRankingRepository,
	service.NewRankingService)

func InitApp() *App {
	wire.Build(
		ioc.InitDB,
		ioc.InitRedis,
		ioc.InitRLockClient,
		ioc.InitLogger,
		ioc.InitKafka,
		ioc.NewConsumers,
		ioc.NewSyncProducer,
		//interactiveSvcProvider,
		//ioc.InitIntrGRPCClient,
		ioc.InitEtcd,
		ioc.InitIntrGRPCClientV1,
		rankingServiceSet,
		ioc.InitJobs,
		ioc.InitRankingJob,

		//events.NewInteractiveReadEventConsumer,
		article.NewKafkaProducer,

		dao.NewUserDAO,
		dao.NewGORMArticleDAO,

		cache.NewUserCache,
		cache.NewRedisArticleCache,
		cache.NewCodeCache,
		repository.NewCodeRepository,
		repository.NewUserRepository,
		repository.NewArticleRepository,
		service.NewUserService,
		service.NewCodeService,
		service.NewArticleService,
		ioc.InitSMSService,
		ioc.NewWechatHandler,
		web.NewUserHandler,
		web.NewArticleHandler,
		web.NewOAuthWechatHandler,
		ijwt.NewRedisJWTHandler,
		ioc.InitOAuthWechatService,
		ioc.InitWeb,
		ioc.InitMiddleware,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
