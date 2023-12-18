package handler

import (
	"douyin/app/gateway/rpc"
	"douyin/idl/pb"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

//rpc Relation(douyin_relation_action_request) returns(douyin_relation_action_response);
//rpc GetFollowerList(douyin_relation_follower_list_request) returns(douyin_relation_follower_list_response);
//rpc GetFollowList(douyin_relation_follow_list_request) returns(douyin_relation_follow_list_response);
//rpc GetFriendList(douyin_relation_friend_list_request) returns(douyin_relation_friend_list_response);

func Relation(ctx *gin.Context) {
	var req pb.DouyinRelationActionRequest

	req.Token = ctx.Query("token")
	ToUserId := ctx.Query("to_user_id")
	req.ToUserId, _ = strconv.ParseInt(ToUserId, 10, 64)
	ActionType, _ := strconv.Atoi(ctx.Query("action_type"))
	req.ActionType = int32(ActionType)

	resp, err := rpc.Relation(ctx, &req)
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

func GetFollowList(ctx *gin.Context) {
	var req pb.DouyinRelationFollowListRequest

	req.Token = ctx.Query("token")
	UserId := ctx.Query("user_id")
	req.UserId, _ = strconv.ParseInt(UserId, 10, 64)

	resp, err := rpc.GetFollowList(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status_code": resp.StatusCode,
			"status_msg":  resp.StatusMsg,
			"err":         err,
		})
		return
	}

	UserList := make([]map[string]interface{}, 0)
	// 循环遍历resp.UserList，并将每个元素转换为map
	for _, user := range resp.UserList {
		UserMap := map[string]interface{}{
			"id":               user.Id,
			"name":             user.Name,
			"follow_count":     user.FollowCount,
			"follower_count":   user.FollowerCount,
			"is_follow":        true,
			"avatar":           user.Avatar,
			"background_image": user.BackgroundImage,
			"signature":        user.Signature,
			"total_favorited":  user.TotalFavorited,
			"work_count":       user.WorkCount,
			"favorite_count":   user.FavoriteCount,
		}
		UserList = append(UserList, UserMap)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status_code": resp.StatusCode,
		"status_msg":  resp.StatusMsg,
		"user_list":   UserList,
	})
}

func GetFollowerList(ctx *gin.Context) {
	var req pb.DouyinRelationFollowerListRequest

	req.Token = ctx.Query("token")
	UserId := ctx.Query("user_id")
	req.UserId, _ = strconv.ParseInt(UserId, 10, 64)

	resp, err := rpc.GetFollowerList(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status_code": resp.StatusCode,
			"status_msg":  resp.StatusMsg,
			"err":         err,
		})
		return
	}

	UserList := make([]map[string]interface{}, 0)
	// 循环遍历resp.UserList，并将每个元素转换为map
	for _, user := range resp.UserList {
		UserMap := map[string]interface{}{
			"id":               user.Id,
			"name":             user.Name,
			"follow_count":     user.FollowCount,
			"follower_count":   user.FollowerCount,
			"is_follow":        user.IsFollow,
			"avatar":           user.Avatar,
			"background_image": user.BackgroundImage,
			"signature":        user.Signature,
			"total_favorited":  user.TotalFavorited,
			"work_count":       user.WorkCount,
			"favorite_count":   user.FavoriteCount,
		}

		UserList = append(UserList, UserMap)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status_code": resp.StatusCode,
		"status_msg":  resp.StatusMsg,
		"user_list":   UserList,
	})
}

func GetFriendList(ctx *gin.Context) {
	var req pb.DouyinRelationFriendListRequest

	req.Token = ctx.Query("token")
	UserId := ctx.Query("user_id")
	req.UserId, _ = strconv.ParseInt(UserId, 10, 64)

	resp, err := rpc.GetFriendList(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status_code": resp.StatusCode,
			"status_msg":  resp.StatusMsg,
			"err":         err,
		})
		return
	}

	UserList := make([]map[string]interface{}, 0)
	// 循环遍历resp.UserList，并将每个元素转换为map
	for _, user := range resp.UserList {
		UserMap := map[string]interface{}{
			"id":               user.Id,
			"name":             user.Name,
			"follow_count":     user.FollowCount,
			"follower_count":   user.FollowerCount,
			"is_follow":        user.IsFollow,
			"avatar":           user.Avatar,
			"background_image": user.BackgroundImage,
			"signature":        user.Signature,
			"total_favorited":  user.TotalFavorited,
			"work_count":       user.WorkCount,
			"favorite_count":   user.FavoriteCount,
		}

		UserList = append(UserList, UserMap)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status_code": resp.StatusCode,
		"status_msg":  resp.StatusMsg,
		"user_list":   UserList,
	})
}
