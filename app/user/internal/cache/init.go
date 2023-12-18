package cache

import (
	"github.com/go-redis/redis/v8"
)

var RDB *redis.NewRing

func InitRedisDB() {
	rdb := redis.NewRing(&redis.RingOptions{
		Addrs:              nil,
		NewClient:          nil,
		HeartbeatFrequency: 0,
		NewConsistentHash:  nil,
		Dialer:             nil,
		OnConnect:          nil,
		Username:           "",
		Password:           "",
		DB:                 0,
		MaxRetries:         0,
		MinRetryBackoff:    0,
		MaxRetryBackoff:    0,
		DialTimeout:        0,
		ReadTimeout:        0,
		WriteTimeout:       0,
		PoolFIFO:           false,
		PoolSize:           0,
		MinIdleConns:       0,
		MaxConnAge:         0,
		PoolTimeout:        0,
		IdleTimeout:        0,
		IdleCheckFrequency: 0,
		TLSConfig:          nil,
		Limiter:            nil,
	})
	RDB = rdb
}
