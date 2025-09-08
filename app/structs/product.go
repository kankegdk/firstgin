package structs

// Product 商品信息表，对应ims_goods表
type Product struct {
	ID                   int     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Weid                 int     `gorm:"column:weid" json:"weid"`
	Ocid                 int     `gorm:"column:ocid" json:"ocid"`
	Sid                  int     `gorm:"column:sid" json:"sid"`           // 店铺ID
	CardTid              int     `gorm:"column:card_tid" json:"card_tid"` // 卡类型id
	Mgid                 int     `gorm:"column:mgid" json:"mgid"`         // 会员组id
	Name                 string  `gorm:"column:name" json:"name"`         // 名称
	Keyword              string  `gorm:"column:keyword" json:"keyword"`
	CatID                int     `gorm:"column:cat_id" json:"cat_id"`                   // 分类ID
	Ptype                int     `gorm:"column:ptype" json:"ptype"`                     // 类型
	BrandID              int     `gorm:"column:brand_id" json:"brand_id"`               // 品牌编号
	Model                string  `gorm:"column:model" json:"model"`                     // 商品型号
	QuantityUnit         string  `gorm:"column:quantity_unit" json:"quantity_unit"`     // 单位
	Location             string  `gorm:"column:location" json:"location"`               // 产品所在地
	Tel                  string  `gorm:"column:tel" json:"tel"`                         // 联系电话
	Quantity             int     `gorm:"column:quantity" json:"quantity"`               // 商品数目
	SaleCount            int     `gorm:"column:sale_count" json:"sale_count"`           // 销量
	SaleCountBase        int     `gorm:"column:sale_count_base" json:"sale_count_base"` // 销量基数
	StockStatusID        int     `gorm:"column:stock_status_id" json:"stock_status_id"` // 库存状态编号
	Image                string  `gorm:"column:image" json:"image"`
	Videotype            int     `gorm:"column:videotype" json:"videotype"` // 1服务器2腾讯视频
	Videoid              string  `gorm:"column:videoid" json:"videoid"`
	Videourl             string  `gorm:"column:videourl" json:"videourl"`
	InstallVideourl      string  `gorm:"column:install_videourl" json:"install_videourl"` // 安装视频地址
	Shipping             int     `gorm:"column:shipping" json:"shipping"`                 // 是否需要运送
	Price                float64 `gorm:"column:price" json:"price"`                       // 商品价格
	IsTimer              int     `gorm:"column:is_timer" json:"is_timer"`                 // 计时
	TimeAmount           int     `gorm:"column:time_amount" json:"time_amount"`           // 服务时长分钟
	OriginalPrice        float64 `gorm:"column:original_price" json:"original_price"`     // 原价
	Costs                float64 `gorm:"column:costs" json:"costs"`                       // 成本
	PointsMethod         int     `gorm:"column:points_method" json:"points_method"`
	Points               int     `gorm:"column:points" json:"points"`                                 // 购买得积分
	PayPoints            int     `gorm:"column:pay_points" json:"pay_points"`                         // 兑换所需积分
	PointsPrice          float64 `gorm:"column:points_price" json:"points_price"`                     // 积分可以抵扣金额
	IsPointsGoods        int     `gorm:"column:is_points_goods" json:"is_points_goods"`               // 是否是积分兑换商品
	CouponID             int     `gorm:"column:coupon_id" json:"coupon_id"`                           // 购买送优惠券id
	CouponNumber         int     `gorm:"column:coupon_number" json:"coupon_number"`                   // 购买送优惠券数量
	DateAvailable        int     `gorm:"column:date_available" json:"date_available"`                 // 供货日期
	Weight               float64 `gorm:"column:weight" json:"weight"`                                 // 重量
	Length               float64 `gorm:"column:length" json:"length"`                                 // 尺寸(长)
	Width                float64 `gorm:"column:width" json:"width"`                                   // 尺寸(宽)
	Height               float64 `gorm:"column:height" json:"height"`                                 // 尺寸(高)
	Subtract             int     `gorm:"column:subtract" json:"subtract"`                             // 是否扣除库存
	Minimum              int     `gorm:"column:minimum" json:"minimum"`                               // 最小起订数目
	Sort                 int     `gorm:"column:sort" json:"sort"`                                     // 排序
	IsTiming             int     `gorm:"column:is_timing" json:"is_timing"`                           // 是否定期服务
	TimingUnit           string  `gorm:"column:timing_unit" json:"timing_unit"`                       // 定期周期
	TimingMum            int     `gorm:"column:timing_mum" json:"timing_mum"`                         // 周期可约次数
	Timesmum             int     `gorm:"column:timesmum" json:"timesmum"`                             // 次数
	Takeeffecttime       int     `gorm:"column:takeeffecttime" json:"takeeffecttime"`                 // 购买后多久生效
	IsTimes              int     `gorm:"column:is_times" json:"is_times"`                             // 0单次服务1次卡2包月包年服务
	Extraprice           float64 `gorm:"column:extraprice" json:"extraprice"`                         // 购物卡赠送金额
	EffectiveTime        int     `gorm:"column:effective_time" json:"effective_time"`                 // 有效时间天
	IsHot                int     `gorm:"column:is_hot" json:"is_hot"`                                 // 热卖
	IsDiscount           int     `gorm:"column:is_discount" json:"is_discount"`                       // 促销
	IsRecommended        int     `gorm:"column:is_recommended" json:"is_recommended"`                 // 推荐
	IsNew                int     `gorm:"column:is_new" json:"is_new"`                                 // 新产品
	IsMemberDiscount     int     `gorm:"column:is_member_discount" json:"is_member_discount"`         // 开启会员组价格
	IsCommission         int     `gorm:"column:is_commission" json:"is_commission"`                   // 是否独立佣金
	IsComment            int     `gorm:"column:is_comment" json:"is_comment"`                         // 是否可评论
	IsCombination        int     `gorm:"column:is_combination" json:"is_combination"`                 // 是否为套装商品:0:否, 1:是
	IsTechnical          int     `gorm:"column:is_technical" json:"is_technical"`                     // 师傅默认产品
	IsSkumore            int     `gorm:"column:is_skumore" json:"is_skumore"`                         // 是否多规格下单
	IsAdditional         int     `gorm:"column:is_additional" json:"is_additional"`                   // 是否支持尾款
	CombinationIDs       string  `gorm:"column:combination_ids" json:"combination_ids"`               // 套装商品id组
	MemberDiscountMethod int     `gorm:"column:member_discount_method" json:"member_discount_method"` // 折扣方式，0.折扣 1.加减金额
	CommissionMethod     int     `gorm:"column:commission_method" json:"commission_method"`           // 分佣方式，0.折扣 1.固定金额
	CommissionPrice      float64 `gorm:"column:commission_price" json:"commission_price"`             // 独立佣金
	Viewed               int     `gorm:"column:viewed" json:"viewed"`                                 // 点击量
	ViewedBase           int     `gorm:"column:viewed_base" json:"viewed_base"`                       // 点击量基数
	ProvinceName         string  `gorm:"column:province_name" json:"province_name"`
	CityName             string  `gorm:"column:city_name" json:"city_name"`
	DistrictName         string  `gorm:"column:district_name" json:"district_name"`
	ProvinceID           int     `gorm:"column:province_id" json:"province_id"`
	CityID               int     `gorm:"column:city_id" json:"city_id"`
	DistrictID           int     `gorm:"column:district_id" json:"district_id"`
	CreateTime           int     `gorm:"column:create_time" json:"create_time"` // 添加时间
	UpdateTime           int     `gorm:"column:update_time" json:"update_time"` // 更新时间
	Status               int     `gorm:"column:status" json:"status"`           // 是否公开
}

