package rpc

import (
	"context"
	"douyin/idl/pb"
)

func PublishVideo(ctx context.Context, req *pb.DouyinPublishActionRequest) (resp *pb.DouyinPublishActionResponse, err error) {
	resp, err = FeedClient.PublishVideo(ctx, req)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

func GetFeedList(ctx context.Context, req *pb.DouyinFeedRequest) (resp *pb.DouyinFeedResponse, err error) {
	resp, err = FeedClient.GetFeedList(ctx, req)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

func Favorite(ctx context.Context, req *pb.DouyinFavoriteActionRequest) (resp *pb.DouyinFavoriteActionResponse, err error) {
	resp, err = FeedClient.Favorite(ctx, req)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

func FavoriteList(ctx context.Context, req *pb.DouyinFavoriteListRequest) (resp *pb.DouyinFavoriteListResponse, err error) {
	resp, err = FeedClient.FavoriteList(ctx, req)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

func GetPublishList(ctx context.Context, req *pb.DouyinPublishListRequest) (resp *pb.DouyinPublishListResponse, err error) {
	resp, err = FeedClient.GetPublishList(ctx, req)
	if err != nil {
		return resp, err
	}
	return resp, nil
}
