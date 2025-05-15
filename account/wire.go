//go:build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/rwpp/RzWeLook/account/grpc"
	"github.com/rwpp/RzWeLook/account/ioc"
	"github.com/rwpp/RzWeLook/account/repository"
	"github.com/rwpp/RzWeLook/account/repository/dao"
	"github.com/rwpp/RzWeLook/account/service"
	"github.com/rwpp/RzWeLook/wego"
)

func Init() *wego.App {
	wire.Build(
		ioc.InitDB,
		ioc.InitLogger,
		ioc.InitEtcdClient,
		ioc.InitGRPCxServer,
		dao.NewCreditGORMDAO,
		repository.NewAccountRepository,
		service.NewAccountService,
		grpc.NewAccountServiceServer,
		wire.Struct(new(wego.App), "GRPCServer"))
	return new(wego.App)
}
