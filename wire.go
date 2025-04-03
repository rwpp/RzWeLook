//go:build wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/rwpp/RzWeLook/internal/ioc"
	"github.com/rwpp/RzWeLook/internal/repository"
	"github.com/rwpp/RzWeLook/internal/repository/cache"
	"github.com/rwpp/RzWeLook/internal/repository/dao"
	"github.com/rwpp/RzWeLook/internal/service"
	"github.com/rwpp/RzWeLook/internal/web"
	ijwt "github.com/rwpp/RzWeLook/internal/web/jwt"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		ioc.InitDB,
		ioc.InitRedis,
		dao.NewUserDAO,
		cache.NewUserCache,
		cache.NewCodeCache,
		repository.NewCodeRepository,
		repository.NewUserRepository,
		service.NewUserService,
		service.NewCodeService,
		ioc.InitSMSService,
		ioc.NewWechatHandler,
		web.NewUserHandler,
		web.NewOAuthWechatHandler,
		ijwt.NewRedisJWTHandler,
		ioc.InitOAuthWechatService,
		ioc.InitWeb,
		ioc.InitMiddleware,
	)
	return new(gin.Engine)
}
