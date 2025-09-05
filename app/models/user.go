package models

import (
	"errors"
	"log"
	"myapi/app/storage"

	"gorm.io/gorm"
)

// User 定义用户模型，用于GORM操作
type User struct {
}

// TableName 设置表名
func (u *User) TableName() string {
	return userTableName
}

// GetAllUsers 获取所有用户
func (u *User) GetAllUsers() ([]map[string]interface{}, error) {
	var users []map[string]interface{}

	gormDB := storage.GetGormDB()
	if gormDB == nil {
		log.Println("GORM连接为空")
		return []map[string]interface{}{}, errors.New("数据库连接失败")
	}

	// 明确打印当前使用的表名，用于调试
	log.Printf("当前使用的表名: %s", userTableName)

	// 使用Unscoped()完全禁用软删除功能，直接查询表中的所有记录
	result := gormDB.Unscoped().Table(userTableName).Select("id,username,mobile").Find(&users)
	return users, result.Error
}

// BeforeSave 钩子函数，在保存前执行
func (u *User) BeforeSave(tx *gorm.DB) (err error) {
	// 这里可以添加保存前的逻辑，比如密码加密等
	return
}

// FindByID 根据ID查找用户，返回map[string]interface{}类型，可以包含数据库中的所有字段
func (u *User) FindByID(id uint) (map[string]interface{}, error) {
	user := make(map[string]interface{})

	gormDB := storage.GetGormDB()
	if gormDB == nil {
		return nil, errors.New("数据库连接失败")
	}

	result := gormDB.Unscoped().Table(userTableName).Where("id = ?", id).First(&user)
	return user, result.Error
}

// FindByUsername 根据用户名查找用户，返回map[string]interface{}类型，可以包含数据库中的所有字段
func (u *User) FindByUsername(username string) (map[string]interface{}, error) {
	user := make(map[string]interface{})

	gormDB := storage.GetGormDB()
	if gormDB == nil {
		return nil, errors.New("数据库连接失败")
	}

	result := gormDB.Unscoped().Table(userTableName).Where("username = ?", username).First(&user)
	return user, result.Error
}

// Create 创建新用户，可以传入map类型的数据
func (u *User) Create(userData map[string]interface{}) error {
	gormDB := storage.GetGormDB()
	if gormDB == nil {
		return errors.New("数据库连接失败")
	}
	return gormDB.Table(userTableName).Create(userData).Error
}

// UpdateByID 根据ID更新用户信息，支持传入map类型的数据更新指定字段
func (u *User) UpdateByID(id uint, updateData map[string]interface{}) error {
	gormDB := storage.GetGormDB()
	if gormDB == nil {
		return errors.New("数据库连接失败")
	}
	return gormDB.Table(userTableName).Where("id = ?", id).Updates(updateData).Error
}

// DeleteByID 根据ID删除用户
func (u *User) DeleteByID(id uint) error {
	gormDB := storage.GetGormDB()
	if gormDB == nil {
		return errors.New("数据库连接失败")
	}
	return gormDB.Unscoped().Table(userTableName).Where("id = ?", id).Delete(nil).Error
}
