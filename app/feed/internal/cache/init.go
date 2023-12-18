package cache

import (
	"douyin/app/feed/internal/db"
	"github.com/go-redis/redis/v8"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

var RDB *redis.Client

func InitRedisDB() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	RDB = rdb
}

func NewDBClient(ctx context.Context) *gorm.DB {
	db := db.DB
	return db.WithContext(ctx)
}
