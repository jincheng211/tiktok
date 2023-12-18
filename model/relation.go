package model

import "time"

// User
type Relation struct {
	ID         int64     `gorm:"column:id;primaryKey;autoIncrement;"`
	FollowID   int64     `gorm:"column:follow_id;"`   // 关注者
	FollowerID int64     `gorm:"column:follower_id;"` // 被关注者
	IsFriend   bool      `gorm:"column:is_friend;default:(-);"`
	CreatedAt  time.Time `gorm:"column:created_at;type:datetime;not null" `
}

func (relation *Relation) TableName() string {
	return "relation"
}
