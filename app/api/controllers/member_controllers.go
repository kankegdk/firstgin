package controllers

import (
	"log"
	"myapi/app/helper"
	"myapi/app/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// LoginByPassword 处理密码登录请求
func LoginByPassword(c *gin.Context) {
	// 1. 解析请求体
	var loginData struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	// 使用公共验证器验证请求
	if !helper.ValidateRequest(c, &loginData) {
		return
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
		Telephone string `json:"telephone" binding:"required,len=11" customvalidate:"mobile"`
		Code      string `json:"code" binding:"required,len=6"`
	}
	// 使用公共验证器验证请求
	if !helper.ValidateRequest(c, &loginData) {
		return
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
		Telephone string `json:"telephone" binding:"required,len=11" customvalidate:"mobile"`
	}
	// 使用公共验证器验证请求
	if !helper.ValidateRequest(c, &smsData) {
		log.Println("发送短信验证码失败: 请求参数验证失败")
		log.Println("请求参数:", smsData)
		return
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
