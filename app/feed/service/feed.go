package service

import (
	"context"
	"douyin/app/feed/internal/db"
	pb "douyin/idl/pb"
	"douyin/pkg/jwt"
	"douyin/pkg/oss"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"strconv"
	"sync"
	"time"
)

type FeedSrv struct {
	pb.UnimplementedFeedServiceServer
}

var FeedSrvIns *FeedSrv
var FeedSrvOnce sync.Once

func GetFeedSrv() *FeedSrv {
	FeedSrvOnce.Do(func() {
		FeedSrvIns = &FeedSrv{}
	})
	return FeedSrvIns
}

func (f *FeedSrv) GetFeedList(ctx context.Context, req *pb.DouyinFeedRequest) (resp *pb.DouyinFeedResponse, err error) {

	conn, err := grpc.Dial("localhost:10002", grpc.WithInsecure())
	if err != nil {
		log.Fatal("Failed to connect: %v", err)
	}
	defer conn.Close()

	userClient := pb.NewUserServiceClient(conn)

	resp = new(pb.DouyinFeedResponse)
	resp.StatusCode = 0 // 状态码，0-成功，其他值-失败
	videoList, err := db.NewFeedDao(ctx).GetFeedList(req)
	if err != nil {
		resp.StatusCode = 500
		resp.StatusMsg = "获取视频流失败"
		return resp, err
	}
	// 假设用户没的路的userid
	var UserID int64
	if req.Token != "" {
		claimn, _ := jwt.ParseToken(req.Token)
		UserID = claimn.UserID
	}
	createTime := time.Now().Unix()
	for _, v := range videoList {
		userReq := pb.DouyinUserRequest{UserId: v.AuthorId}
		userResp, _ := userClient.UserInfo(ctx, &userReq)

		if createTime > v.CreatedAt.Unix() {
			createTime = v.CreatedAt.Unix()
		}

		// 从redis获取最新的评论数
		Count, _ := db.RDB.HGet(ctx, fmt.Sprintf("videoID: %d", v.ID), "CommentCount").Result()
		CommentCount, _ := strconv.ParseInt(Count, 10, 64)
		// 更新评论数
		db.NewFeedDao(ctx).UpdateCommentCount(v.ID, CommentCount)
		//从redis获取最新的点赞数
		Count, _ = db.RDB.HGet(ctx, fmt.Sprintf("videoID: %d", v.ID), "FavoriteCount").Result()
		FavoriteCount, _ := strconv.ParseInt(Count, 10, 64)
		// 更新喜欢人数
		db.NewFeedDao(ctx).UpdateFavoriteCount(v.ID, FavoriteCount)

		video := &pb.Video{
			Id:            v.ID,
			Author:        userResp.User,
			PlayUrl:       oss.GetVideo(v.PlayURL),
			CoverUrl:      oss.GetVideo(v.CoverURL),
			FavoriteCount: FavoriteCount,
			CommentCount:  CommentCount,
			IsFavorite:    db.NewFeedDao(ctx).IsFavorite(UserID, v.ID),
			Title:         v.Title,
		}
		resp.VideoList = append(resp.VideoList, video)
	}

	resp.NextTime = &createTime
	resp.StatusMsg = "获取视频流成功"
	return resp, nil
}

func (f *FeedSrv) GetPublishList(ctx context.Context, req *pb.DouyinPublishListRequest) (resp *pb.DouyinPublishListResponse, err error) {
	conn, err := grpc.Dial("localhost:10002", grpc.WithInsecure())
	if err != nil {
		log.Fatal("Failed to connect: %v", err)
	}
	defer conn.Close()

	userClient := pb.NewUserServiceClient(conn)

	resp = new(pb.DouyinPublishListResponse)
	resp.StatusCode = 0 // 状态码，0-成功，其他值-失败
	videoList, err := db.NewFeedDao(ctx).GetPublishList(req)
	if err != nil {
		resp.StatusCode = 500
		resp.StatusMsg = "获取视频流失败"
		return resp, err
	}

	userReq := pb.DouyinUserRequest{UserId: req.UserId}
	userResp, _ := userClient.UserInfo(ctx, &userReq)

	// 获取总的收到的赞数
	var TotalFavorited int64

	for _, v := range videoList {

		// 从redis获取最新的评论数
		Count, _ := db.RDB.HGet(ctx, fmt.Sprintf("videoID: %d", v.ID), "CommentCount").Result()
		CommentCount, _ := strconv.ParseInt(Count, 10, 64)
		// 更新评论数
		db.NewFeedDao(ctx).UpdateCommentCount(v.ID, CommentCount)
		//从redis获取最新的点赞数
		Count, _ = db.RDB.HGet(ctx, fmt.Sprintf("videoID: %d", v.ID), "FavoriteCount").Result()
		FavoriteCount, _ := strconv.ParseInt(Count, 10, 64)
		// 更新喜欢数
		db.NewFeedDao(ctx).UpdateFavoriteCount(v.ID, FavoriteCount)
		TotalFavorited += FavoriteCount

		video := &pb.Video{
			Id:            v.ID,
			Author:        userResp.User,
			PlayUrl:       oss.GetVideo(v.PlayURL),
			CoverUrl:      oss.GetVideo(v.CoverURL),
			FavoriteCount: FavoriteCount,
			CommentCount:  CommentCount,
			IsFavorite:    db.NewFeedDao(ctx).IsFavorite(req.UserId, v.ID),
			Title:         v.Title,
		}

		resp.VideoList = append(resp.VideoList, video)
	}

	// 将总喜欢数更新到Redis
	db.RDB.HSet(ctx, fmt.Sprintf("user: %d", req.UserId), "TotalFavorited", TotalFavorited)
	resp.StatusMsg = "获取视频流成功"
	return resp, nil
}

