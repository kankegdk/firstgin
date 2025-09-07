package services

import (
	"errors"
	"myapi/app/models"
	"myapi/app/structs"
)

// TuanGoodsService 定义团购商品服务接口
type TuanGoodsService interface {
	GetTuanGoodsByID(id int) *structs.TuanGoods
	GetTuanGoodsByGoodsID(goodsID int) *structs.TuanGoods
	GetTuanFoundByID(id int) *structs.TuanFound
	GetTuanFollowCountByFoundID(foundID int) (int64, error)
	ValidateCanJoinTuan(goodsID int, tuanID int, quantity int) (*structs.TuanGoods, error)
}

// tuanGoodsService 实现TuanGoodsService接口的结构体
type tuanGoodsService struct{}

// NewTuanGoodsService 创建一个新的团购商品服务实例
func NewTuanGoodsService() TuanGoodsService {
	return &tuanGoodsService{}
}

// GetTuanGoodsByID 根据ID获取团购商品
func (s *tuanGoodsService) GetTuanGoodsByID(id int) *structs.TuanGoods {
	return models.GetTuanGoodsByID(id)
}

// GetTuanGoodsByGoodsID 根据商品ID获取团购商品
func (s *tuanGoodsService) GetTuanGoodsByGoodsID(goodsID int) *structs.TuanGoods {
	return models.GetTuanGoodsByGoodsID(goodsID)
}

// GetTuanFoundByID 根据ID获取开团记录
func (s *tuanGoodsService) GetTuanFoundByID(id int) *structs.TuanFound {
	return models.GetTuanFoundByID(id)
}

// GetTuanFollowCountByFoundID 根据开团ID获取参团人数
func (s *tuanGoodsService) GetTuanFollowCountByFoundID(foundID int) (int64, error) {
	return models.GetTuanFollowCountByFoundID(foundID)
}

// ValidateCanJoinTuan 验证是否可平团
func (s *tuanGoodsService) ValidateCanJoinTuan(goodsID int, tuanID int, quantity int) (*structs.TuanGoods, error) {
	// 参数验证
	if goodsID <= 0 {
		return nil, errors.New("商品ID无效")
	}
	if tuanID <= 0 {
		return nil, errors.New("团购ID无效")
	}
	if quantity <= 0 {
		return nil, errors.New("购买数量无效")
	}

	// 调用model层的验证方法
	tuanGoods, err := models.ValidateCanJoinTuan(int64(goodsID), int64(tuanID), int64(quantity))
	if err != nil {
		return nil, err
	}

	return tuanGoods, nil
}
