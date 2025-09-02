package storage

import (
	"context"
	"fmt"
	"log"
	"myapi/app/config"
	"time"

	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client
var redisCtx = context.Background()

// InitRedis 初始化Redis连接
func InitRedis() error {
	// 检查是否启用Redis缓存
	redisEnabled := config.GetBool("redisEnabled", false)
	if !redisEnabled {
		log.Println("Redis缓存功能已禁用")
		return nil
	}
	
	// 直接获取Redis配置
	host := config.GetString("redisHost", "localhost")
	port := config.GetInt("redisPort", 6379)
	password := config.GetString("redisPassword", "")
	db := config.GetInt("redisDb", 0)
	poolSize := config.GetInt("redisPoolSize", 20)
	minIdleConns := config.GetInt("redisMinIdleConns", 5)

	// 创建Redis客户端
	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", host, port),
		Password:     password,
		DB:           db, // 默认DB
		PoolSize:     poolSize,
		MinIdleConns: minIdleConns,
	})

	// 测试连接
	_, err := client.Ping(redisCtx).Result()
	if err != nil {
		return fmt.Errorf("failed to ping redis: %v", err)
	}
	RedisClient = client
	log.Println("Redis连接初始化成功")
	return nil
}

// GetRedis 获取Redis连接实例
func GetRedis() *redis.Client {
	return RedisClient
}

// CloseRedis 关闭Redis连接
func CloseRedis() {
	if RedisClient != nil {
		RedisClient.Close()
		log.Println("Redis连接已关闭")
	}
}

// SetCache 设置缓存，带过期时间
func SetCache(key string, value interface{}, expiration time.Duration) error {
	if RedisClient == nil {
		// Redis未初始化，返回nil错误表示操作被跳过
		return nil
	}
	return RedisClient.Set(redisCtx, key, value, expiration).Err()
}

// GetCache 获取缓存
func GetCache(key string) (string, error) {
	if RedisClient == nil {
		// Redis未初始化，返回空字符串和错误
		return "", fmt.Errorf("redis not initialized")
	}
	return RedisClient.Get(redisCtx, key).Result()
}

// DelCache 删除缓存
func DelCache(key string) error {
	if RedisClient == nil {
		// Redis未初始化，返回nil错误表示操作被跳过
		return nil
	}
	return RedisClient.Del(redisCtx, key).Err()
}
