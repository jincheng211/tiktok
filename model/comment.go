package model

import "time"

// Comment
type Comment struct {
	Content    string    `gorm:"column:content;"`                     // 评论内容
	CreateDate string    `gorm:"column:create_date";`                 // 评论发布日期，格式 mm-dd
	ID         int64     `gorm:"column:id;primaryKey;autoIncrement;"` // 评论id
	UserID     int64     `gorm:"column:user_id;"`                     // 评论用户信息
	VideoID    int64     `gorm:"column:video_id"`
	CreatedAt  time.Time `gorm:"column:created_at;type:datetime;not null;" `
}

func (c *Comment) TableName() string {
	return "Comment"
}
