package models

import (
	"errors"
	"log"
	"myapi/app/storage"
	"myapi/app/structs"
	"time"

	"gorm.io/gorm"
)

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

// GetMemberByID 根据ID获取会员信息
func GetMemberByID(id int) (*structs.Member, error) {
	gormDB := storage.GetGormDB()
	if gormDB == nil {
		return nil, errors.New("数据库连接失败")
	}

	member := &structs.Member{}
	result := gormDB.Table(memberTableName).
		Where("id = ?", id).
		First(member)

	if result.Error != nil {
		log.Printf("根据ID查询会员失败: %v", result.Error)
		if result.RowsAffected == 0 {
			return nil, nil // 用户不存在
		}
		return nil, result.Error
	}

	return member, nil
}

// CreateMember 创建新会员
func CreateMember(member *structs.Member) error {
	gormDB := storage.GetGormDB()
	if gormDB == nil {
		return errors.New("数据库连接失败")
	}

	// 设置默认值
	if member.Regdate == 0 {
		member.Regdate = time.Now().Unix()
	}
	if member.Lastdate == 0 {
		member.Lastdate = time.Now().Unix()
	}
	if member.Status == 0 {
		member.Status = 1 // 默认启用
	}

	result := gormDB.Table(memberTableName).Create(member)
	if result.Error != nil {
		log.Printf("创建会员失败: %v", result.Error)
		return result.Error
	}

	return nil
}

// UpdateMember 更新会员信息
func UpdateMember(member *structs.Member) error {
	gormDB := storage.GetGormDB()
	if gormDB == nil {
		return errors.New("数据库连接失败")
	}

	result := gormDB.Table(memberTableName).
		Where("id = ?", member.ID).
		Updates(member)

	if result.Error != nil {
		log.Printf("更新会员信息失败: %v", result.Error)
		return result.Error
	}

	return nil
}

// GetMemberNameByID 根据会员ID获取昵称
func GetMemberNameByID(uid int) (string, error) {
	member, err := GetMemberByID(uid)
	if err != nil {
		return "", err
	}
	if member == nil {
		return "", nil
	}
	return member.Nickname, nil
}

// GetPIDName 获取上级会员昵称
func GetPIDName(uid int) (string, error) {
	member, err := GetMemberByID(uid)
	if err != nil {
		return "", err
	}
	if member == nil {
		return "", nil
	}

	if member.Pid > 0 {
		pidMember, err := GetMemberByID(member.Pid)
		if err != nil {
			return "", err
		}
		if pidMember != nil {
			return pidMember.Nickname, nil
		}
	}

	return "平台", nil
}

// GetOneLevel 获取一级会员
func GetOneLevel(uid int, isData bool) (interface{}, error) {
	gormDB := storage.GetGormDB()
	if gormDB == nil {
		return nil, errors.New("数据库连接失败")
	}

	if isData {
		var members []structs.Member
		result := gormDB.Table(memberTableName).
			Select("id, nickname, regdate, userpic").
			Where("pid = ?", uid).
			Find(&members)
		if result.Error != nil {
			log.Printf("查询一级会员数据失败: %v", result.Error)
			return nil, result.Error
		}
		return members, nil
	} else {
		var count int64
		result := gormDB.Table(memberTableName).
			Where("pid = ?", uid).
			Count(&count)
		if result.Error != nil {
			log.Printf("查询一级会员数量失败: %v", result.Error)
			return 0, result.Error
		}
		return count, nil
	}
}

// GetTwoLevel 获取二级会员
func GetTwoLevel(uid int, isData bool) (interface{}, error) {
	gormDB := storage.GetGormDB()
	if gormDB == nil {
		return nil, errors.New("数据库连接失败")
	}

	// 查询一级会员ID
	var firstLevelIDs []int
	result := gormDB.Table(memberTableName).
		Select("id").
		Where("pid = ?", uid).
		Pluck("id", &firstLevelIDs)
	if result.Error != nil {
		log.Printf("查询一级会员ID失败: %v", result.Error)
		return 0, result.Error
	}

	if len(firstLevelIDs) == 0 {
		if isData {
			return []structs.Member{}, nil
		} else {
			return 0, nil
		}
	}

	if isData {
		var members []structs.Member
		result := gormDB.Table(memberTableName).
			Select("id, nickname, regdate, userpic").
			Where("pid IN ?", firstLevelIDs).
			Find(&members)
		if result.Error != nil {
			log.Printf("查询二级会员数据失败: %v", result.Error)
			return nil, result.Error
		}
		return members, nil
	} else {
		var count int64
		result := gormDB.Table(memberTableName).
			Where("pid IN ?", firstLevelIDs).
			Count(&count)
		if result.Error != nil {
			log.Printf("查询二级会员数量失败: %v", result.Error)
			return 0, result.Error
		}
		return count, nil
	}
}

