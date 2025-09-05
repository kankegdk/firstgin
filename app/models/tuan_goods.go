package models

import (
	"errors"
	"log"
	"myapi/app/storage"
	"myapi/app/structs"
	"time"
)

// GetTuanGoodsByID 根据ID获取团购商品
func GetTuanGoodsByID(id int) *structs.TuanGoods {
	// 获取共享的gorm连接实例
	gormDB := storage.GetGormDB()
	if gormDB == nil {
		log.Println("GORM连接为空")
		return nil
	}

	tuanGoods := &structs.TuanGoods{}
	result := gormDB.Table(tuanGoodsTableName).
		Where("id = ?", id).
		First(tuanGoods)

	if result.Error != nil {
		log.Printf("根据ID查询团购商品失败: %v", result.Error)
		return nil
	}

	return tuanGoods
}

// GetTuanGoodsByGoodsID 根据商品ID获取团购商品
func GetTuanGoodsByGoodsID(goodsID int) *structs.TuanGoods {
	// 获取共享的gorm连接实例
	gormDB := storage.GetGormDB()
	if gormDB == nil {
		log.Println("GORM连接为空")
		return nil
	}

	tuanGoods := &structs.TuanGoods{}
	result := gormDB.Table(tuanGoodsTableName).
		Where("goods_id = ? AND status = 1", goodsID).
		First(tuanGoods)

	if result.Error != nil {
		log.Printf("根据商品ID查询团购商品失败: %v", result.Error)
		return nil
	}

	return tuanGoods
}

// GetTuanFoundByID 根据ID获取开团记录
func GetTuanFoundByID(id int) *structs.TuanFound {
	// 获取共享的gorm连接实例
	gormDB := storage.GetGormDB()
	if gormDB == nil {
		log.Println("GORM连接为空")
		return nil
	}

	tuanFound := &structs.TuanFound{}
	result := gormDB.Table(tuanFoundTableName).
		Where("id = ?", id).
		First(tuanFound)

	if result.Error != nil {
		log.Printf("根据ID查询开团记录失败: %v", result.Error)
		return nil
	}

	return tuanFound
}

// GetTuanFollowCountByFoundID 根据开团ID获取参团人数
func GetTuanFollowCountByFoundID(foundID int) (int64, error) {
	// 获取共享的gorm连接实例
	gormDB := storage.GetGormDB()
	if gormDB == nil {
		return 0, errors.New("GORM连接为空")
	}

	var count int64
	result := gormDB.Table(tuanFollowTableName).
		Where("found_id = ? AND status = 1", foundID).
		Count(&count)

	if result.Error != nil {
		return 0, result.Error
	}

	return count, nil
}

// ValidateCanJoinTuan 验证是否可平团
func ValidateCanJoinTuan(goodsID int, tuanID int, quantity int) error {
	// 1. 检查团购商品是否存在且有效
	tuanGoods := GetTuanGoodsByGoodsID(goodsID)
	if tuanGoods == nil {
		return errors.New("该商品未开团或已下架")
	}

	// 2. 检查团购是否在有效期内
	now := time.Now().Unix()
	if now < int64(tuanGoods.BeginDate) || now > int64(tuanGoods.EndDate) {
		return errors.New("团购不在有效期内")
	}

	// 3. 检查开团记录是否存在
	tuanFound := GetTuanFoundByID(tuanID)
	if tuanFound == nil {
		return errors.New("开团记录不存在")
	}

	// 4. 检查开团是否有效
	if tuanFound.Status != 1 {
		return errors.New("开团已结束或无效")
	}

	// 5. 检查是否已经人满
	joinCount, err := GetTuanFollowCountByFoundID(tuanID)
	if err != nil {
		return errors.New("获取参团人数失败")
	}
	if joinCount >= int64(tuanGoods.PeopleNum) {
		return errors.New("该团已经满员")
	}

	// 6. 检查购买数量是否超过限制
	if quantity <= 0 {
		return errors.New("购买数量必须大于0")
	}
	if tuanGoods.BuyMax > 0 && quantity > tuanGoods.BuyMax {
		return errors.New("购买数量超过限制")
	}

	// 所有验证通过
	return nil
}