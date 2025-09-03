package structs

// User 定义了用户表结构体
// 对应表 ims_user
// 用于存储平台用户的基本信息和权限配置

type User struct {
	ID         int    `gorm:"column:id;primaryKey;autoIncrement" json:"id"` // id
	Uuid       string `gorm:"column:uuid;default:'';size:68" json:"uuid"`
	Weid       int    `gorm:"column:weid;not null" json:"weid"`
	Lastweid   int    `gorm:"column:lastweid;default:0" json:"lastweid"` // 最后进入的平台
	Uid        int    `gorm:"column:uid" json:"uid"`
	W7uid      int    `gorm:"column:w7uid;not null" json:"w7uid"`
	Did        int    `gorm:"column:did" json:"did"`           // 部门id
	Tid        int    `gorm:"column:tid;default:0" json:"tid"` // 师傅id
	Sid        int    `gorm:"column:sid;default:0" json:"sid"` // 店id
	Ocid       int    `gorm:"column:ocid;default:0" json:"ocid"`
	Tzid       int    `gorm:"column:tzid;default:0" json:"tzid"`                   // 团长id
	Username   string `gorm:"column:username;default:'';size:50" json:"username"`  // 用户名
	Password   string `gorm:"column:password;default:'';size:100" json:"-"`        // 密码不返回给前端
	Salt       string `gorm:"column:salt;default:'';size:10" json:"-"`             // 盐值不返回给前端
	Touxiang   string `gorm:"column:touxiang;default:'';size:255" json:"touxiang"` // 头像
	Qianming   string `gorm:"column:qianming;default:'';size:200" json:"qianming"`
	Title      string `gorm:"column:title;default:'';size:50" json:"title"` // 姓名
	Sex        int8   `gorm:"column:sex;default:0" json:"sex"`
	Mobile     string `gorm:"column:mobile;size:20" json:"mobile"` // 手机号
	RoleId     int    `gorm:"column:role_id;default:0" json:"role_id"`
	Remark     string `gorm:"column:remark;default:'';size:50" json:"remark"` // 备注
	LoginIp    string `gorm:"column:login_ip;size:30" json:"login_ip"`        // 最近登录IP
	LoginTime  int    `gorm:"column:login_time" json:"login_time"`            // 最近登录时间
	Px         int    `gorm:"column:px;default:0" json:"px"`
	Time       int    `gorm:"column:time;default:0" json:"time"`
	Role       string `gorm:"column:role;size:20" json:"role"` // 微擎权限分组
	CreateTime int    `gorm:"column:create_time;default:0" json:"create_time"`
	UpdateTime int    `gorm:"column:update_time;default:0" json:"update_time"`
	Status     int8   `gorm:"column:status;default:1" json:"status"`
}

type UsersRelation struct {
	ID    int    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Weid  int    `gorm:"column:weid;default:0" json:"weid"`
	Ptype string `gorm:"column:ptype;default:'';size:20" json:"ptype"`
	Mid   int    `gorm:"column:mid;not null" json:"mid"` // 会员id
	Uid   int    `gorm:"column:uid;default:0" json:"uid"`
	Reid  int    `gorm:"column:reid;not null" json:"reid"` // 关联id
}

// UuidRelation 定义了UUID关系结构体
// 对应表 ims_uuid_relation
// 用于通过UUID建立用户与其他实体之间的临时关联

type UuidRelation struct {
	ID    int    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Weid  int    `gorm:"column:weid;default:0" json:"weid"`
	Ptype string `gorm:"column:ptype;default:'';size:20" json:"ptype"`
	Uid   int    `gorm:"column:uid;default:0" json:"uid"`
	Uuid  string `gorm:"column:uuid;default:'';size:68" json:"uuid"`
}

// UsersRoles 定义了用户角色结构体
// 对应表 ims_users_roles
// 用于存储用户角色信息及权限配置

type UsersRoles struct {
	ID          int    `gorm:"column:id;primaryKey;autoIncrement" json:"id"` // 编号
	Uuid        string `gorm:"column:uuid;default:'';size:68" json:"uuid"`
	Pid         int    `gorm:"column:pid;default:0" json:"pid"` // 所属父类
	Weid        int    `gorm:"column:weid;default:0" json:"weid"`
	Sid         int    `gorm:"column:sid;default:0" json:"sid"`
	Ocid        int    `gorm:"column:ocid;default:0" json:"ocid"`
	Tzid        int    `gorm:"column:tzid;default:0" json:"tzid"` // 团长id
	Title       string `gorm:"column:title;size:36" json:"title"` // 分组名称
	Datarules   int8   `gorm:"column:datarules;default:0" json:"datarules"`
	Status      int8   `gorm:"column:status" json:"status"`                     // 状态
	Description string `gorm:"column:description;type:text" json:"description"` // 描述
	IsConsole   int8   `gorm:"column:is_console;default:0" json:"is_console"`   // 控制台
	Access      string `gorm:"column:access;type:text" json:"access"`           // 权限节点
	IsAllrole   int8   `gorm:"column:is_allrole;default:0" json:"is_allrole"`
}

// UsersSessions 定义了用户会话结构体
// 对应表 ims_users_sessions
// 用于存储用户会话信息及登录状态

type UsersSessions struct {
	ID         int    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Weid       int    `gorm:"column:weid;default:0" json:"weid"`
	Ptype      string `gorm:"column:ptype;default:'';size:20" json:"ptype"` // 那个平台
	Token      string `gorm:"column:token;not null;size:68" json:"token"`
	Ip         string `gorm:"column:ip;default:'';size:30" json:"ip"`
	Data       string `gorm:"column:data;type:text" json:"data"`
	ExpireTime int    `gorm:"column:expire_time;default:0" json:"expire_time"`
	LastTime   int    `gorm:"column:last_time;default:0" json:"last_time"`
	Status     int8   `gorm:"column:status;default:1" json:"status"` // 状态
	IsError    int8   `gorm:"column:is_error;default:0" json:"is_error"`
	DevStatus  int8   `gorm:"column:dev_status;default:1" json:"dev_status"` // 1正常 0下线
}
