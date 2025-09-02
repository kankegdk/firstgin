package structs

// Member 会员模型，对应member表
type Member struct {
	ID                int     `gorm:"column:id" json:"id"`
	UUID              string  `gorm:"column:uuid" json:"uuid"`
	Weid              int     `gorm:"column:weid" json:"weid"`
	Ocid              int     `gorm:"column:ocid" json:"ocid"`
	Sid               int     `gorm:"column:sid" json:"sid"`
	RegType           string  `gorm:"column:reg_type" json:"reg_type"`
	Username          string  `gorm:"column:username" json:"username"`
	Name              string  `gorm:"column:name" json:"name"`
	Address           string  `gorm:"column:address" json:"address"`
	Password          string  `gorm:"column:password" json:"-"` // 密码不返回给前端
	Salt              string  `gorm:"column:salt" json:"-"` // 盐值不返回给前端
	Status            int     `gorm:"column:status" json:"status"`
	FreeGoodsCourse   int     `gorm:"column:free_goods_course" json:"free_goods_course"`
	Nickname          string  `gorm:"column:nickname" json:"nickname"`
	Xingming          string  `gorm:"column:xingming" json:"xingming"`
	Sex               int     `gorm:"column:sex" json:"sex"`
	Userpic           string  `gorm:"column:userpic" json:"userpic"`
	Pid               int     `gorm:"column:pid" json:"pid"`
	AgentLevel        int     `gorm:"column:agent_level" json:"agent_level"`
	Freeze            float64 `gorm:"column:freeze" json:"freeze"`
	Balance           float64 `gorm:"column:balance" json:"balance"`
	Totleconsumed     float64 `gorm:"column:totleconsumed" json:"totleconsumed"`
	Points            int     `gorm:"column:points" json:"points"`
	WarmCore          float64 `gorm:"column:warm_core" json:"warm_core"`
	NotreceivedPoints int     `gorm:"column:notreceived_points" json:"notreceived_points"`
	CashPoints        int     `gorm:"column:cash_points" json:"cash_points"`
	Wish              int     `gorm:"column:wish" json:"wish"`
	Regdate           int64   `gorm:"column:regdate" json:"regdate"`
	Lastdate          int64   `gorm:"column:lastdate" json:"lastdate"`
	Regip             string  `gorm:"column:regip" json:"regip"`
	Lastip            string  `gorm:"column:lastip" json:"lastip"`
	Loginnum          int     `gorm:"column:loginnum" json:"loginnum"`
	Email             string  `gorm:"column:email" json:"email"`
	Telephone         string  `gorm:"column:telephone" json:"telephone"`
	Gid               int     `gorm:"column:gid" json:"gid"`
	Level             int     `gorm:"column:level" json:"level"`
	Areaid            int     `gorm:"column:areaid" json:"areaid"`
	Message           int     `gorm:"column:message" json:"message"`
	Islock            int     `gorm:"column:islock" json:"islock"`
	Resume            string  `gorm:"column:resume" json:"resume"`
	Country           string  `gorm:"column:country" json:"country"`
	Province          string  `gorm:"column:province" json:"province"`
	City              string  `gorm:"column:city" json:"city"`
	CategoryID        int     `gorm:"column:category_id" json:"category_id"`
	Customtext        string  `gorm:"column:customtext" json:"customtext"`
	Sort              int     `gorm:"column:sort" json:"sort"`
	HasWarmCard       int     `gorm:"column:has_warm_card" json:"has_warm_card"`
}
