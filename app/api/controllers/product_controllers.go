package controllers

import (
	"log"
	"net/http"
	"strconv"

	"myapi/app/helper"
	"myapi/app/services" // 导入服务层

	"github.com/gin-gonic/gin"
)

// GetAllProducts 处理获取所有产品的请求
func GetAllProducts(c *gin.Context) {
	// 1. 创建服务实例
	productService := services.NewProductService()
	// 2. 调用服务层方法
	products := productService.GetAllProducts()

	c.JSON(http.StatusOK, gin.H{
		"message": "获取产品列表成功",
		"data":    products,
		"count":   len(products),
	})
}

// BuyNowInfo 处理立即购买信息请求
func BuyNowInfo(c *gin.Context) {
	// 1. 解析请求参数
	var buyData struct {
		Goods_id   string `json:"goodsId" binding:"required"`
		BuyNumber  string `json:"buyNumber" binding:"required"`
		Msid       string `json:"msid"`       // 秒杀ID，可选
		Tuanid     string `json:"tuanid"`     // 团购ID，可选
		Tecid      string `json:"tecid"`      // 技术ID，可选
		Jointuanid string `json:"jointuanid"` // 参团ID，可选
		Sku        string `json:"sku"`        // SKU，可选
		IsSkumore  string `json:"is_skumore"` // 是否多规格，可选
		Skumore    string `json:"skumore"`    // 多规格详情，可选
	}

	// 使用公共验证器验证请求
	if !helper.ValidateRequest(c, &buyData) {
		return
	}

	// 创建map类型的参数
	params := make(map[string]interface{})
	params["BuyNumber"] = buyData.BuyNumber

	// 解析goods_id为整数
	goodsID, _ := strconv.Atoi(buyData.Goods_id)
	params["GoodsID"] = goodsID
	params["Uid"] = helper.UID(c)
	// 处理秒杀ID
	if buyData.Msid != "" {
		msidInt, err := strconv.ParseInt(buyData.Msid, 10, 64)
		if err == nil && msidInt > 0 {
			params["Msid"] = msidInt
		}
	}

	// 添加其他可选参数
	if buyData.Tuanid != "" {
		params["Tuanid"] = buyData.Tuanid
	}
	if buyData.Tecid != "" {
		params["Tecid"] = buyData.Tecid
	}
	if buyData.Jointuanid != "" {
		params["Jointuanid"] = buyData.Jointuanid
	}
	if buyData.Sku != "" {
		params["Sku"] = buyData.Sku
	}
	if buyData.IsSkumore != "" {
		params["IsSkumore"] = buyData.IsSkumore
	}
	if buyData.Skumore != "" {
		params["Skumore"] = buyData.Skumore
	}
	log.Printf("验证前id: %v", params["Uid"])
	//不能同时有msid和tuanid
	if buyData.Msid != "" && buyData.Tuanid != "" {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "秒杀和团购不能同时购买",
		})
		return
	}
	// 如果Uid=0 且 mid>0就要验证是否登录
	uid, uidOk := params["Uid"].(int64)
	if uidOk && uid == 0 && buyData.Msid != "" {
		msidInt, err := strconv.ParseInt(buyData.Msid, 10, 64)
		if err == nil && msidInt > 0 {
			// 验证登录，使用helper.CheckLogin函数
			isLoggedIn, _, _ := helper.CheckLogin(c)
			if !isLoggedIn {
				c.JSON(http.StatusOK, gin.H{
					"code": 1,
					"msg":  "请先登录",
				})
				return
			}
		}
	}
	params["Uid"] = helper.UID(c)

	log.Printf("验证后id: %v", params["Uid"])
	// 设置weid和ip字段
	weid := helper.Weid(c)
	if weid <= 0 {
		weid = helper.GetWeid() // 使用默认weid
	}
	params["Weid"] = weid
	// 获取客户端IP
	params["Ip"] = helper.GetRealIP(c)
	//键值对形式打印params
	log.Printf("BuyNowInfo params: %v", params)

	// 3. 创建服务实例并调用服务层方法
	productService := services.NewProductService()
	result, err := productService.GetBuyNowInfo(params)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  err.Error(),
		})
		return
	}

	// 4. 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": result,
	})
}
