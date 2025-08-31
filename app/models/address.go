package models

import (
	"log"

	"myapi/app/config"
	"myapi/app/storage"
	"myapi/app/structs"
)

// AddAddress 添加新地址
func AddAddress(data structs.Address) (int64, error) {
	// 获取共享的gorm连接实例
	gormDB := storage.GetGormDB()
	if gormDB == nil {
		log.Println("GORM连接为空")
		return 0, nil
	}

	// 获取表前缀
	cfg := config.GetDatabaseConfig()
	tablePrefix := cfg.GetTablePrefix()
	tableName := tablePrefix + "address"

	// 执行插入操作
	result := gormDB.Table(tableName).Create(&data)
	if result.Error != nil {
		log.Println("添加地址失败:", result.Error)
		return 0, result.Error
	}

	// 返回插入的ID
	return result.RowsAffected, nil
}

// UpdateAddress 更新地址
func UpdateAddress(id int, data structs.Address) error {
	// 获取共享的gorm连接实例
	gormDB := storage.GetGormDB()
	if gormDB == nil {
		log.Println("GORM连接为空")
		return nil
	}

	// 获取表前缀
	cfg := config.GetDatabaseConfig()
	tablePrefix := cfg.GetTablePrefix()
	tableName := tablePrefix + "address"

	// 执行更新操作
	result := gormDB.Table(tableName).Where("id = ?", id).Updates(&data)
	if result.Error != nil {
		log.Println("更新地址失败:", result.Error)
		return result.Error
	}

	return nil
}

// DeleteAddress 删除地址
func DeleteAddress(id int) error {
	// 获取共享的gorm连接实例
	gormDB := storage.GetGormDB()
	if gormDB == nil {
		log.Println("GORM连接为空")
		return nil
	}

	// 获取表前缀
	cfg := config.GetDatabaseConfig()
	tablePrefix := cfg.GetTablePrefix()
	tableName := tablePrefix + "address"

	// 执行删除操作
	result := gormDB.Table(tableName).Where("id = ?", id).Delete(&structs.Address{})
	if result.Error != nil {
		log.Println("删除地址失败:", result.Error)
		return result.Error
	}

	return nil
}

// GetAddressDetail 获取地址详情
func GetAddressDetail(id, uid int) (*structs.Address, error) {
	// 获取共享的gorm连接实例
	gormDB := storage.GetGormDB()
	if gormDB == nil {
		log.Println("GORM连接为空")
		return nil, nil
	}

	// 获取表前缀
	cfg := config.GetDatabaseConfig()
	tablePrefix := cfg.GetTablePrefix()
	tableName := tablePrefix + "address"

	// 查询地址详情
	var address structs.Address
	result := gormDB.Table(tableName).Where("id = ? AND uid = ?", id, uid).First(&address)
	if result.Error != nil {
		log.Println("获取地址详情失败:", result.Error)
		return nil, result.Error
	}

	return &address, nil
}

// GetDefaultAddress 获取默认地址
func GetDefaultAddress(weid, uid int) (*structs.Address, error) {
	// 获取共享的gorm连接实例
	gormDB := storage.GetGormDB()
	if gormDB == nil {
		log.Println("GORM连接为空")
		return nil, nil
	}

	// 获取表前缀
	cfg := config.GetDatabaseConfig()
	tablePrefix := cfg.GetTablePrefix()
	tableName := tablePrefix + "address"

	// 查询默认地址
	var address structs.Address
	result := gormDB.Table(tableName).Where("weid = ? AND uid = ? AND isDefault = ?", weid, uid, 1).First(&address)
	if result.Error != nil {
		log.Println("获取默认地址失败:", result.Error)
		return nil, result.Error
	}

	return &address, nil
}

// GetAddressList 获取地址列表
func GetAddressList(weid, uid int) ([]structs.Address, error) {
	// 获取共享的gorm连接实例
	gormDB := storage.GetGormDB()
	if gormDB == nil {
		log.Println("GORM连接为空")
		return []structs.Address{}, nil
	}

	// 获取表前缀
	cfg := config.GetDatabaseConfig()
	tablePrefix := cfg.GetTablePrefix()
	tableName := tablePrefix + "address"

	// 查询地址列表
	var addresses []structs.Address
	result := gormDB.Table(tableName).Where("weid = ? AND uid = ?", weid, uid).Order("isDefault desc, id desc").Find(&addresses)
	if result.Error != nil {
		log.Println("获取地址列表失败:", result.Error)
		return []structs.Address{}, result.Error
	}

	return addresses, nil
}

// CancelOtherDefaultAddress 取消其他默认地址
func CancelOtherDefaultAddress(weid, uid int) error {
	// 获取共享的gorm连接实例
	gormDB := storage.GetGormDB()
	if gormDB == nil {
		log.Println("GORM连接为空")
		return nil
	}

	// 获取表前缀
	cfg := config.GetDatabaseConfig()
	tablePrefix := cfg.GetTablePrefix()
	tableName := tablePrefix + "address"

	// 取消其他默认地址
	result := gormDB.Table(tableName).Where("weid = ? AND uid = ?", weid, uid).Update("isDefault", 0)
	if result.Error != nil {
		log.Println("取消其他默认地址失败:", result.Error)
		return result.Error
	}

	return nil
}

// SetDefaultAddress 设置默认地址
func SetDefaultAddress(id, weid, uid int) error {
	// 获取共享的gorm连接实例
	gormDB := storage.GetGormDB()
	if gormDB == nil {
		log.Println("GORM连接为空")
		return nil
	}

	// 获取表前缀
	cfg := config.GetDatabaseConfig()
	tablePrefix := cfg.GetTablePrefix()
	tableName := tablePrefix + "address"

	// 设置默认地址
	result := gormDB.Table(tableName).Where("id = ?", id).Update("isDefault", 1)
	if result.Error != nil {
		log.Println("设置默认地址失败:", result.Error)
		return result.Error
	}

	return nil
}