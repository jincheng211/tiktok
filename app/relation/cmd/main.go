package main

import (
	"douyin/app/relation/internal/db"
	"douyin/app/relation/service"
	"douyin/config"
	"douyin/idl/pb"
	"douyin/pkg/discovery"
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
)

func main() {
	config.InitConfig()
	db.InitDB()
	// etcd 地址
	etcdAddress := []string{config.Conf.Etcd.Address}
	// 服务注册
	etcdRegister := discovery.NewRegister(etcdAddress, logrus.New())
	grpcAddress := config.Conf.Services["relation"].Addr[0]
	defer etcdRegister.Stop()
	userNode := discovery.Server{
		Name: config.Conf.Domain["relation"].Name,
		Addr: grpcAddress,
	}
	server := grpc.NewServer()
	defer server.Stop()

	// 绑定service
	pb.RegisterRelationServiceServer(server, service.GetRelationSrv())
	lis, err := net.Listen("tcp", grpcAddress)
	if err != nil {
		panic(err)
	}
	if _, err := etcdRegister.Register(userNode, 10); err != nil {
		panic(fmt.Sprintf("start server failed, err: %v", err))
	}
	logrus.Info("server started listen on ", grpcAddress)
	if err := server.Serve(lis); err != nil {
		panic(err)
	}
}