// GoodsBuynowinfo 商品立即购买信息表，对应ims_goods_buynowinfo表
type GoodsBuynowinfo struct {
	ID         int    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Weid       int    `gorm:"column:weid" json:"weid"`
	Ip         string `gorm:"column:ip" json:"ip"`
	Data       string `gorm:"column:data" json:"data"`
	ExpireTime int64  `gorm:"column:expire_time" json:"expire_time"` // 使用int64以匹配time.Now().Unix()的返回类型
	Status     int    `gorm:"column:status" json:"status"`           // 状态
}

// GoodsImage 商品图片表
type GoodsImage struct {
	ID      int    `gorm:"column:id" json:"id"`
	GoodsID int    `gorm:"column:goods_id" json:"goods_id"`
	Image   string `gorm:"column:image" json:"image"`
	Sort    int    `gorm:"column:sort" json:"sort"`
}

// GoodsDescription 产品信息描述表，对应ims_goods_description表
type GoodsDescription struct {
	ID              int    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	GoodsID         int    `gorm:"column:goods_id" json:"goods_id"`
	Summary         string `gorm:"column:summary" json:"summary"`         // 产品简介
	Description     string `gorm:"column:description" json:"description"` // 商品详情
	MetaDescription string `gorm:"column:meta_description" json:"meta_description"`
	MetaKeyword     string `gorm:"column:meta_keyword" json:"meta_keyword"`
}

