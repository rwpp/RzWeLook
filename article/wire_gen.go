// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/google/wire"
	"github.com/rwpp/RzWeLook/article/events"
	"github.com/rwpp/RzWeLook/article/grpc"
	"github.com/rwpp/RzWeLook/article/ioc"
	"github.com/rwpp/RzWeLook/article/repository"
	"github.com/rwpp/RzWeLook/article/repository/cache"
	"github.com/rwpp/RzWeLook/article/repository/dao"
	"github.com/rwpp/RzWeLook/article/service"
	"github.com/rwpp/RzWeLook/wego"
)

// Injectors from wire.go:

func Init() *wego.App {
	loggerV1 := ioc.InitLogger()
	db := ioc.InitDB(loggerV1)
	articleDAO := dao.NewGORMArticleDAO(db)
	cmdable := ioc.InitRedis()
	articleCache := cache.NewRedisArticleCache(cmdable)
	articleRepository := repository.NewArticleRepository(articleDAO, articleCache, loggerV1)
	syncProducer := ioc.InitProducer()
	producer := events.NewSaramaSyncProducer(syncProducer)
	articleService := service.NewArticleService(articleRepository, loggerV1, producer)
	articleServiceServer := grpc.NewArticleServiceServer(articleService)
	client := ioc.InitEtcdClient()
	server := ioc.InitGRPCxServer(articleServiceServer, client, loggerV1)
	app := &wego.App{
		GRPCServer: server,
	}
	return app
}

// wire.go:

var thirdProvider = wire.NewSet(ioc.InitRedis, ioc.InitLogger, ioc.InitProducer, ioc.InitEtcdClient, ioc.InitDB)
