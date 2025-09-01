package controllers

import (
	"myapi/app/helper"
	"myapi/app/services"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// LoginByPassword 处理密码登录请求
func LoginByPassword(c *gin.Context) {
	// 1. 解析请求体
	var loginData struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBind(&loginData); err != nil {
		// 检查错误信息，确定是哪个字段为空
		if strings.Contains(err.Error(), "Username") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "用户名不能为空"})
			return
		} else if strings.Contains(err.Error(), "Password") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "密码不能为空"})
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
			return
		}
	}

	// 2. 获取客户端IP
	clientIP := helper.GetRealIP(c)

	// 3. 创建服务实例
	memberService := services.NewMemberService()

	// 4. 调用服务层方法进行登录
	member, token, err := memberService.LoginByPassword(loginData.Username, loginData.Password, clientIP)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": err.Error(),
		})
		return
	}

	// 5. 返回登录成功响应
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "登录成功",
		"data": gin.H{
			"token":    token,
			"userInfo": member,
		},
	})
}

// LoginBySmsCode 处理短信验证码登录请求
func LoginBySmsCode(c *gin.Context) {
	// 1. 解析请求体
	var loginData struct {
		Telephone string `json:"telephone" binding:"required,len=11"`
		Code      string `json:"code" binding:"required,len=6"`
	}
	if err := c.ShouldBind(&loginData); err != nil {
		// 检查错误信息，确定是哪个字段为空
		if strings.Contains(err.Error(), "Telephone") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "手机号不能为空"})
			return
		} else if strings.Contains(err.Error(), "Code") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "验证码不能为空"})
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
			return
		}
	}

	// 2. 获取客户端IP
	clientIP := c.ClientIP()

	// 3. 创建服务实例
	memberService := services.NewMemberService()

	// 4. 调用服务层方法进行登录
	member, token, err := memberService.LoginBySmsCode(loginData.Telephone, loginData.Code, clientIP)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": err.Error(),
		})
		return
	}

	// 5. 返回登录成功响应
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "登录成功",
		"data": gin.H{
			"token":    token,
			"userInfo": member,
		},
	})
}

// SendSmsCode 发送短信验证码
func SendSmsCode(c *gin.Context) {
	// 1. 解析请求体
	var smsData struct {
		Telephone string `json:"telephone" binding:"required,len=11"`
	}
	if err := c.ShouldBind(&smsData); err != nil {
		// 检查错误信息，确定是哪个字段为空
		if strings.Contains(err.Error(), "Telephone") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "手机号不能为空"})
			return
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
			return
		}
	}

	// 2. 创建服务实例
	memberService := services.NewMemberService()

	// 3. 调用服务层方法发送短信验证码
	err := memberService.SendSmsCode(smsData.Telephone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "发送验证码失败: " + err.Error(),
		})
		return
	}

	// 4. 返回发送成功响应
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "验证码已发送，请注意查收",
	})
}
