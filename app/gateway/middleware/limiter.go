package middleware

import (
	"context"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

// 限流器(令牌桶)
type limiter struct {
	visitors map[string]*rate.Limiter
}

func NewLimiter() *limiter {
	l := &limiter{visitors: make(map[string]*rate.Limiter)}
	go func() {
		for {
			<-time.After(time.Hour)
			l.visitors = make(map[string]*rate.Limiter)
		}
	}()
	return l
}

func (l *limiter) Visit(ctx context.Context, addr string) error {
	lim, exists := l.visitors[addr]
	if !exists {
		lim = rate.NewLimiter(1, 50000) // 设置每秒最多访问5次
		l.visitors[addr] = lim
	}
	if !lim.Allow() {
		return status.Errorf(codes.ResourceExhausted, "稍等")
	}
	return nil
}

// 限流器
func (l *limiter) Limiting(c *gin.Context) {
	err := l.Visit(c.Request.Context(), c.ClientIP())
	if err != nil {
		c.JSON(429, gin.H{
			"error": err.Error(),
		})
		c.Abort()
		return
	}
	c.Next()
}
