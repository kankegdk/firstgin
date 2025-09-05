package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"myapi/app/helper"

	"github.com/gin-gonic/gin"
)

// TestValidateRequestWithValidPhone 测试使用ValidateRequest验证正确的手机号
func TestValidateRequestWithValidPhone(t *testing.T) {
	// 创建一个测试路由器
	router := gin.Default()

	// 创建一个测试请求处理函数
	router.POST("/test/phone/valid", func(c *gin.Context) {
		// 定义请求结构体
		var requestData struct {
			Telephone string `json:"telephone" binding:"required,len=11"`
			Name      string `json:"name" binding:"required"`
		}

		// 使用ValidateRequest验证请求
		if !helper.ValidateRequest(c, &requestData) {
			return
		}

		// 验证通过，返回成功响应
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "验证成功",
			"data":    requestData,
		})
	})

	// 创建一个测试请求
	reqBody := `{"telephone":"13812345678","name":"张三"}`
	req := httptest.NewRequest("POST", "/test/phone/valid", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// 记录响应
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 验证响应状态码
	if w.Code != http.StatusOK {
		t.Errorf("期望状态码 %d，实际得到 %d", http.StatusOK, w.Code)
	}

	// 解析响应体
	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Errorf("解析响应体失败: %v", err)
	}

	// 验证响应消息
	if message, ok := response["message"].(string); !ok || message != "验证成功" {
		t.Errorf("期望消息 '验证成功'，实际得到 '%v'", response["message"])
	}

	// 验证数据是否正确返回
	if data, ok := response["data"].(map[string]interface{}); ok {
		if telephone, ok := data["Telephone"].(string); !ok || telephone != "13812345678" {
			t.Errorf("期望手机号 '13812345678'，实际得到 '%v'", data["Telephone"])
		}
	} else {
		t.Errorf("响应数据格式错误")
	}
}

// TestValidateRequestWithInvalidPhone 测试使用ValidateRequest验证错误的手机号
func TestValidateRequestWithInvalidPhone(t *testing.T) {
	// 创建一个测试路由器
	router := gin.Default()

	// 创建一个测试请求处理函数
	router.POST("/test/phone/invalid", func(c *gin.Context) {
		// 定义请求结构体
		var requestData struct {
			Telephone string `json:"telephone" binding:"required,len=11"`
			Name      string `json:"name" binding:"required"`
		}

		// 使用ValidateRequest验证请求
		if !helper.ValidateRequest(c, &requestData) {
			return
		}

		// 验证通过，返回成功响应
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "验证成功",
			"data":    requestData,
		})
	})

	// 创建一个测试请求（包含无效的手机号）
	reqBody := `{"telephone":"12345678901","name":"李四"}` // 这不是有效的手机号格式
	req := httptest.NewRequest("POST", "/test/phone/invalid", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// 记录响应
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 验证响应状态码
	if w.Code != http.StatusBadRequest {
		t.Errorf("期望状态码 %d，实际得到 %d", http.StatusBadRequest, w.Code)
	}

	// 解析响应体
	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Errorf("解析响应体失败: %v", err)
	}

	// 验证错误消息
	if errorMsg, ok := response["error"].(string); !ok || errorMsg != "手机号格式错误" {
		t.Errorf("期望错误消息 '手机号格式错误'，实际得到 '%v'", response["error"])
	}
}

// TestValidateRequestWithMap 测试使用map作为参数验证手机号
func TestValidateRequestWithMap(t *testing.T) {
	// 创建一个测试路由器
	router := gin.Default()

	// 创建一个测试请求处理函数
	router.POST("/test/phone/map", func(c *gin.Context) {
		// 定义一个map来接收请求参数
		requestData := make(map[string]interface{})

		// 使用ValidateRequest验证请求
		if !helper.ValidateRequest(c, &requestData) {
			return
		}

		// 验证通过，返回成功响应
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "验证成功",
			"data":    requestData,
		})
	})

	// 创建一个测试请求（包含有效的手机号）
	reqBody := `{"phone":"13812345678","name":"王五"}` // 使用"phone"作为键
	req := httptest.NewRequest("POST", "/test/phone/map", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// 记录响应
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 验证响应状态码
	if w.Code != http.StatusOK {
		t.Errorf("期望状态码 %d，实际得到 %d", http.StatusOK, w.Code)
	}

	// 解析响应体
	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Errorf("解析响应体失败: %v", err)
	}

	// 验证响应消息
	if message, ok := response["message"].(string); !ok || message != "验证成功" {
		t.Errorf("期望消息 '验证成功'，实际得到 '%v'", response["message"])
	}
}

// 运行测试的命令：go test -v myapi/test/validator_test.go
