package services

import (
	"encoding/json"
	"errors"
	"log"
	"myapi/app/config"
	"myapi/app/models"
	"myapi/app/storage"
	"myapi/app/structs"
	"regexp"
	"time"

	"github.com/go-redis/redis/v8"
)

// AddressService 地址服务接口
type AddressService interface {
	AddAddress(data structs.Address) (int64, error)
	UpdateAddress(id int, data structs.Address) error
	DeleteAddress(id, weid, uid int) error
	SetDefaultAddress(id, weid, uid int) error
	GetAddressDetail(id, weid, uid int) (*structs.Address, error)
	GetDefaultAddress(weid, uid int) (*structs.Address, error)
	GetAddressList(weid, uid int) ([]structs.Address, error)
}

// addressService 实现AddressService接口的结构体
type addressService struct{}

// NewAddressService 创建一个新的地址服务实例
func NewAddressService() AddressService {
	return &addressService{}
}

// validateAddressData 验证地址数据
func validateAddressData(data structs.Address) error {
	if data.Name == "" {
		return errors.New("请输入收货人姓名")
	}

	if data.Telephone == "" {
		return errors.New("请输入手机号")
	}

	// 验证手机号格式
	phoneRegex := regexp.MustCompile(`^1[3-9]\d{9}$`)
	if !phoneRegex.MatchString(data.Telephone) {
		return errors.New("请输入正确的手机号")
	}

	if data.ProvinceId == 0 || data.CityId == 0 || data.DistrictId == 0 {
		return errors.New("请选择完整的省市区信息")
	}

	if data.Address == "" {
		return errors.New("请输入详细地址与门牌号")
	}

	return nil
}

// AddAddress 添加地址
func (s *addressService) AddAddress(data structs.Address) (int64, error) {
	// 验证数据
	if err := validateAddressData(data); err != nil {
		return 0, err
	}

	// 如果设置为默认地址，先取消其他默认地址
	if data.IsDefault == 1 {
		if err := models.CancelOtherDefaultAddress(data.Weid, data.Uid); err != nil {
			return 0, err
		}
	}

	// 调用模型添加地址
	return models.AddAddress(data)
}

// UpdateAddress 更新地址
func (s *addressService) UpdateAddress(id int, data structs.Address) error {
	if id <= 0 {
		return errors.New("地址ID不能为空")
	}

	// 验证地址是否存在且属于当前用户
	_, err := models.GetAddressDetail(id, data.Uid)
	if err != nil {
		return errors.New("地址不存在或无权限修改")
	}

	// 验证数据
	if err := validateAddressData(data); err != nil {
		return err
	}

	// 如果设置为默认地址，先取消其他默认地址
	if data.IsDefault == 1 {
		if err := models.CancelOtherDefaultAddress(data.Weid, data.Uid); err != nil {
			return err
		}
	}

	// 调用模型更新地址
	return models.UpdateAddress(id, data)
}

// DeleteAddress 删除地址
func (s *addressService) DeleteAddress(id, weid, uid int) error {
	if id <= 0 {
		return errors.New("地址ID不能为空")
	}

	// 验证地址是否存在且属于当前用户
	_, err := models.GetAddressDetail(id, uid)
	if err != nil {
		return errors.New("地址不存在或无权限删除")
	}

	// 调用模型删除地址
	return models.DeleteAddress(id)
}

// SetDefaultAddress 设置默认地址
func (s *addressService) SetDefaultAddress(id, weid, uid int) error {
	if id <= 0 {
		return errors.New("地址ID不能为空")
	}

	// 验证地址是否存在且属于当前用户
	_, err := models.GetAddressDetail(id, uid)
	if err != nil {
		return errors.New("地址不存在或无权限设置")
	}

	// 先取消所有默认地址
	if err := models.CancelOtherDefaultAddress(weid, uid); err != nil {
		return err
	}

	// 设置新的默认地址
	return models.SetDefaultAddress(id, weid, uid)
}

// GetAddressDetail 获取地址详情
func (s *addressService) GetAddressDetail(id, weid, uid int) (*structs.Address, error) {
	if id <= 0 {
		return nil, errors.New("地址ID不能为空")
	}

	// 从数据库获取地址详情
	return models.GetAddressDetail(id, uid)
}

// GetDefaultAddress 获取默认地址
func (s *addressService) GetDefaultAddress(weid, uid int) (*structs.Address, error) {
	// 构建缓存键
	prefix := config.GetRedisConfig().Prefix
	cacheKey := prefix + ":" + "default_address:" + string(uid)

	// 尝试从缓存获取数据
	cacheData, err := storage.GetCache(cacheKey)
	if err == nil {
		log.Println("从Redis缓存中获取默认地址数据cacheKey", cacheKey)
		// 缓存命中，解析JSON并返回
		var address structs.Address
		if err2 := json.Unmarshal([]byte(cacheData), &address); err2 == nil {
			return &address, nil
		} else {
			log.Println("解析默认地址缓存数据失败:", err2)
		}
	} else if err != redis.Nil {
		// 如果是其他错误（非键不存在），记录日志
		log.Println("获取默认地址Redis缓存失败:", err)
	} else {
		// 键不存在的情况，记录日志但不报错
		log.Println("默认地址Redis缓存键不存在:", cacheKey)
	}

	// 缓存未命中，从数据库获取数据
	address, err := models.GetDefaultAddress(weid, uid)

	// 将结果存入缓存，过期时间设为12小时
	if err == nil && address != nil {
		if data, err := json.Marshal(address); err == nil {
			log.Println("将默认地址数据存入Redis缓存cacheKey", cacheKey)
			storage.SetCache(cacheKey, string(data), time.Hour*12)
		}
	}

	return address, err
}

// GetAddressList 获取地址列表
func (s *addressService) GetAddressList(weid, uid int) ([]structs.Address, error) {
	// 构建缓存键
	prefix := config.GetRedisConfig().Prefix
	cacheKey := prefix + ":" + "address_list:" + string(uid)

	// 尝试从缓存获取数据
	cacheData, err := storage.GetCache(cacheKey)
	if err == nil {
		log.Println("从Redis缓存中获取地址列表数据cacheKey", cacheKey)
		// 缓存命中，解析JSON并返回
		var addresses []structs.Address
		if err2 := json.Unmarshal([]byte(cacheData), &addresses); err2 == nil {
			return addresses, nil
		} else {
			log.Println("解析地址列表缓存数据失败:", err2)
		}
	} else if err != redis.Nil {
		// 如果是其他错误（非键不存在），记录日志
		log.Println("获取地址列表Redis缓存失败:", err)
	} else {
		// 键不存在的情况，记录日志但不报错
		log.Println("地址列表Redis缓存键不存在:", cacheKey)
	}

	// 缓存未命中，从数据库获取数据
	addresses, err := models.GetAddressList(weid, uid)

	// 将结果存入缓存，过期时间设为12小时
	if err == nil {
		if data, err := json.Marshal(addresses); err == nil {
			log.Println("将地址列表数据存入Redis缓存cacheKey", cacheKey)
			storage.SetCache(cacheKey, string(data), time.Hour*12)
		}
	}

	return addresses, err
}