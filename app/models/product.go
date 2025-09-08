package models

import (
	"errors"
	"log"
	"math"
	"myapi/app/helper"
	"myapi/app/storage"
	"myapi/app/structs"
	"strings"
	"time"
)

// GetProductByID 根据ID获取产品
func GetProductByID(id int) *structs.Product {
	// 获取共享的gorm连接实例
	gormDB := storage.GetGormDB()
	if gormDB == nil {
		log.Println("GORM连接为空")
		return nil
	}

	product := &structs.Product{}
	result := gormDB.Table(goodsTableName).
		Where("id = ?", id).
		First(product)

	if result.Error != nil {
		log.Printf("根据ID查询产品失败: %v", result.Error)
		return nil
	}

	return product
}

// GetAllProducts 获取所有产品
func GetAllProducts() []structs.Product {
	// 获取共享的gorm连接实例
	gormDB := storage.GetGormDB()
	if gormDB == nil {
		log.Println("GORM连接为空")
		return nil
	}

	var products []structs.Product
	result := gormDB.Table(goodsTableName).
		Where("status = ?", 1).
		Find(&products)

	if result.Error != nil {
		log.Printf("查询所有产品失败: %v", result.Error)
		return nil
	}

	return products
}

// GetCategoryByID 根据ID获取分类信息
func GetCategoryByID(id int) (*structs.Category, error) {
	gormDB := storage.GetGormDB()
	if gormDB == nil {
		return nil, errors.New("数据库连接失败")
	}

	category := &structs.Category{}
	result := gormDB.Table(categoryTableName).
		Where("id = ?", id).
		First(category)

	if result.Error != nil {
		log.Printf("查询分类信息失败: %v", result.Error)
		return nil, result.Error
	}

	return category, nil
}

// CreateGoodsBuynowinfo 创建商品立即购买信息
func CreateGoodsBuynowinfo(weid int, ip string, data string) (int, error) {
	gormDB := storage.GetGormDB()
	if gormDB == nil {
		return 0, errors.New("数据库连接失败")
	}

	buynowinfo := &structs.GoodsBuynowinfo{
		Weid:       weid,
		Ip:         ip,
		ExpireTime: time.Now().Unix(),
		Data:       data,
		Status:     1,
	}

	result := gormDB.Table(goodsBuynowinfoTableName).Create(buynowinfo)
	if result.Error != nil {
		log.Printf("创建立即购买信息失败: %v", result.Error)
		return 0, result.Error
	}

	return int(buynowinfo.ID), nil
}

