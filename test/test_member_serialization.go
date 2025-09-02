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

	fmt.Println("测试Member对象JSON序列化...")

	// 创建一个模拟的Member对象
	member := &structs.Member{
		ID:        1,
		Username:  "testuser",
		Telephone: "13800138000",
		Nickname:  "测试用户",
		Status:    1,
	}

	// 测试使用修改后的代码生成JWT
	token, err := helper.GenerateJWT(member, config.GetString("jwtPrivateKeyPath", ""), 24*time.Hour)
	if err != nil {
		fmt.Printf("生成JWT失败: %v\n", err)
		return
	}

	fmt.Println("JWT生成成功!")
	fmt.Printf("令牌长度: %d 字符\n", len(token))
	fmt.Println("\n令牌前200个字符:")
	if len(token) > 200 {
		fmt.Println(token[:200] + "...")
	} else {
		fmt.Println(token)
	}

	// 验证JWT并检查Member字段是否正确解析
	claims, err := helper.ParseJWT(token, config.GetString("jwtPublicKeyPath", ""))
	if err != nil {
		fmt.Printf("验证JWT失败: %v\n", err)
		return
	}

	fmt.Println("\nJWT验证成功!")
	fmt.Printf("用户ID: %d\n", claims.UserID)
	fmt.Printf("用户名: %s\n", claims.Username)
	fmt.Printf("手机号: %s\n", claims.Phone)
	fmt.Printf("Member字段是否为map类型: %T\n", claims.Member)
	fmt.Printf("Member字段键值对: %+v\n", claims.Member)

	// 尝试访问map中的值
	if id, ok := claims.Member["id"].(float64); ok {
		fmt.Printf("Member.id: %.0f\n", id)
	}
	if username, ok := claims.Member["username"].(string); ok {
		fmt.Printf("Member.username: %s\n", username)
	}
	if nickname, ok := claims.Member["nickname"].(string); ok {
		fmt.Printf("Member.nickname: %s\n", nickname)
	}
	if status, ok := claims.Member["status"].(float64); ok {
		fmt.Printf("Member.status: %.0f\n", status)
	}

	fmt.Println("\nMember对象已成功转换为键值对(map)格式! 现在JWT中的Member信息将以JSON对象形式存储，前端可以正确解析。")
}
