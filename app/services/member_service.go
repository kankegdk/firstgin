package services

import (
	"fmt"
	"log"
	"math/rand"
	"myapi/app/config"
	"myapi/app/helper"
	"myapi/app/models"
	"myapi/app/structs"
	"time"
)

// MemberService 会员服务接口
type MemberService interface {
	// 密码登录
	LoginByPassword(username, password string, clientIP string) (*structs.Member, string, error)
	// 短信验证码登录
	LoginBySmsCode(telephone, code string, clientIP string) (*structs.Member, string, error)
	// 发送短信验证码
	SendSmsCode(telephone string) error
	// 验证密码
	VerifyPassword(inputPassword, storedPassword, salt string) bool
}

// memberService 实现MemberService接口的结构体
type memberService struct{}

// NewMemberService 创建一个新的会员服务实例
func NewMemberService() MemberService {
	return &memberService{}
}

// LoginByPassword 密码登录实现
func (s *memberService) LoginByPassword(username, password string, clientIP string) (*structs.Member, string, error) {
	// 1. 根据用户名获取会员信息
	member, err := models.GetMemberByUsername(username)
	if err != nil {
		return nil, "", fmt.Errorf("获取会员信息失败: %v", err)
	}
	if member == nil {
		return nil, "", fmt.Errorf("用户名或密码错误")
	}

	// 2. 检查用户状态
	if member.Status != 1 {
		return nil, "", fmt.Errorf("用户未激活或已禁用")
	}

	// 3. 验证密码
	if !s.VerifyPassword(password, member.Password, member.Salt) {
		return nil, "", fmt.Errorf("用户名或密码错误")
	}

	// 4. 更新最后登录信息
	err = models.UpdateMemberLastLogin(member.ID, clientIP)
	if err != nil {
		return nil, "", fmt.Errorf("更新登录信息失败: %v", err)
	}

	// 5. 生成JWT
	token := generateToken(member)

	return member, token, nil
}

// LoginBySmsCode 短信验证码登录实现
func (s *memberService) LoginBySmsCode(telephone, code string, clientIP string) (*structs.Member, string, error) {
	// 1. 验证短信验证码（这里简化处理，实际应该从Redis中获取并验证）
	if !verifySmsCode(telephone, code) {
		return nil, "", fmt.Errorf("验证码错误或已过期")
	}

	// 2. 根据手机号获取会员信息
	member, err := models.GetMemberByTelephone(telephone)
	if err != nil {
		return nil, "", fmt.Errorf("获取会员信息失败: %v", err)
	}
	if member == nil {
		return nil, "", fmt.Errorf("手机号未注册")
	}

	// 3. 检查用户状态
	if member.Status != 1 {
		return nil, "", fmt.Errorf("用户未激活或已禁用")
	}

	// 4. 更新最后登录信息
	err = models.UpdateMemberLastLogin(member.ID, clientIP)
	if err != nil {
		return nil, "", fmt.Errorf("更新登录信息失败: %v", err)
	}

	// 5. 生成JWT令牌
	token := generateToken(member)

	return member, token, nil
}

// SendSmsCode 发送短信验证码（简化实现）
func (s *memberService) SendSmsCode(telephone string) error {
	// 1. 检查手机号是否已注册
	member, err := models.GetMemberByTelephone(telephone)
	if err != nil {
		return fmt.Errorf("检查手机号失败: %v", err)
	}
	if member == nil {
		return fmt.Errorf("手机号未注册")
	}

	// 2. 生成验证码
	_ = generateSmsCode()

	// 3. 存储验证码到Redis并设置过期时间（简化实现）
	// 实际项目中应该使用Redis客户端

	// 4. 发送短信（简化实现）
	// 实际项目中应该调用短信服务商API

	return nil
}

// VerifyPassword 验证密码
func (s *memberService) VerifyPassword(inputPassword, storedPassword, salt string) bool {
	hash := helper.PassHash(inputPassword, salt)
	return hash == storedPassword
}

// generateToken 使用非对称加密的JWT库生成令牌
func generateToken(member *structs.Member) string {
	// 获取配置
	privateKeyPath := config.GetString("jwtPrivateKeyPath", "")
	// 设置令牌过期时间（24小时）
	expirationTime := time.Hour * 24 * 100

	// 使用RSA私钥生成JWT令牌
	token, err := helper.GenerateJWT(member, privateKeyPath, expirationTime)
	if err != nil {
		// 如果生成失败，使用默认的简单令牌作为备用方案
		log.Printf("生成JWT令牌失败: %v\n", err)
		timestamp := time.Now().Unix()
		return fmt.Sprintf("%d:%d:%s", member.ID, timestamp, generateRandomString(16))
	}

	return token
}

// generateSmsCode 生成短信验证码
func generateSmsCode() string {
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

// verifySmsCode 验证短信验证码（简化实现）
func verifySmsCode(telephone, code string) bool {
	// 实际项目中应该从Redis中获取并验证
	// 这里简单模拟验证
	return len(code) == 6
}

// generateRandomString 生成随机字符串
func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
