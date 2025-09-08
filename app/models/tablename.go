package models

import (
	"myapi/app/config"
)

// 全局变量存储完整表名
var (
	// 商品相关表
	goodsTableName                string
	goodsBuynowinfoTableName      string
	goodsSkuTableName             string
	goodsSkuValueTableName        string
	goodsImageTableName           string
	goodsDescriptionTableName     string
	goodsDiscountTableName        string
	categoryTableName             string
	miaoshaGoodsTableName         string
	tuanGoodsTableName            string
	MiaoshaGoodsSkuValueTableName string
	TuanGoodsSkuValueTableName    string

	// 团购相关表
	tuanFoundTableName  string
	tuanFollowTableName string

	// 其他表
	adTableName                     string
	addressTableName                string
	memberTableName                 string
	operatingcityTableName          string
	orderTableName                  string
	userTableName                   string
	usersRelationTableName          string
	uuidRelationTableName           string
	operatingcityIncomelogTableName string
	tuanGoodsSkuValueTableName      string
)

// init函数在包初始化时执行，配置表前缀
func init() {
	// 获取表前缀
	tablePrefix := config.GetString("dbPrefix", "")

	// 商品相关表
	goodsTableName = tablePrefix + "goods"
	goodsBuynowinfoTableName = tablePrefix + "goods_buynowinfo"
	goodsSkuTableName = tablePrefix + "goods_sku"
	goodsSkuValueTableName = tablePrefix + "goods_sku_value"
	goodsImageTableName = tablePrefix + "goods_image"
	goodsDescriptionTableName = tablePrefix + "goods_description"
	goodsDiscountTableName = tablePrefix + "goods_discount"
	categoryTableName = tablePrefix + "category"
	miaoshaGoodsTableName = tablePrefix + "miaosha_goods"
	tuanGoodsTableName = tablePrefix + "tuan_goods"
	MiaoshaGoodsSkuValueTableName = tablePrefix + "miaosha_goods_sku_value"

	// 团购相关表
	tuanFoundTableName = tablePrefix + "tuan_found"
	tuanFollowTableName = tablePrefix + "tuan_follow"
	tuanGoodsSkuValueTableName = tablePrefix + "tuan_goods_sku_value"

	// 其他表
	adTableName = tablePrefix + "ad"
	addressTableName = tablePrefix + "address"
	memberTableName = tablePrefix + "member"
	operatingcityTableName = tablePrefix + "operatingcity"
	orderTableName = tablePrefix + "order"
	userTableName = tablePrefix + "user"
	usersRelationTableName = tablePrefix + "users_relation"
	uuidRelationTableName = tablePrefix + "uuid_relation"
	operatingcityIncomelogTableName = tablePrefix + "operatingcity_incomelog"
}
