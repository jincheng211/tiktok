package service

import (
	"context"
	"douyin/app/message/internal/db"
	pb "douyin/idl/pb"
	"fmt"
	"sync"
)

type MessageSrv struct {
	pb.UnimplementedMessageServiceServer
}

var MessageSrvIns *MessageSrv
var MessageSrvOnce sync.Once

func GetMessageSrv() *MessageSrv {
	MessageSrvOnce.Do(func() {
		MessageSrvIns = &MessageSrv{}
	})
	return MessageSrvIns
}

func (m *MessageSrv) Message(ctx context.Context, req *pb.DouyinMessageActionRequest) (*pb.DouyinMessageActionResponse, error) {
	resp := new(pb.DouyinMessageActionResponse)
	resp.StatusCode = 0
	err := db.NewMessageDao(ctx).Message(req)
	fmt.Println("content", req.Content)
	if err != nil {
		resp.StatusCode = 500
		resp.StatusMsg = "service Message发送错误"
		return resp, err
	}

	resp.StatusMsg = "消息发送成功"
	return resp, nil
}

func (m *MessageSrv) GetMessageList(ctx context.Context, req *pb.DouyinMessageChatRequest) (*pb.DouyinMessageChatResponse, error) {
	resp := new(pb.DouyinMessageChatResponse)
	resp.StatusCode = 0
	messageList, err := db.NewMessageDao(ctx).GetMessageList(req)
	if err != nil {
		resp.StatusCode = 500
		resp.StatusMsg = "service Messagelist获取失败"
		return resp, err
	}

	resp.StatusMsg = "service Messagelist获取成功"
	for _, m := range messageList {
		message := &pb.Message{
			Id:         m.ID,
			ToUserId:   m.ToUserID,
			FromUserId: m.FromUserID,
			Content:    m.Content,
			CreateTime: m.CreateTime,
		}
		resp.MessageList = append(resp.MessageList, message)
	}
	return resp, nil
}
