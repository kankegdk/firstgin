package helper

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"myapi/app/config"
	"myapi/app/structs"

	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims 定义JWT声明结构
type JWTClaims struct {
	UserID   int                 `json:"user_id"`
	Username string              `json:"username"`
	Phone    string              `json:"phone"`
	Member   map[string]interface{} `json:"member,omitempty"`
	jwt.RegisteredClaims
}

// 生成RSA密钥对并保存到文件
func GenerateRSAKeyPair(privateKeyPath, publicKeyPath string) error {
	// 生成2048位RSA私钥
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("生成RSA私钥失败: %v", err)
	}

	// 确保目录存在
	privateKeyDir := filepath.Dir(privateKeyPath)
	if err := os.MkdirAll(privateKeyDir, 0700); err != nil {
		return fmt.Errorf("创建私钥目录失败: %v", err)
	}

	publicKeyDir := filepath.Dir(publicKeyPath)
	if err := os.MkdirAll(publicKeyDir, 0755); err != nil {
		return fmt.Errorf("创建公钥目录失败: %v", err)
	}

	// 保存私钥到文件
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: privateKeyBytes,
		},
	)
	if err := os.WriteFile(privateKeyPath, privateKeyPEM, 0600); err != nil {
		return fmt.Errorf("保存私钥文件失败: %v", err)
	}

	// 保存公钥到文件
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return fmt.Errorf("编码公钥失败: %v", err)
	}
	publicKeyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: publicKeyBytes,
		},
	)
	if err := os.WriteFile(publicKeyPath, publicKeyPEM, 0644); err != nil {
		return fmt.Errorf("保存公钥文件失败: %v", err)
	}

	return nil
}

// 从文件加载RSA私钥
func LoadPrivateKey(privateKeyPath string) (*rsa.PrivateKey, error) {
	log.Printf("尝试加载私钥: %s", privateKeyPath)

	// 检查文件是否存在
	if _, err := os.Stat(privateKeyPath); err != nil {
		log.Printf("私钥文件不存在: %v", err)
		return nil, fmt.Errorf("私钥文件不存在: %v", err)
	}

	privateKeyData, err := os.ReadFile(privateKeyPath)
	if err != nil {
		log.Printf("读取私钥文件失败: %v", err)
		return nil, fmt.Errorf("读取私钥文件失败: %v", err)
	}
	log.Printf("成功读取私钥文件，大小: %d 字节", len(privateKeyData))

	block, rest := pem.Decode(privateKeyData)
	if block == nil {
		log.Printf("无法解码PEM数据，剩余字节: %d", len(rest))
		return nil, fmt.Errorf("无效的私钥数据: 无法解码PEM格式")
	}
	log.Printf("成功解码PEM数据，类型: %s", block.Type)

	// 获取配置的JWT_SECRET作为解密密钥
	cfg := config.GetConfig()
	encryptionKey := cfg.JWTSecret
	if encryptionKey == "" {
		log.Printf("警告: JWT_SECRET为空，使用默认密钥")
		encryptionKey = "default_encryption_key"
	} else {
		log.Printf("使用配置的JWT_SECRET作为解密密钥")
	}

	// 注释掉敏感日志
	// log.Printf("JWT_SECRET: %s", encryptionKey)

	var privateKeyBytes []byte

	// 根据私钥类型处理
	if block.Type == "ENCRYPTED RSA PRIVATE KEY" {
		log.Printf("检测到加密的私钥，尝试解密...")
		// 解密私钥数据

		privateKeyBytes, err = decryptData(block.Bytes, encryptionKey)
		if err != nil {
			log.Printf("解密私钥失败: %v", err)
			return nil, fmt.Errorf("解密私钥失败: %v", err)
		}
		log.Printf("成功解密私钥数据，大小: %d 字节", len(privateKeyBytes))
	} else if block.Type == "RSA PRIVATE KEY" {
		log.Printf("检测到未加密的私钥，直接使用...")
		// 未加密的私钥，直接使用
		privateKeyBytes = block.Bytes
	} else {
		log.Printf("不支持的私钥类型: %s", block.Type)
		return nil, fmt.Errorf("不支持的私钥类型: %s", block.Type)
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBytes)
	if err != nil {
		log.Printf("解析私钥失败: %v", err)
		// 尝试其他格式的解析
		parsedKey, err2 := x509.ParsePKCS8PrivateKey(privateKeyBytes)
		if err2 != nil {
			log.Printf("尝试解析PKCS8格式失败: %v", err2)
			return nil, fmt.Errorf("解析私钥失败: %v (PKCS8也失败: %v)", err, err2)
		}
		// 检查是否是RSA密钥
		if rsaKey, ok := parsedKey.(*rsa.PrivateKey); ok {
			log.Printf("成功解析PKCS8格式的RSA私钥")
			return rsaKey, nil
		}
		return nil, fmt.Errorf("解析的密钥不是RSA类型")
	}

	log.Printf("成功解析RSA私钥")
	return privateKey, nil
}