// GoodsDiscount 数量折扣表，对应ims_goods_discount表
type GoodsDiscount struct {
	ID       int     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	GoodsID  int     `gorm:"column:goods_id;index" json:"goods_id"`
	Quantity int     `gorm:"column:quantity" json:"quantity"`
	Price    float64 `gorm:"column:price" json:"price"`
}

// Category 商品分类表，对应ims_xm_mallv3_category表
type Category struct {
	ID                int    `gorm:"column:id;primaryKey;autoIncrement;unsigned" json:"id"` // ID
	Weid              int    `gorm:"column:weid" json:"weid"`
	Pid               int    `gorm:"column:pid;index" json:"pid"`                         // 上级分类ID
	Title             string `gorm:"column:title" json:"title"`                           // 标题
	Goodsnum          int    `gorm:"column:goodsnum" json:"goodsnum"`                     // 商品数量
	IsBinding         int    `gorm:"column:is_binding" json:"is_binding"`                 // 选定师傅才能下单
	IsOrdercount      int    `gorm:"column:is_ordercount" json:"is_ordercount"`           // 是否统计订单
	IsStoragelocation int    `gorm:"column:is_storagelocation" json:"is_storagelocation"` // 是否用存放位置
	Ptype             int    `gorm:"column:ptype" json:"ptype"`                           // 类型
	ServicetimePtype  int    `gorm:"column:servicetime_ptype" json:"servicetime_ptype"`
	Deliverymode      string `gorm:"column:deliverymode" json:"deliverymode"`         // 交付方式
	SubmitOrderTxt    string `gorm:"column:submit_order_txt" json:"submit_order_txt"` // 下单按钮文字
	Ordergoodsremark  string `gorm:"column:ordergoodsremark" json:"ordergoodsremark"`
	Order1remark      string `gorm:"column:order1remark" json:"order1remark"`
	Image             string `gorm:"column:image" json:"image"`
	MetaKeyword       string `gorm:"column:meta_keyword" json:"meta_keyword"`
	MetaDescription   string `gorm:"column:meta_description" json:"meta_description"`
	Sort              int    `gorm:"column:sort" json:"sort"`
	CreateTime        int    `gorm:"column:create_time" json:"create_time"`
	UpdateTime        int    `gorm:"column:update_time" json:"update_time"`
	Status            int    `gorm:"column:status" json:"status"`
	IsShowHome        int    `gorm:"column:is_show_home" json:"is_show_home"` // 是否在首页展示：0-不展示 1-展示
}

