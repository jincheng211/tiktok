package rpc

import (
	"context"
	"douyin/idl/pb"
)

func UserRegister(ctx context.Context, req *pb.DouyinUserRegisterRequest) (resp *pb.DouyinUserRegisterResponse, err error) {
	resp, err = UserClient.UserRegister(ctx, req)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func UserLogin(ctx context.Context, req *pb.DouyinUserLoginRequest) (resp *pb.DouyinUserLoginResponse, err error) {
	resp, err = UserClient.UserLogin(ctx, req)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func UserInfo(ctx context.Context, req *pb.DouyinUserRequest) (resp *pb.DouyinUserResponse, err error) {
	resp, err = UserClient.UserInfo(ctx, req)
	if err != nil {
		return resp, err
	}

	return resp, nil

}
