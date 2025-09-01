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

	"myapi/app/config"
)

func main() {
	// 初始化配置（使用与主应用相同的配置加载逻辑）
	config.Init()
	cfg := config.GetConfig()

	// 获取密钥保存路径（从配置中读取）
	privateKeyPath := cfg.JWTPrivateKeyPath
	publicKeyPath := cfg.JWTPublicKeyPath

	// 确保密钥文件所在目录存在
	privateKeyDir := filepath.Dir(privateKeyPath)
	if err := os.MkdirAll(privateKeyDir, 0755); err != nil {
		panic("无法创建私钥目录: " + err.Error())
	}

	// 获取加密密钥
	encryptionKey := cfg.JWTSecret
	if encryptionKey == "" {
		encryptionKey = "default_encryption_key" // 默认加密密钥
		fmt.Println("警告: JWT_SECRET为空，使用默认密钥")
	}
	// 注释掉敏感日志
	// println("加密密钥: " + encryptionKey)

	// 生成2048位RSA密钥对
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic("生成RSA密钥失败: " + err.Error())
	}

	// 保存私钥到文件（使用环境配置的密钥进行加密）
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)

	// 使用环境配置的密钥加密私钥
	encryptedPrivateKey, err := encryptData(privateKeyBytes, encryptionKey)
	if err != nil {
		panic("加密私钥失败: " + err.Error())
	}

	privateKeyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "ENCRYPTED RSA PRIVATE KEY",
			Bytes: encryptedPrivateKey,
		},
	)
	if err := os.WriteFile(privateKeyPath, privateKeyPEM, 0600); err != nil {
		panic("保存私钥失败: " + err.Error())
	}

	// 保存公钥到文件
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		panic("序列化公钥失败: " + err.Error())
	}
	publicKeyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: publicKeyBytes,
		},
	)
	if err := os.WriteFile(publicKeyPath, publicKeyPEM, 0644); err != nil {
		panic("保存公钥失败: " + err.Error())
	}

	println("RSA密钥对生成成功!")
	println("私钥: " + privateKeyPath + " (已加密)")
	println("公钥: " + publicKeyPath)
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
