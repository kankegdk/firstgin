package services

import (
	"fmt"
	"log"
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
	GetBuyNowInfo(params map[string]interface{}) (*structs.BuyNowInfoData, error)
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
func (s *productService) GetBuyNowInfo(params map[string]interface{}) (*structs.BuyNowInfoData, error) {
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
	e, err := models.CartGoods(cartParams)
	log.Println("cart 产品", e)
	if err != nil {
		return nil, fmt.Errorf("获取商品信息失败: %v", err)
	}
	return nil, nil
	// // 检查商品状态
	// if product.Status != 1 {
	// 	return nil, errors.New("商品已下架")
	// }

	// // 处理积分商品
	// if product.Points > 0 {
	// 	// 在实际项目中，这里需要验证用户积分是否足够
	// }

	// // 处理秒杀商品
	// if msid, err := strconv.ParseInt(params.Msid, 10, 64); err == nil && msid > 0 {
	// 	// 在实际项目中，这里需要验证秒杀活动是否有效
	// }

	// // 处理团购商品
	// if tuanid, err := strconv.ParseInt(params.Tuanid, 10, 64); err == nil && tuanid > 0 {
	// 	// 在实际项目中，这里需要验证团购活动是否有效
	// }

	// // 获取分类信息
	// category, err := models.GetCategoryByID(product.CatID)
	// if err != nil {
	// 	// 分类信息不是必须的，可以继续处理
	// }

	// // 构建返回数据
	// result := &structs.BuyNowInfoData{
	// 	Product:  product,
	// 	Category: category,
	// 	// 默认配送方式为快递
	// 	ShippingType: "express",
	// }

	// // 转换为JSON字符串
	// dataJSON, err := json.Marshal(params)
	// if err != nil {
	// 	return nil, fmt.Errorf("转换为JSON字符串失败: %v", err)
	// }

	// // 记录购买信息
	// recordId, err := models.CreateGoodsBuynowinfo(params.Weid, params.Ip, string(dataJSON))
	// if err == nil {
	// 	result.RecordId = recordId
	// }

	// return result, nil
}
