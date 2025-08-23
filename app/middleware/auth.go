package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"myapi/app/logs"
)

// AuthMiddleware JWT认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		// 1. 从请求头中获取token
		authHeader := c.GetHeader("Authorization")
		logs.Debug("Authorization header:", authHeader, "11111")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "请求头中Authorization为空",
			})
			c.Abort() // 中止后续处理
			return
		}

		// 2. 检查token格式（应该是 "Bearer <token>"）
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "Authorization格式错误，应为Bearer <token>",
			})
			c.Abort()
			return
		}

		// 3. 解析token（这里简化处理，实际项目中应该使用j-go等库验证JWT）
		token := parts[1]
		if token != "secret-token" { // 这里用硬编码演示，实际应该验证JWT签名
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "无效的token",
			})
			c.Abort()
			return
		}

		// 4. Token验证通过，可以在这里将用户信息存入上下文供后续使用
		// 例如：c.Set("userID", userID)
		c.Set("user", gin.H{
			"id":   1,
			"name": "管理员",
		})

		// 5. 继续处理下一个中间件或路由 handler
		c.Next()
	}
}

// OptionalAuth 可选认证中间件（验证失败也不中止请求）
func OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) == 2 && parts[0] == "Bearer" && parts[1] == "secret-token" {
				c.Set("user", gin.H{"id": 1, "name": "管理员"})
			}
		}
		c.Next()
	}
}
