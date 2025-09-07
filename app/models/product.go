package models

import (
	"errors"
	"log"
	"myapi/app/storage"
	"myapi/app/structs"
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
	// 安全地进行类型断言
	skuStr, ok := sku.(string)
	if !ok && sku != nil {
		log.Println("SKU参数类型错误，无法转换为字符串")
		return nil, errors.New("SKU参数无效")
	}

	// log.Printf("确认商品信息存在，准备处理各种活动信息")
	extraData["quantity"] = params["quantity"]
	log.Println("团购价格处理前的价格: ,", extraData["price"], "数量:", extraData["quantity"])

	// 处理团购信息
	var tuanid = params["tuanid"].(int64)
	if tuanid > 0 {
		// 正确传递字符串给可变参数函数
		tuanGoods, err := ValidateCanJoinTuan(params["GoodsID"].(int64), int64(tuanid), params["quantity"].(int64), skuStr)
		if err != nil {
			log.Printf("团购验证失败: %v，直接返回错误", err)
			return nil, err
		}
		extraData["price"] = tuanGoods.Price
		log.Println("团购价格处理后的价格: ", extraData["price"])
		//如果是团购就返回团购相关信息
		return extraData, nil
	} else {
		log.Println("未提供有效的团购ID参数，跳过团购信息处理")
	}

	log.Printf("秒杀ID参数:%T,%v", params["msid"], params["msid"])

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

	return nil, nil

	/*
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


		// 处理多SKU情况
		if isSkumoreParam, ok := params["is_skumore"]; ok && isSkumoreParam != "" {
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
		quantity := 0.0
		quantityInt := 0
		if q, ok := extraData["quantity"].(float64); ok {
			quantity = q
			quantityInt = int(q)
		} else if q, ok := extraData["quantity"].(int); ok {
			quantity = float64(q)
			quantityInt = q
		} else {
			log.Printf("警告: quantity字段类型未知(%T)，使用默认值0", extraData["quantity"])
		}
		total := price * quantity
		totalPayPoints := 0
		totalPointsPrice := 0
		log.Printf("计算商品总价: 单价%.2f * 数量%d = %.2f", price, quantityInt, total)

		// 安全获取PayPoints
		if payPointsValue, ok := extraData["PayPoints"]; ok && payPointsValue != nil {
			log.Printf("尝试获取支付积分值: %v", payPointsValue)
			if pp, ok := payPointsValue.(int); ok {
				totalPayPoints = pp * quantityInt
				log.Printf("支付积分为int类型，计算总支付积分: %d * %d = %d", pp, quantityInt, totalPayPoints)
			} else if pp, ok := payPointsValue.(float64); ok {
				totalPayPoints = int(pp) * quantityInt
				log.Printf("支付积分为float64类型，计算总支付积分: %d * %d = %d", int(pp), quantityInt, totalPayPoints)
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
	*/
}
