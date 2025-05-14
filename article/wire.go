package article

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

var thirdProvider = wire.NewSet(
	ioc.InitRedis,
	ioc.InitLogger,
	ioc.InitProducer,
	ioc.InitEtcdClient,
	ioc.InitDB,
)

func Init() *wego.App {
	wire.Build(
		thirdProvider,
		events.NewSaramaSyncProducer,
		cache.NewRedisArticleCache,
		dao.NewGORMArticleDAO,
		repository.NewArticleRepository,
		service.NewArticleService,
		grpc.NewArticleServiceServer,
		ioc.InitGRPCxServer,
		wire.Struct(new(wego.App), "GRPCServer"),
	)
	return new(wego.App)
}