// GetThreeLevel 获取三级会员
func GetThreeLevel(uid int, isData bool) (interface{}, error) {
	gormDB := storage.GetGormDB()
	if gormDB == nil {
		return nil, errors.New("数据库连接失败")
	}

	// 查询一级会员ID
	var firstLevelIDs []int
	result := gormDB.Table(memberTableName).
		Select("id").
		Where("pid = ?", uid).
		Pluck("id", &firstLevelIDs)
	if result.Error != nil {
		log.Printf("查询一级会员ID失败: %v", result.Error)
		return 0, result.Error
	}

	if len(firstLevelIDs) == 0 {
		if isData {
			return []structs.Member{}, nil
		} else {
			return 0, nil
		}
	}

	// 查询二级会员ID
	var secondLevelIDs []int
	result = gormDB.Table(memberTableName).
		Select("id").
		Where("pid IN ?", firstLevelIDs).
		Pluck("id", &secondLevelIDs)
	if result.Error != nil {
		log.Printf("查询二级会员ID失败: %v", result.Error)
		return 0, result.Error
	}

	if len(secondLevelIDs) == 0 {
		if isData {
			return []structs.Member{}, nil
		} else {
			return 0, nil
		}
	}

	if isData {
		var members []structs.Member
		result := gormDB.Table(memberTableName).
			Select("id, nickname, regdate, userpic").
			Where("pid IN ?", secondLevelIDs).
			Find(&members)
		if result.Error != nil {
			log.Printf("查询三级会员数据失败: %v", result.Error)
			return nil, result.Error
		}
		return members, nil
	} else {
		var count int64
		result := gormDB.Table(memberTableName).
			Where("pid IN ?", secondLevelIDs).
			Count(&count)
		if result.Error != nil {
			log.Printf("查询三级会员数量失败: %v", result.Error)
			return 0, result.Error
		}
		return count, nil
	}
}

// IsTelephoneRegistered 检查手机号是否已注册
func IsTelephoneRegistered(telephone string) (bool, error) {
	member, err := GetMemberByTelephone(telephone)
	if err != nil {
		return false, err
	}
	return member == nil, nil // 返回true表示未注册
}

// IsMemberChecked 检查会员是否已审核通过
func IsMemberChecked(uid int) (bool, error) {
	member, err := GetMemberByID(uid)
	if err != nil {
		return false, err
	}
	if member == nil {
		return false, nil
	}
	return member.Status == 1, nil
}

// BindTelephone 绑定手机号到会员
func BindTelephone(uid int, telephone string) error {
	gormDB := storage.GetGormDB()
	if gormDB == nil {
		return errors.New("数据库连接失败")
	}

	result := gormDB.Table(memberTableName).
		Where("id = ?", uid).
		Update("telephone", telephone)

	if result.Error != nil {
		log.Printf("绑定手机号失败: %v", result.Error)
		return result.Error
	}

	return nil
}

// GetMemberWithGroupInfo 获取会员信息并包含会员组信息
func GetMemberWithGroupInfo(member *structs.Member) (*structs.Member, error) {
	// 在实际实现中，这里应该查询会员组信息并添加到会员对象中
	// 由于没有看到会员组相关的模型，这里暂时返回原会员对象
	return member, nil
}
func GetUserPoints(uid int) (int, error) {
	gormDB := storage.GetGormDB()
	if gormDB == nil {
		return 0, errors.New("数据库连接失败")
	}

	var points int
	result := gormDB.Table(memberTableName).
		Select("points").
		Where("id = ?", uid).
		Pluck("points", &points)

	if result.Error != nil {
		log.Printf("查询用户积分失败: %v", result.Error)
		return 0, result.Error
	}

	return points, nil
}
