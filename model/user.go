package model

import (
	"github.com/CocaineCong/grpc-todolist/consts"
	"golang.org/x/crypto/bcrypt"
)

// User
type User struct {
	Avatar          string `gorm:"column:avatar;"`                      // 用户头像
	BackgroundImage string `gorm:"column:background_image;"`            // 用户个人页顶部大图
	FavoriteCount   int64  `gorm:"column:favorite_count;default:(-);"`  // 喜欢数
	FollowCount     int64  `gorm:"column:follow_count;default:(-);"`    // 关注总数
	FollowerCount   int64  `gorm:"column:follower_count;default:(-);"`  // 粉丝总数
	ID              int64  `gorm:" column:id;primaryKey;unique;"`       // 用户id
	IsFollow        bool   `gorm:"column:is_follow;default:(-);"`       // true-已关注，false-未关注
	Name            string `gorm:"column:name;NOT NULL;"`               // 用户名称
	Signature       string `gorm:"column:signature;"`                   // 个人简介
	TotalFavorited  string `gorm:"column:total_favorited;default:(-);"` // 获赞数量
	WorkCount       int64  `gorm:"column:work_count;default:(-);"`      // 作品数
	Password        string `gorm:"column:password;"`                    // 密码
}

func (user *User) TableName() string {
	return "user"
}

// 密码加密
func (user *User) SetPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), consts.PassWordCost)
	if err != nil {
		return err
	}
	user.Password = string(bytes)
	return nil
}

// 检验密码
func (user *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}
