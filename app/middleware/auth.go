package middleware

import (
	"net/http"
	"strconv"
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
			// 如果是简单令牌格式（向后兼容）
			if isSimpleToken(tokenString) {
				// 从简单令牌中提取用户ID
				userID := extractUserIDFromSimpleToken(tokenString)
				if userID > 0 {
					// 设置用户信息到上下文中
					c.Set("userID", userID)
					c.Set("user", gin.H{
						"id": userID,
					})
					c.Next()
					return
				}
			}

			logs.Debug("JWT验证失败:", err)
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "无效的token",
			})
			c.Abort()
			return
		}
		// log.Println("claims:", claims)

		// 4. Token验证通过，将用户信息存入上下文
		c.Set("userID", claims.UserID)
		c.Set("user", gin.H{
			"id": claims.UserID,
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
			if len(parts) == 2 && parts[0] == "Bearer" {
				tokenString := parts[1]
				cfg := config.GetConfig()
				publicKeyPath := cfg.JWTPublicKeyPath

				// 尝试解析JWT令牌
				claims, err := helper.ParseJWT(tokenString, publicKeyPath)
				if err == nil {
					// JWT验证成功
					c.Set("userID", claims.UserID)
					c.Set("user", gin.H{
						"id": claims.UserID,
					})
				} else if isSimpleToken(tokenString) {
					// 尝试解析简单令牌格式（向后兼容）
					userID := extractUserIDFromSimpleToken(tokenString)
					if userID > 0 {
						c.Set("userID", userID)
						c.Set("user", gin.H{
							"id": userID,
						})
					}
				}
			}
		}
		c.Next()
	}
}

// 判断是否为简单令牌格式（向后兼容）
func isSimpleToken(token string) bool {
	// 简单令牌格式: userID:timestamp:randomString
	parts := strings.Split(token, ":")
	return len(parts) == 3
}

// 从简单令牌中提取用户ID（向后兼容）
func extractUserIDFromSimpleToken(token string) int {
	parts := strings.Split(token, ":")
	if len(parts) != 3 {
		return 0
	}

	// 尝试解析第一个部分为用户ID
	userID, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0
	}

	return userID
}
