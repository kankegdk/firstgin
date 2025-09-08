package controllers

import (
	"log"
	"net/http"

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
		Goods_id   int64  `json:"goodsId" binding:"required|gt=0"`
		BuyNumber  int64  `json:"buyNumber" binding:"required|gt=0"`
		Msid       int64  `json:"msid" `       // 秒杀ID，可选
		Tuanid     int64  `json:"tuanid" `     // 团购ID，可选
		Tecid      int64  `json:"tecid" `      // 技术ID，可选
		Jointuanid int64  `json:"jointuanid" ` // 参团ID，可选
		Sku        string `json:"sku"`         // SKU，可选
		IsSkumore  int    `json:"is_skumore"`  // 是否多规格，可选
		Skumore    int    `json:"skumore"`     // 多规格详情，可选
	}

	// 使用公共验证器验证请求
	if !helper.ValidateRequest(c, &buyData) {
		return
	}

	// 创建map类型的参数
	params := make(map[string]interface{})
	params["Msid"] = buyData.Msid
	params["Tuanid"] = buyData.Tuanid
	params["Tecid"] = buyData.Tecid
	params["Jointuanid"] = buyData.Jointuanid
	params["Sku"] = buyData.Sku
	params["IsSkumore"] = buyData.IsSkumore
	params["Skumore"] = buyData.Skumore
	params["BuyNumber"] = buyData.BuyNumber
	params["GoodsID"] = buyData.Goods_id

	// 解析goods_id为整数
	params["Uid"] = helper.UID(c)
	//不能同时有msid和tuanid
	if buyData.Msid != 0 && buyData.Tuanid != 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"msg":  "秒杀和团购不能同时购买",
		})
		return
	}
	// 如果Uid=0 且 mid>0就要验证是否登录,必须先登录才能下单
	// if params["Uid"] == 0 && buyData.Msid > 0 {
	// 	// 验证登录，使用helper.CheckLogin函数
	// 	log.Printf("验证前id: %v", params["Uid"])
	// 	isLoggedIn, _, _ := helper.CheckLogin(c)
	// 	if !isLoggedIn {
	// 		c.JSON(http.StatusOK, gin.H{
	// 			"code": 1,
	// 			"msg":  "请先登录",
	// 		})
	// 		return
	// 	}
	// }
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
