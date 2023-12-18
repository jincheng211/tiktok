package handler

import (
	"douyin/app/gateway/middleware"
	"douyin/app/gateway/rpc"
	"douyin/idl/pb"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

//rpc Comment(douyin_comment_action_request) returns(douyin_comment_action_response);
//rpc CommentList(douyin_comment_list_request) returns(douyin_comment_list_response);

func Comment(ctx *gin.Context) {
	// 使用Prometheus
	middleware.PrometheusCli.WithLabelValues("POST", "/douyin/comment/action/")

	var req pb.DouyinCommentActionRequest
	req.Token = ctx.Query("token")
	VideoId := ctx.Query("video_id")
	req.VideoId, _ = strconv.ParseInt(VideoId, 10, 64)
	ActionType, _ := strconv.Atoi(ctx.Query("action_type"))
	req.ActionType = int32(ActionType)
	if req.ActionType == 1 {
		CommentText := ctx.Query("comment_text")
		req.CommentText = CommentText
	} else {
		CommentId, _ := strconv.ParseInt(ctx.Query("comment_id"), 10, 64)
		req.CommentId = CommentId
	}
	resp, err := rpc.Comment(ctx, &req)
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
		"comment": map[string]interface{}{
			"id": resp.Comment.Id,
			"user": map[string]interface{}{
				"id":               resp.Comment.User.Id,
				"name":             resp.Comment.User.Name,
				"follow_count":     resp.Comment.User.FollowCount,
				"follower_count":   resp.Comment.User.FollowerCount,
				"is_follow":        resp.Comment.User.IsFollow,
				"avatar":           resp.Comment.User.Avatar,
				"background_image": resp.Comment.User.BackgroundImage,
				"signature":        resp.Comment.User.Signature,
				"total_favorited":  resp.Comment.User.TotalFavorited,
				"work_count":       resp.Comment.User.WorkCount,
				"favorite_count":   resp.Comment.User.FavoriteCount,
			},
			"content":     resp.Comment.Content,
			"create_date": resp.Comment.CreateDate,
		},
	})
}

func CommentList(ctx *gin.Context) {
	var req pb.DouyinCommentListRequest
	req.Token = ctx.Query("token")
	VideoId := ctx.Query("video_id")
	req.VideoId, _ = strconv.ParseInt(VideoId, 10, 64)

	resp, err := rpc.CommentList(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status_code": resp.StatusCode,
			"status_msg":  resp.StatusMsg,
			"err":         err,
		})
		return
	}

	CommentList := make([]map[string]interface{}, 0)

	// 循环遍历resp.VideoList，并将每个元素转换为map
	for _, c := range resp.CommentList {
		UserMap := map[string]interface{}{
			"id":               c.User.Id,
			"name":             c.User.Name,
			"follow_count":     c.User.FollowCount,
			"follower_count":   c.User.FollowerCount,
			"is_follow":        c.User.IsFollow,
			"avatar":           c.User.Avatar,
			"background_image": c.User.BackgroundImage,
			"signature":        c.User.Signature,
			"total_favorited":  c.User.TotalFavorited,
			"work_count":       c.User.WorkCount,
			"favorite_count":   c.User.FavoriteCount,
		}

		CommentMap := map[string]interface{}{
			"id":          c.Id,
			"user":        UserMap,
			"content":     c.Content,
			"create_date": c.CreateDate,
		}
		CommentList = append(CommentList, CommentMap)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status_code":  resp.StatusCode,
		"status_msg":   resp.StatusMsg,
		"comment_list": CommentList,
	})
}
