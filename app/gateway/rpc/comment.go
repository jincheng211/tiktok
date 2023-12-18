package rpc

import (
	"context"
	"douyin/idl/pb"
)

func Comment(ctx context.Context, req *pb.DouyinCommentActionRequest) (resp *pb.DouyinCommentActionResponse, err error) {
	resp, err = CommentClient.Comment(ctx, req)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

func CommentList(ctx context.Context, req *pb.DouyinCommentListRequest) (resp *pb.DouyinCommentListResponse, err error) {
	resp, err = CommentClient.CommentList(ctx, req)
	if err != nil {
		return resp, err
	}
	return resp, nil
}
