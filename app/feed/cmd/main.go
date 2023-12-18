package main

import (
	"douyin/app/feed/internal/db"
	"douyin/app/feed/service"
	"douyin/config"
	"douyin/idl/pb"
	"douyin/pkg/discovery"
	"douyin/pkg/oss"
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
)

func main() {
	config.InitConfig()
	db.InitDB()
	oss.Init_oss()

	// etcd 地址
	etcdAddress := []string{config.Conf.Etcd.Address}
	// 服务注册
	etcdRegister := discovery.NewRegister(etcdAddress, logrus.New())
	grpcAddress := config.Conf.Services["feed"].Addr[0]
	defer etcdRegister.Stop()
	feedNode := discovery.Server{
		Name: config.Conf.Domain["feed"].Name,
		Addr: grpcAddress,
	}
	server := grpc.NewServer()
	defer server.Stop()

	// 绑定service
	pb.RegisterFeedServiceServer(server, service.GetFeedSrv())
	lis, err := net.Listen("tcp", grpcAddress)
	if err != nil {
		panic(err)
	}
	if _, err := etcdRegister.Register(feedNode, 10); err != nil {
		panic(fmt.Sprintf("start server failed, err: %v", err))
	}
	logrus.Info("server started listen on ", grpcAddress)
	if err := server.Serve(lis); err != nil {
		panic(err)
	}
}
