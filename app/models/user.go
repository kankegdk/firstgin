package models

// User 用户模型，定义了用户数据的结构
// 使用json标签来指定JSON序列化时的字段名
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// GetUserByID 模拟从数据库获取用户的方法
// 这是一个模型层的方法，负责具体的业务数据操作
func GetUserByID(id int) *User {
	// 这里为了演示，直接返回模拟数据
	// 实际项目中这里会是数据库查询逻辑
	return &User{
		ID:    id,
		Name:  "张三",
		Email: "zhangsan@example.com",
	}
}

// GetAllUsers 模拟获取所有用户
func GetAllUsers() []User {
	return []User{
		{ID: 1, Name: "张三", Email: "zhangsan@example.com"},
		{ID: 2, Name: "李四", Email: "lisi@example.com"},
		{ID: 3, Name: "王五", Email: "wangwu@example.com"},
	}
}
