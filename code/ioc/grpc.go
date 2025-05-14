package ioc

import (
	grpc2 "github.com/rwpp/RzWeLook/code/grpc"
	"github.com/rwpp/RzWeLook/pkg/grpcx"
	"github.com/rwpp/RzWeLook/pkg/logger"
	"github.com/spf13/viper"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
)

func InitGRPCxServer(codeServiceServer *grpc2.CodeServiceServer,
	ecli *clientv3.Client,
	l logger.LoggerV1) *grpcx.Server {
	type Config struct {
		Port     int    `yaml:"port"`
		EtcdAddr string `yaml:"etcdAddr"`
		EtcdTTL  int64  `yaml:"etcdTTL"`
	}
	var cfg Config
	err := viper.UnmarshalKey("grpc.server", &cfg)
	if err != nil {
		panic(err)
	}
	server := grpc.NewServer()
	codeServiceServer.Register(server)
	return &grpcx.Server{
		Server:     server,
		Port:       cfg.Port,
		Name:       "code",
		L:          l,
		EtcdTTL:    cfg.EtcdTTL,
		EtcdClient: ecli,
	}
}
