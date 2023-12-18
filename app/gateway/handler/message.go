package handler

import (
	"douyin/app/gateway/rpc"
	"douyin/idl/pb"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

//rpc Message(douyin_message_action_request) returns(douyin_message_action_response);
//rpc GetMessageList (douyin_message_chat_request) returns(douyin_message_chat_response);

func Message(ctx *gin.Context) {
	var req pb.DouyinMessageActionRequest
	req.Token = ctx.Query("token")
	ToUserId := ctx.Query("to_user_id")
	req.ToUserId, _ = strconv.ParseInt(ToUserId, 10, 64)
	ActionType, _ := strconv.Atoi(ctx.Query("action_type"))
	req.ActionType = int32(ActionType)
	req.Content = ctx.Query("content")
	resp, err := rpc.Message(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status_code": resp.StatusCode,
			"status_msg":  resp.StatusMsg,
			"err":         err,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status_code": resp.StatusCode,
		"status_msg":  resp.StatusMsg,
	})
}

func GetMessageList(ctx *gin.Context) {
	var req pb.DouyinMessageChatRequest
	req.Token = ctx.Query("token")
	ToUserId := ctx.Query("to_user_id")
	req.ToUserId, _ = strconv.ParseInt(ToUserId, 10, 64)

	resp, err := rpc.GetMessageList(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status_code": resp.StatusCode,
			"status_msg":  resp.StatusMsg,
			"err":         err,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status_code":   resp.StatusCode,
		"status_msg":    resp.StatusMsg,
		"messsage_list": resp.MessageList,
	})
}
