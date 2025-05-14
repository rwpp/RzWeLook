package ioc

import (
	grpc2 "github.com/rwpp/RzWeLook/interactive/grpc"
	"github.com/rwpp/RzWeLook/pkg/grpcx"
	"github.com/rwpp/RzWeLook/pkg/logger"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func InitGRPCxServer(l logger.LoggerV1, intrServer *grpc2.InteractiveServiceServer) *grpcx.Server {
	type Config struct {
		Port      int      `yaml:"port"`
		EtcdAddrs []string `yaml:"etcdAddrs"`
	}
	var cfg Config
	err := viper.UnmarshalKey("grpc.server", &cfg)
	if err != nil {
		panic(err)
	}
	server := grpc.NewServer()
	intrServer.Register(server)
	return &grpcx.Server{
		Server:    server,
		Port:      cfg.Port,
		EtcdAddrs: cfg.EtcdAddrs,
		Name:      "interactive",
		L:         l,
	}
}
