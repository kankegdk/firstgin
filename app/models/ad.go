package models

import (
	"log"

	"myapi/app/config"
	"myapi/app/storage"
	"myapi/app/structs"
)

// GetAllAds 从数据库获取所有广告的方法
// pageUrl 参数用于筛选特定页面的广告，如果为空则返回所有广告
func GetAllAds(pageUrl string) []structs.Ad {
	// 获取共享的gorm连接实例
	gormDB := storage.GetGormDB()
	if gormDB == nil {
		log.Println("GORM连接为空")
		return []structs.Ad{}
	}

	// 获取表前缀
	cfg := config.GetDatabaseConfig()
	tablePrefix := cfg.GetTablePrefix()
	tableName := tablePrefix + "ad"

	// 使用gorm查询数据
	var ads []structs.Ad
	query := gormDB.Table(tableName).Order("sort ASC")
	// 如果提供了pageUrl参数，则添加条件筛选
	if pageUrl != "" {
		query = query.Where("page_url = ?", pageUrl)
	}
	result := query.Find(&ads)
	if result.Error != nil {
		log.Printf("查询广告数据失败: %v", result.Error)
		return []structs.Ad{}
	}

	return ads
}
