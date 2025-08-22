package controllers

import (
	"net/http"

	// 导入模型层

	"github.com/gin-gonic/gin"
)

// GetUser 处理获取单个用户的请求
func GetAllProducts(c *gin.Context) {
	// 1. 从URL参数中获取id
	//idStr := c.Param("id")
	// 2. 将字符串id转换为整数

	c.JSON(http.StatusOK, gin.H{
		"message": "获取产品列表成功",
		//"data":    ["id":1],
		"count": 100,
	})
}
