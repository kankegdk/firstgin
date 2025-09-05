package structs

// MiaoshaGoods 秒杀商品表，对应ims_miaosha_goods表
type MiaoshaGoods struct {
	ID           int     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Weid         int     `gorm:"column:weid" json:"weid"`
	Sid          int     `gorm:"column:sid" json:"sid"`
	Title        string  `gorm:"column:title" json:"title"`
	Ocid         int     `gorm:"column:ocid" json:"ocid"`
	GoodsID      int     `gorm:"column:goods_id" json:"goods_id"`
	Price        float64 `gorm:"column:price" json:"price"`
	SaleCount    int     `gorm:"column:sale_count" json:"sale_count"`
	EndDate      int     `gorm:"column:end_date" json:"end_date"`
	BeginDate    int     `gorm:"column:begin_date" json:"begin_date"`
	BuyMax       int     `gorm:"column:buy_max" json:"buy_max"`
	MemberBuyMax int     `gorm:"column:member_buy_max" json:"member_buy_max"`
	BuyLimit     int     `gorm:"column:buy_limit" json:"buy_limit"`
	IsLevel      int     `gorm:"column:is_level" json:"is_level"`
	CreateTime   int     `gorm:"column:create_time" json:"create_time"`
	UpdateTime   int     `gorm:"column:update_time" json:"update_time"`
	Sort         int     `gorm:"column:sort" json:"sort"`
	Status       int     `gorm:"column:status" json:"status"`
}

// MiaoshaGoodsSkuValue 秒杀商品SKU值表，对应ims_miaosha_goods_sku_value表
type MiaoshaGoodsSkuValue struct {
	ID       int     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	MsID     int     `gorm:"column:ms_id" json:"ms_id"`
	GoodsID  int     `gorm:"column:goods_id" json:"goods_id"`
	Sku      string  `gorm:"column:sku" json:"sku"`
	Image    string  `gorm:"column:image" json:"image"`
	Quantity int     `gorm:"column:quantity" json:"quantity"`
	Price    float64 `gorm:"column:price" json:"price"`
}

// MiaoshaTime 秒杀时间段表，对应ims_miaosha_time表
type MiaoshaTime struct {
	ID         int    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Weid       int    `gorm:"column:weid" json:"weid"`
	BeginTime  string `gorm:"column:begin_time" json:"begin_time"`
	EndTime    string `gorm:"column:end_time" json:"end_time"`
	Sort       int    `gorm:"column:sort" json:"sort"`
	CreateTime int    `gorm:"column:create_time" json:"create_time"`
	UpdateTime int    `gorm:"column:update_time" json:"update_time"`
	Status     int    `gorm:"column:status" json:"status"`
}
