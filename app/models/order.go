package models

import (
	"errors"
	"log"

	"myapi/app/storage"
	"myapi/app/structs"
)

// CheckMiaoshaMemberBuyMax 检查用户是否超过秒杀商品的购买限制
// uid 用户ID，miaosha 秒杀活动信息
func CheckMiaoshaMemberBuyMax(uid int64, miaosha structs.MiaoshaGoods) (int, error) {
	// 获取共享的gorm连接实例
	gormDB := storage.GetGormDB()
	if gormDB == nil {
		log.Println("GORM连接为空")
		return 0, errors.New("数据库连接失败")
	}

	// 查询用户购买该秒杀商品的订单数量
	var orderCount int64
	result := gormDB.Table(orderTableName).
		Where("uid = ?", uid).
		Where("ms_id = ?", miaosha.ID).
		Count(&orderCount)

	if result.Error != nil {
		log.Printf("查询用户秒杀订单数量失败: %v", result.Error)
		return 0, result.Error
	}

	// 检查是否超过购买限制
	if orderCount >= int64(miaosha.MemberBuyMax) {
		return 1, nil // 超过限制，返回1
	} else {
		return 0, nil // 未超过限制，返回0
	}
}
