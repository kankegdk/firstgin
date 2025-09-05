package models

import (
	"errors"
	"log"

	"myapi/app/helper"
	"myapi/app/storage"
	"myapi/app/structs"

	"gorm.io/gorm"
)

// GetAdminID 根据会员ID和类型获取管理员ID
func GetAdminID(mid int, ptype string) (int, error) {
	gormDB := storage.GetGormDB()
	if gormDB == nil {
		return 0, errors.New("数据库连接失败")
	}

	var relation structs.UsersRelation
	result := gormDB.Table(usersRelationTableName).
		Where("weid = ? AND ptype = ? AND mid = ?", helper.GetWeid(), ptype, mid).
		Order("id DESC").
		First(&relation)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return 0, nil
		}
		log.Printf("获取管理员ID失败: %v", result.Error)
		return 0, result.Error
	}

	return relation.Uid, nil
}

// GetReID 根据会员ID和类型获取关联ID
func GetReID(mid int, ptype string) (int, error) {
	gormDB := storage.GetGormDB()
	if gormDB == nil {
		return 0, errors.New("数据库连接失败")
	}

	var relation structs.UsersRelation
	result := gormDB.Table(usersRelationTableName).
		Where("weid = ? AND ptype = ? AND mid = ?", helper.GetWeid(), ptype, mid).
		Order("id DESC").
		First(&relation)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// 如果在users_relation中找不到，尝试从uuid_relation中获取
			return GetUuidReID(mid, ptype)
		}
		log.Printf("获取关联ID失败: %v", result.Error)
		return 0, result.Error
	}

	return relation.Reid, nil
}

// GetRelaArray 获取关联数组
func GetRelaArray(mid int, ptype string) (map[string]interface{}, error) {
	gormDB := storage.GetGormDB()
	if gormDB == nil {
		return nil, errors.New("数据库连接失败")
	}

	var relation structs.UsersRelation
	result := gormDB.Table(usersRelationTableName).
		Where("weid = ? AND ptype = ? AND mid = ?", helper.GetWeid(), ptype, mid).
		Order("id DESC").
		First(&relation)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// 如果在users_relation中找不到，尝试从uuid_relation中获取
			return GetUuidRelaArray(mid, ptype)
		}
		log.Printf("获取关联数组失败: %v", result.Error)
		return nil, result.Error
	}

	// 转换为map返回
	return map[string]interface{}{
		"id":    relation.ID,
		"weid":  relation.Weid,
		"ptype": relation.Ptype,
		"mid":   relation.Mid,
		"uid":   relation.Uid,
		"reid":  relation.Reid,
	}, nil
}

// GetUid 根据关联ID和类型获取用户ID
func GetUid(reid int, ptype string) (int, error) {
	gormDB := storage.GetGormDB()
	if gormDB == nil {
		return 0, errors.New("数据库连接失败")
	}

	var relation structs.UsersRelation
	result := gormDB.Table(usersRelationTableName).
		Where("weid = ? AND ptype = ? AND reid = ?", helper.GetWeid(), ptype, reid).
		Order("id DESC").
		First(&relation)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// 如果在users_relation中找不到，尝试从uuid_relation中获取
			return GetUuidUid(reid, ptype)
		}
		log.Printf("获取用户ID失败: %v", result.Error)
		return 0, result.Error
	}

	return relation.Mid, nil
}

// GetAdminMid 根据管理员ID获取会员ID
func GetAdminMid(uid int) (int, error) {
	gormDB := storage.GetGormDB()
	if gormDB == nil {
		return 0, errors.New("数据库连接失败")
	}

	var relation structs.UsersRelation
	result := gormDB.Table(usersRelationTableName).
		Where("weid = ? AND uid = ? AND ptype = ?", helper.GetWeid(), uid, "admin").
		Order("id DESC").
		First(&relation)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return 0, nil
		}
		log.Printf("获取管理员会员ID失败: %v", result.Error)
		return 0, result.Error
	}

	return relation.Mid, nil
}

// GetReIDByID 根据用户信息获取关联ID
func GetReIDByID(userInfo map[string]interface{}, ptype string) (int, error) {
	// 检查用户ID是否存在
	id, idExists := userInfo["id"].(int)
	if !idExists || id == 0 {
		return 0, nil
	}

	// TODO: 实现Users::find方法调用
	// nowUserInfo := Users::find(userInfo["id"])
	// 暂时跳过更新时间检查

	// 根据不同的类型返回不同的关联ID
	if ptype == "store" {
		sid, sidExists := userInfo["sid"].(int)
		if sidExists && sid > 0 {
			return sid, nil
		}
	}

	if ptype == "tuanzhang" {
		tzid, tzidExists := userInfo["tzid"].(int)
		if tzidExists && tzid > 0 {
			return tzid, nil
		}
	}

	if ptype == "operatingcity" {
		ocid, ocidExists := userInfo["ocid"].(int)
		if ocidExists && ocid > 0 {
			return ocid, nil
		}
	}

	if ptype == "technical" {
		tid, tidExists := userInfo["tid"].(int)
		if tidExists && tid > 0 {
			return tid, nil
		}
		// TODO: 实现Technical::where方法调用
	}

	return 0, nil
}

