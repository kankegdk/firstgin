package services

import (
	"myapi/app/models"
)

// ProductService 产品服务接口
type ProductService interface {
	GetProductByID(id int) *models.Product
	GetAllProducts() []models.Product
	CreateProduct(product *models.Product) error
	UpdateProduct(id int, product *models.Product) error
	DeleteProduct(id int) error
}

// productService 实现ProductService接口的结构体
type productService struct{}

// NewProductService 创建一个新的产品服务实例
func NewProductService() ProductService {
	return &productService{}
}

// GetProductByID 根据ID获取产品
func (s *productService) GetProductByID(id int) *models.Product {
	// 调用模型层方法
	return models.GetProductByID(id)
}

// GetAllProducts 获取所有产品
func (s *productService) GetAllProducts() []models.Product {
	// 调用模型层方法
	return models.GetAllProducts()
}

// CreateProduct 创建产品
func (s *productService) CreateProduct(product *models.Product) error {
	// 实际项目中，这里会包含业务逻辑验证
	// 然后调用模型层的创建方法
	// 这里只是一个示例
	return nil
}

// UpdateProduct 更新产品
func (s *productService) UpdateProduct(id int, product *models.Product) error {
	// 实际项目中，这里会包含业务逻辑验证
	// 然后调用模型层的更新方法
	// 这里只是一个示例
	return nil
}

// DeleteProduct 删除产品
func (s *productService) DeleteProduct(id int) error {
	// 实际项目中，这里会包含业务逻辑验证
	// 然后调用模型层的删除方法
	// 这里只是一个示例
	return nil
}
