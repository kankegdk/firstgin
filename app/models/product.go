package models

import (
	"errors"
	"fmt"
	"log"
	"math"
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

	// 处理商品价格和库存，确保类型安全
	price := 0.0
	stores := 0
	stock := 0
	label := ""
	sku := params["sku"]
	isSkumore := 0

	// 获取参数
	id, ok := params["id"].(int)
	if !ok {
		log.Println("商品ID参数无效或类型错误，直接返回错误")
		return nil, errors.New("商品ID无效")
	}
	extraData := make(map[string]interface{})
	result := gormDB.Table(goodsTableName).
		Where("id = ?", id).
		Limit(1).
		Scan(&extraData)

	if result.Error != nil {
		log.Printf("查询商品信息失败: %v，直接返回错误", result.Error)
		return nil, result.Error
	}
	// 检查查询是否返回了数据
	if len(extraData) == 0 {
		log.Printf("未找到ID为%d的商品信息，直接返回错误", id)
		return nil, errors.New("未找到商品信息")
	}
	log.Printf("确认商品信息存在，准备处理各种活动信息")
	extraData["quantity"] = params["quantity"]

	// 处理团购信息
	if tuanid, ok := params["tuanid"].(int64); ok && tuanid > 0 {
		log.Printf("检测到团购ID参数: %v，开始查询团购信息", tuanid)

		var tuan structs.TuanGoods
		tuanResult := gormDB.Table(tuanGoodsTableName).
			Where("id = ?", tuanid).
			Where("goods_id = ?", id).
			First(&tuan)
		if tuanResult.Error == nil {
			log.Printf("成功获取团购信息，将团购信息添加到返回结果中")
			extraData["tuan"] = tuan
		} else {
			log.Printf("查询团购信息失败: %v，不添加团购信息", tuanResult.Error)
			return nil, errors.New("团购信息不存在")
		}
	} else {
		log.Println("未提供有效的团购ID参数，跳过团购信息处理")
	}

	// 处理秒杀信息
	if msid, ok := params["msid"].(int64); ok && msid > 0 {
		log.Printf("检测到秒杀ID参数: %v，开始查询秒杀信息", msid)
		var miaosha structs.MiaoshaGoods
		miaoshaResult := gormDB.Table(miaoshaGoodsTableName).
			Where("id = ?", msid).
			Where("goods_id = ?", id).
			First(&miaosha)
		if miaoshaResult.Error == nil {
			log.Printf("成功获取秒杀信息，检查秒杀购买限制")
			// 在实际项目中，需要实现检查秒杀购买限制的逻辑
			if miaosha.MemberBuyMax > 0 {
				log.Printf("秒杀商品设置了购买限制: %d", miaosha.MemberBuyMax)
				isMemberBuyMax, err := CheckMiaoshaMemberBuyMax(int64(params["Uid"].(float64)), miaosha)
				if err != nil {
					log.Printf("检查秒杀购买限制失败: %v，直接返回错误", err)
					return nil, err
				}
				log.Printf("秒杀购买限制检查结果: %v，将结果添加到返回数据", isMemberBuyMax)
				extraData["is_member_buy_max"] = isMemberBuyMax
			}
			log.Println("将秒杀信息添加到返回结果中")
			extraData["miaosha"] = miaosha
		} else {
			log.Printf("查询秒杀信息失败: %v，不添加秒杀信息", miaoshaResult.Error)
			return nil, errors.New("秒杀信息不存在")
		}
	} else {
		log.Println("未提供有效的秒杀ID参数，跳过秒杀信息处理")
	}

	extraData["quantity"] = params["quantity"]
	log.Printf("购买数量已设置到返回结果: %v", extraData["quantity"])

	// 处理SKU信息
	if sku, ok := params["sku"].(string); ok && sku != "" {
		log.Printf("检测到SKU参数: %s，开始查询SKU详情", sku)
		var goodsSkuValue structs.GoodsSkuValue
		skuResult := gormDB.Table(goodsSkuTableName).
			Where("goods_id = ? AND sku = ?", id, sku).
			First(&goodsSkuValue)
		if skuResult.Error == nil {
			log.Printf("成功获取SKU信息，更新价格、库存和图片")
			extraData["Price"] = goodsSkuValue.Price
			extraData["quantity"] = goodsSkuValue.Quantity
			if goodsSkuValue.Image != "" {
				log.Printf("SKU有自定义图片，更新图片信息")
				extraData["Image"] = goodsSkuValue.Image
			} else {
				log.Println("SKU没有自定义图片，保持原有图片")
			}
		} else {
			log.Printf("查询SKU信息失败: %v，不更新SKU相关信息", skuResult.Error)
		}
	} else {
		log.Println("未提供有效的SKU参数，跳过SKU详情处理")
	}

	// 处理多SKU信息
	if isSkumore, ok := params["is_skumore"]; ok {
		log.Printf("检测到多SKU标记: %v，添加到返回结果", isSkumore)
		extraData["is_skumore"] = isSkumore
	} else {
		log.Println("未提供多SKU标记，跳过")
	}
	if skumore, ok := params["skumore"]; ok {
		log.Printf("检测到多SKU详情数据，添加到返回结果")
		extraData["skumore"] = skumore
	} else {
		log.Println("未提供多SKU详情数据，跳过")
	}

	log.Println("开始处理商品价格和库存信息，初始化相关变量")

	// 安全地获取数量字段
	if quantityValue, ok := extraData["quantity"]; ok && quantityValue != nil {
		log.Printf("尝试获取数量字段值: %v", quantityValue)
		if q, ok := quantityValue.(int); ok {
			stores = q
			stock = q
			log.Printf("数量字段为int类型，库存设置为: %d", stores)
		} else if q, ok := quantityValue.(float64); ok {
			stores = int(q)
			stock = int(q)
			log.Printf("数量字段为float64类型，转换为int后库存设置为: %d", stores)
		} else {
			log.Printf("数量字段类型不支持，保持默认值: %d", stores)
		}
	} else {
		log.Println("数量字段不存在或为空，保持默认值")
	}
	skumorequantity := 0
	log.Printf("库存信息处理完成，当前库存: %d", stores)

	// 处理团购价格
	if tuan, ok := extraData["tuan"]; ok && tuan != nil {
		log.Printf("检测到团购信息，开始处理团购价格")
		tuanData := tuan.(structs.TuanGoods)
		if tuanData.Price > 0 {
			price = tuanData.Price
			log.Printf("团购价格设置成功: %f", price)
		} else {
			log.Printf("团购价格无效或为0，保持原有价格")
		}
	} else {
		log.Println("没有团购信息，跳过团购价格处理")
	}

	// 处理秒杀价格
	if miaosha, ok := extraData["miaosha"]; ok && miaosha != nil {
		log.Printf("检测到秒杀信息，开始处理秒杀价格")
		miaoshaData := miaosha.(structs.MiaoshaGoods)
		if miaoshaData.Price > 0 {
			price = miaoshaData.Price
			log.Printf("秒杀价格设置成功: %f", price)
		} else {
			log.Printf("秒杀价格无效或为0，保持原有价格")
		}
	} else {
		log.Println("没有秒杀信息，跳过秒杀价格处理")
	}

	// 处理多SKU情况
	if isSkumoreParam, ok := params["is_skumore"]; ok && isSkumoreParam != "" {
		log.Printf("检测到多SKU标记参数: %v", isSkumoreParam)
		isSkumoreStr := fmt.Sprintf("%v", isSkumoreParam)
		if isSkumoreStr == "1" {
			log.Println("商品为多SKU类型，重置相关参数")
			isSkumore = 1
			price = 0
			stores = 0
			skumorequantity = 0
			label = ""

			// 处理skumore参数
			if skumore, ok := params["skumore"]; ok && skumore != "" {
				log.Printf("检测到多SKU详情参数，添加到返回结果")
				extraData["skumore"] = skumore
			} else {
				log.Println("未提供多SKU详情参数，不添加skumore字段")
			}
		} else {
			log.Printf("多SKU标记不为1，不进行多SKU处理")
		}
	} else if skuParam, ok := params["sku"]; ok && fmt.Sprintf("%v", skuParam) != "" {
		log.Printf("检测到SKU参数: %v，设置SKU和label", skuParam)
		sku = fmt.Sprintf("%v", skuParam)
	} else {
		log.Println("既没有多SKU标记也没有SKU参数，跳过特殊SKU处理")
	}

	// 根据不同活动类型查询对应SKU
	var goodsSkuValue structs.GoodsSkuValue
	query := gormDB.Table(goodsSkuTableName).Where("goods_id = ?", id)
	log.Println("开始根据活动类型查询对应SKU信息")

	if tuan && params["sku"] != "" {
		log.Printf("检测到团购ID参数: %v，查询团购SKU", tuanid)
		// 查询团购SKU
		var tuanSkuValue structs.TuanGoodsSkuValue
		query = gormDB.Table(TuanGoodsSkuValueTableName).Where("goods_id = ? AND tuan_id = ?", id, tuanid)
		tuanSkuResult := query.First(&tuanSkuValue)
		if tuanSkuResult.Error == nil {
			log.Printf("成功获取团购SKU信息，更新库存: %d，价格: %f", tuanSkuValue.Quantity, tuanSkuValue.Price)
			stores = tuanSkuValue.Quantity
			price = tuanSkuValue.Price
			if tuanSkuValue.Image != "" {
				log.Printf("团购SKU有自定义图片，更新图片信息")
				extraData["Image"] = tuanSkuValue.Image
			}
		} else {
			log.Printf("查询团购SKU失败: %v", tuanSkuResult.Error)
		}
	}
	if miaosha, ok := extraData["miaosha"]; ok && miaosha != nil && params["sku"] != "" {
		log.Printf("检测到秒杀ID参数: %v，查询秒杀SKU", extraData["miaosha"])
		// 查询秒杀SKU
		var miaoshaSkuValue structs.MiaoshaGoodsSkuValue
		query = gormDB.Table(MiaoshaGoodsSkuValueTableName).Where("goods_id = ? AND ms_id = ?", id, params["misd"])
		miaoshaSkuResult := query.First(&miaoshaSkuValue)
		if miaoshaSkuResult.Error == nil {
			log.Printf("成功获取秒杀SKU信息，更新库存: %d，价格: %f", miaoshaSkuValue.Quantity, miaoshaSkuValue.Price)
			stores = miaoshaSkuValue.Quantity
			price = miaoshaSkuValue.Price
			if miaoshaSkuValue.Image != "" {
				log.Printf("秒杀SKU有自定义图片，更新图片信息")
				extraData["Image"] = miaoshaSkuValue.Image
			}
		} else {
			log.Printf("查询秒杀SKU失败: %v", miaoshaSkuResult.Error)
		}
	} else {
		log.Println("没有团购或秒杀ID，查询普通SKU")
		// 查询普通SKU
		// 分割sku字符串并构建查询条件
		skuArray := strings.Split(sku, ",")
		for _, s := range skuArray {
			if s != "" {
				query = query.Where("FIND_IN_SET(sku, ?)", s)
			}
		}
		skuResult := query.First(&goodsSkuValue)
		if skuResult.Error == nil {
			log.Printf("成功获取普通SKU信息，更新库存: %d，价格: %f", goodsSkuValue.Quantity, goodsSkuValue.Price)
			stores = goodsSkuValue.Quantity
			price = goodsSkuValue.Price
			if goodsSkuValue.Image != "" {
				log.Printf("普通SKU有自定义图片，更新图片信息")
				extraData["Image"] = goodsSkuValue.Image
			}
		} else {
			log.Printf("查询普通SKU失败: %v", skuResult.Error)
		}
	}
	log.Printf("SKU查询处理完成，当前库存: %d，价格: %f", stores, price)

	// 应用数量折扣

	var discounts []structs.GoodsDiscount

	discountResult := gormDB.Table(goodsDiscountTableName).
		Where("goods_id = ? AND quantity >= ?", id, extraData["quantity"]).
		Order("quantity DESC, price ASC").
		Find(&discounts)
	if discountResult.Error != nil {
		log.Printf("查询数量折扣失败: %v", discountResult.Error)
	} else {
		if len(discounts) > 0 {
			// 假设price字段存储的是百分比，例如95表示95%折扣
			originalPrice := price
			price = math.Round(price*discounts[0].Price*10) / 100
			log.Printf("应用数量折扣成功: 原价%.2f，折扣率%.2f，折后价%.2f", originalPrice, discounts[0].Price, price)
		} else {
			log.Printf("未找到适用的数量折扣，保持原价")
		}
	}

	// 处理积分相关字段 - 安全检查
	log.Println("开始处理积分相关字段安全检查")
	pointsValue, pointsExists := extraData["Points"]
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

	// 设置商品基本信息
	log.Println("设置商品基本信息到返回结果")
	extraData["GoodsID"] = id
	extraData["Quantity"] = extraData["quantity"]
	extraData["Price"] = price
	log.Printf("商品基本信息已设置: ID=%d，数量=%d，价格=%.2f", id, extraData["quantity"], price)

	// 计算总价 - 安全获取字段值
	total := price * float64(extraData["quantity"].(int))
	totalPayPoints := 0
	totalPointsPrice := 0
	log.Printf("计算商品总价: 单价%.2f * 数量%d = %.2f", price, extraData["quantity"].(int), total)

	// 安全获取PayPoints
	if payPointsValue, ok := extraData["PayPoints"]; ok && payPointsValue != nil {
		log.Printf("尝试获取支付积分值: %v", payPointsValue)
		if pp, ok := payPointsValue.(int); ok {
			totalPayPoints = pp * extraData["quantity"].(int)
			log.Printf("支付积分为int类型，计算总支付积分: %d * %d = %d", pp, extraData["quantity"].(int), totalPayPoints)
		} else if pp, ok := payPointsValue.(float64); ok {
			totalPayPoints = int(pp) * extraData["quantity"].(int)
			log.Printf("支付积分为float64类型，计算总支付积分: %d * %d = %d", int(pp), extraData["quantity"].(int), totalPayPoints)
		} else {
			log.Printf("支付积分类型不支持，保持默认值: %d", totalPayPoints)
		}
	} else {
		log.Println("支付积分字段不存在或为空，保持默认值")
	}

	// 安全获取PointsPrice
	if pointsPriceValue, ok := extraData["PointsPrice"]; ok && pointsPriceValue != nil {
		log.Printf("尝试获取积分价格值: %v", pointsPriceValue)
		if pp, ok := pointsPriceValue.(int); ok {
			totalPointsPrice = pp * extraData["quantity"].(int)
			log.Printf("积分价格为int类型，计算总积分价格: %d * %d = %d", pp, extraData["quantity"].(int), totalPointsPrice)
		} else if pp, ok := pointsPriceValue.(float64); ok {
			totalPointsPrice = int(pp) * extraData["quantity"].(int)
			log.Printf("积分价格为float64类型，计算总积分价格: %d * %d = %d", int(pp), extraData["quantity"].(int), totalPointsPrice)
		} else {
			log.Printf("积分价格类型不支持，保持默认值: %d", totalPointsPrice)
		}
	} else {
		log.Println("积分价格字段不存在或为空，保持默认值")
	}

	// 计算返积分 - 安全获取字段值
	var totalReturnPoints float64
	pointsValue = extraData["Points"]
	log.Printf("开始计算返积分")
	if pointsValue != nil {
		var points int
		log.Printf("尝试获取积分值: %v", pointsValue)
		if pInt, ok := pointsValue.(int); ok {
			points = pInt
			log.Printf("积分值为int类型，值为: %d", points)
		} else if pFloat, ok := pointsValue.(float64); ok {
			points = int(pFloat)
			log.Printf("积分值为float64类型，转换为int后的值为: %d", points)
		} else {
			log.Printf("积分值类型不支持，保持默认值: %d", points)
			return extraData, nil
		}

		if points > 0 {
			log.Printf("积分值大于0，准备计算返积分")
			// 安全获取PointsMethod
			pointsMethod := 0
			if methodValue, ok := extraData["PointsMethod"]; ok && methodValue != nil {
				log.Printf("尝试获取积分方法值: %v", methodValue)
				if mInt, ok := methodValue.(int); ok {
					pointsMethod = mInt
					log.Printf("积分方法为int类型，值为: %d", pointsMethod)
				} else if mFloat, ok := methodValue.(float64); ok {
					pointsMethod = int(mFloat)
					log.Printf("积分方法为float64类型，转换为int后的值为: %d", pointsMethod)
				} else {
					log.Printf("积分方法类型不支持，保持默认值: %d", pointsMethod)
				}
			} else {
				log.Println("积分方法字段不存在或为空，保持默认值")
			}

			if pointsMethod == 1 {
				log.Printf("积分方法为按数量计算，每件商品返%d积分", points)
				totalReturnPoints = float64(points * extraData["quantity"].(int))
				log.Printf("按数量计算返积分: %d * %d = %.2f", points, extraData["quantity"].(int), totalReturnPoints)
			} else {
				log.Printf("积分方法为按金额比例计算，比例为%d%%", points)
				// 假设points字段存储的是百分比
				totalReturnPoints = total * (float64(points) / 100)
				log.Printf("按金额比例计算返积分: %.2f * %d%% = %.2f", total, points, totalReturnPoints)
			}
		} else {
			log.Printf("积分值为0或负数，不计算返积分")
		}
	} else {
		log.Println("积分值为nil，不计算返积分")
	}

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
}
