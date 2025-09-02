package models

import (
	"errors"
	"log"
	"myapi/app/config"
	"myapi/app/storage"
	"myapi/app/structs"
	"time"

	"gorm.io/gorm"
)

// 全局变量存储完整表名
var memberTableName string

// init函数在包初始化时执行，只配置一次表前缀
func init() {
	// 获取表前缀
	tablePrefix := config.GetString("dbPrefix", "")
	memberTableName = tablePrefix + "member"
}

// GetMemberByUsername 根据用户名获取会员信息
func GetMemberByUsername(username string) (*structs.Member, error) {
	// 获取共享的gorm连接实例
	gormDB := storage.GetGormDB()
	if gormDB == nil {
		log.Println("GORM连接为空")
		return nil, errors.New("数据库连接失败")
	}

	member := &structs.Member{}
	result := gormDB.Table(memberTableName).
		Where("username = ?", username).
		First(member)

	if result.Error != nil {
		log.Printf("根据用户名查询会员失败: %v", result.Error)
		if result.RowsAffected == 0 {
			return nil, nil // 用户不存在
		}
		return nil, result.Error
	}

	return member, nil
}

// GetMemberByTelephone 根据手机号获取会员信息
func GetMemberByTelephone(telephone string) (*structs.Member, error) {
	// 获取共享的gorm连接实例
	gormDB := storage.GetGormDB()
	if gormDB == nil {
		log.Println("GORM连接为空")
		return nil, errors.New("数据库连接失败")
	}

	member := &structs.Member{}
	result := gormDB.Table(memberTableName).
		Where("telephone = ?", telephone).
		First(member)

	if result.Error != nil {
		log.Printf("根据手机号查询会员失败: %v", result.Error)
		if result.RowsAffected == 0 {
			return nil, nil // 用户不存在
		}
		return nil, result.Error
	}

	return member, nil
}

// UpdateMemberLastLogin 更新会员最后登录信息
func UpdateMemberLastLogin(id int, lastIp string) error {
	// 获取共享的gorm连接实例
	gormDB := storage.GetGormDB()
	if gormDB == nil {
		log.Println("GORM连接为空")
		return errors.New("数据库连接失败")
	}

	// 构建更新数据
	updateData := map[string]interface{}{
		"lastdate": time.Now().Unix(),
		"lastip":   lastIp,
	}

	// 执行更新操作
	result := gormDB.Table(memberTableName).
		Where("id = ?", id).
		Updates(updateData).
		UpdateColumn("loginnum", gorm.Expr("loginnum + ?", 1))

	if result.Error != nil {
		log.Printf("更新会员最后登录信息失败: %v", result.Error)
		return result.Error
	}

	return nil
}
