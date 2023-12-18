package handler

import (
	"bytes"
	"douyin/app/gateway/rpc"
	pb "douyin/idl/pb"
	"douyin/pkg/jwt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strconv"
)

func PublishVideo(ctx *gin.Context) {
	var req pb.DouyinPublishActionRequest

	req.Token = ctx.PostForm("token")
	req.Title = ctx.PostForm("title")
	file, _ := ctx.FormFile("data")

	// 获取电影数据
	video, err := file.Open()
	if err != nil {
		panic(err)
	}
	defer video.Close()
	buf := bytes.NewBuffer(nil)
	if _, err = io.Copy(buf, video); err != nil {
		panic(err)
	}

	req.Data = buf.Bytes()
	resp, err := rpc.PublishVideo(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
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

func GetFeedList(ctx *gin.Context) {
	var req pb.DouyinFeedRequest
	req.LatestTime, _ = strconv.ParseInt(ctx.Query("latest_time"), 10, 64)
	req.Token = ctx.Query("token")

	resp, err := rpc.GetFeedList(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status_code": resp.StatusCode,
			"status_msg":  resp.StatusMsg,
			"err":         err,
		})
		return
	}

	videoList := make([]map[string]interface{}, 0)

	// 循环遍历resp.VideoList，并将每个元素转换为map
	for _, v := range resp.VideoList {
		authorMap := map[string]interface{}{
			"id":               v.Author.Id,
			"name":             v.Author.Name,
			"follow_count":     v.Author.FollowCount,
			"follower_count":   v.Author.FollowerCount,
			"is_follow":        v.Author.IsFollow,
			"avatar":           v.Author.Avatar,
			"background_image": v.Author.BackgroundImage,
			"signature":        v.Author.Signature,
			"total_favorited":  v.Author.TotalFavorited,
			"work_count":       v.Author.WorkCount,
			"favorite_count":   v.Author.FavoriteCount,
		}

		videoMap := map[string]interface{}{
			"id":             v.Id,
			"author":         authorMap,
			"play_url":       v.PlayUrl,
			"cover_url":      v.CoverUrl,
			"favorite_count": v.FavoriteCount,
			"comment_count":  v.CommentCount,
			"is_favorite":    v.IsFavorite,
			"title":          v.Title,
		}
		videoList = append(videoList, videoMap)
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status_code": resp.StatusCode,
		"status_msg":  resp.StatusMsg,
		"next_time":   resp.NextTime,
		"video_list":  videoList,
	})
}

func GetPublishList(ctx *gin.Context) {
	var req pb.DouyinPublishListRequest
	req.Token = ctx.Query("token")
	UserId := ctx.Query("user_id")
	req.UserId, _ = strconv.ParseInt(UserId, 10, 64)
	resp, err := rpc.GetPublishList(ctx, &req)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status_code": resp.StatusCode,
			"status_msg":  resp.StatusMsg,
			"err":         err,
		})
		return
	}

	videoList := make([]map[string]interface{}, 0)
	// 循环遍历resp.VideoList，并将每个元素转换为map
	for _, v := range resp.VideoList {
		authorMap := map[string]interface{}{
			"id":               v.Author.Id,
			"name":             v.Author.Name,
			"follow_count":     v.Author.FollowCount,
			"follower_count":   v.Author.FollowerCount,
			"is_follow":        v.Author.IsFollow,
			"avatar":           v.Author.Avatar,
			"background_image": v.Author.BackgroundImage,
			"signature":        v.Author.Signature,
			"total_favorited":  v.Author.TotalFavorited,
			"work_count":       v.Author.WorkCount,
			"favorite_count":   v.Author.FavoriteCount,
		}

		videoMap := map[string]interface{}{
			"id":             v.Id,
			"author":         authorMap,
			"play_url":       v.PlayUrl,
			"cover_url":      v.CoverUrl,
			"favorite_count": v.FavoriteCount,
			"comment_count":  v.CommentCount,
			"is_favorite":    v.IsFavorite,
			"title":          v.Title,
		}
		videoList = append(videoList, videoMap)
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status_code": resp.StatusCode,
		"status_msg":  resp.StatusMsg,
		"video_list":  videoList,
	})
}

func Favorite(ctx *gin.Context) {
	var req pb.DouyinFavoriteActionRequest
	req.Token = ctx.Query("token")
	videoId := ctx.Query("video_id")
	actioType := ctx.Query("action_type")
	req.VideoId, _ = strconv.ParseInt(videoId, 10, 64)
	i, _ := strconv.Atoi(actioType)
	req.ActionType = int32(i)
	resp, err := rpc.Favorite(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status_code": resp.StatusCode,
			"status_msg":  resp.StatusMsg,
			"error":       err,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status_code": resp.StatusCode,
		"status_msg":  resp.StatusMsg,
	})
}

func FavoriteList(ctx *gin.Context) {
	var req pb.DouyinFavoriteListRequest
	req.Token = ctx.Query("token")
	userId := ctx.Query("user_id")
	req.UserId, _ = strconv.ParseInt(userId, 10, 64)
	_, err := jwt.ParseToken(req.Token)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status_code": 500,
			"status_msg":  "token解析失败",
			"error":       err,
		})
		return
	}
	resp, err := rpc.FavoriteList(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status_code": resp.StatusCode,
			"status_msg":  resp.StatusMsg,
			"error":       err,
		})
		return
	}
	videoList := make([]map[string]interface{}, 0)

	// 循环遍历resp.VideoList，并将每个元素转换为map
	for _, v := range resp.VideoList {
		authorMap := map[string]interface{}{
			"id":               v.Author.Id,
			"name":             v.Author.Name,
			"follow_count":     v.Author.FollowCount,
			"follower_count":   v.Author.FollowerCount,
			"is_follow":        v.Author.IsFollow,
			"avatar":           v.Author.Avatar,
			"background_image": v.Author.BackgroundImage,
			"signature":        v.Author.Signature,
			"total_favorited":  v.Author.TotalFavorited,
			"work_count":       v.Author.WorkCount,
			"favorite_count":   v.Author.FavoriteCount,
		}

		videoMap := map[string]interface{}{
			"id":             v.Id,
			"author":         authorMap,
			"play_url":       v.PlayUrl,
			"cover_url":      v.CoverUrl,
			"favorite_count": v.FavoriteCount,
			"comment_count":  v.CommentCount,
			"is_favorite":    v.IsFavorite,
			"title":          v.Title,
		}
		videoList = append(videoList, videoMap)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status_code": resp.StatusCode,
		"status_msg":  resp.StatusMsg,
		"video_list":  videoList,
	})
}
