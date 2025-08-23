package app

import (
	"myapi/app/admin"
	"myapi/app/api"
	"myapi/app/config"

	"github.com/gin-gonic/gin"
)

// SetupRouter 配置并返回路由引擎
// 这是一个公开的函数，将在main.go中被调用
func SetupRouter() *gin.Engine {
	// 使用gin.New()而不是gin.Default()，避免加载默认的日志中间件
	r := gin.New()
	// 仅添加Recovery中间件来处理panic
	r.Use(gin.Recovery())

	cnf := config.GetConfig()
	// fmt.Println(cnf.AppName)

	// dbcfg := config.GetDatabaseConfig()
	// fmt.Printf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
	// 	dbcfg.User, dbcfg.Password, dbcfg.Host, dbcfg.Port, dbcfg.DBName)

	// 设置API路由
	api.SetupAPIRoutes(r, cnf.AppName)

	// 设置admin路由
	admin.SetupAdminRouter(r, cnf.BackendAppName)
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": " hello world",
		})
	})
	return r
}
