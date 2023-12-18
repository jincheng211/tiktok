package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			// 捕获函数执行过程中出现的 panic 错误。
			if r := recover(); r != nil {
				c.JSON(http.StatusOK, gin.H{
					"code": 404,
					"msg":  fmt.Sprintf("%s", r),
				})

				// 终止请求处理流程，确保在中间件函数内部处理完错误后不再继续处理后续的中间件和路由处理函数。
				c.Abort()
			}
		}()
		c.Next()
	}
}
