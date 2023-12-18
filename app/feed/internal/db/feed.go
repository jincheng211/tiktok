package db

import (
	"context"
	"douyin/idl/pb"
	"douyin/model"
	"douyin/pkg/ffmpeg"
	"douyin/pkg/jwt"
	"douyin/pkg/oss"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type FeedDao struct {
	*gorm.DB
}

func NewFeedDao(ctx context.Context) *FeedDao {
	return &FeedDao{NewDBClient(ctx)}
}

// 获取视频信息
func (dao *FeedDao) GetFeedList(req *pb.DouyinFeedRequest) ([]*model.Video, error) {
	var videoList []*model.Video

	// 第一次获取视频信息
	if req.LatestTime == 0 {
		err := dao.Model(&model.Video{}).Where("created_at < ?", time.Now()).Order("created_at DESC").
			Limit(5).Find(&videoList).Error
		if err != nil {
			return nil, err
		}
		return videoList, nil
	}

	// 下滑获得更新的数据
	err := dao.Model(&model.Video{}).Where("created_at < ?", time.Unix(req.LatestTime, 0).Format("2006-01-02 15:04:05")).Order("created_at ASC").
		Limit(5).Find(&videoList).Error
	if err != nil {
		return nil, err
	}
	return videoList, nil
}

// 获取发布视频流
func (dao *FeedDao) GetPublishList(req *pb.DouyinPublishListRequest) ([]*model.Video, error) {
	var videoList []*model.Video
	err := dao.Model(&model.Video{}).Where("author_id = ?", req.UserId).Find(&videoList).Error
	if err != nil {
		return nil, err
	}
	return videoList, nil
}

// 发布视频
func (dao *FeedDao) PublishVideo(req *pb.DouyinPublishActionRequest) error {
	var video model.Video
	var videoCount int64

	// 解析token数据
	clamin, err := jwt.ParseToken(req.Token)
	if err != nil {
		return err
	}

	// 开启数据库事务
	tx := dao.Begin()

	// 获取视频数
	dao.Model(&model.Video{}).Count(&videoCount)

	// 获取视频封面
	cover, err := ffmpeg.GetVideoCover(req.Data)
	if err != nil {
		// 封面获取失败
		tx.Rollback()
		return err
	}

	// 上传视频以及封面
	err = oss.PutVideo(req.Data, cover, time.Now().Unix(), clamin.UserID)
	if err != nil {
		// 上传视频失败，回滚事务
		tx.Rollback()
		return err
	}

	// 上传视频
	video = model.Video{
		AuthorId:      clamin.UserID,
		CommentCount:  0,
		CoverURL:      "douyinVideoList/" + strconv.FormatInt(clamin.UserID, 10) + "/" + strconv.FormatInt(time.Now().Unix(), 10) + ".jpg",
		FavoriteCount: 0,
		ID:            videoCount + 10001,
		IsFavorite:    false,
		PlayURL:       "douyinVideoList/" + strconv.FormatInt(clamin.UserID, 10) + "/" + strconv.FormatInt(time.Now().Unix(), 10) + ".mp4",
		Title:         req.Title,
		CreatedAt:     time.Now(),
	}

	// 插入视频数据到数据库
	if err := tx.Create(&video).Error; err != nil {
		// 插入数据失败，回滚事务
		tx.Rollback()
		return err
	}

	// 提交事务
	tx.Commit()

	// 成功上传视频和保存数据到数据库
	return nil
}

func (dao *FeedDao) Favorite(req *pb.DouyinFavoriteActionRequest) (int64, error) {
	claimn, _ := jwt.ParseToken(req.Token)
	favorite := model.FavoriteVideo{
		VideoID: req.VideoId,
		UserID:  claimn.UserID,
	}
	// 1-点赞
	if req.ActionType == 1 {
		if err := dao.Model(&model.FavoriteVideo{}).Create(&favorite).Error; err != nil {
			return -1, err
		}
		// 点赞数+1
		var like model.Video
		err := dao.Model(&model.Video{}).Where("id = ?", req.VideoId).Find(&like).Error
		if err != nil {
			return -1, err
		}
		like.FavoriteCount++
		err = dao.Model(&model.Video{}).Where("id = ?", req.VideoId).Update("favorite_count", like.FavoriteCount).Error
		if err != nil {
			return -1, err
		}

		return like.FavoriteCount, nil
	}

	// 取消点赞
	if err := dao.Where(&favorite).Delete(&model.FavoriteVideo{}).Error; err != nil {
		return -1, err
	}
	// 点赞数-1
	var like model.Video
	err := dao.Model(&model.Video{}).Where("id = ?", req.VideoId).Find(&like).Error
	if err != nil {
		return -1, err
	}
	like.FavoriteCount--
	err = dao.Model(&model.Video{}).Where("id = ?", req.VideoId).Update("favorite_count", like.FavoriteCount).Error
	if err != nil {
		return -1, err
	}
	return like.FavoriteCount, nil
}

// 获取喜欢视频信息
func (dao *FeedDao) FavoriteList(req *pb.DouyinFavoriteListRequest) ([]*model.Video, error) {
	var favoriteVideos []*model.FavoriteVideo
	if err := dao.Model(&model.FavoriteVideo{}).Where("user_id = ?", req.UserId).Find(&favoriteVideos).Error; err != nil {
		return nil, err
	}

	var videoIDs []int64
	for _, fav := range favoriteVideos {
		videoIDs = append(videoIDs, fav.VideoID)
	}

	var videoList []*model.Video
	if err := dao.Model(&model.Video{}).Find(&videoList, videoIDs).Error; err != nil {
		return nil, err
	}
	return videoList, nil
}

func (dao *FeedDao) GetVideoCount(req *pb.DouyinPublishActionRequest) (VideoCount int64, err error) {
	// 解析token数据
	clamin, _ := jwt.ParseToken(req.Token)
	err = dao.Model(&model.Video{}).Where("author_id = ?", clamin.UserID).Count(&VideoCount).Error
	if err != nil {
		return 0, err
	}
	return VideoCount, nil
}

func (dao *FeedDao) UpdateCommentCount(VideoID, CommentCount int64) error {
	err := dao.Model(&model.Video{}).Where("id = ?", VideoID).Update("comment_count", CommentCount).Error
	if err != nil {
		return err
	}
	return nil
}

func (dao *FeedDao) UpdateFavoriteCount(VideoID, FavoriteCount int64) error {
	err := dao.Model(&model.Video{}).Where("id = ?", VideoID).Update("favorite_count", FavoriteCount).Error
	if err != nil {
		return err
	}
	return nil
}

func (dao *FeedDao) IsFavorite(UserID, VideoID int64) bool {
	var like model.FavoriteVideo // 假设你有一个名为 Like 的模型结构

	// 构建查询条件
	err := dao.Model(&model.FavoriteVideo{}).Where("video_id = ? AND user_id = ?", VideoID, UserID).First(&like).Error
	if err != nil {
		// 处理错误，可能是未找到记录或其他错误
		return false
	} else {
		// 找到匹配的记录
		return true
	}
}

func (dao *FeedDao) FavoriteVideoCount(req *pb.DouyinFavoriteActionRequest) (int64, error) {
	var FavoriteVideoCount int64
	claimn, _ := jwt.ParseToken(req.Token)
	if err := dao.Model(model.FavoriteVideo{}).Where("user_id = ?", claimn.UserID).Count(&FavoriteVideoCount).Error; err != nil {
		return -1, err
	}
	return FavoriteVideoCount, nil
}
