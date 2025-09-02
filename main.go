package main

import (
	"log"
	"myapi/app/admin"
	"myapi/app/api"
	"myapi/app/config"
	"myapi/app/storage"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化配置
	config.Init()

	log.Println("正在启动服务器...")

	// 初始化MySQL连接
	if err := storage.InitMySQL(); err != nil {
		log.Fatalf("初始化MySQL连接失败: %v", err)
	}
	defer storage.CloseMySQL()

	// 初始化Redis连接
	if err := storage.InitRedis(); err != nil {
		log.Fatalf("初始化Redis连接失败: %v", err)
	}
	defer storage.CloseRedis()

	// 创建路由引擎但不设置路由
	router := gin.New()
	router.Use(gin.Recovery())
	// 先应用自定义的访问日志中间件

	// 然后再设置路由
	appName := config.GetString("appName", "api")
	backendAppName := config.GetString("backendAppName", "admin")
	api.SetupAPIRoutes(router, appName)
	admin.SetupAdminRouter(router, backendAppName)

	// 添加根路由
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello, World!"})
	})

	// 启动服务器
	go func() {
		serverPort := config.GetString("serverPort", "8080")
		if err := router.Run(":" + serverPort); err != nil {
			log.Fatalf("启动服务器失败: %v", err)
		}
	}()

	// 等待中断信号优雅地关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("正在关闭服务器...")
}
