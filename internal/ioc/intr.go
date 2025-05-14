package ioc

import (
	"github.com/fsnotify/fsnotify"
	intrv1 "github.com/rwpp/RzWeLook/api/proto/gen/intr/v1"
	"github.com/rwpp/RzWeLook/interactive/service"
	"github.com/rwpp/RzWeLook/internal/web/client"
	"github.com/spf13/viper"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitEtcd() *clientv3.Client {
	var cfg clientv3.Config
	err := viper.UnmarshalKey("etcd", &cfg)
	if err != nil {
		panic(err)
	}
	cli, err := clientv3.New(cfg)
	if err != nil {
		panic(err)
	}
	return cli
}

// grpc客户端
func InitIntrGRPCClientV1(client *clientv3.Client) intrv1.InteractiveServiceClient {
	type Config struct {
		Secure bool
		Name   string
	}
	var cfg Config
	err := viper.UnmarshalKey("grpc.client.intr", &cfg)
	if err != nil {
		panic(err)
	}
	bd, err := resolver.NewBuilder(client)
	if err != nil {
		panic(err)
	}
	opts := []grpc.DialOption{grpc.WithResolvers(bd)}
	if cfg.Secure {
		//加载证书启动https
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	cc, err := grpc.Dial("etcd:///service/"+cfg.Name, opts...)
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
