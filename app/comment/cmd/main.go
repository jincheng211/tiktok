package main

import (
	"douyin/app/comment/internal/cache"
	"douyin/app/comment/internal/db"
	"douyin/app/comment/service"
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
	cache.InitRedisDB()
	oss.Init_oss()

	// etcd 地址
	etcdAddress := []string{config.Conf.Etcd.Address}
	// 服务注册
	etcdRegister := discovery.NewRegister(etcdAddress, logrus.New())
	grpcAddress := config.Conf.Services["comment"].Addr[0]
	defer etcdRegister.Stop()

	feedNode := discovery.Server{
		Name: config.Conf.Domain["comment"].Name,
		Addr: grpcAddress,
	}

	server := grpc.NewServer()
	defer server.Stop()

	// 绑定service
	pb.RegisterCommentServiceServer(server, service.GetCommentSrv())
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
