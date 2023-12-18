package service

import (
	"context"
	"douyin/app/relation/internal/db"
	"douyin/idl/pb"
	"douyin/pkg/jwt"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"sync"
)

type RelationSrv struct {
	pb.UnimplementedRelationServiceServer
}

var RelationSrvIns *RelationSrv
var RelationSrvOnce sync.Once

// 单例模式
// https://www.python100.com/html/112030.html(单例模式作用)
func GetRelationSrv() *RelationSrv {
	RelationSrvOnce.Do(func() {
		RelationSrvIns = &RelationSrv{}
	})
	return RelationSrvIns
}

func (r *RelationSrv) Relation(ctx context.Context, req *pb.DouyinRelationActionRequest) (*pb.DouyinRelationActionResponse, error) {
	resp := new(pb.DouyinRelationActionResponse)
	resp.StatusCode = 0
	claimn, _ := jwt.ParseToken(req.Token)
	if req.ActionType == 1 {
		err := db.NewRelationDao(ctx).Follow(req)
		if err != nil {
			resp.StatusCode = 500
			resp.StatusMsg = "service关注操作失败!"
			return resp, err
		}

		// 关注数存入redis
		followCount, err := db.NewRelationDao(ctx).GetFollowCount(claimn.UserID)
		db.RDB.HSet(ctx, fmt.Sprintf("user: %d", claimn.UserID), "FollowCount", followCount)
		resp.StatusMsg = "service关注操作成功!"
		return resp, nil
	}

	err := db.NewRelationDao(ctx).UnFollow(req)
	if err != nil {
		resp.StatusCode = 500
		resp.StatusMsg = "service取消关注操作失败!"
		return resp, err
	}

	// 关注数存入redis
	followCount, err := db.NewRelationDao(ctx).GetFollowCount(claimn.UserID)
	db.RDB.HSet(ctx, fmt.Sprintf("user: %d", claimn.UserID), "FollowCount", followCount)
	resp.StatusMsg = "service取消关注操作成功!"
	return resp, nil
}

// 关注列表
func (r *RelationSrv) GetFollowList(ctx context.Context, req *pb.DouyinRelationFollowListRequest) (*pb.DouyinRelationFollowListResponse, error) {
	resp := new(pb.DouyinRelationFollowListResponse)
	resp.StatusCode = 0

	followlist, err := db.NewRelationDao(ctx).GetFollowList(req)

	if err != nil {
		resp.StatusCode = 500
		resp.StatusMsg = "service获取关注列表失败!"
		return resp, err
	}

	conn, err := grpc.Dial("localhost:10002", grpc.WithInsecure())
	if err != nil {
		log.Fatal("Failed to connect: %v", err)
	}
	defer conn.Close()
	userClient := pb.NewUserServiceClient(conn)

	for _, follow := range followlist {
		userReq := pb.DouyinUserRequest{UserId: follow.FollowerID}
		userResp, _ := userClient.UserInfo(ctx, &userReq)
		resp.UserList = append(resp.UserList, userResp.User)
	}

	// 关注数存入redis
	followCount, err := db.NewRelationDao(ctx).GetFollowCount(req.UserId)
	db.RDB.HSet(ctx, fmt.Sprintf("user: %d", req.UserId), "FollowCount", followCount)

	resp.StatusMsg = "service获取关注列表成功!"
	return resp, nil
}

// 粉丝
func (r *RelationSrv) GetFollowerList(ctx context.Context, req *pb.DouyinRelationFollowerListRequest) (*pb.DouyinRelationFollowerListResponse, error) {
	resp := new(pb.DouyinRelationFollowerListResponse)
	resp.StatusCode = 0

	followerlist, err := db.NewRelationDao(ctx).GetFollowerList(req)
	if err != nil {
		resp.StatusCode = 500
		resp.StatusMsg = "service获取关注列表操作失败!"
		return resp, err
	}
	fmt.Println(followerlist)

	conn, err := grpc.Dial("localhost:10002", grpc.WithInsecure())
	if err != nil {
		log.Fatal("Failed to connect: %v", err)
	}
	defer conn.Close()
	userClient := pb.NewUserServiceClient(conn)

	for _, follower := range followerlist {
		userReq := pb.DouyinUserRequest{UserId: follower.FollowID}
		userResp, _ := userClient.UserInfo(ctx, &userReq)
		if db.NewRelationDao(ctx).IsFriend(req.UserId, userResp.User.Id) == true {
			userResp.User.IsFollow = true
		}
		resp.UserList = append(resp.UserList, userResp.User)
	}

	// 粉丝数存入redis
	followerCount, err := db.NewRelationDao(ctx).GetFollowerCount(req.UserId)
	db.RDB.HSet(ctx, fmt.Sprintf("user: %d", req.UserId), "FollowerCount", followerCount)

	resp.StatusMsg = "service操作成功!"
	return resp, nil
}

func (r *RelationSrv) GetFriendList(ctx context.Context, req *pb.DouyinRelationFriendListRequest) (*pb.DouyinRelationFriendListResponse, error) {
	resp := new(pb.DouyinRelationFriendListResponse)
	resp.StatusCode = 0

	friendlist, err := db.NewRelationDao(ctx).GetFriendList(req)
	if err != nil {
		resp.StatusCode = 500
		resp.StatusMsg = "service赞操作失败!"
		return resp, err
	}

	conn, err := grpc.Dial("localhost:10002", grpc.WithInsecure())
	if err != nil {
		log.Fatal("Failed to connect: %v", err)
	}
	defer conn.Close()
	userClient := pb.NewUserServiceClient(conn)

	for _, friend := range friendlist {
		userReq := pb.DouyinUserRequest{UserId: friend.FollowerID}
		userResp, _ := userClient.UserInfo(ctx, &userReq)
		resp.UserList = append(resp.UserList, userResp.User)
	}

	resp.StatusMsg = "service操作成功!"
	return resp, nil
}
