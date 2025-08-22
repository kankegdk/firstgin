package controllers

import (
	"net/http"
	"strconv"

	"myapi/app/models" // 导入模型层

	"github.com/gin-gonic/gin"
)

// GetUser 处理获取单个用户的请求
func GetUser(c *gin.Context) {
	// 1. 从URL参数中获取id
	idStr := c.Param("id")
	// 2. 将字符串id转换为整数
	id, err := strconv.Atoi(idStr)
	if err != nil {
		// 如果转换失败，返回错误响应
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的用户ID",
		})
		return
	}

	// 3. 调用模型层的方法获取用户数据
	user := models.GetUserByID(id)

	// 4. 返回JSON响应
	c.JSON(http.StatusOK, gin.H{
		"message": "获取用户成功",
		"data":    user,
	})
}

// GetAllUsers 处理获取所有用户的请求
func GetAllUsers(c *gin.Context) {
	// 直接调用模型层方法
	users := models.GetAllUsers()

	c.JSON(http.StatusOK, gin.H{
		"message": "获取用户列表成功",
		"data":    users,
		"count":   len(users),
	})
}

// CreateUser 处理创建用户的请求 (示例)
func CreateUser(c *gin.Context) {
	// 这里可以学习如何绑定JSON请求体
	// 在实际项目中，你会在这里解析请求数据，验证，然后调用 models.CreateUser(...)
	c.JSON(http.StatusOK, gin.H{
		"message": "用户创建功能待实现",
	})
}
