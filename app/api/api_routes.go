package api

import (
	"myapi/app/api/controllers"
	"myapi/app/middleware"

	"github.com/gin-gonic/gin"
)

// SetupAPIRoutes 配置API路由
func SetupAPIRoutes(r *gin.Engine, apiPrefix string) {
	// 定义API路由组，应用认证中间件
	api := r.Group(apiPrefix)
	api.Use(middleware.AuthMiddleware())
	{
		// 用户相关路由
		users := api.Group("/u")
		{
			users.GET("/", controllers.GetAllUsers) // GET /api/users
			users.GET("/:id", controllers.GetUser)  // GET /api/users/123
			users.POST("/", controllers.CreateUser) // POST /api/users
		}

		// 产品相关路由
		products := api.Group("/products")
		{
			products.GET("/", controllers.GetAllProducts)
		}
	}

	// 公共API路由组，不需要认证
	public := r.Group(apiPrefix)
	{
		// 定义健康检查路由
		public.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "OK"})
		})
		// 定义ping路由
		public.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "OK", "message": "pong"})
		})
		// 广告相关路由
		public.GET("/ads", controllers.GetAllAds)
	}
}
