package controllers

import (
	"log"
	"net/http"
	"strconv"

	"myapi/app/helper"
	"myapi/app/services"
	"myapi/app/structs"

	"github.com/gin-gonic/gin"
)

// AddAddress 处理添加新地址的请求
func AddAddress(c *gin.Context) {
	// 1. 创建服务实例
	addressService := services.NewAddressService()

	// 2. 从请求中获取数据
	var address structs.Address
	if err := c.ShouldBindJSON(&address); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 从上下文中获取用户信息
	uid := helper.UID(c)
	if uid == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}
	weid := helper.Weid(c)

	// 设置用户ID和weid
	address.Uid = uid
	address.Weid = weid

	// 4. 调用服务层方法添加地址
	addressId, err := addressService.AddAddress(address)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 5. 返回JSON响应
	c.JSON(http.StatusOK, gin.H{
		"message": "地址添加成功",
		"data": gin.H{
			"address_id": addressId,
		},
	})
}

// UpdateAddress 处理更新地址的请求
func UpdateAddress(c *gin.Context) {
	// 1. 创建服务实例
	addressService := services.NewAddressService()

	// 2. 从请求中获取地址ID
	addressId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的地址ID"})
		return
	}

	// 3. 从请求中获取数据
	var address structs.Address
	if err := c.ShouldBindJSON(&address); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 从上下文中获取用户信息
	uid := helper.UID(c)
	if uid == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}
	weid := helper.Weid(c)

	// 设置用户ID和weid
	address.Uid = uid
	address.Weid = weid

	// 5. 调用服务层方法更新地址
	if err := addressService.UpdateAddress(addressId, address); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 6. 返回JSON响应
	c.JSON(http.StatusOK, gin.H{
		"message": "地址修改成功",
		"data":    addressId,
	})
}

// DeleteAddress 处理删除地址的请求
func DeleteAddress(c *gin.Context) {
	// 1. 创建服务实例
	addressService := services.NewAddressService()

	// 2. 从请求中获取地址ID
	addressId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的地址ID"})
		return
	}

	// 3. 从上下文中获取用户信息
	uid := helper.UID(c)
	if uid == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}
	weid := helper.Weid(c)

	// 4. 调用服务层方法删除地址
	if err := addressService.DeleteAddress(addressId, weid, uid); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 5. 返回JSON响应
	c.JSON(http.StatusOK, gin.H{"message": "地址删除成功"})
}

// SetDefaultAddress 处理设置默认地址的请求
func SetDefaultAddress(c *gin.Context) {
	// 1. 创建服务实例
	addressService := services.NewAddressService()

	// 2. 从请求中获取地址ID
	addressId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的地址ID"})
		return
	}

	uid := helper.UID(c)
	if uid == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}
	weid := helper.Weid(c)

	// 4. 调用服务层方法设置默认地址
	if err := addressService.SetDefaultAddress(addressId, weid, uid); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 5. 返回JSON响应
	c.JSON(http.StatusOK, gin.H{
		"message": "默认地址设置成功",
		"data":    gin.H{"address_id": addressId},
	})
}

// GetAddressDetail 处理获取地址详情的请求
func GetAddressDetail(c *gin.Context) {
	// 1. 创建服务实例
	addressService := services.NewAddressService()

	// 2. 从请求中获取地址ID
	addressId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的地址ID"})
		return
	}

	uid := helper.UID(c)
	if uid == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	weid, weidExists := c.Get("weid")
	if !weidExists {
		weid = 1 // 默认值
	}

	// 转换类型
	uidInt := uid
	weidInt := 1
	if weidIntVal, ok := weid.(int); ok {
		weidInt = weidIntVal
	}

	// 4. 调用服务层方法获取地址详情
	address, err := addressService.GetAddressDetail(addressId, weidInt, uidInt)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 5. 返回JSON响应
	c.JSON(http.StatusOK, gin.H{"data": address})
}

// GetDefaultAddress 处理获取默认地址的请求
func GetDefaultAddress(c *gin.Context) {
	// 1. 创建服务实例
	addressService := services.NewAddressService()

	uid := helper.UID(c)
	if uid == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}
	weid := helper.Weid(c)

	// 3. 调用服务层方法获取默认地址
	address, err := addressService.GetDefaultAddress(weid, uid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if address == nil {
		c.JSON(http.StatusOK, gin.H{"data": nil, "message": "暂无地址信息"})
		return
	}

	// 4. 返回JSON响应
	c.JSON(http.StatusOK, gin.H{"data": address})
}

// GetAddressList 处理获取地址列表的请求
func GetAddressList(c *gin.Context) {
	// 1. 创建服务实例
	addressService := services.NewAddressService()

	uid := helper.UID(c)
	if uid == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}
	weid := helper.Weid(c)
	log.Println("weid:", weid)
	log.Println("uid:", uid)

	// 3. 调用服务层方法获取地址列表
	addresses, err := addressService.GetAddressList(weid, uid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 4. 返回JSON响应
	c.JSON(http.StatusOK, gin.H{"data": addresses})
}
