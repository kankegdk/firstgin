package services

import (
	"encoding/json"
	"log"
	"myapi/app/config"
	"myapi/app/models"
	"myapi/app/storage"
	"myapi/app/structs"
	"time"

	"github.com/go-redis/redis/v8"
)

// AdService 广告服务接口
type AdService interface {
	GetAllAds(pageUrl string) []structs.Ad
}

// adService 实现AdService接口的结构体
type adService struct{}

// NewAdService 创建一个新的广告服务实例
func NewAdService() AdService {
	return &adService{}
}

// GetAllAds 获取所有广告
func (s *adService) GetAllAds(pageUrl string) []structs.Ad {
	// 如果pageUrl有值，使用Redis缓存
	if pageUrl != "" {
		// 从配置中获取Redis前缀
		prefix := config.GetRedisConfig().Prefix

		// 构建缓存键
		cacheKey := prefix + ":" + "ads:" + pageUrl

		// 尝试从缓存获取数据
		cacheData, err := storage.GetCache(cacheKey)
		if err == nil {
			//这里加一个日志
			log.Println("从Redis缓存中获取广告数据cacheKey", cacheKey)
			// 缓存命中，解析JSON并返回
			var ads []structs.Ad
			if err2 := json.Unmarshal([]byte(cacheData), &ads); err2 == nil {
				return ads
			} else {
				log.Println("解析缓存数据失败:", err2)
			}
		} else if err != redis.Nil {
			// 如果是其他错误（非键不存在），记录日志
			log.Println("获取Redis缓存失败:", err)
		} else {
			// 键不存在的情况，记录日志但不报错
			log.Println("Redis缓存键不存在:", cacheKey)
		}

		// 缓存未命中，从数据库获取数据
		ads := models.GetAllAds(pageUrl)

		// 将结果存入缓存，过期时间设为24小时
		if data, err := json.Marshal(ads); err == nil {
			log.Println("将广告数据存入Redis缓存cacheKey", cacheKey)
			storage.SetCache(cacheKey, string(data), time.Hour*24)
		}

		return ads
	}
	// 等于空就去取数据库，直接返回空数组
	return []structs.Ad{}
	// 如果pageUrl为空，直接从数据库获取
	//return models.GetAllAds(pageUrl)
}
