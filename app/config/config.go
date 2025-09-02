package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/joho/godotenv" // 用于从.env文件加载环境变量
)

var (
	config map[string]interface{}
	once   sync.Once // 确保配置只初始化一次
)

// Init 初始化配置
func Init() {
	once.Do(func() {
		// 初始化配置map
		config = make(map[string]interface{})

		// 自动扫描并加载env目录下的所有.env文件
		envDir, _ := filepath.Abs("env")
		fileInfos, err := os.ReadDir(envDir)
		if err != nil {
			log.Printf("无法读取env目录: %v", err)
		} else {
			for _, fileInfo := range fileInfos {
				if !fileInfo.IsDir() && filepath.Ext(fileInfo.Name()) == ".env" {
					envFile := filepath.Join(envDir, fileInfo.Name())
					// 加载.env文件到map
					envMap, err := godotenv.Read(envFile)
					if err != nil {
						log.Printf("未找到或无法加载配置文件 %s: %v", envFile, err)
					} else {
						log.Printf("成功加载配置文件: %s", envFile)
						// 将配置合并到全局config map中
						for key, value := range envMap {
							// 尝试将值转换为合适的类型
							if boolValue, err := strconv.ParseBool(value); err == nil {
								config[key] = boolValue
							} else if intValue, err := strconv.Atoi(value); err == nil {
								config[key] = intValue
							} else if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
								config[key] = floatValue
							} else {
								// 保留为字符串
								config[key] = value
							}
						}
					}
				}
			}
		}
		log.Println("配置加载完成")
	})
}

// GetConfig 获取全局配置map（单例模式）
func GetConfig() map[string]interface{} {
	if config == nil {
		Init()
	}
	return config
}

// IsDevelopment 判断是否是开发环境
func IsDevelopment() bool {
	return GetString("environment", "") == "development"
}

// IsProduction 判断是否是生产环境
func IsProduction() bool {
	return GetString("environment", "") == "production"
}

// GetServerAddress 获取服务器监听地址
func GetServerAddress() string {
	return ":" + GetString("serverPort", "8080")
}

// GetString 获取字符串类型的配置值
func GetString(key string, defaultValue string) string {
	if config == nil {
		Init()
	}

	value, exists := config[key]
	if !exists {
		return defaultValue
	}

	// 根据值的实际类型进行转换
	switch v := value.(type) {
	case string:
		return v
	default:
		// 尝试将其他类型转换为字符串
		return fmt.Sprintf("%v", v)
	}
}

// GetBool 获取布尔类型的配置值
func GetBool(key string, defaultValue bool) bool {
	if config == nil {
		Init()
	}

	value, exists := config[key]
	if !exists {
		return defaultValue
	}

	// 尝试直接获取布尔值
	if boolValue, ok := value.(bool); ok {
		return boolValue
	}

	// 尝试将字符串转换为布尔值
	if strValue, ok := value.(string); ok {
		if boolValue, err := strconv.ParseBool(strValue); err == nil {
			return boolValue
		}
	}

	return defaultValue
}

// GetInt 获取整数类型的配置值
func GetInt(key string, defaultValue int) int {
	if config == nil {
		Init()
	}

	value, exists := config[key]
	if !exists {
		return defaultValue
	}

	// 尝试直接获取整数值
	if intValue, ok := value.(int); ok {
		return intValue
	}

	// 尝试将其他数值类型转换为int
	if floatValue, ok := value.(float64); ok {
		return int(floatValue)
	}

	// 尝试将字符串转换为int
	if strValue, ok := value.(string); ok {
		if intValue, err := strconv.Atoi(strValue); err == nil {
			return intValue
		}
	}

	return defaultValue
}

// GetFloat 获取浮点类型的配置值
func GetFloat(key string, defaultValue float64) float64 {
	if config == nil {
		Init()
	}

	value, exists := config[key]
	if !exists {
		return defaultValue
	}

	// 尝试直接获取浮点值
	if floatValue, ok := value.(float64); ok {
		return floatValue
	}

	// 尝试将整数类型转换为float64
	if intValue, ok := value.(int); ok {
		return float64(intValue)
	}

	// 尝试将字符串转换为float64
	if strValue, ok := value.(string); ok {
		if floatValue, err := strconv.ParseFloat(strValue, 64); err == nil {
			return floatValue
		}
	}

	return defaultValue
}
