package controllers

import (
	"net/http"
	"strconv"

	"myapi/app/services"
	"myapi/app/structs"

	"github.com/gin-gonic/gin"
)

// GetProduct 处理获取单个产品的请求
func GetProduct(c *gin.Context) {
	// 1. 从URL参数中获取id
	idStr := c.Param("id")
	// 2. 将字符串id转换为整数
	id, err := strconv.Atoi(idStr)
	if err != nil {
		// 如果转换失败，返回错误响应
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的产品ID",
		})
		return
	}

	// 3. 创建服务实例
	productService := services.NewProductService()
	// 4. 调用服务层方法
	product := productService.GetProductByID(id)

	// 5. 返回JSON响应
	c.JSON(http.StatusOK, gin.H{
		"message": "获取产品成功",
		"data":    product,
	})
}

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

// CreateProduct 处理创建产品的请求
func CreateProduct(c *gin.Context) {
	// 1. 解析请求体
	var product structs.Product
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 2. 创建服务实例
	productService := services.NewProductService()
	// 3. 调用服务层方法
	err := productService.CreateProduct(&product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "创建产品失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "创建产品成功",
		"data":    product,
	})
}

// UpdateProduct 处理更新产品的请求
func UpdateProduct(c *gin.Context) {
	// 实现逻辑类似CreateProduct
	c.JSON(http.StatusOK, gin.H{
		"message": "更新产品功能待实现",
	})
}

// DeleteProduct 处理删除产品的请求
func DeleteProduct(c *gin.Context) {
	// 实现逻辑类似GetProduct
	c.JSON(http.StatusOK, gin.H{
		"message": "删除产品功能待实现",
	})
}
