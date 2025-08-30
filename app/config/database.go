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
	Charset      string
	TablePrefix  string
	MaxLifetime  int
	MaxOpenConns int
	MaxIdleConns int
}

// GetDatabaseConfig 从主配置中获取数据库配置
func GetDatabaseConfig() *DatabaseConfig {

	// 这里可以解析DATABASE_URL或者使用单独的环境变量
	return &DatabaseConfig{
		Host:         getEnv("DB_HOST", "localhost"),
		Port:         getEnvAsInt("DB_PORT", 3306),
		User:         getEnv("DB_USER", "root"),
		Password:     getEnv("DB_PASSWORD", ""),
		DBName:       getEnv("DB_NAME", "mydb"),
		SSLMode:      getEnv("DB_SSL_MODE", "disable"),
		Charset:      getEnv("DB_CHARSET", "utf8mb4"),
		TablePrefix:  getEnv("DB_PREFIX", ""),
		MaxLifetime:  getEnvAsInt("DB_MaxLifetime", 30),
		MaxIdleConns: getEnvAsInt("DB_MaxIdleConns", 100),
		MaxOpenConns: getEnvAsInt("DB_MaxOpenConns", 10),
	}
}

// GetConnectionString 获取数据库连接字符串
func (d *DatabaseConfig) GetConnectionString() string {
	// 对于MySQL数据库，构建DSN格式的连接字符串
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		d.User, d.Password, d.Host, d.Port, d.DBName, d.Charset)
}

// GetTablePrefix 获取表前缀
func (d *DatabaseConfig) GetTablePrefix() string {
	return d.TablePrefix
}
