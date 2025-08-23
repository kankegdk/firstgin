package controllers

import (
	"net/http"

	"myapi/app/services" // 导入服务层

	"github.com/gin-gonic/gin"
)

// GetAllProducts 处理获取所有产品的请求
func GetAllProducts(c *gin.Context) {
	// 1. 创建服务实例
	productService := services.NewProductService()
	// 2. 调用服务层方法
	products := productService.GetAllProducts()

	c.JSON(http.StatusOK, gin.H{
		"message": "获取产品列表成功",
		"data":    products,
		"count":   len(products),
	})
}
