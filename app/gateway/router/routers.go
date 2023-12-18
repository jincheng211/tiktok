package router

import (
	"douyin/app/gateway/handler"
	"douyin/app/gateway/middleware"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func NewRouter() *gin.Engine {
	ginRouter := gin.Default()
	l := middleware.NewLimiter()
	EsHookLog := middleware.EsHookLog()
	middleware.PrometheusInit()

	ginRouter.Use(middleware.Cors(), middleware.ErrorMiddleware(), middleware.ElasticMiddleware(EsHookLog))

	// 限流 熔断
	ginRouter.Use(middleware.HystrixMiddleware(), l.Limiting)

	store := cookie.NewStore([]byte("something-very-secret"))
	ginRouter.Use(sessions.Sessions("mysession", store))

	// 添加Prometheus的注册器到Gin
	ginRouter.GET("/metrics", gin.WrapH(promhttp.Handler()))

	douyin := ginRouter.Group("/douyin")
	{
		douyin.GET("ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, "success")
		})

		// 用户服务
		douyin.POST("/user/register/", handler.UserRegister)
		douyin.POST("/user/login/", handler.UserLogin)
		douyin.GET("/user/", handler.UserInfo)

		// feed
		douyin.GET("/feed/", handler.GetFeedList)
		douyin.POST("/publish/action/", handler.PublishVideo)
		douyin.GET("/publish/list/", handler.GetPublishList)
		douyin.POST("/favorite/action/", handler.Favorite)
		douyin.GET("/favorite/list/", handler.FavoriteList)

		// comment
		douyin.POST("/comment/action/", handler.Comment)
		douyin.GET("/comment/list/", handler.CommentList)

		// message
		douyin.POST("/message/action/", handler.Message)
		douyin.GET("/message/chat/", handler.GetMessageList)

		// relation
		douyin.POST("/relation/action/", handler.Relation)
		douyin.GET("/relation/follow/list/", handler.GetFollowList)
		douyin.GET("/relation/follower/list/", handler.GetFollowerList)
		douyin.GET("/relation/friend/list/", handler.GetFriendList)

		// 设置静态文件服务
		douyin.Static("/static", "../douyin/static")
	}

	return ginRouter
}
