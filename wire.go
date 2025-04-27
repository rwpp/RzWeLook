//go:build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/rwpp/RzWeLook/internal/events/article"
	"github.com/rwpp/RzWeLook/internal/ioc"
	"github.com/rwpp/RzWeLook/internal/repository"
	"github.com/rwpp/RzWeLook/internal/repository/cache"
	"github.com/rwpp/RzWeLook/internal/repository/dao"
	"github.com/rwpp/RzWeLook/internal/service"
	"github.com/rwpp/RzWeLook/internal/web"
	ijwt "github.com/rwpp/RzWeLook/internal/web/jwt"
)

var rankingServiceSet = wire.NewSet(
	repository.NewRankingRepository,
	cache.NewRankingCache,
	service.NewRankingService)

func InitApp() *App {
	wire.Build(
		ioc.InitDB,
		ioc.InitRedis,
		ioc.InitLogger,
		ioc.InitKafka,
		ioc.NewConsumers,
		ioc.NewSyncProducer,
		rankingServiceSet,
		ioc.InitJobs,
		ioc.InitRankingJob,

		article.NewInteractiveReadEventConsumer,
		article.NewKafkaProducer,

		dao.NewUserDAO,
		dao.NewGORMArticleDAO,
		dao.NewGORMInteractiveDAO,
		cache.NewRedisInteractiveCache,
		cache.NewUserCache,
		cache.NewRedisArticleCache,
		cache.NewCodeCache,
		repository.NewCodeRepository,
		repository.NewUserRepository,
		repository.NewArticleRepository,
		repository.NewCachedInteractiveRepository,
		service.NewUserService,
		service.NewCodeService,
		service.NewArticleService,
		service.NewInteractiveService,
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