// GoodsAttribute 商品属性表，对应ims_goods_attribute表
type GoodsAttribute struct {
	ID               int `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	GoodsID          int `gorm:"column:goods_id" json:"goods_id"`
	AttributeValueID int `gorm:"column:attribute_value_id" json:"attribute_value_id"`
}

// GoodsBrand 商品品牌表，对应ims_goods_brand表
type GoodsBrand struct {
	ID      int `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	GoodsID int `gorm:"column:goods_id" json:"goods_id"`
	BrandID int `gorm:"column:brand_id" json:"brand_id"`
}

// BuyNowInfoParams 立即购买参数
type BuyNowInfoParams struct {
	GoodsID    int    `json:"goodsId" binding:"required"`
	Msid       string `json:"msid"`
	Tuanid     string `json:"tuanid"`
	Tecid      string `json:"tecid"`
	Jointuanid string `json:"jointuanid"`
	Sku        string `json:"sku"`
	IsSkumore  string `json:"is_skumore"`
	Skumore    string `json:"skumore"`
	BuyNumber  string `json:"buyNumber"`
	Weid       int    `json:"weid"`
	Ip         string `json:"ip"`
}

// BuyNowInfoData 立即购买返回数据
type BuyNowInfoData struct {
	*Product                        // 使用指针类型，与代码中的使用一致
	IsCombination     string        `json:"is_combination"`
	Pic               string        `json:"pic"`
	GoodsLength       int           `json:"goodslength"`
	Jointuanid        string        `json:"jointuanid"`
	Miaosha           interface{}   `json:"miaosha"`
	Category          *Category     `json:"category"` // 使用指针类型，避免空结构体
	Deliverymode      string        `json:"deliverymode"`
	Deliverymodearray []interface{} `json:"deliverymodearray"`
	Buynowinfoid      int           `json:"buynowinfoid"`
	ShippingType      string        `json:"shipping_type"` // 配送方式
	RecordId          int           `json:"record_id"`     // 立即购买记录ID
}

