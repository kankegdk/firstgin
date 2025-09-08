package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"myapi/app/helper"
	"myapi/app/models"
	"myapi/app/structs"
)

// ProductService 定义产品服务接口
type ProductService interface {
	GetProductByID(id int) *structs.Product
	GetAllProducts() []structs.Product
	CreateProduct(product *structs.Product) error
	UpdateProduct(product *structs.Product) error
	DeleteProduct(id int) error
	GetBuyNowInfo(params map[string]interface{}) (map[string]interface{}, error)
}

// productService 实现ProductService接口的结构体
type productService struct{}

// NewProductService 创建一个新的产品服务实例
func NewProductService() ProductService {
	return &productService{}
}

// GetProductByID 根据ID获取产品
func (s *productService) GetProductByID(id int) *structs.Product {
	return models.GetProductByID(id)
}

// GetAllProducts 获取所有产品
func (s *productService) GetAllProducts() []structs.Product {
	return models.GetAllProducts()
}

// CreateProduct 创建产品
func (s *productService) CreateProduct(product *structs.Product) error {
	// 实际实现
	return nil
}

// UpdateProduct 更新产品
func (s *productService) UpdateProduct(product *structs.Product) error {
	// 实际实现
	return nil
}

// DeleteProduct 删除产品
func (s *productService) DeleteProduct(id int) error {
	// 实际实现
	return nil
}

// GetBuyNowInfo 获取立即购买信息
func (s *productService) GetBuyNowInfo(params map[string]interface{}) (map[string]interface{}, error) {
	// 参数验证
	var goodsID = params["GoodsID"]

	// 构建查询参数
	cartParams := make(map[string]interface{})
	cartParams["GoodsID"] = goodsID
	cartParams["id"] = goodsID
	cartParams["sku"] = params["Sku"]
	cartParams["quantity"] = params["BuyNumber"]
	cartParams["msid"] = params["Msid"]
	cartParams["tuanid"] = params["Tuanid"]
	cartParams["tecid"] = params["Tecid"]
	cartParams["jointuanid"] = params["Jointuanid"]
	cartParams["is_skumore"] = params["IsSkumore"]
	cartParams["skumore"] = params["Skumore"]
	cartParams["uid"] = params["Uid"]
	// return nil, nil
	// // 获取商品详情
	product, err := models.CartGoods(cartParams)
	log.Println("cart 产品", product)
	if err != nil {
		return nil, fmt.Errorf("获取商品信息失败: %v", err)
	}
	status, _ := helper.ToInt(product["status"])
	if status != 1 {
		return nil, errors.New("商品已下架")
	}

	// 处理积分商品
	if is_points_goods, err := helper.ToInt(product["is_points_goods"]); true || err == nil && is_points_goods == 1 {
		// 在实际项目中，这里需要验证用户积分是否足够
		points, _ := helper.ToInt(product["pay_points"])
		uid, _ := helper.ToInt(params["Uid"])
		if points > 0 {
			// 验证用户积分是否足够
			userPoints, err := models.GetUserPoints(uid)

			log.Println("用户积分", userPoints, "需要积分：", points)
			if err != nil {
				return nil, fmt.Errorf("获取用户积分失败: %v", err)
			}
			if userPoints < points {
				return nil, errors.New("积分不足")
			}
		}
	}

	infodata := make(map[string]interface{})
	infodata["shopList"] = product
	infodata["sid"] = product["sid"]
	infodata["ptype"] = 1

	// 转换为JSON字符串
	dataJSON, err := json.Marshal(product)
	if err != nil {
		return nil, fmt.Errorf("转换为JSON字符串失败: %v", err)
	}

	// 记录购买信息
	recordId, err := models.CreateGoodsBuynowinfo(params["Weid"].(int), params["Ip"].(string), string(dataJSON))

	// 构建返回数据

	if err == nil {
		infodata["recordId"] = recordId
	}

	return infodata, nil
}
