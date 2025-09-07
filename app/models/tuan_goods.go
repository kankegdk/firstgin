package models

import (
	"errors"
	"log"
	"myapi/app/storage"
	"myapi/app/structs"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
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

// GetTuanGoodsByGoodsIDAndTuanID 根据商品ID和团购ID获取团购商品
func GetTuanGoodsByGoodsIDAndTuanID(goodsID int, tuanID int) *structs.TuanGoods {
	// 获取共享的gorm连接实例
	gormDB := storage.GetGormDB()
	if gormDB == nil {
		log.Println("GORM连接为空")
		return nil
	}

	tuanGoods := &structs.TuanGoods{}
	result := gormDB.Table(tuanGoodsTableName).
		Where("id = ? AND goods_id = ? AND status = 1", tuanID, goodsID).
		First(tuanGoods)

	if result.Error != nil {
		log.Printf("根据商品ID和团购ID查询团购商品失败: %v", result.Error)
		return nil
	}

	return tuanGoods
}

// GetTuanGoodsByGoodsID 根据商品ID获取团购商品（保留原有函数用于兼容性）
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
		Where("found_id = ? ", foundID).
		Where("status IN (?,?,?)", 0, 1, 2).
		Count(&count)

	if result.Error != nil {
		return 0, result.Error
	}

	return count, nil
}

// ValidateCanJoinTuan 验证是否可平团
func ValidateCanJoinTuan(goodsID int64, tuanID int64, quantity int64, sku ...string) (*structs.TuanGoods, error) {
	// 1. 检查团购商品是否存在且有效
	tuanGoods := GetTuanGoodsByGoodsIDAndTuanID(int(goodsID), int(tuanID))
	if tuanGoods == nil {
		return nil, errors.New("该商品未开团或已下架")
	}
	if tuanGoods.BuyLimit < int(quantity) {
		return nil, errors.New("该团购已结束或无效，最多购买" + strconv.Itoa(tuanGoods.BuyLimit) + "件")
	}

	// 2. 检查团购是否在有效期内
	now := time.Now().Unix()
	if now < int64(tuanGoods.BeginDate) || now > int64(tuanGoods.EndDate) {
		return nil, errors.New("团购不在有效期内")
	}

	// 3. 检查开团记录是否存在
	// tuanFound := GetTuanFoundByID(tuanID)
	// if tuanFound == nil {
	// 	return nil, errors.New("开团记录不存在")
	// }

	// // 4. 检查开团是否有效
	// if tuanFound.Status != 1 {
	// 	return nil, errors.New("开团已结束或无效")
	// }

	// 5. 检查是否已经人满
	joinCount, err := GetTuanFollowCountByFoundID(int(tuanID))
	if err != nil {
		return nil, errors.New("获取参团人数失败")
	}
	if joinCount >= int64(tuanGoods.PeopleNum) {
		return nil, errors.New("该团已经满员")
	}

	// 4. 检查购买数量是否超过限制
	if tuanGoods.BuyLimit < int(quantity) {
		return nil, errors.New("该团购已结束或无效，最多购买" + strconv.Itoa(tuanGoods.BuyLimit) + "件")
	}
	if tuanGoods.BuyMax > 0 && int(quantity) > tuanGoods.BuyMax {
		return nil, errors.New("购买数量超过限制" + strconv.Itoa(tuanGoods.BuyMax))
	}

	if len(sku) > 0 && sku[0] != "" {
		gormDB := storage.GetGormDB()
		if gormDB == nil {
			return nil, errors.New("GORM连接为空")
		}
		var tuanGoodsSkuValue structs.TuanGoodsSkuValue
		var query *gorm.DB

		// 确保第一个元素不是空字符串
		if sku[0] != "" {
			log.Println("开始检查SKU值: ", sku[0])
			query = gormDB.Table(tuanGoodsSkuValueTableName).Where("goods_id = ? AND tuan_id = ?", goodsID, tuanID)

			// 分割sku字符串并添加条件
			skuValues := strings.Split(sku[0], ",")
			for _, vo := range skuValues {
				if vo != "" {
					query = query.Where("FIND_IN_SET(?, sku)", vo)
				}
			}

			result := query.First(&tuanGoodsSkuValue)
			if result.Error != nil {
				return nil, errors.New("未找到对应的SKU信息")
			}

			// 找到了SKU信息，更新价格
			tuanGoods.Price = tuanGoodsSkuValue.Price

			// 检查库存是否足够
			if tuanGoodsSkuValue.Quantity < int(quantity) {
				return nil, errors.New("SKU库存不足")
			}
		} else {
			log.Println("SKU值为空字符串，跳过SKU检查")
		}

	}

	// 所有验证通过
	return tuanGoods, nil
}
