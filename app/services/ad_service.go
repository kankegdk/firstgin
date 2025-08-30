package services

import (
	"myapi/app/models"
)

// AdService 广告服务接口
type AdService interface {
	GetAllAds(pageUrl string) []models.Ad
}

// adService 实现AdService接口的结构体
type adService struct{}

// NewAdService 创建一个新的广告服务实例
func NewAdService() AdService {
	return &adService{}
}

// GetAllAds 获取所有广告
func (s *adService) GetAllAds(pageUrl string) []models.Ad {
	// 调用模型层方法，并传递pageUrl参数
	return models.GetAllAds(pageUrl)
}