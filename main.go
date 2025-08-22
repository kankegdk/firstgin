package main

import (
	"myapi/app"
	"myapi/app/config"
)

func main() {
	config.Init()
	// 通过app包下的SetupRouter函数获取配置好的路由引擎
	router := app.SetupRouter()

	// 启动服务器
	router.Run() // 默认监听 0.0.0.0:8080
}
