package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Index 处理索引请求
func Index(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "hello world",
	})
}
