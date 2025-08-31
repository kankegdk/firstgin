package structs

// Address 用户收货地址模型
type Address struct {
	ID              int    `json:"id"`
	Weid            int    `json:"weid"`
	Uid             int    `json:"uid"`
	Name            string `json:"name"`
	Telephone       string `json:"telephone"`
	Address         string `json:"address"`
	RegionName      string `json:"region_name"`
	ProvinceName    string `json:"province_name"`
	CityName        string `json:"city_name"`
	DistrictName    string `json:"district_name"`
	IsBindingaddress int    `json:"is_bindingaddress"`
	IsDefault       int    `json:"is_default"`
	Street          int    `json:"street"`
	CityId          int    `json:"city_id"`
	DistrictId      int    `json:"district_id"`
	ProvinceId      int    `json:"province_id"`
	Longitude       string `json:"longitude"`
	Latitude        string `json:"latitude"`
}