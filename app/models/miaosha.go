package models

import (
	"errors"
	"log"
	"myapi/app/storage"
	"myapi/app/structs"
	"strings"
	"time"

	"gorm.io/gorm"
)

// GetMiaoshaGoodsByID 根据ID获取秒杀商品信息
func GetMiaoshaGoodsByID(id int) *structs.MiaoshaGoods {
	// 获取共享的gorm连接实例
	gormDB := storage.GetGormDB()
	if gormDB == nil {
		log.Println("GORM连接为空")
		return nil
	}

	miaoshaGoods := &structs.MiaoshaGoods{}
	result := gormDB.Table(miaoshaGoodsTableName).
		Where("id = ?", id).
		First(miaoshaGoods)

	if result.Error != nil {
		log.Printf("根据ID查询秒杀商品失败: %v", result.Error)
		return nil
	}

	return miaoshaGoods
}

// ValidateCanBuyMiaosha 验证是否可以购买秒杀商品
func ValidateCanBuyMiaosha(goodsID int64, msID int64, quantity int64, sku ...string) (*structs.MiaoshaGoods, error) {
	// 1. 检查秒杀商品是否存在且有效
	miaoshaGoods := GetMiaoshaGoodsByID(int(msID))
	if miaoshaGoods == nil {
		return nil, errors.New("该秒杀商品不存在")
	}

	// 2. 检查秒杀是否在有效期内
	now := time.Now().Unix()
	if now < int64(miaoshaGoods.BeginDate) {
		return nil, errors.New("秒杀活动还未开始")
	}
	if now > int64(miaoshaGoods.EndDate) {
		return nil, errors.New("秒杀活动已结束")
	}

	// 3. 检查秒杀商品状态是否有效
	if miaoshaGoods.Status != 1 {
		return nil, errors.New("秒杀商品已下架")
	}

	// 4. 检查购买数量是否超过限制
	if miaoshaGoods.BuyLimit < int(quantity) {
		return nil, errors.New("购买数量超过限制")
	}

	if miaoshaGoods.BuyMax > 0 && int(quantity) > miaoshaGoods.BuyMax {
		return nil, errors.New("购买数量超过单次购买上限")
	}

	// 5. 处理SKU信息
	if len(sku) > 0 && sku[0] != "" {
		gormDB := storage.GetGormDB()
		if gormDB == nil {
			return nil, errors.New("数据库连接失败")
		}

		var query *gorm.DB
		query = gormDB.Table(MiaoshaGoodsSkuValueTableName).Where("goods_id = ? AND ms_id = ?", goodsID, msID)

		// 分割sku字符串并添加条件
		skuArray := strings.Split(sku[0], ",")
		for _, s := range skuArray {
			if s != "" {
				query = query.Where("FIND_IN_SET( ?,sku)", s)
			}
		}

		var miaoshaGoodsSkuValue structs.MiaoshaGoodsSkuValue
		result := query.First(&miaoshaGoodsSkuValue)
		if result.Error != nil {
			log.Printf("查询秒杀SKU失败: %v", result.Error)
			return nil, errors.New("未找到匹配的秒杀SKU")
		}

		// 找到了SKU信息，更新价格
		miaoshaGoods.Price = miaoshaGoodsSkuValue.Price

		// 检查库存是否足够
		if miaoshaGoodsSkuValue.Quantity < int(quantity) {
			return nil, errors.New("秒杀商品库存不足")
		}
	}

	return miaoshaGoods, nil
}

// GetMiaoshaGoodsSkuByMsIDAndGoodsID 根据秒杀ID和商品ID获取秒杀SKU信息
func GetMiaoshaGoodsSkuByMsIDAndGoodsID(msID int, goodsID int, sku string) (*structs.MiaoshaGoodsSkuValue, error) {
	// 获取共享的gorm连接实例
	gormDB := storage.GetGormDB()
	if gormDB == nil {
		log.Println("GORM连接为空")
		return nil, errors.New("数据库连接失败")
	}

	var miaoshaGoodsSkuValue structs.MiaoshaGoodsSkuValue
	result := gormDB.Table(MiaoshaGoodsSkuValueTableName).
		Where("ms_id = ? AND goods_id = ? AND sku = ?", msID, goodsID, sku).
		First(&miaoshaGoodsSkuValue)

	if result.Error != nil {
		log.Printf("查询秒杀SKU信息失败: %v", result.Error)
		return nil, result.Error
	}

	return &miaoshaGoodsSkuValue, nil
}
