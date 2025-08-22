package config

import (
	"fmt"
)

// redisConfig 数据库配置
type RedisConfig struct {
	Host         string
	Port         int
	User         string
	Password     string
	DBName       string
	PoolSize     int
	MinIdleConns int
}

func GetRedisConfig() *RedisConfig {

	// 这里可以解析DATABASE_URL或者使用单独的环境变量
	return &RedisConfig{
		Host:         getEnv("REDIS_HOST", "localhost"),
		Port:         getEnvAsInt("REDIS_PORT", 6379),
		Password:     getEnv("REDIS_PASSWORD", "postgres"),
		DBName:       getEnv("REDIS_DB", ""),
		PoolSize:     getEnvAsInt("REDIS_POOL_SIZE", 20),
		MinIdleConns: getEnvAsInt("REDIS_MIN_IDLE_CONNS", 5),
	}
}

// GetConnectionString 获取数据库连接字符串
func (d *RedisConfig) GetConnectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s PoolSize=%s MinIdleConns=%s",
		d.Host, d.Port, d.User, d.Password, d.DBName, d.PoolSize, d.MinIdleConns)
}
