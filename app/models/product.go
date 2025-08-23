package models

// Product 产品模型
type Product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	// 其他字段...
}

// GetProductByID 根据ID获取产品
func GetProductByID(id int) *Product {
	// 实际实现，目前返回模拟数据
	return &Product{
		ID:    id,
		Name:  "示例产品",
		Price: 99.99,
	}
}

// GetAllProducts 获取所有产品
func GetAllProducts() []Product {
	// 实际实现，目前返回模拟数据
	return []Product{
		{ID: 1, Name: "产品1", Price: 99.99},
		{ID: 2, Name: "产品2", Price: 199.99},
	}
}
