package db

import (
	"context"
	"douyin/idl/pb"
	"douyin/model"
	"douyin/pkg/jwt"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type MessageDao struct {
	*gorm.DB
}

func NewMessageDao(ctx context.Context) *MessageDao {
	return &MessageDao{NewDBClient(ctx)}
}

func (dao *MessageDao) Message(req *pb.DouyinMessageActionRequest) error {
	claimn, _ := jwt.ParseToken(req.Token)
	message := model.Message{
		Content:    req.Content,
		CreateTime: time.Now().Unix(),
		FromUserID: claimn.UserID,
		ToUserID:   req.ToUserId,
	}
	fmt.Println(message.Content)
	err := dao.Model(model.Message{}).Create(&message).Error
	if err != nil {
		return err
	}
	return nil
}

func (dao *MessageDao) GetMessageList(req *pb.DouyinMessageChatRequest) ([]*model.Message, error) {
	var messageList []*model.Message
	claimn, _ := jwt.ParseToken(req.Token)
	err := dao.Model(&model.Message{}).Where("(from_user_id = ? AND to_user_id = ?) OR (from_user_id = ? AND to_user_id = ?) ", claimn.UserID, req.ToUserId, req.ToUserId, claimn.UserID).
		Find(&messageList).Error
	if err != nil {
		return nil, err
	}
	return messageList, nil
}
