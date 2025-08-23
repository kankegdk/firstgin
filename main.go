package main

import (
	"log"
	"myapi/app/admin"
	"myapi/app/api"
	"myapi/app/config"
	"myapi/app/logs"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化配置
	config.Init()
	cnf := config.GetConfig()

	// 初始化日志记录器
	logger, err := logs.NewLogger(
		cnf.AccessLogPath, // 从配置中获取访问日志文件路径
		cnf.DebugLogPath,  // 从配置中获取调试日志文件路径
	)
	if err != nil {
		log.Fatalf("初始化日志记录器失败: %v", err)
	}
	// 设置全局日志实例
	logs.GlobalLogger = logger
	defer logger.Close()

	// 重定向标准输出和标准错误到调试日志
	logger.RedirectStdoutStderr()

	log.Println("正在启动服务器...")

	// 创建路由引擎但不设置路由
	router := gin.New()
	router.Use(gin.Recovery())
	// 先应用自定义的访问日志中间件
	router.Use(logger.GinLogger())
	
	// 然后再设置路由
	api.SetupAPIRoutes(router, cnf.AppName)
	admin.SetupAdminRouter(router, cnf.BackendAppName)
	
	// 添加根路由
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello, World!"})
	})

	// 启动服务器
	go func() {
		if err := router.Run(":" + cnf.ServerPort); err != nil {
			log.Fatalf("启动服务器失败: %v", err)
		}
	}()

	// 等待中断信号优雅地关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("正在关闭服务器...")
}