// CartGoods 获取购物车商品信息
func CartGoods(params map[string]interface{}) (map[string]interface{}, error) {
	log.Println("开始执行CartGoods函数，参数:", params)
	gormDB := storage.GetGormDB()
	if gormDB == nil {
		log.Println("数据库连接失败，直接返回错误")
		return nil, errors.New("数据库连接失败")
	}

	sku := params["sku"]

	extraData := make(map[string]interface{})
	result := gormDB.Table(goodsTableName).
		Where("id = ? ", params["GoodsID"]).
		Limit(1).
		Scan(&extraData)

	if result.Error != nil {
		log.Printf("查询商品信息失败: %v，直接返回错误", result.Error)
		return nil, result.Error
	}
	// 检查查询是否返回了数据
	if len(extraData) == 0 {
		log.Printf("未找到ID为%d的商品信息，直接返回错误", params["GoodsID"])
		return nil, errors.New("未找到商品信息")
	}

	var tuanid = params["tuanid"].(int64)
	var msid = params["msid"].(int64)

	// 安全地进行类型断言
	skuStr, ok := sku.(string)
	if !ok && sku != nil {
		log.Println("SKU参数类型错误，无法转换为字符串")
		return nil, errors.New("SKU参数无效")
	}

	// log.Printf("确认商品信息存在，准备处理各种活动信息")
	extraData["quantity"] = params["quantity"]
	//如果没有团购，秒杀，但是又sku
	if skuStr != "" && (tuanid == 0 && msid == 0) {
		// 检查SKU是否存在
		var SkuValueresult structs.GoodsSkuValue
		var query = gormDB.Table(goodsSkuValueTableName).
			Where("goods_id = ? ", params["GoodsID"])
		// 分割sku字符串并添加条件
		skuValues := strings.Split(skuStr, ",")
		for _, vo := range skuValues {
			if vo != "" {
				query = query.Where("FIND_IN_SET(?, sku)", vo)
			}
		}

		result := query.First(&SkuValueresult)
		//如果有就替换价格
		if result.Error == nil {
			extraData["price"] = SkuValueresult.Price
			extraData["skuvid"] = SkuValueresult.ID
			log.Println("SKU存在，价格:", SkuValueresult.Price)
		}
	}
	log.Println("团购价格处理前的价格: ,", extraData["price"], "数量:", extraData["quantity"])

	// 处理团购信息
	if tuanid > 0 {
		// 正确传递字符串给可变参数函数
		tuanGoods, err := ValidateCanJoinTuan(params["GoodsID"].(int64), int64(tuanid), params["quantity"].(int64), skuStr)
		if err != nil {
			log.Printf("团购验证失败: %v，直接返回错误", err)
			return nil, err
		}
		extraData["price"] = tuanGoods.Price
		extraData["tuanid"] = tuanid
		log.Println("团购价格处理后的价格: ", extraData["price"])
		//如果是团购就返回团购相关信息
		return extraData, nil
	} else {
		log.Println("未提供有效的团购ID参数，跳过团购信息处理")
	}

	//如果有秒杀
	if params["msid"].(int64) > 0 {
		// 验证是否可以购买秒杀商品
		miaoshaGoods, err := ValidateCanBuyMiaosha(params["GoodsID"].(int64), params["msid"].(int64), params["quantity"].(int64), skuStr)
		if err != nil {
			log.Printf("秒杀验证失败: %v，直接返回错误", err)
			return nil, err
		}
		// 设置秒杀价格
		extraData["price"] = miaoshaGoods.Price
		log.Println("秒杀价格处理后的价格: ", extraData["price"])
		// 存储秒杀ID
		extraData["ms_id"] = params["msid"]
		//如果是秒杀就返回秒杀相关信息
		return extraData, nil
	} else {
		log.Println("未提供有效的秒杀ID参数，跳过秒杀信息处理")
	}

	//折扣价格处理
	var discounts []structs.GoodsDiscount
	discountResult := gormDB.Table(goodsDiscountTableName).
		Where("goods_id = ? AND quantity <= ?", params["GoodsID"], extraData["quantity"]).
		Order("quantity DESC, price ASC").
		Find(&discounts)
	if discountResult.Error != nil {
		log.Printf("查询数量折扣失败: %v", discountResult.Error)
	} else {
		if len(discounts) > 0 {
			// 假设price字段存储的是百分比，例如95表示95%折扣
			originalPrice, _ := helper.ToFloat64(extraData["price"]) // 显式转换为float64类型

			extraData["price"] = math.Round(originalPrice*discounts[0].Price/100*100) / 100

			log.Printf("应用数量折扣成功: 原价%.2f，折扣率%.2f，折后价%.2f", originalPrice, discounts[0].Price, extraData["price"].(float64))
		} else {
			log.Printf("未找到适用的数量折扣，保持原价")
		}
	}

	// 安全检查Points字段
	pointsValue, pointsExists := extraData["points"]
	if pointsExists && pointsValue != nil {
		log.Printf("检测到积分字段，当前值: %v", pointsValue)
		if p, ok := pointsValue.(int); ok && p < 0 {
			extraData["Points"] = 0
			log.Printf("积分值为负数(%d)，已纠正为0", p)
		} else if p, ok := pointsValue.(float64); ok && p < 0 {
			extraData["Points"] = 0
			log.Printf("积分值为负数(%.2f)，已纠正为0", p)
		} else {
			log.Printf("积分值有效，无需纠正")
		}
	} else {
		log.Println("积分字段不存在或为空，跳过安全检查")
	}

	// 计算总价并添加到extraData中
	price, _ := helper.ToFloat64(extraData["price"])
	quantity, _ := helper.ToFloat64(extraData["quantity"])
	extraData["totalprice"] = price * quantity
	log.Printf("计算商品总价: 单价%.2f * 数量%.2f = %.2f", price, quantity, extraData["totalprice"].(float64))

	// 安全获取PayPoints
	log.Printf("开始计算支付积分%v", extraData)
	if payPointsValue, ok := extraData["pay_points"]; ok && payPointsValue != nil {
		log.Printf("尝试获取支付积分值: %v", payPointsValue)
		p, _ := helper.ToFloat64(payPointsValue)
		totalPayPoints := p * quantity
		extraData["PayPoints"] = totalPayPoints
		log.Printf("支付积分为int类型，计算总支付积分: %.2f * %.2f = %.2f", p, quantity, totalPayPoints)

	} else {
		log.Println("支付积分字段不存在或为空，保持默认值")
	}

	// 安全获取PointsPrice
	if pointsPriceValue, ok := extraData["points_price"]; ok && pointsPriceValue != nil {
		log.Printf("尝试获取积分价格值: %v", pointsPriceValue)
		p, _ := helper.ToFloat64(pointsPriceValue)
		totalPointsPrice := p * quantity
		extraData["PointsPrice"] = totalPointsPrice
		log.Printf("积分价格为int类型，计算总积分价格: %.2f * %.2f = %.2f", p, quantity, totalPointsPrice)

	} else {
		log.Println("积分价格字段不存在或为空，保持默认值")
	}

	return nil, nil

	/*
		total := price * quantity
		totalPayPoints := 0
		totalPointsPrice := 0
		log.Printf("计算商品总价: 单价%.2f * 数量%d = %.2f", price, quantityInt, total)

		// 更新商品重量 - 安全获取字段值
		if weightValue, ok := extraData["Weight"]; ok && weightValue != nil {
			log.Printf("尝试获取商品重量值: %v", weightValue)
			if w, ok := weightValue.(float64); ok {
				weightTotal := w * float64(extraData["quantity"].(int))
				extraData["Weight"] = weightTotal
				log.Printf("重量为float64类型，计算总重量: %.2f * %d = %.2f", w, extraData["quantity"].(int), weightTotal)
			} else if w, ok := weightValue.(int); ok {
				weightTotal := float64(w) * float64(extraData["quantity"].(int))
				extraData["Weight"] = weightTotal
				log.Printf("重量为int类型，计算总重量: %d * %d = %.2f", w, extraData["quantity"].(int), weightTotal)
			} else {
				log.Printf("重量类型不支持，不更新重量信息")
			}
		} else {
			log.Println("重量字段不存在或为空，不更新重量信息")
		}

		// 设置额外数据
		log.Println("设置商品额外数据到返回结果")
		extraData["price"] = price
		extraData["total"] = total
		extraData["totalPayPoints"] = totalPayPoints
		extraData["totalPointsPrice"] = totalPointsPrice
		extraData["totalReturnPoints"] = totalReturnPoints
		extraData["stock"] = stock
		extraData["stores"] = stores
		extraData["label"] = label
		extraData["sku"] = sku
		extraData["is_skumore"] = isSkumore
		extraData["skumorequantity"] = skumorequantity

		log.Printf("额外数据已设置: 总价=%.2f，总支付积分=%d，总积分价格=%d，总返积分=%.2f，库存=%d", total, totalPayPoints, totalPointsPrice, totalReturnPoints, stock)

		// 处理积分商品 - 安全获取字段值
		isPointsGoods := 0
		log.Println("开始处理积分商品标记")
		if isPointsValue, ok := extraData["IsPointsGoods"]; ok && isPointsValue != nil {
			log.Printf("尝试获取积分商品标记值: %v", isPointsValue)
			if ipg, ok := isPointsValue.(int); ok {
				isPointsGoods = ipg
				log.Printf("积分商品标记为int类型，值为: %d", isPointsGoods)
			} else if ipg, ok := isPointsValue.(float64); ok {
				isPointsGoods = int(ipg)
				log.Printf("积分商品标记为float64类型，转换为int后的值为: %d", isPointsGoods)
			} else {
				log.Printf("积分商品标记类型不支持，保持默认值: %d", isPointsGoods)
			}
		} else {
			log.Println("积分商品标记字段不存在或为空，保持默认值")
		}

		if isPointsGoods == 1 {
			log.Println("该商品为积分商品，将总价和积分价格设置为0")
			extraData["total"] = 0.0
			extraData["totalPointsPrice"] = 0.0
		} else {
			log.Println("该商品不是积分商品，保持原价格设置")
		}

		log.Println("CartGoods函数执行完成，返回商品信息")
		return extraData, nil
	*/
}
