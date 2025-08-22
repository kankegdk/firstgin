package app

import (
	"myapi/app/config"
	"myapi/app/controllers" // 导入控制器
	"myapi/app/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRouter 配置并返回路由引擎
// 这是一个公开的函数，将在main.go中被调用
func SetupRouter() *gin.Engine {
	r := gin.Default()
	cnf := config.GetConfig()
	// fmt.Println(cnf.AppName)

	// dbcfg := config.GetDatabaseConfig()
	// fmt.Printf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
	// 	dbcfg.User, dbcfg.Password, dbcfg.Host, dbcfg.Port, dbcfg.DBName)

	// 定义API路由,这个需要验证路由
	api := r.Group(cnf.AppName)
	api.Use(middleware.AuthMiddleware())
	{
		// 用户相关路由
		users := api.Group("/u")
		{
			users.GET("/", controllers.GetAllUsers) // GET /api/users
			users.GET("/:id", controllers.GetUser)  // GET /api/users/123
			users.POST("/", controllers.CreateUser) // POST /api/users
		}

		// 你可以在这里继续添加其他路由组，例如：
		products := api.Group("/products")
		{
			products.GET("/", controllers.GetAllProducts)
		}
	}

	public := r.Group(cnf.AppName)
	{
		// 定义一个健康检查路由
		public.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "OK"})
		})
		// 定义一个健康检查路由
		public.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "OK", "message": "pong"})
		})

	}
	return r
}
