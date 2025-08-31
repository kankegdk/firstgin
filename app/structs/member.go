package structs

// Member 会员模型，对应ims_xm_mallv3_member表
type Member struct {
	ID                 int     `json:"id"`
	UUID               string  `json:"uuid"`
	Weid               int     `json:"weid"`
	Ocid               int     `json:"ocid"`
	Sid                int     `json:"sid"`
	RegType            string  `json:"reg_type"`
	Username           string  `json:"username"`
	Name               string  `json:"name"`
	Address            string  `json:"address"`
	Password           string  `json:"-"` // 密码不返回给前端
	Salt               string  `json:"-"` // 盐值不返回给前端
	Status             int     `json:"status"`
	FreeGoodsCourse    int     `json:"free_goods_course"`
	Nickname           string  `json:"nickname"`
	Xingming           string  `json:"xingming"`
	Sex                int     `json:"sex"`
	Userpic            string  `json:"userpic"`
	Pid                int     `json:"pid"`
	AgentLevel         int     `json:"agent_level"`
	Freeze             float64 `json:"freeze"`
	Balance            float64 `json:"balance"`
	Totleconsumed      float64 `json:"totleconsumed"`
	Points             int     `json:"points"`
	WarmCore           float64 `json:"warm_core"`
	NotreceivedPoints  int     `json:"notreceived_points"`
	CashPoints         int     `json:"cash_points"`
	Wish               int     `json:"wish"`
	Regdate            int64   `json:"regdate"`
	Lastdate           int64   `json:"lastdate"`
	Regip              string  `json:"regip"`
	Lastip             string  `json:"lastip"`
	Loginnum           int     `json:"loginnum"`
	Email              string  `json:"email"`
	Telephone          string  `json:"telephone"`
	Gid                int     `json:"gid"`
	Level              int     `json:"level"`
	Areaid             int     `json:"areaid"`
	Message            int     `json:"message"`
	Islock             int     `json:"islock"`
	Resume             string  `json:"resume"`
	Country            string  `json:"country"`
	Province           string  `json:"province"`
	City               string  `json:"city"`
	CategoryID         int     `json:"category_id"`
	Customtext         string  `json:"customtext"`
	Sort               int     `json:"sort"`
	HasWarmCard        int     `json:"has_warm_card"`
}