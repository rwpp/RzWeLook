package ioc

//import (
//	userv1 "github.com/rwpp/RzWeLook/api/proto/gen/user/v1"
//	"github.com/spf13/viper"
//	"google.golang.org/grpc"
//	"google.golang.org/grpc/credentials/insecure"
//)

//func InitUserRpcClient() userv1.UserServiceClient {
//	type config struct {
//		Addr string `yaml:"addr"`
//	}
//	var cfg config
//	err := viper.UnmarshalKey("grpc.client.user", &cfg)
//	if err != nil {
//		panic(err)
//	}
//	conn, err := grpc.Dial(cfg.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
//	if err != nil {
//		panic(err)
//	}
//	client := userv1.NewUserServiceClient(conn)
//	return client
//}
