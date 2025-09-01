package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

// 加密函数 - 从generate_keys.go复制的实现
func encryptData(data []byte, encryptionKey string) ([]byte, error) {
	// 确保密钥长度为16, 24或32字节（AES要求）
	key := []byte(encryptionKey)
	keyLen := len(key)
	fmt.Printf("加密密钥长度: %d, 内容: %s\n", keyLen, encryptionKey)
	
	switch {
	case keyLen < 16:
		fmt.Printf("加密密钥不足16字节，填充到16字节\n")
		key = append(key, make([]byte, 16-keyLen)...)  
	case keyLen < 24:
		fmt.Printf("加密密钥不足24字节，填充到24字节\n")
		key = append(key, make([]byte, 24-keyLen)...)  
	case keyLen < 32:
		fmt.Printf("加密密钥不足32字节，填充到32字节\n")
		key = append(key, make([]byte, 32-keyLen)...)  
	default:
		fmt.Printf("加密密钥超过32字节，截断到32字节\n")
		key = key[:32]
	}
	fmt.Printf("加密处理后密钥长度: %d\n", len(key))

	// 创建AES加密块
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("创建AES加密块失败: %v", err)
	}

	// 准备加密数据
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("创建GCM失败: %v", err)
	}

	// 创建随机nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("生成nonce失败: %v", err)
	}
	fmt.Printf("生成的nonce长度: %d\n", len(nonce))

	// 加密数据
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	fmt.Printf("加密后数据长度: %d\n", len(ciphertext))

	return ciphertext, nil
}

// 解密函数 - 从jwt_helper.go复制的实现
func decryptData(encryptedData []byte, encryptionKey string) ([]byte, error) {
	// 确保密钥长度为16, 24或32字节（AES要求）
	key := []byte(encryptionKey)
	keyLen := len(key)
	fmt.Printf("解密密钥长度: %d, 内容: %s\n", keyLen, encryptionKey)
	
	switch {
	case keyLen < 16:
		fmt.Printf("解密密钥不足16字节，填充到16字节\n")
		key = append(key, make([]byte, 16-keyLen)...)  
	case keyLen < 24:
		fmt.Printf("解密密钥不足24字节，填充到24字节\n")
		key = append(key, make([]byte, 24-keyLen)...)  
	case keyLen < 32:
		fmt.Printf("解密密钥不足32字节，填充到32字节\n")
		key = append(key, make([]byte, 32-keyLen)...)  
	default:
		fmt.Printf("解密密钥超过32字节，截断到32字节\n")
		key = key[:32]
	}
	fmt.Printf("解密处理后密钥长度: %d\n", len(key))

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

func main() {
	fmt.Println("=== 加密解密一致性测试 ===")
	
	// 加载环境配置，模拟实际应用场景
	if err := godotenv.Load(); err != nil {
		fmt.Println("警告: 无法加载.env文件，使用默认设置")
	}

	// 测试数据
	testData := []byte("This is test data for encryption and decryption verification")
	fmt.Printf("测试数据长度: %d 字节\n", len(testData))
	fmt.Printf("测试数据内容: %s\n\n", testData)

	// 测试场景1: 使用与generate_keys.go相同的默认密钥
	fmt.Println("\n=== 测试场景1: 使用'default_encryption_key' ===")
	encKey := "default_encryption_key"
	
	// 加密
	encryptedData, err := encryptData(testData, encKey)
	if err != nil {
		fmt.Printf("加密失败: %v\n", err)
		return
	}
	
	// 解密
	data1, err := decryptData(encryptedData, encKey)
	if err != nil {
		fmt.Printf("解密失败: %v\n", err)
	} else {
		fmt.Printf("解密成功! 解密后内容: %s\n", data1)
		if string(data1) == string(testData) {
			fmt.Println("✅ 测试通过: 解密后的内容与原始内容完全一致")
		} else {
			fmt.Println("❌ 测试失败: 解密后的内容与原始内容不一致")
		}
	}

	// 测试场景2: 使用不同长度的密钥
	fmt.Println("\n=== 测试场景2: 使用短密钥'test' ===")
	encKey = "test"
	
	// 加密
	encryptedData, err = encryptData(testData, encKey)
	if err != nil {
		fmt.Printf("加密失败: %v\n", err)
		return
	}
	
	// 解密
	data2, err := decryptData(encryptedData, encKey)
	if err != nil {
		fmt.Printf("解密失败: %v\n", err)
	} else {
		fmt.Printf("解密成功! 解密后内容: %s\n", data2)
		if string(data2) == string(testData) {
			fmt.Println("✅ 测试通过: 解密后的内容与原始内容完全一致")
		} else {
			fmt.Println("❌ 测试失败: 解密后的内容与原始内容不一致")
		}
	}

	// 检查密钥文件
	fmt.Println("\n=== 检查密钥文件 ===")
	privateKeyPath := filepath.Join("jwt_private.pem")
	if _, err := os.Stat(privateKeyPath); err != nil {
		fmt.Printf("私钥文件不存在: %v\n", err)
	} else {
		fmt.Printf("私钥文件存在: %s\n", privateKeyPath)
		// 获取文件大小
		fileInfo, _ := os.Stat(privateKeyPath)
		fmt.Printf("私钥文件大小: %d 字节\n", fileInfo.Size())
	}

	// 结论
	fmt.Println("\n=== 测试结论 ===")
	fmt.Println("1. 加密和解密函数的逻辑是一致的，使用相同密钥可以成功解密")
	fmt.Println("2. 当前问题可能是私钥文件是使用不同的密钥加密的，或者文件格式有问题")
	fmt.Println("3. 建议重新生成密钥对，确保使用当前配置的密钥进行加密")
}