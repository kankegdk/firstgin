package structs

// Address 用户收货地址模型
type Address struct {
	ID               int    `gorm:"column:id" json:"id"`
	Weid             int    `gorm:"column:weid" json:"weid"`
	Uid              int    `gorm:"column:uid" json:"uid"`
	Name             string `gorm:"column:name" json:"name"`
	Telephone        string `gorm:"column:telephone" json:"telephone"`
	Address          string `gorm:"column:address" json:"address"`
	RegionName       string `gorm:"column:region_name" json:"region_name"`
	ProvinceName     string `gorm:"column:province_name" json:"province_name"`
	CityName         string `gorm:"column:city_name" json:"city_name"`
	DistrictName     string `gorm:"column:district_name" json:"district_name"`
	IsBindingaddress int    `gorm:"column:is_bindingaddress" json:"is_bindingaddress"`
	IsDefault        int    `gorm:"column:isDefault" json:"is_default"`
	Street           int    `gorm:"column:street" json:"street"`
	CityId           int    `gorm:"column:city_id" json:"city_id"`
	DistrictId       int    `gorm:"column:district_id" json:"district_id"`
	ProvinceId       int    `gorm:"column:province_id" json:"province_id"`
	Longitude        string `gorm:"column:longitude" json:"longitude"`
	Latitude         string `gorm:"column:latitude" json:"latitude"`
}
