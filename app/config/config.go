package config

import (
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/joho/godotenv" // 用于从.env文件加载环境变量
)

// Config 全局配置结构体
type Config struct {
	AppName        string
	Environment    string // development, production, testing
	ServerPort     string
	JWTSecret      string
	DebugMode      bool
	DatabaseURL    string
	BackendAppName string
	AccessLogPath  string
	DebugLogPath   string
}

var (
	config *Config
	once   sync.Once // 确保配置只初始化一次
)

// Init 初始化配置
func Init() {
	once.Do(func() {
		// 1. 加载.env文件（如果存在）
		// 在生产环境中通常不会使用.env文件，而是直接设置环境变量
		godotenv.Load(".env")

		// 2. 初始化配置
		config = &Config{
		AppName:        getEnv("APP_NAME", "MyGinApp"),
		Environment:    getEnv("ENVIRONMENT", "development"),
		ServerPort:     getEnv("SERVER_PORT", "8080"),
		BackendAppName: getEnv("BackendAppName", "admin"),
		JWTSecret:      getEnv("JWT_SECRET", "your-default-secret-key-change-in-production"),
		DebugMode:      getEnvAsBool("DEBUG_MODE", true),
		DatabaseURL:    getEnv("DATABASE_URL", "host=localhost user=postgres dbname=mydb sslmode=disable"),
		AccessLogPath:  getEnv("ACCESS_LOG_PATH", "logs/access.log"),
		DebugLogPath:   getEnv("DEBUG_LOG_PATH", "logs/debug.log"),
	}

		log.Printf("配置初始化完成: %s (%s)", config.AppName, config.Environment)
	})
}

// GetConfig 获取全局配置实例（单例模式）
func GetConfig() *Config {
	if config == nil {
		Init()
	}
	return config
}

// IsDevelopment 判断是否是开发环境
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

// IsProduction 判断是否是生产环境
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

// GetServerAddress 获取服务器监听地址
func (c *Config) GetServerAddress() string {
	return ":" + c.ServerPort
}

// Helper function to get environment variable or return default value
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// Helper function to get environment variable as boolean
func getEnvAsBool(key string, defaultValue bool) bool {
	strValue := getEnv(key, "")
	if strValue == "" {
		return defaultValue
	}

	value, err := strconv.ParseBool(strValue)
	if err != nil {
		return defaultValue
	}
	return value
}

// Helper function to get environment variable as integer
func getEnvAsInt(key string, defaultValue int) int {
	strValue := getEnv(key, "")
	if strValue == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(strValue)
	if err != nil {
		return defaultValue
	}
	return value
}