// SetLogout 注销用户关联
func SetLogout(uid int) error {
	gormDB := storage.GetGormDB()
	if gormDB == nil {
		return errors.New("数据库连接失败")
	}

	// TODO: 实现Users::find方法调用
	// userInfo := Users::find(uid)
	// 假设已经获取了userInfo
	userInfo := map[string]interface{}{"weid": helper.GetWeid(), "id": uid}

	// 删除用户关联
	result := gormDB.Table(usersRelationTableName).
		Where("weid = ? AND uid = ?", userInfo["weid"], userInfo["id"]).
		Delete(nil)

	if result.Error != nil {
		log.Printf("删除用户关联失败: %v", result.Error)
		return result.Error
	}

	// 删除师傅关联
	if tid, exists := userInfo["tid"].(int); exists && tid > 0 {
		result = gormDB.Table(usersRelationTableName).
			Where("weid = ? AND reid = ?", userInfo["weid"], tid).
			Delete(nil)

		if result.Error != nil {
			log.Printf("删除师傅关联失败: %v", result.Error)
		}
	}

	// 删除团长关联
	if tzid, exists := userInfo["tzid"].(int); exists && tzid > 0 {
		result = gormDB.Table(usersRelationTableName).
			Where("weid = ? AND reid = ?", userInfo["weid"], tzid).
			Delete(nil)

		if result.Error != nil {
			log.Printf("删除团长关联失败: %v", result.Error)
		}
	}

	// 删除城市代理关联
	if ocid, exists := userInfo["ocid"].(int); exists && ocid > 0 {
		result = gormDB.Table(usersRelationTableName).
			Where("weid = ? AND reid = ?", userInfo["weid"], ocid).
			Delete(nil)

		if result.Error != nil {
			log.Printf("删除城市代理关联失败: %v", result.Error)
		}
	}

	// 删除店铺关联
	if sid, exists := userInfo["sid"].(int); exists && sid > 0 {
		result = gormDB.Table(usersRelationTableName).
			Where("weid = ? AND reid = ?", userInfo["weid"], sid).
			Delete(nil)

		if result.Error != nil {
			log.Printf("删除店铺关联失败: %v", result.Error)
		}
	}

	return nil
}

// GetUuidReID 从UUID关系中获取关联ID
func GetUuidReID(uid int, ptype string) (int, error) {
	relaArray, err := GetUuidRelaArray(uid, ptype)
	if err != nil {
		return 0, err
	}

	if reid, exists := relaArray["reid"].(int); exists {
		return reid, nil
	}

	return 0, nil
}

// GetUuidRelaArray 从UUID关系中获取关联数组
func GetUuidRelaArray(uid int, ptype string) (map[string]interface{}, error) {
	gormDB := storage.GetGormDB()
	if gormDB == nil {
		return nil, errors.New("数据库连接失败")
	}

	var relation structs.UuidRelation
	result := gormDB.Table(uuidRelationTableName).
		Where("weid = ? AND ptype = ? AND uid = ? AND uuid IS NOT NULL AND uuid <> ''", helper.GetWeid(), ptype, uid).
		Order("id DESC").
		First(&relation)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		log.Printf("获取UUID关联数组失败: %v", result.Error)
		return nil, result.Error
	}

	res := make(map[string]interface{})

	// TODO: 实现Users::where方法调用
	// users := Users::where('uuid', relation.Uuid)->find()
	// 暂时设置默认值
	res["uid"] = 0

	// 根据不同的类型获取不同的关联实体
	var reid int = 0

	if ptype == "store" {
		// TODO: 实现Store::where方法调用
	} else if ptype == "tuanzhang" {
		// TODO: 实现Tuanzhang::where方法调用
	} else if ptype == "operatingcity" {
		// TODO: 实现Operatingcity::where方法调用
	} else if ptype == "technical" {
		// TODO: 实现Technical::where方法调用
	}

	res["reid"] = reid

	// 创建用户关系记录
	newRelation := structs.UsersRelation{
		Weid:  helper.GetWeid(),
		Ptype: ptype,
		Reid:  reid,
		Uid:   res["uid"].(int),
		Mid:   uid,
	}

	result = gormDB.Table(usersRelationTableName).Create(&newRelation)
	if result.Error != nil {
		log.Printf("创建用户关系记录失败: %v", result.Error)
	}

	// 删除UUID关系记录
	result = gormDB.Table(uuidRelationTableName).
		Where("weid = ? AND ptype = ? AND uid = ?", helper.GetWeid(), ptype, uid).
		Delete(nil)

	if result.Error != nil {
		log.Printf("删除UUID关系记录失败: %v", result.Error)
	}

	return res, nil
}

// GetUuidUid 从UUID关系中获取用户ID
func GetUuidUid(reid int, ptype string) (int, error) {
	gormDB := storage.GetGormDB()
	if gormDB == nil {
		return 0, errors.New("数据库连接失败")
	}

	var uuid string = ""

	// 根据不同的类型获取不同的关联实体UUID
	if ptype == "store" {
		// TODO: 实现Store::find方法调用
	} else if ptype == "tuanzhang" {
		// TODO: 实现Tuanzhang::find方法调用
	} else if ptype == "operatingcity" {
		// TODO: 实现Operatingcity::find方法调用
	} else if ptype == "technical" {
		// TODO: 实现Technical::find方法调用
	}

	if uuid == "" {
		return 0, nil
	}

	var relation structs.UuidRelation
	result := gormDB.Table(uuidRelationTableName).
		Where("weid = ? AND uuid = ?", helper.GetWeid(), uuid).
		Order("id DESC").
		First(&relation)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return 0, nil
		}
		log.Printf("获取UUID用户ID失败: %v", result.Error)
		return 0, result.Error
	}

	return relation.Uid, nil
}
