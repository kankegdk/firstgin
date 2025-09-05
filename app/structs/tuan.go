package structs

// TuanGoods 团购商品表，对应ims_tuan_goods表
type TuanGoods struct {
	ID           int     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Weid         int     `gorm:"column:weid" json:"weid"`
	Sid          int     `gorm:"column:sid" json:"sid"`
	Title        string  `gorm:"column:title" json:"title"`
	Ocid         int     `gorm:"column:ocid" json:"ocid"`
	GoodsID      int     `gorm:"column:goods_id" json:"goods_id"`
	SaleCount    int     `gorm:"column:sale_count" json:"sale_count"`
	PeopleNum    int     `gorm:"column:people_num" json:"people_num"`
	RobotNum     int     `gorm:"column:robot_num" json:"robot_num"`
	AutoInitiate int     `gorm:"column:auto_initiate" json:"auto_initiate"`
	IsLuckyDraw  int     `gorm:"column:is_luckydraw" json:"is_luckydraw"`
	LuckyDrawNum int     `gorm:"column:luckydraw_num" json:"luckydraw_num"`
	TimeLimit    int     `gorm:"column:time_limit" json:"time_limit"`
	Price        float64 `gorm:"column:price" json:"price"`
	TuanMinPrice float64 `gorm:"column:tuan_min_price" json:"tuan_min_price"`
	TuanMaxPrice float64 `gorm:"column:tuan_max_price" json:"tuan_max_price"`
	EndDate      int     `gorm:"column:end_date" json:"end_date"`
	BeginDate    int     `gorm:"column:begin_date" json:"begin_date"`
	BuyMax       int     `gorm:"column:buy_max" json:"buy_max"`
	BuyLimit     int     `gorm:"column:buy_limit" json:"buy_limit"`
	IsLevel      int     `gorm:"column:is_level" json:"is_level"`
	CreateTime   int     `gorm:"column:create_time" json:"create_time"`
	UpdateTime   int     `gorm:"column:update_time" json:"update_time"`
	Sort         int     `gorm:"column:sort" json:"sort"`
	Status       int     `gorm:"column:status" json:"status"`
}

// TuanGoodsSkuValue 团购商品SKU值表，对应ims_tuan_goods_sku_value表
type TuanGoodsSkuValue struct {
	ID       int     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	TuanID   int     `gorm:"column:tuan_id" json:"tuan_id"`
	GoodsID  int     `gorm:"column:goods_id" json:"goods_id"`
	Sku      string  `gorm:"column:sku" json:"sku"`
	Image    string  `gorm:"column:image" json:"image"`
	Quantity int     `gorm:"column:quantity" json:"quantity"`
	Price    float64 `gorm:"column:price" json:"price"`
}

// TuanFollow 参团记录表，对应ims_tuan_follow表
type TuanFollow struct {
	ID          int    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Uid         int    `gorm:"column:uid" json:"uid"`
	Nickname    string `gorm:"column:nickname" json:"nickname"`
	Avatar      string `gorm:"column:avatar" json:"avatar"`
	JoinTime    int    `gorm:"column:join_time" json:"join_time"`
	OrderID     int    `gorm:"column:order_id" json:"order_id"`
	FoundID     int    `gorm:"column:found_id" json:"found_id"`
	TuanID      int    `gorm:"column:tuan_id" json:"tuan_id"`
	IsHead      int    `gorm:"column:is_head" json:"is_head"`
	IsRobot     int    `gorm:"column:is_robot" json:"is_robot"`
	PayTime     int    `gorm:"column:pay_time" json:"pay_time"`
	IsRefund    int    `gorm:"column:is_refund" json:"is_refund"`
	Status      int    `gorm:"column:status" json:"status"`
	TuanEndTime int    `gorm:"column:tuan_end_time" json:"tuan_end_time"`
}

// TuanFound 开团记录表，对应ims_tuan_found表
type TuanFound struct {
	ID           int    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Weid         int    `gorm:"column:weid" json:"weid"`
	Sid          int    `gorm:"column:sid" json:"sid"`
	Ocid         int    `gorm:"column:ocid" json:"ocid"`
	Sn           string `gorm:"column:sn" json:"sn"`
	FoundTime    int    `gorm:"column:found_time" json:"found_time"`
	FoundEndTime int    `gorm:"column:found_end_time" json:"found_end_time"`
	Uid          int    `gorm:"column:uid" json:"uid"`
	TuanID       int    `gorm:"column:tuan_id" json:"tuan_id"`
	Nickname     string `gorm:"column:nickname" json:"nickname"`
	Avatar       string `gorm:"column:avatar" json:"avatar"`
	Join         int    `gorm:"column:join" json:"join"`
	Need         int    `gorm:"column:need" json:"need"`
	PayTime      int    `gorm:"column:pay_time" json:"pay_time"`
	Status       int    `gorm:"column:status" json:"status"`
	TuanEndTime  int    `gorm:"column:tuan_end_time" json:"tuan_end_time"`
}

