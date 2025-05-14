package ioc

import (
	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rwpp/RzWeLook/interactive/repository/dao"
	"github.com/rwpp/RzWeLook/pkg/ginx"
	"github.com/rwpp/RzWeLook/pkg/gormx/connpool"
	"github.com/rwpp/RzWeLook/pkg/logger"
	events2 "github.com/rwpp/RzWeLook/pkg/migrator/events"
	"github.com/rwpp/RzWeLook/pkg/migrator/events/fixer"
	"github.com/rwpp/RzWeLook/pkg/migrator/scheduler"
	"github.com/spf13/viper"
)

const topic = "migrator_interactives"

func InitFixDataConsumer(l logger.LoggerV1,
	src SrcDB,
	dst DstDB,
	client sarama.Client) *fixer.Consumer[dao.Interactive] {
	res, err := fixer.NewConsumer[dao.Interactive](client, l,
		topic, src, dst)
	if err != nil {
		panic(err)
	}
	return res
}
func InitMigradatorProducer(p sarama.SyncProducer) events2.Producer {
	return events2.NewSaramaProducer(p, topic)
}
func InitMigratorServer(l logger.LoggerV1, src SrcDB, dst DstDB, pool *connpool.DoubleWritePool, producer events2.Producer) *ginx.Server {
	intrSch := scheduler.NewScheduler[dao.Interactive](l, src, dst, pool, producer)
	engine := gin.Default()
	ginx.InitCounter(prometheus.CounterOpts{
		Namespace: "RzWeLook",
		Subsystem: "webook_intr",
		Name:      "http_biz_code",
		Help:      "HTTP 请求",
	})
	intrSch.RegisterRoutes(engine.Group("/migrator"))
	//todo读不到addr
	addr := viper.GetString("migrator.web.addr")
	return &ginx.Server{
		Engine: engine,
		Addr:   addr,
	}
}
