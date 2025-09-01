package middleware

import (
	"log"
	"net/http"
	"strings"

	"myapi/app/config"
	"myapi/app/helper"
	"myapi/app/logs"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware JWT认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		// 1. 从请求头中获取token
		authHeader := c.GetHeader("Authorization")
		logs.Debug("Authorization header:", authHeader)
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

		// 3. 解析并验证JWT令牌
		tokenString := parts[1]
		cfg := config.GetConfig()
		publicKeyPath := cfg.JWTPublicKeyPath

		// 使用RSA公钥解析和验证JWT
		claims, err := helper.ParseJWT(tokenString, publicKeyPath)
		if err != nil {
			logs.Debug("JWT验证失败:", err)
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "无效的token",
			})
			c.Abort()
			return
		}

		// 4. Token验证通过，将用户信息存入上下文
		c.Set("userID", claims.UserID)
		c.Set("user", claims.Member)
		// 安全地从map中获取weid值并进行类型断言
		if weid, ok := claims.Member["weid"]; ok {
			c.Set("weid", weid)
			log.Println("weid:", weid)
		} else {
			log.Println("claims.Member 中没有 weid 字段")
		}
		if sid, ok := claims.Member["sid"]; ok {
			c.Set("sid", sid)
			log.Println("sid:", sid)
		} else {
			log.Println("claims.Member 中没有 sid 字段")
		}
		log.Println("claims:", claims.Member)
		// 5. 继续处理下一个中间件或路由 handler
		c.Next()
	}
}