func (f *FeedSrv) PublishVideo(ctx context.Context, req *pb.DouyinPublishActionRequest) (resp *pb.DouyinPublishActionResponse, err error) {
	resp = new(pb.DouyinPublishActionResponse)
	resp.StatusCode = 0 // 状态码，0-成功，其他值-失败;

	// 获取用户信息
	claimn, err := jwt.ParseToken(req.Token)
	if err != nil {
		resp.StatusCode = 500
		resp.StatusMsg = "用户信息获取失败"
		fmt.Println("(publish) token解析err:", err)
		return resp, err
	}

	_, err = db.RDB.HGetAll(ctx, fmt.Sprintf("user:%d", claimn.UserID)).Result()
	if err != nil {
		resp.StatusCode = 500
		resp.StatusMsg = "用户信息获取失败"
		fmt.Println("(publish) redis get userinfo err:", err)
		return resp, err
	}

	err = db.NewFeedDao(ctx).PublishVideo(req)
	if err != nil {
		resp.StatusCode = 500
		resp.StatusMsg = "上传视频失败"
		fmt.Println("(publish) 上传video err:", err)
		return resp, err
	}

	videoCount, err := db.NewFeedDao(ctx).GetVideoCount(req)
	if err != nil {
		resp.StatusCode = 500
		resp.StatusMsg = "service 获取视频个数失败"
		fmt.Println("(publish) 上传video err:", err)
		return resp, err
	}
	// videoCount结果存到redis
	db.RDB.HSet(ctx, fmt.Sprintf("user: %d", claimn.UserID), "WorkCount", videoCount)

	resp.StatusMsg = "上传视频成功"
	return resp, nil
}

func (f *FeedSrv) Favorite(ctx context.Context, req *pb.DouyinFavoriteActionRequest) (resp *pb.DouyinFavoriteActionResponse, err error) {
	resp = new(pb.DouyinFavoriteActionResponse)
	resp.StatusCode = 0

	FavoriteCount, err := db.NewFeedDao(ctx).Favorite(req)
	// redis存某条视频的总点赞数
	if FavoriteCount > -1 {
		db.RDB.HSet(ctx, fmt.Sprintf("videoID: %d", req.VideoId), "FavoriteCount", FavoriteCount)
	}
	FavoriteVideoCount, err := db.NewFeedDao(ctx).FavoriteVideoCount(req)
	// 用户的点赞数
	if FavoriteCount > -1 {
		db.RDB.HSet(ctx, fmt.Sprintf("user: %d", req.VideoId), "FavoriteCount", FavoriteVideoCount)
	}

	if err != nil {
		resp.StatusCode = 500
		resp.StatusMsg = "service赞操作失败!"
		return resp, err
	}
	resp.StatusMsg = "service操作成功!"
	return resp, nil
}

func (f *FeedSrv) FavoriteList(ctx context.Context, req *pb.DouyinFavoriteListRequest) (resp *pb.DouyinFavoriteListResponse, err error) {
	resp = new(pb.DouyinFavoriteListResponse)
	resp.StatusCode = 0

	conn, err := grpc.Dial("localhost:10002", grpc.WithInsecure())
	if err != nil {
		log.Fatal("Failed to connect: %v", err)
	}
	defer conn.Close()

	userClient := pb.NewUserServiceClient(conn)
	videoList, err := db.NewFeedDao(ctx).FavoriteList(req)

	if err != nil {
		resp.StatusCode = 500
		resp.StatusMsg = "sevice获取喜欢列表失败"
		return resp, err
	}

	for _, v := range videoList {
		userReq := pb.DouyinUserRequest{UserId: v.AuthorId}
		userResp, _ := userClient.UserInfo(ctx, &userReq)
		video := &pb.Video{
			Id:            v.ID,
			Author:        userResp.User,
			PlayUrl:       oss.GetVideo(v.PlayURL),
			CoverUrl:      oss.GetVideo(v.CoverURL),
			FavoriteCount: v.FavoriteCount,
			CommentCount:  v.CommentCount,
			IsFavorite:    v.IsFavorite,
			Title:         v.Title,
		}
		resp.VideoList = append(resp.VideoList, video)
	}

	resp.StatusMsg = "sevice获取喜欢列表成功"
	return resp, nil
}
