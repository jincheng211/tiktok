package model

import (
	"gorm.io/gorm"
	"time"
)

// Video
type Video struct {
	ID            int64          `gorm:"column:id;primaryKey;unique;"`
	AuthorId      int64          `gorm:"column:author_id;default:(-);"`
	CommentCount  int64          `gorm:"column:comment_count;default:(-);"`
	FavoriteCount int64          `gorm:"column:favorite_count;default:(-);"`
	IsFavorite    bool           `gorm:"column:is_favorite;default:(-);"`
	CoverURL      string         `gorm:"column:cover_url;default:(-);"`
	PlayURL       string         `gorm:"column:play_url;default:(-);"`
	Title         string         `gorm:"column:title;default:(-);"`
	CreatedAt     time.Time      `gorm:"column:created_at;type:datetime;not null" `
	UpdatedAt     time.Time      `gorm:"column:updated_at;type:datetime;"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at;type:datetime;" `
}

type FavoriteVideo struct {
	ID      int64 `gorm:"column:id;primaryKey;unique;auto_increment:1;"`
	VideoID int64 `gorm:"column:video_id;"`
	UserID  int64 `gorm:"column:user_id;"`
}

func (v *Video) TableName() string {
	return "video"
}
