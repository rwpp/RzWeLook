package ioc

import (
	"github.com/fsnotify/fsnotify"
	intrv1 "github.com/rwpp/RzWeLook/api/proto/gen/intr/v1"
	"github.com/rwpp/RzWeLook/bff/web/client"
	"github.com/rwpp/RzWeLook/interactive/service"
	"github.com/spf13/viper"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// grpc客户端
func InitIntrGRPCClientV1(client *clientv3.Client) intrv1.InteractiveServiceClient {
	type Config struct {
		Target string `json:"target"`
		Secure bool   `json:"secure"`
	}
	var cfg Config
	err := viper.UnmarshalKey("grpc.client.intr", &cfg)
	if err != nil {
		panic(err)
	}
	rs, err := resolver.NewBuilder(client)
	if err != nil {
		panic(err)
	}
	opts := []grpc.DialOption{grpc.WithResolvers(rs)}
	if !cfg.Secure {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	cc, err := grpc.Dial(cfg.Target, opts...)
	if err != nil {
		panic(err)
	}
	return intrv1.NewInteractiveServiceClient(cc)

}

// 流量控制客户端
func InitIntrGRPCClient(svc service.InteractiveService) intrv1.InteractiveServiceClient {
	type Config struct {
		Addr      string
		Secure    bool
		Threshold int32
	}
	var cfg Config
	err := viper.UnmarshalKey("grpc.client.intr", &cfg)
	if err != nil {
		panic(err)
	}
	var opts []grpc.DialOption
	if cfg.Secure {

	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	cc, err := grpc.Dial(cfg.Addr, opts...)
	if err != nil {
		panic(err)
	}
	remote := intrv1.NewInteractiveServiceClient(cc)
	local := client.NewInteractiveServiceAdapter(svc)
	res := client.NewGreyScaleInteractiveServiceClient(remote, local)
	viper.OnConfigChange(func(e fsnotify.Event) {
		var cfg Config
		err := viper.UnmarshalKey("grpc.client.intr", &cfg)
		if err != nil {

		}
		res.UpdateThreshold(cfg.Threshold)
	})
	return res
}
