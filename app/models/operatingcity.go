package models

import (
	"errors"
	"fmt"
	"log"
	"time"

	"myapi/app/config"
	"myapi/app/storage"
	"myapi/app/structs"

	"gorm.io/gorm"
)

// 全局变量存储完整表名
var operatingcityTableName string
var operatingcityIncomelogTableName string

// init函数在包初始化时执行，只配置一次表前缀
func init() {
	// 获取表前缀
	tablePrefix := config.GetString("dbPrefix", "")
	operatingcityTableName = tablePrefix + "operatingcity"
	operatingcityIncomelogTableName = tablePrefix + "operatingcity_incomelog"
}

// GetOperatingcityByID 根据ID获取运营城市信息
func GetOperatingcityByID(id int) ([]map[string]interface{}, error) {
	// 获取共享的gorm连接实例

	var operatingcity []map[string]interface{}
	gormDB := storage.GetGormDB()
	if gormDB == nil {
		log.Println("GORM连接为空")
		return []map[string]interface{}{}, errors.New("数据库连接失败")
	}

	result := gormDB.Table(operatingcityTableName).Where("id = ?", id).First(&operatingcity)
	if result.Error != nil {
		log.Printf("查询运营城市数据失败: %v", result.Error)
		return nil, result.Error
	}

	return operatingcity, nil
}

// GetTitle 根据ID获取运营城市标题
func GetTitle(id int) (string, error) {
	city, err := GetOperatingcityByID(id)
	if err != nil {
		return "", err
	}
	if len(city) == 0 {
		return "", errors.New("未找到城市数据")
	}
	// 从map中获取Title字段，需要类型断言
	if title, ok := city[0]["title"].(string); ok {
		return title, nil
	}
	return "", errors.New("无法获取城市标题")
}

// GetCityID 根据ID获取城市ID
func GetCityID(id int) (int, error) {
	city, err := GetOperatingcityByID(id)
	if err != nil {
		return 0, err
	}
	if len(city) == 0 {
		return 0, errors.New("未找到城市数据")
	}
	// 从map中获取CityID字段，需要类型断言
	if cityID, ok := city[0]["city_id"].(int); ok {
		return cityID, nil
	}
	// 处理可能的数值类型转换问题
	if cityIDFloat, ok := city[0]["city_id"].(float64); ok {
		return int(cityIDFloat), nil
	}
	return 0, errors.New("无法获取城市ID")
}

// GetCityName 根据ID获取城市名称
func GetCityName(id int) (string, error) {
	city, err := GetOperatingcityByID(id)
	if err != nil {
		return "", err
	}
	if len(city) == 0 {
		return "", errors.New("未找到城市数据")
	}
	// 从map中获取CityName字段，需要类型断言
	if cityName, ok := city[0]["city_name"].(string); ok {
		return cityName, nil
	}
	return "", errors.New("无法获取城市名称")
}

// GetSettings 获取运营城市设置
func GetSettings(id int, cache int) (map[string]interface{}, error) {
	// 获取共享的gorm连接实例
	gormDB := storage.GetGormDB()
	if gormDB == nil {
		log.Println("GORM连接为空")
		return nil, errors.New("数据库连接失败")
	}

	var operatingcity structs.Operatingcity
	result := gormDB.Table(operatingcityTableName).Where("id = ?", id).First(&operatingcity)
	if result.Error != nil {
		log.Printf("查询运营城市数据失败: %v", result.Error)
		return map[string]interface{}{}, nil
	}

	// 这里假设settings是JSON字符串，在实际应用中需要根据实际格式进行解析
	settings := make(map[string]interface{})
	settings["id"] = operatingcity.ID
	// 实际应用中需要解析operatingcity.Settings字符串到settings map

	return settings, nil
}

// GetOperatingcityByAreaType 获取指定区域类型和名称的运营城市
func GetOperatingcityByAreaType(areaType int, areaName string) ([]structs.Operatingcity, error) {
	// 获取共享的gorm连接实例
	gormDB := storage.GetGormDB()
	if gormDB == nil {
		log.Println("GORM连接为空")
		return nil, errors.New("数据库连接失败")
	}

	var cities []structs.Operatingcity
	var query *gorm.DB

	switch areaType {
	case 1:
		query = gormDB.Table(operatingcityTableName).Where("areatype = ? AND province_name = ? AND status = ?", areaType, areaName, 1)
	case 2:
		query = gormDB.Table(operatingcityTableName).Where("areatype = ? AND city_name = ? AND status = ?", areaType, areaName, 1)
	case 3:
		query = gormDB.Table(operatingcityTableName).Where("areatype = ? AND district_name = ? AND status = ?", areaType, areaName, 1)
	default:
		return nil, errors.New("无效的区域类型")
	}

	result := query.Order("id asc").Find(&cities)
	if result.Error != nil {
		log.Printf("查询运营城市数据失败: %v", result.Error)
		return nil, result.Error
	}

	return cities, nil
}

