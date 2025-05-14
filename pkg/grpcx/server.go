package grpcx

import (
	"context"
	"github.com/rwpp/RzWeLook/pkg/logger"
	"github.com/rwpp/RzWeLook/pkg/netx"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
	"google.golang.org/grpc"
	"net"
	"strconv"
	"time"
)

type Server struct {
	*grpc.Server
	Port int
	// ETCD 服务注册租约 TTL
	EtcdTTL     int64
	EtcdClient  *clientv3.Client
	etcdManager endpoints.Manager
	etcdKey     string
	cancel      func()
	Name        string
	L           logger.LoggerV1
}

// Serve 启动服务器并且阻塞
func (s *Server) Serve() error {
	// 初始化一个控制整个过程的 ctx
	// 你也可以考虑让外面传进来，这样的话就是 main 函数自己去控制了
	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel
	port := strconv.Itoa(s.Port)
	l, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}
	// 要先确保启动成功，再注册服务
	err = s.register(ctx, port)
	if err != nil {
		return err
	}
	return s.Server.Serve(l)
}

func (s *Server) register(ctx context.Context, port string) error {
	cli := s.EtcdClient
	serviceName := "service/" + s.Name
	em, err := endpoints.NewManager(cli,
		serviceName)
	if err != nil {
		return err
	}
	s.etcdManager = em
	ip := netx.GetOutboundIP()
	s.etcdKey = serviceName + "/" + ip
	addr := ip + ":" + port
	leaseResp, err := cli.Grant(ctx, s.EtcdTTL)
	// 开启续约
	ch, err := cli.KeepAlive(ctx, leaseResp.ID)
	if err != nil {
		return err
	}
	go func() {
		// 可以预期，当我们的 cancel 被调用的时候，就会退出这个循环
		for chResp := range ch {
			s.L.Debug("续约：", logger.String("resp", chResp.String()))
		}
	}()
	// metadata 我们这里没啥要提供的
	return em.AddEndpoint(ctx, s.etcdKey,
		endpoints.Endpoint{Addr: addr}, clientv3.WithLease(leaseResp.ID))
}

func (s *Server) Close() error {
	s.cancel()
	if s.etcdManager != nil {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		err := s.etcdManager.DeleteEndpoint(ctx, s.etcdKey)
		if err != nil {
			return err
		}
	}
	err := s.EtcdClient.Close()
	if err != nil {
		return err
	}
	s.Server.GracefulStop()
	return nil
}

//type Server struct {
//	*grpc.Server
//	Addr      string
//	Port      int
//	EtcdAddrs []string
//	Name      string
//	L         logger.LoggerV1
//	kaCancel  func()
//	em        endpoints.Manager
//	client    *etcdv3.Client
//	key       string
//}

//func (s *Server) Serve() error {
//	l, err := net.Listen("tcp", ":"+strconv.Itoa(s.Port))
//	if err != nil {
//		return err
//	}
//	err = s.register()
//	if err != nil {
//		return err
//	}
//	// 就是直接启动，我现在要嵌入服务注册过程
//	// 这边会阻塞，类似与 gin.Run
//	return s.Server.Serve(l)
//
//}

//func (s *Server) register() error {
//	//初始化客户端
//	client, err := etcdv3.New(etcdv3.Config{
//		Endpoints: s.EtcdAddrs,
//	})
//	if err != nil {
//		return err
//	}
//	s.client = client
//	// endpoint 以服务为维度。一个服务一个 Manager
//	em, err := endpoints.NewManager(client, "service/"+s.Name)
//	if err != nil {
//		return err
//	}
//	//获取自己的ip
//	addr := netx.GetOutboundIP() + ":" + strconv.Itoa(s.Port)
//	key := "service/" + s.Name + "/" + addr
//	s.key = key
//	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
//	// 你可以做成配置的
//	var ttl int64 = 30
//	leaseResp, err := client.Grant(ctx, ttl)
//	//创建续约
//	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
//	defer cancel()
//	//注册服务
//	err = em.AddEndpoint(ctx, key, endpoints.Endpoint{
//		Addr: addr,
//	}, etcdv3.WithLease(leaseResp.ID))
//
//	kaCtx, kaCancel := context.WithCancel(context.Background())
//	s.kaCancel = kaCancel
//	//启动租约
//	ch, err := client.KeepAlive(kaCtx, leaseResp.ID)
//	if err != nil {
//		return err
//	}
//	go func() {
//		for kaResp := range ch {
//			// 正常就是打印一下 DEBUG 日志啥的
//			s.L.Debug(kaResp.String())
//		}
//	}()
//	return nil
//}
//
//func (s *Server) Close() error {
//	if s.kaCancel != nil {
//		s.kaCancel()
//	}
//	if s.em != nil {
//		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
//		defer cancel()
//		err := s.em.DeleteEndpoint(ctx, s.key)
//		if err != nil {
//			return err
//		}
//	}
//	if s.client != nil {
//		err := s.client.Close()
//		if err != nil {
//			return err
//		}
//	}
//	s.Server.GracefulStop()
//	return nil
//}
