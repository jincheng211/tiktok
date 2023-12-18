package middleware

import (
	"github.com/afex/hystrix-go/hystrix"
	"github.com/gin-gonic/gin"
)

func ConfigureHystrixCommand(commandName string, Hystrix hystrix.CommandConfig) {
	hystrix.ConfigureCommand(commandName, Hystrix)
}

func HystrixMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 在这里进行Hystrix命令的配置

		Hystrix := hystrix.CommandConfig{
			Timeout:                1000,   // 设置超时时间
			MaxConcurrentRequests:  100000, // 最大并发请求数
			ErrorPercentThreshold:  10,     // 错误百分比阈值
			SleepWindow:            5000,   // 熔断持续时间
			RequestVolumeThreshold: 100000, // 阈值请求量
		}
		ConfigureHystrixCommand("my_command", Hystrix)

		// 继续处理请求
		c.Next()
	}
}