// TransportExtend 运费模板扩展表，对应ims_transport_extend表
type TransportExtend struct {
	ID            int     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Weid          int     `gorm:"column:weid" json:"weid"`
	AreaID        string  `gorm:"column:area_id" json:"area_id"`
	TopAreaID     string  `gorm:"column:top_area_id" json:"top_area_id"`
	AreaName      string  `gorm:"column:area_name" json:"area_name"`
	Snum          int     `gorm:"column:snum" json:"snum"`
	Sprice        float64 `gorm:"column:sprice" json:"sprice"`
	Xnum          int     `gorm:"column:xnum" json:"xnum"`
	Xprice        float64 `gorm:"column:xprice" json:"xprice"`
	IsDefault     int     `gorm:"column:is_default" json:"is_default"`
	Sort          int     `gorm:"column:sort" json:"sort"`
	Status        int     `gorm:"column:status" json:"status"`
	TransportID   int     `gorm:"column:transport_id" json:"transport_id"`
	TransportTitle string  `gorm:"column:transport_title" json:"transport_title"`
}

// Tuanzhang 团长表，对应ims_tuanzhang表
type Tuanzhang struct {
	ID             int     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Uuid           string  `gorm:"column:uuid" json:"uuid"`
	Weid           int     `gorm:"column:weid" json:"weid"`
	Ocid           int     `gorm:"column:ocid" json:"ocid"`
	Uid            int     `gorm:"column:uid" json:"uid"`
	Answer         string  `gorm:"column:answer" json:"answer"`
	Tel            string  `gorm:"column:tel" json:"tel"`
	CommunityTitle string  `gorm:"column:community_title" json:"community_title"`
	Title          string  `gorm:"column:title" json:"title"`
	Level          int     `gorm:"column:level" json:"level"`
	Touxiang       string  `gorm:"column:touxiang" json:"touxiang"`
	ServiceTimes   int     `gorm:"column:service_times" json:"service_times"`
	CateIds        string  `gorm:"column:cate_ids" json:"cate_ids"`
	TotalIncome    float64 `gorm:"column:total_income" json:"total_income"`
	Income         float64 `gorm:"column:income" json:"income"`
	Points         int     `gorm:"column:points" json:"points"`
	Email          string  `gorm:"column:email" json:"email"`
	Introduction   string  `gorm:"column:introduction" json:"introduction"`
	IdCart         string  `gorm:"column:id_cart" json:"id_cart"`
	Imageurl1      string  `gorm:"column:imageurl1" json:"imageurl1"`
	Imageurl2      string  `gorm:"column:imageurl2" json:"imageurl2"`
	Imageurl3      string  `gorm:"column:imageurl3" json:"imageurl3"`
	Imageurl4      string  `gorm:"column:imageurl4" json:"imageurl4"`
	ProvinceName   string  `gorm:"column:province_name" json:"province_name"`
	CityName       string  `gorm:"column:city_name" json:"city_name"`
	DistrictName   string  `gorm:"column:district_name" json:"district_name"`
	RegionName     string  `gorm:"column:region_name" json:"region_name"`
	Longitude      string  `gorm:"column:longitude" json:"longitude"`
	Latitude       string  `gorm:"column:latitude" json:"latitude"`
	Dizhi          string  `gorm:"column:dizhi" json:"dizhi"`
	HouseNumber    string  `gorm:"column:house_number" json:"house_number"`
	Customtext     string  `gorm:"column:customtext" json:"customtext"`
	CreateTime     int     `gorm:"column:create_time" json:"create_time"`
	IsBusiness     int     `gorm:"column:is_business" json:"is_business"`
	IsStore        int     `gorm:"column:is_store" json:"is_store"`
	Sort           int     `gorm:"column:sort" json:"sort"`
	Status         int     `gorm:"column:status" json:"status"`
}

// TuanzhangIncomeLog 团长收入明细表，对应ims_tuanzhang_incomelog表
type TuanzhangIncomeLog struct {
	ID             int     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Weid           int     `gorm:"column:weid" json:"weid"`
	Tzid           int     `gorm:"column:tzid" json:"tzid"`
	Uuid           string  `gorm:"column:uuid" json:"uuid"`
	Level          int     `gorm:"column:level" json:"level"`
	OrderID        int     `gorm:"column:order_id" json:"order_id"`
	Ptype          int     `gorm:"column:ptype" json:"ptype"`
	OrderNumAlias  string  `gorm:"column:order_num_alias" json:"order_num_alias"`
	BuyerID        int     `gorm:"column:buyer_id" json:"buyer_id"`
	Income         float64 `gorm:"column:income" json:"income"`
	ReturnPercent  float64 `gorm:"column:return_percent" json:"return_percent"`
	Percentremark  string  `gorm:"column:percentremark" json:"percentremark"`
	Description    string  `gorm:"column:description" json:"description"`
	OrderTotal     float64 `gorm:"column:order_total" json:"order_total"`
	PayTime        int     `gorm:"column:pay_time" json:"pay_time"`
	CreateTime     int     `gorm:"column:create_time" json:"create_time"`
	MonthTime      string  `gorm:"column:month_time" json:"month_time"`
	YearTime       string  `gorm:"column:year_time" json:"year_time"`
	OrderStatusID  int     `gorm:"column:order_status_id" json:"order_status_id"`
}

// TuanzhangLevel 团长等级表，对应ims_tuanzhang_level表
type TuanzhangLevel struct {
	ID            int     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Weid          int     `gorm:"column:weid" json:"weid"`
	Title         string  `gorm:"column:title" json:"title"`
	ReturnPercent float64 `gorm:"column:return_percent" json:"return_percent"`
	Sort          int     `gorm:"column:sort" json:"sort"`
	Status        int     `gorm:"column:status" json:"status"`
}
