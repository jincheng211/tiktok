package rpc

import (
	"context"
	"douyin/idl/pb"
)

//rpc Message(douyin_message_action_request) returns(douyin_message_action_response);
//rpc GetMessageList (douyin_message_chat_request) returns(douyin_message_chat_response);

func Message(ctx context.Context, req *pb.DouyinMessageActionRequest) (resp *pb.DouyinMessageActionResponse, err error) {
	resp, err = MessageClient.Message(ctx, req)
	if err != nil {
		return resp, err
	}
	return resp, nil
}

func GetMessageList(ctx context.Context, req *pb.DouyinMessageChatRequest) (resp *pb.DouyinMessageChatResponse, err error) {
	resp, err = MessageClient.GetMessageList(ctx, req)
	if err != nil {
		return resp, err
	}
	return resp, nil
}
