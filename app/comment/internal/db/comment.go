package db

import (
	"context"
	"douyin/idl/pb"
	"douyin/model"
	"douyin/pkg/jwt"
	"gorm.io/gorm"
	"time"
)

type CommentDao struct {
	*gorm.DB
}

func NewCommentDao(ctx context.Context) *CommentDao {
	return &CommentDao{NewDBClient(ctx)}
}

func (dao *CommentDao) Comment(req *pb.DouyinCommentActionRequest) (int64, error) {
	// 评论
	claimn, _ := jwt.ParseToken(req.Token)
	comment := model.Comment{
		Content:    req.CommentText,
		CreateDate: time.Now().Format("2006-01-02 15:04:05"),
		UserID:     claimn.UserID,
		VideoID:    req.VideoId,
		CreatedAt:  time.Now(),
	}
	err := dao.Model(&model.Comment{}).Create(&comment).Error
	if err != nil {
		return 0, err
	}
	return comment.ID, nil
}

// 删除评论
func (dao *CommentDao) DeleteComment(req *pb.DouyinCommentActionRequest) error {
	var comment model.Comment
	result := dao.First(&comment, req.CommentId)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil
		} else {
			return result.Error
		}
	}
	result = dao.Model(&model.Comment{}).Delete(&comment)
	if result.Error != nil {
		// 处理错误
		return result.Error
	}
	return nil
}

// 获取评论
func (dao *CommentDao) CommentList(req *pb.DouyinCommentListRequest) ([]*model.Comment, error) {
	var commentList []*model.Comment
	err := dao.Model(&model.Comment{}).Where("video_id = ?", req.VideoId).Order("created_at DESC").Find(&commentList).Error
	if err != nil {
		return nil, err
	}
	return commentList, nil
}

// 获取评论count
func (dao *CommentDao) CommentCount(VideoID int64) (int64, error) {
	var CommentCount int64
	err := dao.Model(&model.Comment{}).Where("video_id = ?", VideoID).Count(&CommentCount).Error
	if err != nil {
		return -1, err
	}
	return CommentCount, nil
}