// SetIncome 处理订单的城市代理收入分配
func SetIncome(orderInfo map[string]interface{}) error {
	// 处理省级城市代理
	if provinceName, ok := orderInfo["shipping_province_name"].(string); ok && provinceName != "" {
		cities, err := GetOperatingcityByAreaType(1, provinceName)
		if err != nil {
			log.Printf("获取省级城市代理失败: %v", err)
		}
		for _, city := range cities {
			if err := Calculate(orderInfo, city); err != nil {
				log.Printf("计算省级城市代理收入失败: %v", err)
			}
		}
	}

	// 处理市级城市代理
	if cityName, ok := orderInfo["shipping_city_name"].(string); ok && cityName != "" {
		cities, err := GetOperatingcityByAreaType(2, cityName)
		if err != nil {
			log.Printf("获取市级城市代理失败: %v", err)
		}
		for _, city := range cities {
			if err := Calculate(orderInfo, city); err != nil {
				log.Printf("计算市级城市代理收入失败: %v", err)
			}
		}
	}

	// 处理区级城市代理
	if districtName, ok := orderInfo["shipping_district_name"].(string); ok && districtName != "" {
		cities, err := GetOperatingcityByAreaType(3, districtName)
		if err != nil {
			log.Printf("获取区级城市代理失败: %v", err)
		}
		for _, city := range cities {
			if err := Calculate(orderInfo, city); err != nil {
				log.Printf("计算区级城市代理收入失败: %v", err)
			}
		}
	}

	return nil
}

// Calculate 计算并分配城市代理收入
func Calculate(orderInfo map[string]interface{}, operatingcity structs.Operatingcity) error {
	// 获取共享的gorm连接实例
	gormDB := storage.GetGormDB()
	if gormDB == nil {
		log.Println("GORM连接为空")
		return errors.New("数据库连接失败")
	}

	if operatingcity.Status != 1 {
		return nil
	}

	// 获取返点比例，这里需要调用其他模块的方法，暂时使用默认值
	percent := 0.0
	// TODO: 实现OperatingcityLevel::getPercent方法

	if percent <= 0 {
		// TODO: 从配置中获取默认返点比例
		return nil
	}

	// 计算收入
	income := 0.0
	// TODO: 实现佣金计算逻辑，需要调用其他模块的方法

	if income <= 0 {
		return nil
	}

	// 检查佣金是否超额
	// TODO: 实现订单佣金总额检查逻辑

	// 更新订单佣金总额
	// TODO: 实现订单佣金总额更新逻辑

	// 检查收入日志是否已存在
	var incomelog structs.OperatingcityIncomelog
	orderID := 0
	if id, ok := orderInfo["id"].(int); ok {
		orderID = id
	}
	weid := 0
	if w, ok := orderInfo["weid"].(int); ok {
		weid = w
	}

	result := gormDB.Table(operatingcityIncomelogTableName).
		Where("ocid = ? AND areatype = ? AND weid = ? AND order_id = ?",
			operatingcity.ID, operatingcity.Areatype, weid, orderID).
		First(&incomelog)

	if result.Error != nil {
		// 更新城市代理收入
		gormDB.Table(operatingcityTableName).
			Where("id = ?", operatingcity.ID).
			Update("income", gorm.Expr("income + ?", income)).
			Update("total_income", gorm.Expr("total_income + ?", income))

		// 创建收入日志
		currentTime := time.Now()
		incomedata := structs.OperatingcityIncomelog{
			Ocid:          operatingcity.ID,
			Ptype:         1,
			Areatype:      operatingcity.Areatype,
			Weid:          weid,
			OrderID:       orderID,
			BuyerID:       0,
			Income:        income,
			ReturnPercent: percent,
			PayTime:       int(currentTime.Unix()),
			CreateTime:    int(currentTime.Unix()),
			MonthTime:     fmt.Sprintf("%d-%02d", currentTime.Year(), currentTime.Month()),
			YearTime:      fmt.Sprintf("%d", currentTime.Year()),
			OrderStatusID: 2, // 已付款
		}

		// 设置订单号
		if orderNumAlias, ok := orderInfo["order_num_alias"].(string); ok {
			incomedata.OrderNumAlias = orderNumAlias
		}

		// 设置买家ID
		if uid, ok := orderInfo["uid"].(int); ok {
			incomedata.BuyerID = uid
		}

		// 设置订单总额
		if total, ok := orderInfo["total"].(float64); ok {
			incomedata.OrderTotal = total
			incomedata.PercentRemark = fmt.Sprintf("%.2f%% x %.2f%%", total, percent)
		}

		// 创建收入日志
		result = gormDB.Table(operatingcityIncomelogTableName).Create(&incomedata)
		if result.Error != nil {
			log.Printf("创建城市代理收入日志失败: %v", result.Error)
			return result.Error
		}
	}

	return nil
}

// Conversion 转换运营城市数据格式
func Conversion(vo map[string]interface{}) map[string]interface{} {
	// 处理结束时间
	if endTime, ok := vo["end_time"].(int); ok {
		if endTime == 0 {
			vo["end_time"] = "永久有效"
		} else {
			// 转换时间格式
			t := time.Unix(int64(endTime), 0)
			vo["end_time"] = t.Format("2006-01-02 15:04:05")
		}
	}

	// 处理创建时间
	if createTime, ok := vo["create_time"].(int); ok {
		t := time.Unix(int64(createTime), 0)
		vo["create_time_str"] = t.Format("2006-01-02 15:04:05")
	}

	// 处理更新时间
	if updateTime, ok := vo["update_time"].(int); ok {
		t := time.Unix(int64(updateTime), 0)
		vo["update_time_str"] = t.Format("2006-01-02 15:04:05")
	}

	// TODO: 处理分类、等级、区域类型等字段的转换
	// 需要调用其他模块的方法

	return vo
}
