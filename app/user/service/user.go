package service

import (
	"context"
	"douyin/app/user/internal/db"
	"douyin/idl/pb"
	"douyin/pkg/e"
	"douyin/pkg/jwt"
	"fmt"
	"strconv"
	"sync"
)

type UserSrv struct {
	pb.UnimplementedUserServiceServer
}

var UserSrvIns *UserSrv
var UserSrvOnce sync.Once

// 单例模式
// https://www.python100.com/html/112030.html(单例模式作用)
func GetUserSrv() *UserSrv {
	UserSrvOnce.Do(func() {
		UserSrvIns = &UserSrv{}
	})
	return UserSrvIns
}

func (u *UserSrv) UserRegister(ctx context.Context, req *pb.DouyinUserRegisterRequest) (*pb.DouyinUserRegisterResponse, error) {
	resp := new(pb.DouyinUserRegisterResponse)
	resp.StatusCode = e.SUCCESS
	id, err := db.NewUserDao(ctx).UserRegister(req)
	if err != nil {
		resp.StatusCode = e.ERROR
		resp.StatusMsg = "用户注册失败"
		return resp, err
	}
	// 用户注册成功
	resp.StatusMsg = "用户注册成功"
	resp.UserId = id
	return resp, nil
}

func (u *UserSrv) UserLogin(ctx context.Context, req *pb.DouyinUserLoginRequest) (resp *pb.DouyinUserLoginResponse, err error) {

	resp = new(pb.DouyinUserLoginResponse)
	resp.StatusCode = e.SUCCESS
	r, err := db.NewUserDao(ctx).GetUserLoginInfo(req)

	if err != nil {
		resp.StatusCode = e.ERROR
		return nil, err
	}

	userInfo := map[string]interface{}{
		"Id":              r.ID,
		"Name":            r.Name,
		"FollowCount":     r.FollowCount,
		"FollowerCount":   r.FollowerCount,
		"IsFollow":        false,
		"Avatar":          r.Avatar,
		"BackgroundImage": r.BackgroundImage,
		"Signature":       r.Signature,
		"TotalFavorited":  r.TotalFavorited,
		"WorkCount":       r.WorkCount,
		"FavoriteCount":   r.FavoriteCount,
	}

	// 将用户信息存入redis
	err = db.RDB.HMSet(ctx, fmt.Sprintf("user: %d", r.ID), userInfo).Err()
	if err != nil {
		panic(err)
	}

	token, err := jwt.GenerateToken(r.ID)
	if err != nil {
		panic(err)
	}

	resp = &pb.DouyinUserLoginResponse{
		StatusCode: e.SUCCESS,
		StatusMsg:  "登录成功",
		UserId:     r.ID,
		Token:      token,
	}
	return
}

func (u *UserSrv) UserInfo(ctx context.Context, req *pb.DouyinUserRequest) (*pb.DouyinUserResponse, error) {
	resp := new(pb.DouyinUserResponse)
	resp.StatusCode = 0 // 状态码，0-成功，其他值-失败
	user, _ := db.NewUserDao(ctx).GetUserInfo(req)

	// redis更新数据
	WorkCount, _ := db.RDB.HGet(ctx, fmt.Sprintf("user: %d", req.UserId), "WorkCount").Int64()
	if WorkCount > user.WorkCount {
		db.NewUserDao(ctx).UpdateUserVideoCount(req, WorkCount)
	}

	FavoriteCount, _ := db.RDB.HGet(ctx, fmt.Sprintf("user: %d", req.UserId), "FavoriteCount").Int64()
	if FavoriteCount > user.FavoriteCount {
		db.NewUserDao(ctx).UpdateFavoriteVideoCount(req, FavoriteCount)
	}

	// string类型
	TotalFavoritedCount, _ := db.RDB.HGet(ctx, fmt.Sprintf("user: %d", req.UserId), "TotalFavorited").Result()
	if TotalFavoritedCount > user.TotalFavorited {
		db.NewUserDao(ctx).UpdateFavoriteCount(req, TotalFavoritedCount)
	}
	TotalFavorited, _ := strconv.ParseInt(TotalFavoritedCount, 10, 64)

	FollowCount, _ := db.RDB.HGet(ctx, fmt.Sprintf("user: %d", req.UserId), "FollowCount").Int64()
	if FollowCount > user.FollowCount {
		db.NewUserDao(ctx).UpdateFollowCount(req, FollowCount)
	}

	FollowerCount, _ := db.RDB.HGet(ctx, fmt.Sprintf("user: %d", req.UserId), "FollowerCount").Int64()
	if FollowerCount > user.FollowerCount {
		db.NewUserDao(ctx).UpdateFollowerCount(req, FollowerCount)
	}

	resp.User = &pb.User{
		Id:              user.ID,
		Name:            user.Name,
		FollowCount:     FollowCount,
		FollowerCount:   FollowerCount,
		IsFollow:        false,
		Avatar:          user.Avatar,
		BackgroundImage: user.BackgroundImage,
		Signature:       user.Signature,
		TotalFavorited:  TotalFavorited,
		WorkCount:       WorkCount,
		FavoriteCount:   FavoriteCount,
	}

	resp.StatusMsg = "获取信息成功"
	return resp, nil
}