// 使用环境配置的密钥解密数据
func decryptData(encryptedData []byte, encryptionKey string) ([]byte, error) {
	// 确保密钥长度为16, 24或32字节（AES要求）
	key := []byte(encryptionKey)
	keyLen := len(key)
	switch {
	case keyLen < 16:
		// 不足16字节，填充到16字节
		key = append(key, make([]byte, 16-keyLen)...)
	case keyLen < 24:
		// 不足24字节，填充到24字节
		key = append(key, make([]byte, 24-keyLen)...)
	case keyLen < 32:
		// 不足32字节，填充到32字节
		key = append(key, make([]byte, 32-keyLen)...)
	default:
		// 超过32字节，截断到32字节
		key = key[:32]
	}

	// 创建AES解密块
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// 提取nonce
	nonceSize := gcm.NonceSize()
	if len(encryptedData) < nonceSize {
		return nil, fmt.Errorf("加密数据太短")
	}

	nonce, ciphertext := encryptedData[:nonceSize], encryptedData[nonceSize:]

	// 解密数据
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("解密失败: %v", err)
	}

	return plaintext, nil
}

// 从文件加载RSA公钥
func LoadPublicKey(publicKeyPath string) (*rsa.PublicKey, error) {
	publicKeyData, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("读取公钥文件失败: %v", err)
	}

	block, _ := pem.Decode(publicKeyData)
	if block == nil || block.Type != "RSA PUBLIC KEY" {
		return nil, fmt.Errorf("无效的公钥数据")
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("解析公钥失败: %v", err)
	}

	rsaPublicKey, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("不是RSA公钥")
	}

	return rsaPublicKey, nil
}

// GenerateJWT 使用RSA私钥生成JWT令牌
func GenerateJWT(member *structs.Member, privateKeyPath string, expirationTime time.Duration) (string, error) {
	// 加载私钥
	privateKey, err := LoadPrivateKey(privateKeyPath)
	// log.Println("privateKey:", privateKey, "path:", privateKeyPath)
	if err != nil {
		// 如果加载失败，尝试生成新的密钥对
		if strings.Contains(err.Error(), "读取私钥文件失败") {
			publicKeyPath := strings.Replace(privateKeyPath, "private.pem", "public.pem", 1)
			if genErr := GenerateRSAKeyPair(privateKeyPath, publicKeyPath); genErr != nil {
				return "", fmt.Errorf("生成密钥对失败: %v", genErr)
			}
			// 重新加载私钥
			if privateKey, err = LoadPrivateKey(privateKeyPath); err != nil {
				return "", fmt.Errorf("加载新生成的私钥失败: %v", err)
			}
		} else {
			return "", err
		}
	}

	// 设置令牌过期时间
	expiration := time.Now().Add(expirationTime)

	// 将Member结构体转换为map[string]interface{}
	memberMap := make(map[string]interface{})
	memberJSON, err := json.Marshal(member)
	if err != nil {
		log.Printf("转换Member到JSON失败: %v", err)
		// 如果转换失败，手动构建map
		memberMap["id"] = member.ID
		memberMap["username"] = member.Username
		memberMap["telephone"] = member.Telephone
		memberMap["nickname"] = member.Nickname
		memberMap["status"] = member.Status
	} else {
		// 将JSON转换为map
		if err := json.Unmarshal(memberJSON, &memberMap); err != nil {
			log.Printf("JSON转换为map失败: %v", err)
			// 如果转换失败，手动构建map
			memberMap["id"] = member.ID
			memberMap["username"] = member.Username
			memberMap["telephone"] = member.Telephone
			memberMap["nickname"] = member.Nickname
			memberMap["status"] = member.Status
		}
	}

	// 创建声明
	claims := &JWTClaims{
		UserID:   member.ID,
		Username: member.Username,
		Phone:    member.Telephone,
		Member:   memberMap,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiration),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	// 创建令牌
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	// 使用私钥签名令牌
	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return "", fmt.Errorf("签名令牌失败: %v", err)
	}

	return tokenString, nil
}

// ParseJWT 使用RSA公钥解析JWT令牌
func ParseJWT(tokenString string, publicKeyPath string) (*JWTClaims, error) {
	// 加载公钥
	publicKey, err := LoadPublicKey(publicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("加载公钥失败: %v", err)
	}

	// 解析令牌
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("意外的签名方法: %v", token.Header["alg"])
		}
		return publicKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("解析令牌失败: %v", err)
	}

	// 验证令牌是否有效
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("无效的令牌")
}
