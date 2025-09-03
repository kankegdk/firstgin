package controllers

import (
	"net/http"

	"myapi/app/services"

	"github.com/gin-gonic/gin"
)

// GetUser 处理获取单个用户的请求
func GetUser(c *gin.Context) {
	// 1. 从URL参数中获取id
	// idStr := c.Param("id")
	// // 2. 将字符串id转换为整数
	// id, err := strconv.Atoi(idStr)
	// if err != nil {
	// 	// 如果转换失败，返回错误响应
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"error": "无效的用户ID",
	// 	})
	// 	return
	// }

	// // 3. 创建服务实例
	// userService := services.NewUserService()
	// // 4. 调用服务层方法
	// user, err := userService.GetUserByID(id)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{
	// 		"error": "获取用户失败",
	// 	})
	// 	return
	// }

	// // 5. 返回JSON响应
	// c.JSON(http.StatusOK, gin.H{
	// 	"message": "获取用户成功",
	// 	"data":    user,
	// })
}

// GetAllUsers 处理获取所有用户的请求
func GetAllUsers(c *gin.Context) {
	// 1. 创建服务实例
	userService := services.NewUserService()
	// 2. 调用服务层方法
	users, err := userService.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "获取用户列表失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "admin 获取用户",
		"data":    users,
		"count":   len(users),
	})
}

// CreateUser 处理创建用户的请求
func CreateUser(c *gin.Context) {
	// // 1. 解析请求体
	// var user models.User
	// if err := c.ShouldBindJSON(&user); err != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"error": err.Error(),
	// 	})
	// 	return
	// }

	// // 2. 创建服务实例
	// userService := services.NewUserService()
	// // 3. 调用服务层方法
	// err := userService.CreateUser(&user)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{
	// 		"error": "创建用户失败",
	// 	})
	// 	return
	// }

	// c.JSON(http.StatusOK, gin.H{
	// 	"message": "创建用户成功",
	// 	"data":    user,
	// })
}

// UpdateUser 处理更新用户的请求
func UpdateUser(c *gin.Context) {
	// 实现逻辑类似CreateUser
	c.JSON(http.StatusOK, gin.H{
		"message": "更新用户功能待实现",
	})
}

// DeleteUser 处理删除用户的请求
func DeleteUser(c *gin.Context) {
	// 实现逻辑类似GetUser
	c.JSON(http.StatusOK, gin.H{
		"message": "删除用户功能待实现",
	})
}
