package helper

import (
	"fmt"
	"net/http"
	"strings"

	"myapi/app/config"

	"github.com/gin-gonic/gin"
)

// CheckLogin 验证用户是否已登录
// 返回值：
// - bool: 用户是否已登录
// - *JWTClaims: 用户的JWT声明信息，如果未登录则为nil
// - error: 验证过程中遇到的错误，如果验证成功则为nil
func CheckLogin(c *gin.Context) (bool, *JWTClaims, error) {
	// 从请求头中获取Authorization
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return false, nil, fmt.Errorf("请求头中Authorization为空")
	}

	// 检查Authorization格式
	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		return false, nil, fmt.Errorf("Authorization格式错误，应为Bearer <token>")
	}

	// 解析JWT令牌
	tokenString := parts[1]
	// 使用配置中的公钥路径
	publicKeyPath := config.GetString("jwtPublicKeyPath", "config/rsa_public.pem")

	// 使用RSA公钥解析和验证JWT
	claims, err := ParseJWT(tokenString, publicKeyPath)
	if err != nil {
		return false, nil, fmt.Errorf("无效的token: %v", err)
	}

	// 如果解析成功，将用户信息存入上下文
	c.Set("userID", claims.UserID)
	c.Set("user", claims.Member)
	// 安全地从map中获取weid值并进行类型断言
	if weid, ok := claims.Member["weid"]; ok {
		c.Set("weid", weid)
	}
	if sid, ok := claims.Member["sid"]; ok {
		c.Set("sid", sid)
	}

	return true, claims, nil
}

// MustLogin 检查用户是否登录，如果未登录则直接返回错误响应
// 这个函数更符合控制器中当前的使用方式
func MustLogin(c *gin.Context) bool {
	isLoggedIn, _, _ := CheckLogin(c)
	if !isLoggedIn {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "请先登录",
		})
		return false
	}
	return true
}

// GetCurrentUser 获取当前登录用户的信息
func GetCurrentUser(c *gin.Context) (map[string]interface{}, bool) {
	user, exists := c.Get("user")
	if !exists {
		// 如果上下文中没有用户信息，尝试重新验证
		isLoggedIn, claims, _ := CheckLogin(c)
		if !isLoggedIn || claims == nil {
			return nil, false
		}
		return claims.Member, true
	}

	userMap, ok := user.(map[string]interface{})
	return userMap, ok
}

// GetCurrentUserID 获取当前登录用户的ID
func GetCurrentUserID(c *gin.Context) (int, bool) {
	userID, exists := c.Get("userID")
	if !exists {
		// 如果上下文中没有用户ID，尝试重新验证
		isLoggedIn, claims, _ := CheckLogin(c)
		if !isLoggedIn || claims == nil {
			return 0, false
		}
		return claims.UserID, true
	}

	id, ok := userID.(int)
	return id, ok
}
