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
	cfg := config.GetRedisConfig()

	// 创建Redis客户端
	client := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password:     cfg.Password,
		DB:           cfg.DBName, // 默认DB
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
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
	return RedisClient.Set(redisCtx, key, value, expiration).Err()
}

// GetCache 获取缓存
func GetCache(key string) (string, error) {
	return RedisClient.Get(redisCtx, key).Result()
}

// DelCache 删除缓存
func DelCache(key string) error {
	return RedisClient.Del(redisCtx, key).Err()
}
