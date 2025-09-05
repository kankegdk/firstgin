package middleware

import (
	"log"
	"net/http"

	"myapi/app/helper"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware JWT认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 使用helper包中的CheckLogin函数进行验证
		isLoggedIn, _, err := helper.CheckLogin(c)
		if !isLoggedIn {
			log.Println("JWT验证失败:", err)
			// 根据错误信息返回不同的错误提示
			errMsg := "无效的token"
			if err != nil {
				errMsg = err.Error()
			}
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": errMsg,
			})
			c.Abort() // 中止后续处理
			return
		}

		// 验证通过，继续处理下一个中间件或路由 handler
		c.Next()
	}
}
