package admin

import (
	"myapi/app/admin/controllers"
	"myapi/app/middleware"

	"github.com/gin-gonic/gin"
)

// SetupAdminRouter 配置并返回admin路由
func SetupAdminRouter(r *gin.Engine, apiPrefix string) {
	// 创建admin路由组（不使用认证中间件）
	admin := r.Group(apiPrefix)

	// 管理员索引路由
	admin.GET("/", controllers.Index)

	// 注释掉认证中间件，允许直接访问
	// admin.Use(middleware.AuthMiddleware())
	admin.Use(middleware.AuthMiddleware())
	{
		// 管理员用户相关路由
		users := admin.Group("/users")
		{
			users.GET("/", controllers.GetAllUsers)
			users.GET("/:id", controllers.GetUser)
			users.POST("/", controllers.CreateUser)
			users.PUT("/:id", controllers.UpdateUser)
			users.DELETE("/:id", controllers.DeleteUser)
		}

		// 管理员产品相关路由
		products := admin.Group("/products")
		{
			products.GET("/", controllers.GetAllProducts)
			products.GET("/:id", controllers.GetProduct)
			products.POST("/", controllers.CreateProduct)
			products.PUT("/:id", controllers.UpdateProduct)
			products.DELETE("/:id", controllers.DeleteProduct)
		}

	}

}
