package main

import (
	"fmt"
	"time"

	"myapi/app/config"
	"myapi/app/helper"
	"myapi/app/structs"
)

func main() {
	// 初始化配置
	config.Init()
	cfg := config.GetConfig()

	fmt.Println("配置信息:")
	fmt.Printf("JWTSecret: %s\n", cfg.JWTSecret)
	fmt.Printf("JWTPrivateKeyPath: %s\n", cfg.JWTPrivateKeyPath)
	fmt.Printf("JWTPublicKeyPath: %s\n", cfg.JWTPublicKeyPath)

	// 创建一个模拟的Member对象
	member := &structs.Member{
		ID:        1,
		Username:  "testuser",
		Telephone: "13800138000",
	}

	// 测试使用配置中指定路径的私钥生成JWT
	fmt.Println("\n测试生成JWT...")
	token, err := helper.GenerateJWT(member, cfg.JWTPrivateKeyPath, 24*time.Hour)
	if err != nil {
		fmt.Printf("生成JWT失败: %v\n", err)
		return
	}

	fmt.Println("JWT生成成功!")
	fmt.Printf("令牌长度: %d 字符\n", len(token))
	fmt.Println("\n令牌前100个字符:")
	if len(token) > 100 {
		fmt.Println(token[:100] + "...")
	} else {
		fmt.Println(token)
	}

	// 测试验证JWT
	fmt.Println("\n测试验证JWT...")
	claims, err := helper.ParseJWT(token, cfg.JWTPublicKeyPath)
	if err != nil {
		fmt.Printf("验证JWT失败: %v\n", err)
		return
	}

	fmt.Println("JWT验证成功!")
	fmt.Printf("用户ID: %d\n", claims.UserID)
	fmt.Printf("用户名: %s\n", claims.Username)
	fmt.Printf("手机号: %s\n", claims.Phone)
	fmt.Println("\n配置文件路径问题已修复，密钥文件现在正确保存在配置指定的位置！")
}