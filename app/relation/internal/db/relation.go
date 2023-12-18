package db

import (
	"context"
	"douyin/idl/pb"
	"douyin/model"
	"douyin/pkg/jwt"
	"gorm.io/gorm"
)

type RelationDao struct {
	*gorm.DB
}

func NewRelationDao(ctx context.Context) *RelationDao {
	return &RelationDao{NewDBClient(ctx)}
}

// 关注
func (dao *RelationDao) Follow(req *pb.DouyinRelationActionRequest) (err error) {
	claimn, _ := jwt.ParseToken(req.Token)
	// 查询对方是否关注了自己
	var follow model.Relation
	result := dao.Model(&model.Relation{}).Where("follow_id = ? AND follower_id = ?", req.ToUserId, claimn.UserID).First(&follow)
	// 关注
	if result.Error != nil {
		// (没有查询到) 对方没有关注
		if result.Error == gorm.ErrRecordNotFound {
			relation := model.Relation{
				FollowID:   claimn.UserID,
				FollowerID: req.ToUserId,
				IsFriend:   false,
			}
			err = dao.Model(&model.Relation{}).Create(&relation).Error
			// 创建错误
			if err != nil {
				return err
			}
			return nil
		}
		// 其他错误
		return err
	}

	// 查询到了,则粉丝变朋友
	follow.IsFriend = true
	err = dao.Model(&model.Relation{}).Where("follow_id = ? AND follower_id = ?", req.ToUserId, claimn.UserID).Save(&follow).Error
	if err != nil {
		return err
	}
	relation := model.Relation{
		FollowID:   claimn.UserID,
		FollowerID: req.ToUserId,
		IsFriend:   true,
	}
	err = dao.Model(&model.Relation{}).Create(&relation).Error
	if err != nil {
		return err
	}
	return nil
}

// 取关
func (dao *RelationDao) UnFollow(req *pb.DouyinRelationActionRequest) (err error) {
	claimn, _ := jwt.ParseToken(req.Token)
	// 查询对方是否关注了自己
	var follow model.Relation
	result := dao.Model(&model.Relation{}).Where("follow_id = ? AND follower_id = ?", req.ToUserId, claimn.UserID).First(&follow)
	// 取关
	if result.Error != nil {
		// (没有查询到) 对方没有关注
		if result.Error == gorm.ErrRecordNotFound {
			err = dao.Where("follow_id = ? AND follower_id = ?", claimn.UserID, req.ToUserId).Delete(&model.Relation{}).Error
			if err != nil {
				return err
			}
			return nil
		} else {
			return err
		}
	}
	// 对方关注了 朋友变粉丝
	follow.IsFriend = false
	err = dao.Model(&model.Relation{}).Where("follow_id = ? AND follower_id = ?", req.ToUserId, claimn.UserID).Save(&follow).Error
	if err != nil {
		return err
	}

	err = dao.Where("follow_id = ? AND follower_id = ?", claimn.UserID, req.ToUserId).Delete(&model.Relation{}).Error
	if err != nil {
		return err
	}
	return nil
}

// 获取关注列表
func (dao *RelationDao) GetFollowList(req *pb.DouyinRelationFollowListRequest) (FollowerList []*model.Relation, err error) {
	err = dao.Model(&model.Relation{}).Where("follow_id = ?", req.UserId).Find(&FollowerList).Error
	return FollowerList, err
}

// 获取关注人数
func (dao *RelationDao) GetFollowCount(UserID int64) (FollowCount int64, err error) {
	if err = dao.Model(&model.Relation{}).Where("follow_id = ?", UserID).Count(&FollowCount).Error; err != nil {
		return -1, err
	}
	return FollowCount, nil
}

// 获取粉丝列表
func (dao *RelationDao) GetFollowerList(req *pb.DouyinRelationFollowerListRequest) (FollowList []*model.Relation, err error) {
	err = dao.Model(&model.Relation{}).Where("follower_id = ?", req.UserId).Find(&FollowList).Error
	return FollowList, err
}

// 获取关注人数
func (dao *RelationDao) GetFollowerCount(UserID int64) (FollowerCount int64, err error) {
	if err = dao.Model(&model.Relation{}).Where("follower_id = ?", UserID).Count(&FollowerCount).Error; err != nil {
		return -1, err
	}
	return FollowerCount, nil
}

// 获取朋友列表
func (dao *RelationDao) GetFriendList(req *pb.DouyinRelationFriendListRequest) (FriendList []*model.Relation, err error) {
	err = dao.Model(&model.Relation{}).Where("follow_id = ? AND is_friend = ?", req.UserId, true).Find(&FriendList).Error
	return FriendList, err
}

// 判断有没有关注粉丝
func (dao *RelationDao) IsFriend(UserID, ToUserId int64) bool {
	var IsFriend model.Relation
	err := dao.Model(model.Relation{}).Where("follower_id = ? AND follow_id = ?", ToUserId, UserID).First(&IsFriend).Error
	if err != nil {
		return false
	}
	return true
}
