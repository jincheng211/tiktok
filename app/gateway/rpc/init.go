package rpc

import (
	"context"
	"douyin/config"
	"douyin/idl/pb"
	"douyin/pkg/discovery"
	"fmt"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
	"log"
	"time"
)

var (
	Register   *discovery.Resolver
	ctx        context.Context
	CancelFunc context.CancelFunc

	UserClient     pb.UserServiceClient
	FeedClient     pb.FeedServiceClient
	CommentClient  pb.CommentServiceClient
	MessageClient  pb.MessageServiceClient
	RelationClient pb.RelationServiceClient
)

func Init() {
	Register = discovery.NewResolver([]string{config.Conf.Etcd.Address}, logrus.New())
	resolver.Register(Register)
	ctx, CancelFunc = context.WithTimeout(context.Background(), 3*time.Second)

	defer Register.Close()

	// 每次写完微服务在这初始化
	initClient(config.Conf.Domain["user"].Name, &UserClient)
	initClient(config.Conf.Domain["feed"].Name, &FeedClient)
	initClient(config.Conf.Domain["comment"].Name, &CommentClient)
	initClient(config.Conf.Domain["message"].Name, &MessageClient)
	initClient(config.Conf.Domain["relation"].Name, &RelationClient)

}

func initClient(serviceName string, client interface{}) {
	conn, err := connectServer(serviceName)

	if err != nil {
		panic(err)
	}

	// 这里也要改
	switch c := client.(type) {
	case *pb.UserServiceClient:
		*c = pb.NewUserServiceClient(conn)
	case *pb.FeedServiceClient:
		*c = pb.NewFeedServiceClient(conn)
	case *pb.CommentServiceClient:
		*c = pb.NewCommentServiceClient(conn)
	case *pb.MessageServiceClient:
		*c = pb.NewMessageServiceClient(conn)
	case *pb.RelationServiceClient:
		*c = pb.NewRelationServiceClient(conn)

	default:
		panic("unsupported client type")
	}
}

func connectServer(serviceName string) (conn *grpc.ClientConn, err error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	addr := fmt.Sprintf("%s:///%s", Register.Scheme(), serviceName)

	if config.Conf.Services[serviceName].LoadBalance {
		log.Printf("load balance enabled for %s\n", serviceName)
		opts = append(opts, grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, "round_robin")))
	}

	conn, err = grpc.DialContext(ctx, addr, opts...)
	return
}
