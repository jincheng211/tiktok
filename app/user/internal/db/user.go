package db

import (
	"context"
	pb "douyin/idl/pb"
	"douyin/model"
	"errors"
	"github.com/CocaineCong/grpc-todolist/pkg/util/logger"
	"gorm.io/gorm"
	"strconv"
)

type UserDao struct {
	*gorm.DB
}

func NewUserDao(ctx context.Context) *UserDao {
	return &UserDao{NewDBClient(ctx)}
}

// register
func (dao *UserDao) UserRegister(req *pb.DouyinUserRegisterRequest) (int64, error) {
	var user model.User
	var count int64
	var userCount int64

	dao.Model(&model.User{}).Where("name = ?", req.Username).Count(&count)
	if count != 0 {
		return 0, errors.New("Username exit")
	}

	// 获取用户数量(确定id)
	err := dao.Model(&model.User{}).Count(&userCount).Error
	if err != nil {
		logger.LogrusObj.Error("get User count Error:" + err.Error())
		return 0, err
	}

	user = model.User{
		TotalFavorited:  strconv.Itoa(0),
		WorkCount:       0,
		ID:              userCount + 100001,
		Name:            req.Username,
		Signature:       "抖音越来越好!",
		Avatar:          "http://127.0.0.1:8080/douyin/static/avatar/1.png",
		BackgroundImage: "http://127.0.0.1:8080/douyin/static/background_image/1.jpeg",
		FavoriteCount:   0,
		FollowerCount:   0,
		FollowCount:     0,
		IsFollow:        true,
	}

	_ = user.SetPassword(req.Password)
	if err := dao.Model(&model.User{}).Create(&user).Error; err != nil {
		logger.LogrusObj.Error("Insert User Error:" + err.Error())
		return 0, err
	}
	return user.ID, nil
}

// GetUserInfo登录信息
func (dao *UserDao) GetUserLoginInfo(req *pb.DouyinUserLoginRequest) (r *model.User, err error) {
	err = dao.Model(&model.User{}).Where("name", req.Username).First(&r).Error
	return
}

// UserInfo 获取用户
func (dao *UserDao) GetUserInfo(req *pb.DouyinUserRequest) (r *model.User, err error) {
	err = dao.Model(&model.User{}).Where("id", req.UserId).First(&r).Error
	return
}

// UserInfo 更新用户视频数
func (dao *UserDao) UpdateUserVideoCount(req *pb.DouyinUserRequest, WorkCount int64) (err error) {
	// 增加VideoCount字段的值
	err = dao.Model(&model.User{}).Where("id = ?", req.UserId).Update("work_count", WorkCount).Error
	return
}

// UserInfo 更新用户喜欢数
func (dao *UserDao) UpdateFavoriteVideoCount(req *pb.DouyinUserRequest, FavoriteCount int64) (err error) {
	// 增加VideoCount字段的值
	err = dao.Model(&model.User{}).Where("id = ?", req.UserId).Update("favorite_count", FavoriteCount).Error
	return
}

// 更新用户的获赞数
func (dao *UserDao) UpdateFavoriteCount(req *pb.DouyinUserRequest, TotalFavorited string) error {
	err := dao.Model(&model.User{}).Where("id = ?", req.UserId).Update("total_favorited", TotalFavorited).Error
	return err
}

// 更新用户的获赞数
func (dao *UserDao) UpdateFollowCount(req *pb.DouyinUserRequest, FollowCount int64) error {
	err := dao.Model(&model.User{}).Where("id = ?", req.UserId).Update("follow_count", FollowCount).Error
	return err
}

// 更新用户的获赞数
func (dao *UserDao) UpdateFollowerCount(req *pb.DouyinUserRequest, FollowerCount int64) error {
	err := dao.Model(&model.User{}).Where("id = ?", req.UserId).Update("follower_count", FollowerCount).Error
	return err
}

// 获取用户的视频数
//func (dao *UserDao) GetUserVideoCount(req *pb.DouyinUserRequest) (workCount int64, err error) {
//	var user model.User
//	err = dao.Model(&model.User{}).Where("id = ?", req.UserId).First(&user).Error
//	if err != nil {
//		return user.WorkCount, err
//	}
//	return user.WorkCount, nil
//}
