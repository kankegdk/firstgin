package services

import (
	"myapi/app/models"
)

// UserService 用户服务接口
type UserService interface {
	//GetUserByID(id int) *models.User
	GetAllUsers() ([]map[string]interface{}, error)
	//CreateUser(user *models.User) error
}

// userService 实现UserService接口的结构体
type userService struct{}

// NewUserService 创建一个新的用户服务实例
func NewUserService() UserService {
	return &userService{}
}

// // GetUserByID 根据ID获取用户
// func (s *userService) GetUserByID(id int) *models.User {
// 	// 调用模型层方法
// 	return models.GetUserByID(id)
// }

// GetAllUsers 获取所有用户
func (s *userService) GetAllUsers() ([]map[string]interface{}, error) {
	// 调用模型层方法
	user := &models.User{}
	return user.GetAllUsers()
}

// // CreateUser 创建用户
// func (s *userService) CreateUser(user *models.User) error {
// 	// 实际项目中，这里会包含业务逻辑验证
// 	// 然后调用模型层的创建方法
// 	// 这里只是一个示例
// 	return nil
// }
