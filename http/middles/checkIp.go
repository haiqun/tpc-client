package middles

import (
	"github.com/gin-gonic/gin"
	"tcp-tls-project/http/controllers"
)

// 是否开启当前中间件
const checkIpStatus bool = false

/**
 * 文件名和函数名称保持一致
 */
func CheckIp() gin.HandlerFunc {
	if !checkIpStatus {
		return EmptyMiddles
	}

	controller := controllers.Controller{}

	return func(context *gin.Context) {
		if context.ClientIP() == "127.0.0.1" {
			controller.BadResponse(10000, "", context)
			// 此处结束
			context.Abort()
			return
		}
		context.Next()
	}
}