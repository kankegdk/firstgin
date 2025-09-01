package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("=== 重新生成RSA密钥对 ===")
	
	// 设置密钥保存路径
	privateKeyPath := "jwt_private.pem"
	publicKeyPath := "jwt_public.pem"

	// 确保keys目录存在
	keysDir := filepath.Dir(privateKeyPath)
	if keysDir != "." {
		if err := os.MkdirAll(keysDir, 0755); err != nil {
			panic("无法创建keys目录: " + err.Error())
		}
	}

	// 加载环境配置
	if err := godotenv.Load(); err != nil {
		fmt.Println("警告: 无法加载.env文件，使用默认设置")
	}

	// 从环境变量获取加密密钥，如果没有则使用默认值
	encryptionKey := os.Getenv("JWT_SECRET")
	if encryptionKey == "" {
		encryptionKey = "default_encryption_key"
		fmt.Println("警告: 未设置JWT_SECRET环境变量，使用默认密钥 'default_encryption_key'")
	} else {
		fmt.Printf("使用环境变量JWT_SECRET作为加密密钥\n")
	}
	fmt.Printf("加密密钥: %s\n", encryptionKey)
	fmt.Printf("加密密钥长度: %d 字节\n\n", len(encryptionKey))

	// 生成2048位RSA密钥对
	fmt.Println("正在生成2048位RSA密钥对...")
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic("生成RSA密钥失败: " + err.Error())
	}
	fmt.Println("RSA密钥对生成成功!")

	// 序列化私钥
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	fmt.Printf("私钥序列化后大小: %d 字节\n", len(privateKeyBytes))

	// 使用环境配置的密钥加密私钥
	fmt.Println("\n正在使用配置的密钥加密私钥...")
	encryptedPrivateKey, err := encryptData(privateKeyBytes, encryptionKey)
	if err != nil {
		panic("加密私钥失败: " + err.Error())
	}
	fmt.Printf("私钥加密后大小: %d 字节\n", len(encryptedPrivateKey))

	// 保存加密后的私钥到文件
	privateKeyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "ENCRYPTED RSA PRIVATE KEY",
			Bytes: encryptedPrivateKey,
		},
	)
	if err := os.WriteFile(privateKeyPath, privateKeyPEM, 0600); err != nil {
		panic("保存私钥失败: " + err.Error())
	}
	fmt.Printf("私钥已保存到: %s (已加密)\n", privateKeyPath)

	// 保存公钥到文件
	fmt.Println("\n正在保存公钥...")
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		panic("序列化公钥失败: " + err.Error())
	}
	fmt.Printf("公钥序列化后大小: %d 字节\n", len(publicKeyBytes))
	
	publicKeyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: publicKeyBytes,
		},
	)
	if err := os.WriteFile(publicKeyPath, publicKeyPEM, 0644); err != nil {
		panic("保存公钥失败: " + err.Error())
	}
	fmt.Printf("公钥已保存到: %s\n", publicKeyPath)

	// 验证生成的密钥对
	fmt.Println("\n=== 验证生成的密钥对 ===")
	fmt.Println("重新生成密钥对成功! 现在可以使用以下命令测试新生成的密钥:")
	fmt.Println("go run test_private_key.go")
	fmt.Println("\n重要提示：")
	fmt.Println("1. 确保在生产环境中设置安全的JWT_SECRET环境变量")
	fmt.Println("2. 备份好生成的密钥文件，丢失后将无法解密现有令牌")
}

// 使用环境配置的密钥加密数据
func encryptData(data []byte, encryptionKey string) ([]byte, error) {
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

	// 创建AES加密块
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// 准备加密数据
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// 创建随机nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// 加密数据
	ciphertext := gcm.Seal(nonce, nonce, data, nil)

	return ciphertext, nil
}