// GoodsCombination 商品套装组合表，对应ims_goods_combination表
type GoodsCombination struct {
	ID       int `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	GoodsID  int `gorm:"column:goods_id" json:"goods_id"`   // 商品id
	ParentID int `gorm:"column:parent_id" json:"parent_id"` // 所属套装主id
	Numbers  int `gorm:"column:numbers" json:"numbers"`     // 数量
}

// GoodsCourse 课程表，对应ims_goods_course表
type GoodsCourse struct {
	ID         int    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Title      string `gorm:"column:title" json:"title"`
	Image      string `gorm:"column:image" json:"image"`
	SourceType int    `gorm:"column:source_type" json:"source_type"` // 1图文2音频3视频
	GoodsID    int    `gorm:"column:goods_id" json:"goods_id"`
	Url        string `gorm:"column:url" json:"url"`
	Content    string `gorm:"column:content" json:"content"`
	IsTry      int    `gorm:"column:is_try" json:"is_try"` // 试看
	Sort       int    `gorm:"column:sort" json:"sort"`
	CreateTime int    `gorm:"column:create_time" json:"create_time"`
	UpdateTime int    `gorm:"column:update_time" json:"update_time"`
	Status     int    `gorm:"column:status" json:"status"`
}

// GoodsGiftcardCommission 购物卡分佣表，对应ims_goods_giftcard_commission表
type GoodsGiftcardCommission struct {
	ID               int     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	CardTid          int     `gorm:"column:card_tid" json:"card_tid"`
	CommissionMethod int     `gorm:"column:commission_method" json:"commission_method"`
	Roletype         string  `gorm:"column:roletype" json:"roletype"`
	ReturnPercent    float64 `gorm:"column:return_percent" json:"return_percent"`
}

// GoodsGiftcardType 购物卡类型表，对应ims_goods_giftcard_type表
type GoodsGiftcardType struct {
	ID               int     `gorm:"column:id;primaryKey;autoIncrement;unsigned" json:"id"`
	Ptype            int     `gorm:"column:ptype" json:"ptype"`
	Weid             int     `gorm:"column:weid" json:"weid"`
	Ocid             int     `gorm:"column:ocid" json:"ocid"`
	Sid              int     `gorm:"column:sid" json:"sid"`
	Name             string  `gorm:"column:name" json:"name"`           // 名称
	BuyPrice         float64 `gorm:"column:buy_price" json:"buy_price"` // 购买价格
	Facevalue        float64 `gorm:"column:facevalue" json:"facevalue"` // 面值
	PointsMethod     int     `gorm:"column:points_method" json:"points_method"`
	Points           int     `gorm:"column:points" json:"points"`                       // 购买得积分
	CommissionMethod int     `gorm:"column:commission_method" json:"commission_method"` // 佣金结算方式
	Color            string  `gorm:"column:color" json:"color"`                         // 颜色
	UseGoods         int     `gorm:"column:use_goods" json:"use_goods"`                 // 是否适用商品
	CatIDs           string  `gorm:"column:cat_ids" json:"cat_ids"`                     // 适用品类
	GoodsIDs         string  `gorm:"column:goods_ids" json:"goods_ids"`                 // 适用产品编号
	ConditionType    int     `gorm:"column:condition_type" json:"condition_type"`       // 使用门槛0无门槛1有门槛
	MinPrice         float64 `gorm:"column:min_price;unsigned" json:"min_price"`        // 最低消费金额
	ExpireType       int     `gorm:"column:expire_type;unsigned" json:"expire_type"`    // 到期类型(10领取后生效 20固定时间)
	ExpireDay        int     `gorm:"column:expire_day;unsigned" json:"expire_day"`      // 领取后生效-有效天数
	StartTime        int     `gorm:"column:start_time;unsigned" json:"start_time"`      // 固定时间-开始时间
	EndTime          int     `gorm:"column:end_time;unsigned" json:"end_time"`          // 固定时间-结束时间
	TotalNum         int     `gorm:"column:total_num" json:"total_num"`                 // 发放总数量(-1为不限制)
	ReceiveNum       int     `gorm:"column:receive_num;unsigned" json:"receive_num"`    // 已领取数量
	Sort             int     `gorm:"column:sort;unsigned" json:"sort"`                  // 排序方式(数字越小越靠前)
	IsDelete         int     `gorm:"column:is_delete;unsigned" json:"is_delete"`        // 软删除
	CreateTime       int     `gorm:"column:create_time;unsigned" json:"create_time"`    // 创建时间
	UpdateTime       int     `gorm:"column:update_time;unsigned" json:"update_time"`    // 更新时间
	Status           int     `gorm:"column:status" json:"status"`                       // 状态
}

// GoodsMemberDiscount 会员分组价格表，对应ims_goods_member_discount表
type GoodsMemberDiscount struct {
	ID             int     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	GoodsID        int     `gorm:"column:goods_id;index" json:"goods_id"`
	Mgid           int     `gorm:"column:mgid" json:"mgid"`                       // 会员组
	DiscountMethod int     `gorm:"column:discount_method" json:"discount_method"` // 折扣方式，0.折扣 1.加减金额
	Addsubtract    int     `gorm:"column:addsubtract" json:"addsubtract"`         // 0:加1减
	Price          float64 `gorm:"column:price;type:decimal(15,2)" json:"price"`
}

// GoodsQuantityUnit 数量单位表，对应ims_goods_quantity_unit表
type GoodsQuantityUnit struct {
	ID         int    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Weid       int    `gorm:"column:weid" json:"weid"`
	Ptype      int    `gorm:"column:ptype" json:"ptype"`
	Title      string `gorm:"column:title" json:"title"`
	CreateTime int    `gorm:"column:create_time" json:"create_time"`
	UpdateTime int    `gorm:"column:update_time" json:"update_time"`
	Sort       int    `gorm:"column:sort" json:"sort"`
	Status     int    `gorm:"column:status" json:"status"`
}

// GoodsReply 评论表，对应ims_goods_reply表
type GoodsReply struct {
	ID                   int    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`                // 评论ID
	Uid                  int    `gorm:"column:uid" json:"uid"`                                       // 用户ID
	Oid                  int    `gorm:"column:oid" json:"oid"`                                       // 订单ID
	Unique               string `gorm:"column:unique" json:"unique"`                                 // 唯一id
	GoodsID              int    `gorm:"column:goods_id" json:"goods_id"`                             // 商品id
	ReplyType            string `gorm:"column:reply_type;index" json:"reply_type"`                   // 某种商品类型(普通商品、秒杀商品）
	ProductScore         int    `gorm:"column:product_score;index" json:"product_score"`             // 商品分数
	ServiceScore         int    `gorm:"column:service_score;index" json:"service_score"`             // 服务分数
	Comment              string `gorm:"column:comment" json:"comment"`                               // 评论内容
	Pics                 string `gorm:"column:pics" json:"pics"`                                     // 评论图片
	AddTime              int    `gorm:"column:add_time;index" json:"add_time"`                       // 评论时间
	MerchantReplyContent string `gorm:"column:merchant_reply_content" json:"merchant_reply_content"` // 管理员回复内容
	MerchantReplyTime    int    `gorm:"column:merchant_reply_time" json:"merchant_reply_time"`       // 管理员回复时间
	IsDel                int    `gorm:"column:is_del;index" json:"is_del"`                           // 0未删除1已删除
	IsReply              int    `gorm:"column:is_reply" json:"is_reply"`                             // 0未回复1已回复
	Nickname             string `gorm:"column:nickname" json:"nickname"`                             // 用户名称
	Avatar               string `gorm:"column:avatar" json:"avatar"`                                 // 用户头像
	Suk                  string `gorm:"column:suk" json:"suk"`                                       // 规格名称
	Status               int    `gorm:"column:status" json:"status"`                                 // 评论状态
}

