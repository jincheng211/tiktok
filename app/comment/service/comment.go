package service

import (
	"context"
	"douyin/app/comment/internal/db"
	pb "douyin/idl/pb"
	"douyin/pkg/jwt"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"sync"
	"time"
)

type CommentSrv struct {
	pb.UnimplementedCommentServiceServer
}

var CommentSrvIns *CommentSrv
var CommentSrvOnce sync.Once

func GetCommentSrv() *CommentSrv {
	CommentSrvOnce.Do(func() {
		CommentSrvIns = &CommentSrv{}
	})
	return CommentSrvIns
}

func (c *CommentSrv) Comment(ctx context.Context, req *pb.DouyinCommentActionRequest) (*pb.DouyinCommentActionResponse, error) {
	resp := new(pb.DouyinCommentActionResponse)
	resp.StatusCode = 0
	// 链接微服务
	conn, err := grpc.Dial("localhost:10002", grpc.WithInsecure())
	if err != nil {
		log.Fatal("Failed to connect: %v", err)
	}
	defer conn.Close()

	// 查询信息
	claimn, _ := jwt.ParseToken(req.Token)
	userClient := pb.NewUserServiceClient(conn)
	uerReq := pb.DouyinUserRequest{UserId: claimn.UserID}
	userResp, _ := userClient.UserInfo(ctx, &uerReq)

	// 写评论
	if req.ActionType == 1 {
		id, err := db.NewCommentDao(ctx).Comment(req)

		if err != nil {
			resp.StatusCode = 500
			resp.StatusMsg = "sevice评论失败"
			return resp, err
		}
		resp.Comment = &pb.Comment{
			Id:         id,
			User:       userResp.User,
			Content:    req.CommentText,
			CreateDate: time.Now().Format("2006-01-02 15:04:05"),
		}
		resp.StatusMsg = "sevice comment评论成功"

	} else {
		// 删除评论
		err := db.NewCommentDao(ctx).DeleteComment(req)
		if err != nil {
			return resp, err
		}
		resp.Comment = &pb.Comment{
			Id:         req.CommentId,
			User:       userResp.User,
			Content:    "",
			CreateDate: time.Now().Format("2006-01-02 15:04:05"),
		}
		// 使用redis增加评论
		resp.StatusMsg = "sevice comment评论删除"
	}

	// 使用redis存储视频评论数
	CommentCount, _ := db.NewCommentDao(ctx).CommentCount(req.VideoId)
	db.RDB.HSet(ctx, fmt.Sprintf("videoID: %d", req.VideoId), "CommentCount", CommentCount)
	return resp, nil
}

func (c *CommentSrv) CommentList(ctx context.Context, req *pb.DouyinCommentListRequest) (*pb.DouyinCommentListResponse, error) {
	resp := new(pb.DouyinCommentListResponse)
	resp.StatusCode = 0

	CommentList, err := db.NewCommentDao(ctx).CommentList(req)
	if err != nil {
		resp.StatusCode = 500
		resp.StatusMsg = "sevice commentList获取失败"
		return resp, err
	}

	resp.StatusMsg = "sevice commentList获取成功"
	for _, c := range CommentList {
		conn, err := grpc.Dial("localhost:10002", grpc.WithInsecure())
		if err != nil {
			log.Fatal("Failed to connect: %v", err)
		}
		defer conn.Close()

		userClient := pb.NewUserServiceClient(conn)
		uerReq := pb.DouyinUserRequest{UserId: c.UserID}
		userResp, _ := userClient.UserInfo(ctx, &uerReq)

		comment := &pb.Comment{
			Id:         c.ID,
			User:       userResp.User,
			Content:    c.Content,
			CreateDate: c.CreateDate,
		}

		resp.CommentList = append(resp.CommentList, comment)
	}

	return resp, nil
}
