package controllers

import (
	"log"
	"net/http"

	"myapi/app/helper"
	"myapi/app/models"   // 导入模型层
	"myapi/app/services" // 导入服务层

	"github.com/gin-gonic/gin"
)

// GetUser 处理获取单个用户的请求
func GetUser(c *gin.Context) {
	// 1. 从中间件设置的上下文中获取userID
	userID := helper.UID(c)
	log.Printf("userID: %v", userID)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "未授权访问",
		})
		return
	}
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "没有用户信息",
		})
		return
	}
	// 5. 返回JSON响应
	c.JSON(http.StatusOK, gin.H{
		"message": "获取用户成功",
		"data":    user,
	})
}

// GetAllUsers 处理获取所有用户的请求
func GetAllUsers(c *gin.Context) {
	// 1. 创建服务实例
	userService := services.NewUserService()
	// 2. 调用服务层方法
	users := userService.GetAllUsers()

	c.JSON(http.StatusOK, gin.H{
		"message": "获取用户列表成功",
		"data":    users,
		"count":   len(users),
	})
}

// CreateUser 处理创建用户的请求
func CreateUser(c *gin.Context) {
	// 1. 解析请求体
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 2. 创建服务实例
	userService := services.NewUserService()
	// 3. 调用服务层方法
	err := userService.CreateUser(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "创建用户失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "创建用户成功",
		"data":    user,
	})
}
