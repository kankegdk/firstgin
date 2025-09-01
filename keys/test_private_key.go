package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"myapi/app/config"
	"myapi/app/helper"
)

// 自定义解密函数，添加详细调试
func customDecryptData(encryptedData []byte, encryptionKey string) ([]byte, error) {
	// 确保密钥长度为16, 24或32字节（AES要求）
	key := []byte(encryptionKey)
	keyLen := len(key)
	fmt.Printf("原始密钥长度: %d, 内容: %s\n", keyLen, encryptionKey)

	switch {
	case keyLen < 16:
		// 不足16字节，填充到16字节
		fmt.Printf("密钥不足16字节，填充到16字节\n")
		key = append(key, make([]byte, 16-keyLen)...)
	case keyLen < 24:
		// 不足24字节，填充到24字节
		fmt.Printf("密钥不足24字节，填充到24字节\n")
		key = append(key, make([]byte, 24-keyLen)...)
	case keyLen < 32:
		// 不足32字节，填充到32字节
		fmt.Printf("密钥不足32字节，填充到32字节\n")
		key = append(key, make([]byte, 32-keyLen)...)
	default:
		// 超过32字节，截断到32字节
		fmt.Printf("密钥超过32字节，截断到32字节\n")
		key = key[:32]
	}
	fmt.Printf("处理后密钥长度: %d\n", len(key))

	// 创建AES解密块
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("创建AES解密块失败: %v", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("创建GCM失败: %v", err)
	}

	// 提取nonce
	nonceSize := gcm.NonceSize()
	fmt.Printf("nonceSize: %d\n", nonceSize)
	if len(encryptedData) < nonceSize {
		return nil, fmt.Errorf("加密数据太短: %d 字节", len(encryptedData))
	}

	nonce, ciphertext := encryptedData[:nonceSize], encryptedData[nonceSize:]
	fmt.Printf("提取的nonce长度: %d, ciphertext长度: %d\n", len(nonce), len(ciphertext))

	// 解密数据
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("解密失败: %v\n\t检查加密和解密使用的密钥是否一致", err)
	}

	fmt.Printf("解密成功，明文长度: %d\n", len(plaintext))
	return plaintext, nil
}

// 手动测试解密私钥
func manualDecryptTest(privateKeyPath string, encryptionKey string) (*rsa.PrivateKey, error) {
	fmt.Println("\n=== 手动解密测试 ===")
	fmt.Printf("使用密钥: %s\n", encryptionKey)

	// 读取私钥文件
	privateKeyData, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("读取私钥文件失败: %v", err)
	}
	fmt.Printf("成功读取私钥文件，大小: %d 字节\n", len(privateKeyData))

	// 解码PEM
	block, _ := pem.Decode(privateKeyData)
	if block == nil {
		return nil, fmt.Errorf("无效的PEM数据")
	}
	fmt.Printf("PEM类型: %s, 数据大小: %d 字节\n", block.Type, len(block.Bytes))

	// 尝试解密
	privateKeyBytes, err := customDecryptData(block.Bytes, encryptionKey)
	if err != nil {
		return nil, fmt.Errorf("解密失败: %v", err)
	}

	// 尝试解析私钥
	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBytes)
	if err != nil {
		fmt.Printf("解析PKCS1私钥失败: %v，尝试PKCS8格式\n", err)
		parsedKey, err2 := x509.ParsePKCS8PrivateKey(privateKeyBytes)
		if err2 != nil {
			return nil, fmt.Errorf("解析PKCS8私钥也失败: %v", err2)
		}
		if rsaKey, ok := parsedKey.(*rsa.PrivateKey); ok {
			fmt.Println("成功解析为PKCS8格式的RSA私钥")
			return rsaKey, nil
		}
		return nil, fmt.Errorf("解析的不是RSA私钥")
	}

	fmt.Println("成功解析为PKCS1格式的RSA私钥")
	return privateKey, nil
}

func main() {
	// 初始化配置
	config.Init()

	// 获取配置
	cfg := config.GetConfig()
	fmt.Println("配置信息:")
	fmt.Printf("JWTSecret: %s\n", cfg.JWTSecret)
	fmt.Printf("JWTPrivateKeyPath: %s\n", cfg.JWTPrivateKeyPath)
	fmt.Printf("JWTPublicKeyPath: %s\n", cfg.JWTPublicKeyPath)

	// 检查文件是否存在
	privateKeyPath := cfg.JWTPrivateKeyPath
	if _, err := os.Stat(privateKeyPath); err != nil {
		fmt.Printf("错误: 私钥文件不存在: %v\n", err)
		// 尝试使用相对路径
		absPath, err := filepath.Abs("jwt_private.pem")
		if err != nil {
			fmt.Printf("获取绝对路径失败: %v\n", err)
			return
		}
		fmt.Printf("尝试使用相对路径: %s\n", absPath)
		privateKeyPath = absPath
	}

	// 测试使用不同的密钥解密
	// 1. 使用配置中的密钥
	privateKey, err := manualDecryptTest(privateKeyPath, cfg.JWTSecret)
	if err != nil {
		fmt.Printf("\n使用配置密钥解密失败: %v\n", err)
		// 2. 尝试使用默认密钥
		fmt.Println("\n尝试使用默认密钥 'default_encryption_key' 解密...")
		privateKey, err = manualDecryptTest(privateKeyPath, "default_encryption_key")
		if err != nil {
			fmt.Printf("使用默认密钥解密也失败: %v\n", err)
			return
		}
	}

	if privateKey == nil {
		fmt.Println("警告: 私钥加载成功但返回nil")
	} else {
		fmt.Println("\n私钥加载成功!")
		fmt.Printf("私钥类型: %T\n", privateKey)
		fmt.Printf("私钥位数: %d\n", privateKey.N.BitLen())
	}

	// 如果解密成功，测试JWT生成
	if privateKey != nil {
		fmt.Println("\n测试生成JWT...")
		token, err := helper.GenerateJWT(1, privateKeyPath, 24*time.Hour)
		if err != nil {
			fmt.Printf("生成JWT失败: %v\n", err)
		} else {
			fmt.Println("JWT生成成功:")
			fmt.Println(token)
		}
	}
}