// GoodsSku 商品SKU表，对应ims_goods_sku表
type GoodsSku struct {
	ID       int    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	GoodsID  int    `gorm:"column:goods_id" json:"goods_id"`
	Name     string `gorm:"column:name" json:"name"`
	Item     string `gorm:"column:item" json:"item"`
	Ptype    string `gorm:"column:ptype" json:"ptype"`
	Required int    `gorm:"column:required" json:"required"`
}

// GoodsSkuValue 商品SKU值表，对应ims_goods_sku_value表
type GoodsSkuValue struct {
	ID       int     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	GoodsID  int     `gorm:"column:goods_id" json:"goods_id"`
	Sku      string  `gorm:"column:sku" json:"sku"`
	Image    string  `gorm:"column:image" json:"image"`
	Quantity int     `gorm:"column:quantity" json:"quantity"`
	Price    float64 `gorm:"column:price;" json:"price"`
}

// GoodsTimeDiscount 时间段价格表，对应ims_goods_time_discount表
type GoodsTimeDiscount struct {
	ID             int     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	GoodsID        int     `gorm:"column:goods_id;index" json:"goods_id"`
	DiscountMethod int     `gorm:"column:discount_method" json:"discount_method"` // 折扣方式，0.折扣 1.加减金额
	BeginTime      string  `gorm:"column:begin_time" json:"begin_time"`
	EndTime        string  `gorm:"column:end_time" json:"end_time"`
	Addsubtract    int     `gorm:"column:addsubtract" json:"addsubtract"` // 0:加1减
	Price          float64 `gorm:"column:price;type:decimal(15,2)" json:"price"`
}

// GoodsToCategory 商品对应分类表，对应ims_goods_to_category表
type GoodsToCategory struct {
	GoodsID    int `gorm:"column:goods_id;primaryKey" json:"goods_id"`
	CategoryID int `gorm:"column:category_id;primaryKey" json:"category_id"`
}
