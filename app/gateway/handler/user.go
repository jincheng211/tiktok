package handler

import (
	"douyin/app/gateway/rpc"
	"douyin/idl/pb"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func UserRegister(ctx *gin.Context) {
	var userReq pb.DouyinUserRegisterRequest
	userReq.Username = ctx.Query("username")
	userReq.Password = ctx.Query("password")

	r, err := rpc.UserRegister(ctx, &userReq)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status_code": 500,
			"status_msg":  "http 注册错误!",
			"err":         err,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status_code": http.StatusOK,
		"status_msg":  "http 用户注册成功!",
		"user_id":     r.UserId,
		"token":       r.Token,
	})
}

func UserLogin(ctx *gin.Context) {
	var userReq pb.DouyinUserLoginRequest
	userReq.Username = ctx.Query("username")
	userReq.Password = ctx.Query("password")

	userResp, err := rpc.UserLogin(ctx, &userReq)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status_code": 500,
			"status_msg":  "http 登录错误!",
			"err":         err,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status_code": 0,
		"status_msg":  "http 登录成功!",
		"user_id":     userResp.UserId,
		"token":       userResp.Token,
	})
}

func UserInfo(ctx *gin.Context) {
	var userReq pb.DouyinUserRequest
	userReq.UserId, _ = strconv.ParseInt(ctx.Query("user_id"), 10, 64)
	userReq.Token = ctx.Query("token")

	userResp, err := rpc.UserInfo(ctx, &userReq)
	if err != nil {
		fmt.Println("err", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status_code": 500,
			"status_msg":  userResp.StatusMsg,
			"err":         err,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status_code": 0,
		"status_msg":  userResp.StatusMsg,
		"user":        userResp.User,
	})
}
