package config

import (
	"fmt"
)

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host         string
	Port         int
	User         string
	Password     string
	DBName       string
	SSLMode      string
	MaxLifetime  int
	MaxOpenConns int
	MaxIdleConns int
}

// GetDatabaseConfig 从主配置中获取数据库配置
func GetDatabaseConfig() *DatabaseConfig {

	// 这里可以解析DATABASE_URL或者使用单独的环境变量
	return &DatabaseConfig{
		Host:         getEnv("DB_HOST", "localhost"),
		Port:         getEnvAsInt("DB_PORT", 5432),
		User:         getEnv("DB_USER", "postgres"),
		Password:     getEnv("DB_PASSWORD", ""),
		DBName:       getEnv("DB_NAME", "mydb"),
		SSLMode:      getEnv("DB_SSL_MODE", "disable"),
		MaxLifetime:  getEnvAsInt("DB_MaxLifetime", 30),
		MaxIdleConns: getEnvAsInt("DB_MaxIdleConns", 100),
		MaxOpenConns: getEnvAsInt("DB_MaxOpenConns", 10),
	}
}

// GetConnectionString 获取数据库连接字符串
func (d *DatabaseConfig) GetConnectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.DBName, d.SSLMode)
}
