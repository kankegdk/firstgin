package controllers

import (
	"net/http"

	"myapi/app/services"

	"github.com/gin-gonic/gin"
)

// GetAllAds 处理获取所有广告的请求
func GetAllAds(c *gin.Context) {
	// 1. 创建服务实例
	adService := services.NewAdService()
	// 2. 从请求中获取pageUrl参数
	pageUrl := c.Query("page_url")
	// 3. 调用服务层方法获取广告数据，传递pageUrl参数
	ads := adService.GetAllAds(pageUrl)

	// 4. 返回JSON响应
	c.JSON(http.StatusOK, gin.H{
		"message": "获取广告列表成功",
		"data":    ads,
	})
}