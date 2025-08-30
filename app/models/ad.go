package models

import (
	"log"
	"time"

	"myapi/app/config"
	"myapi/app/storage"
)

// Ad 广告模型，定义了广告数据的结构
// 使用json标签来指定JSON序列化时的字段名
type Ad struct {
	ID               int       `json:"id"`
	Weid             int       `json:"weid"`
	Ocid             int       `json:"ocid"`
	Sid              int       `json:"sid"`
	Ptype            int       `json:"ptype"`
	Url              string    `json:"url"`
	Title            string    `json:"title"`
	Pic              string    `json:"pic"`
	Sort             int       `json:"sort"`
	Status           int       `json:"status"`
	ValidPeriodStart time.Time `json:"valid_period_start"`
	ValidPeriodEnd   time.Time `json:"valid_period_end"`
	PageUrl          string    `json:"page_url"`
}

// GetAllAds 从数据库获取所有广告的方法
// pageUrl 参数用于筛选特定页面的广告，如果为空则返回所有广告
func GetAllAds(pageUrl string) []Ad {
	// 获取共享的gorm连接实例
	gormDB := storage.GetGormDB()
	if gormDB == nil {
		log.Println("GORM连接为空")
		return []Ad{}
	}

	// 获取表前缀
	cfg := config.GetDatabaseConfig()
	tablePrefix := cfg.GetTablePrefix()
	tableName := tablePrefix + "ad"

	// 使用gorm查询数据
	var ads []Ad
	query := gormDB.Table(tableName).Order("sort ASC")
	
	// 如果提供了pageUrl参数，则添加条件筛选
	if pageUrl != "" {
		query = query.Where("page_url = ?", pageUrl)
	}
	
	result := query.Find(&ads)
	if result.Error != nil {
		log.Printf("查询广告数据失败: %v", result.Error)
		return []Ad{}
	}

	return ads
}
