//go:build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/rwpp/RzWeLook/bff/ioc"
	"github.com/rwpp/RzWeLook/bff/web"
	ijwt "github.com/rwpp/RzWeLook/bff/web/jwt"
	"github.com/rwpp/RzWeLook/wego"
)

func InitApp() *wego.App {
	wire.Build(
		ioc.InitLogger,
		ioc.InitRedis,
		ioc.InitEtcdClient,

		web.NewArticleHandler,
		web.NewUserHandler,
		web.NewRewardHandler,

		ijwt.NewRedisJWTHandler,

		ioc.InitUserClient,
		ioc.InitIntrGRPCClientV1,
		ioc.InitRewardClient,
		ioc.InitCodeClient,
		ioc.InitArticleClient,
		ioc.InitGinServer,
		wire.Struct(new(wego.App), "WebServer"),
	)
	return new(wego.App)
}
