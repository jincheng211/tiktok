package rpc

import (
	"context"
	"douyin/idl/pb"
)

//rpc Relation(douyin_relation_action_request) returns(douyin_relation_action_response);
//rpc GetFollowerList(douyin_relation_follower_list_request) returns(douyin_relation_follower_list_response);
//rpc GetFollowList(douyin_relation_follow_list_request) returns(douyin_relation_follow_list_response);
//rpc GetFriendList(douyin_relation_friend_list_request) returns(douyin_relation_friend_list_response);

func Relation(ctx context.Context, req *pb.DouyinRelationActionRequest) (resp *pb.DouyinRelationActionResponse, err error) {
	resp, err = RelationClient.Relation(ctx, req)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

func GetFollowerList(ctx context.Context, req *pb.DouyinRelationFollowerListRequest) (resp *pb.DouyinRelationFollowerListResponse, err error) {
	resp, err = RelationClient.GetFollowerList(ctx, req)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

func GetFollowList(ctx context.Context, req *pb.DouyinRelationFollowListRequest) (resp *pb.DouyinRelationFollowListResponse, err error) {
	resp, err = RelationClient.GetFollowList(ctx, req)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

func GetFriendList(ctx context.Context, req *pb.DouyinRelationFriendListRequest) (resp *pb.DouyinRelationFriendListResponse, err error) {
	resp, err = RelationClient.GetFriendList(ctx, req)
	if err != nil {
		return resp, err
	}
	return resp, nil
}
