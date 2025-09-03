package structs

// Operatingcity 运营城市模型，与数据库表 ims_operatingcity 对应
type Operatingcity struct {
	ID           int     `gorm:"column:id" json:"id"`
	UUID         string  `gorm:"column:uuid" json:"uuid"`
	Weid         int     `gorm:"column:weid" json:"weid"`
	Uid          int     `gorm:"column:uid" json:"uid"`
	Level        int8    `gorm:"column:level" json:"level"`
	Areatype     int8    `gorm:"column:areatype" json:"areatype"`
	Title        string  `gorm:"column:title" json:"title"`
	Income       float64 `gorm:"column:income" json:"income"`
	TotalIncome  float64 `gorm:"column:total_income" json:"total_income"`
	CateIDs      string  `gorm:"column:cate_ids" json:"cate_ids"`
	RegionName   string  `gorm:"column:region_name" json:"region_name"`
	Latitude     string  `gorm:"column:latitude" json:"latitude"`
	Longitude    string  `gorm:"column:longitude" json:"longitude"`
	ProvinceID   int     `gorm:"column:province_id" json:"province_id"`
	ProvinceName string  `gorm:"column:province_name" json:"province_name"`
	CityID       int     `gorm:"column:city_id" json:"city_id"`
	CityName     string  `gorm:"column:city_name" json:"city_name"`
	DistrictID   int     `gorm:"column:district_id" json:"district_id"`
	DistrictName string  `gorm:"column:district_name" json:"district_name"`
	AreaName     string  `gorm:"column:area_name" json:"area_name"`
	HouseNumber  string  `gorm:"column:house_number" json:"house_number"`
	Tel          string  `gorm:"column:tel" json:"tel"`
	CustomText   string  `gorm:"column:customtext" json:"customtext"`
	Settings     string  `gorm:"column:settings" json:"settings"`
	CreateTime   int     `gorm:"column:create_time" json:"create_time"`
	UpdateTime   int     `gorm:"column:update_time" json:"update_time"`
	EndTime      int     `gorm:"column:end_time" json:"end_time"`
	IsDefault    int8    `gorm:"column:is_default" json:"is_default"`
	Sort         int     `gorm:"column:sort" json:"sort"`
	Status       int8    `gorm:"column:status" json:"status"`
}

// OperatingcityIncomelog 运营城市收入日志模型，与数据库表 ims_operatingcity_incomelog 对应
type OperatingcityIncomelog struct {
	ID            int     `gorm:"column:id" json:"id"`
	Weid          int     `gorm:"column:weid" json:"weid"`
	Ocid          int     `gorm:"column:ocid" json:"ocid"`
	Level         int     `gorm:"column:level" json:"level"`
	OrderID       int     `gorm:"column:order_id" json:"order_id"`
	Ptype         int     `gorm:"column:ptype" json:"ptype"`
	Areatype      int8    `gorm:"column:areatype" json:"areatype"`
	OrderNumAlias string  `gorm:"column:order_num_alias" json:"order_num_alias"`
	BuyerID       int     `gorm:"column:buyer_id" json:"buyer_id"`
	Income        float64 `gorm:"column:income" json:"income"`
	ReturnPercent float64 `gorm:"column:return_percent" json:"return_percent"`
	PercentRemark string  `gorm:"column:percentremark" json:"percentremark"`
	Description   string  `gorm:"column:description" json:"description"`
	OrderTotal    float64 `gorm:"column:order_total" json:"order_total"`
	PayTime       int     `gorm:"column:pay_time" json:"pay_time"`
	CreateTime    int     `gorm:"column:create_time" json:"create_time"`
	MonthTime     string  `gorm:"column:month_time" json:"month_time"`
	YearTime      string  `gorm:"column:year_time" json:"year_time"`
	OrderStatusID int     `gorm:"column:order_status_id" json:"order_status_id"`
}

// OperatingcityLevel 运营城市等级模型，与数据库表 ims_operatingcity_level 对应
type OperatingcityLevel struct {
	ID            int     `gorm:"column:id" json:"id"`
	Weid          int     `gorm:"column:weid" json:"weid"`
	Title         string  `gorm:"column:title" json:"title"`
	ReturnPercent float64 `gorm:"column:return_percent" json:"return_percent"`
	Sort          int     `gorm:"column:sort" json:"sort"`
	Status        int8    `gorm:"column:status" json:"status"`
}
