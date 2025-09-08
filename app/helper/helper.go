package helper

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func UID(c *gin.Context) int {
	uid, exists := c.Get("userID")
	if !exists {
		return 0
	}
	if uidInt, ok := uid.(int); ok {
		return uidInt
	}
	return 0
}

// Only 从请求中获取指定的字段，类似于PHP中的only方法
func Only(c *gin.Context, fields string) map[string]interface{} {
	result := make(map[string]interface{})
	if fields == "" {
		return result
	}

	fieldList := strings.Split(fields, ",")
	for _, field := range fieldList {
		field = strings.TrimSpace(field)
		if field != "" {
			// 获取POST请求中的字段值
			value := c.PostForm(field)
			if value == "" {
				// 如果POST中没有，尝试从query中获取
				value = c.Query(field)
			}
			if value != "" {
				result[field] = value
			}
		}
	}

	return result
}

// Weid 获取微信公众号ID，适配Gin框架的实现
func Weid(c *gin.Context) int {
	// 尝试从不同的地方获取weid
	weid, exists := c.Get("weid")
	log.Println("helper weid:", weid)
	if !exists {
		return 0
	}
	// 尝试直接转换为int
	if weidInt, ok := weid.(int); ok {
		return weidInt
	}
	// 尝试从float64类型转换
	if weidFloat, ok := weid.(float64); ok {
		return int(weidFloat)
	}
	log.Printf("helper weid%v 类型转换失败:类型%T", weid, weid)
	return 0
}

// Weid 获取微信公众号ID，适配Gin框架的实现
func Sid(c *gin.Context) int {
	// 尝试从不同的地方获取weid
	sid, exists := c.Get("sid")
	log.Println("helper weid:", sid)
	if !exists {
		return 0
	}
	// 尝试直接转换为int
	if sidInt, ok := sid.(int); ok {
		return sidInt
	}
	// 尝试从float64类型转换
	if sidFloat, ok := sid.(float64); ok {
		return int(sidFloat)
	}
	log.Printf("helper sid%v 类型转换失败:类型%T", sid, sid)
	return 0
}

// Ocid 获取ocid参数值
func Ocid(c *gin.Context) int {
	ocid := 0
	if ocidStr := c.Query("ocid"); ocidStr != "" {
		if val, err := strconv.Atoi(ocidStr); err == nil {
			ocid = val
		}
	}
	return ocid
}

// Tzid 获取tzid参数值
func Tzid(c *gin.Context) int {
	tzid := 0
	// 尝试从query参数获取tzid
	if tzidStr := c.Query("tzid"); tzidStr != "" {
		if val, err := strconv.Atoi(tzidStr); err == nil {
			tzid = val
		}
	}
	// 尝试从上下文获取tz_id
	if tzid == 0 {
		if tzID, exists := c.Get("tz_id"); exists {
			if val, ok := tzID.(int); ok {
				tzid = val
			} else if valStr, ok := tzID.(string); ok {
				if val, err := strconv.Atoi(valStr); err == nil {
					tzid = val
				}
			}
		}
	}
	return tzid
}

// GetClient 获取客户端类型
func GetClient(c *gin.Context) string {
	ptype := c.Query("from")
	if ptype == "" {
		ptype = c.DefaultQuery("param.from", "")
	}
	if ptype == "" {
		ptype = "wxapp"
	}
	return ptype
}

// GetWeid 获取当前站点ID（不需要Gin上下文的版本）
// 这是一个公共方法，适配PHP中的weid()函数
func GetWeid() int {
	// TODO: 实现从配置或其他来源获取weid的逻辑
	// 当前暂时返回默认值1
	return 1
}

// AreaConversion 地区名称转换
func AreaConversion(data string) string {
	switch data {
	case "北京市":
		data = "北京"
	case "上海市":
		data = "上海"
	case "天津市":
		data = "天津"
	case "重庆市":
		data = "重庆"
	case "广西壮族自治区":
		data = "广西"
	}
	return data
}

// EncryptTel 加密手机号，保留前三位和后四位
func EncryptTel(tel string) string {
	if len(tel) <= 7 {
		return tel
	}
	// 保留前三位和后四位，中间用****替换
	return tel[:3] + "****" + tel[len(tel)-4:]
}

// GetDomainName 获取当前域名
func GetDomainName(c *gin.Context) string {
	// 尝试从X-Forwarded-Host获取
	host := c.GetHeader("X-Forwarded-Host")
	if host != "" {
		return host
	}

	// 尝试从X-Forwarded-Server获取
	host = c.GetHeader("X-Forwarded-Server")
	if host != "" {
		return host
	}

	// 获取请求的Host
	host = c.Request.Host
	if host != "" {
		return host
	}

	// 最后尝试从ServerName获取
	return c.Request.URL.Host
}

// GetServerIP 获取服务器IP地址
func GetServerIP(c *gin.Context) string {
	domain := GetDomainName(c)
	if domain != "" {
		// 在实际生产环境中，可以使用net包解析域名获取IP
		// 这里为简化实现，直接返回域名
		return domain
	}
	return "127.0.0.1"
}

// IsHTTPS 判断当前请求是否为HTTPS
func IsHTTPS(c *gin.Context) bool {
	// 检查HTTPS请求头
	if c.GetHeader("HTTPS") == "on" || strings.ToLower(c.GetHeader("HTTPS")) != "off" {
		return true
	}

	// 检查X-Forwarded-Proto头
	if c.GetHeader("X-Forwarded-Proto") == "https" {
		return true
	}

	// 检查HTTP_FRONT_END_HTTPS头
	if c.GetHeader("HTTP_FRONT_END_HTTPS") == "on" || strings.ToLower(c.GetHeader("HTTP_FRONT_END_HTTPS")) != "off" {
		return true
	}

	// 检查请求URL是否为HTTPS
	if c.Request.URL.Scheme == "https" {
		return true
	}

	// 检查请求端口是否为443
	if c.Request.URL.Port() == "443" {
		return true
	}

	return false
}

// GetHost 获取带协议的主机名
func GetHost(c *gin.Context) string {
	protocol := "http"
	if IsHTTPS(c) {
		protocol = "https"
	}
	return protocol + "://" + GetDomainName(c)
}

// GetHTTPSHost 获取HTTPS协议的主机名
func GetHTTPSHost(c *gin.Context) string {
	return "https://" + GetDomainName(c)
}

// GetURL 获取当前完整URL
func GetURL(c *gin.Context) string {
	return GetHost(c) + c.Request.RequestURI
}

// GetRealIP 获取客户端真实IP地址
func GetRealIP(c *gin.Context) string {
	// 检查X-Forwarded-For头
	forwarded := c.GetHeader("X-Forwarded-For")
	if forwarded != "" {
		// 多个IP地址时，第一个通常是真实IP
		ips := strings.Split(forwarded, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}
	// 如果没有X-Forwarded-For头，使用RemoteAddr
	ip := getClientIP(c)
	return ip
}

// 获取客户端IP（优先IPv4）
func getClientIP(c *gin.Context) string {
	ip := c.ClientIP()
	// 检查是否是IPv6地址（包含冒号）且不是::1以外的地址
	if strings.Contains(ip, ":") && ip != "::1" {
		// 尝试从X-Forwarded-For头中获取IPv4地址
		xff := c.Request.Header.Get("X-Forwarded-For")
		if xff != "" {
			ips := strings.Split(xff, ",")
			for _, v := range ips {
				v = strings.TrimSpace(v)
				if !strings.Contains(v, ":") {
					return v
				}
			}
		}
	}
	// 如果是::1（本地回环地址），返回127.0.0.1
	if ip == "::1" {
		return "127.0.0.1"
	}
	return ip
}

// VersionIncreasing 版本号递增
func VersionIncreasing(version string) string {
	if version == "" {
		return ""
	}

	// 分割版本号
	parts := strings.Split(version, ".")
	if len(parts) == 0 {
		return version
	}

	// 将最后一部分转换为数字并递增
	lastPart, err := strconv.Atoi(parts[len(parts)-1])
	if err != nil {
		return version
	}
	parts[len(parts)-1] = strconv.Itoa(lastPart + 1)

	// 重新组合版本号
	return strings.Join(parts, ".")
}

// SetURL 格式化URL，确保带有协议头
func SetURL(path string) string {
	if path == "" {
		return ""
	}

	// 检查是否已经包含协议头
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		return path
	}

	// 添加默认的http协议头
	return "http://" + path
}

// SMSReplace 移除短信内容中的括号
func SMSReplace(str string) string {
	str = strings.ReplaceAll(str, "【", "")
	str = strings.ReplaceAll(str, "】", "")
	str = strings.ReplaceAll(str, "[", "")
	str = strings.ReplaceAll(str, "]", "")
	return str
}

// WithdrawStatus 获取提现状态的中文描述
func WithdrawStatus(status int) string {
	switch status {
	case 0:
		return "未处理"
	case 1:
		return "已结算"
	case 2:
		return "驳回"
	default:
		return "未知状态"
	}
}

// Sex 根据数字获取性别中文描述
func Sex(sex int) string {
	switch sex {
	case 0:
		return "保密"
	case 1:
		return "男"
	case 2:
		return "女"
	default:
		return "未知"
	}
}

// SexArray 获取性别数组
func SexArray() []map[string]interface{} {
	return []map[string]interface{}{
		{"val": 0, "key": "保密"},
		{"val": 1, "key": "男"},
		{"val": 2, "key": "女"},
	}
}

// ArrayEmpty 检查数组是否全为空
func ArrayEmpty(val map[string]interface{}) bool {
	for _, v := range val {
		if v != nil && v != "" {
			return false
		}
	}
	return true
}

// XmStrToTime 将字符串时间转换为时间戳
func XmStrToTime(thistime string) int64 {
	t, err := time.Parse("2006-01-02 15:04:05", thistime)
	if err != nil {
		// 尝试其他常见格式
		t, err = time.Parse("2006-01-02", thistime)
		if err != nil {
			return 0
		}
	}

	timestamp := t.Unix()
	if timestamp < 0 {
		timestamp = 0
	}
	return timestamp
}

// GetWeekRecentlyDay 获取最近的指定星期几的日期
func GetWeekRecentlyDay(week int) string {
	// PHP中周日是0，而Go中周日是0，所以这里不需要调整
	now := time.Now()

	for i := 0; i < 7; i++ {
		thistime := now.AddDate(0, 0, i)
		if int(thistime.Weekday()) == week {
			return thistime.Format("2006-01-02")
		}
	}

	return ""
}

// GetDayRecentlyDay 获取最近的指定日期的日期
func GetDayRecentlyDay(day int) string {
	now := time.Now()

	for i := 0; i < 31; i++ {
		thistime := now.AddDate(0, 0, i)
		if thistime.Day() == day {
			return thistime.Format("2006-01-02")
		}
	}

	return ""
}

// TimeFormat 格式化时间为指定格式
func TimeFormat(timestamp interface{}, format string) string {
	return _time(timestamp, format)
}

// TimeMDHi 格式化时间为月-日 时:分格式
func TimeMDHi(timestamp interface{}) string {
	return _time(timestamp, "01-02 15:04")
}

// TimeYMD 格式化时间为年-月-日格式
func TimeYMD(timestamp interface{}) string {
	return _time(timestamp, "2006-01-02")
}

// TimeMD 格式化时间为月-日格式
func TimeMD(timestamp interface{}) string {
	return _time(timestamp, "01-02")
}

// TimeY 格式化时间为年份格式
func TimeY(timestamp interface{}) string {
	return _time(timestamp, "2006")
}

// _time 内部时间格式化函数
func _time(timestamp interface{}, format string) string {
	if timestamp == nil || timestamp == "" {
		return ""
	}

	var t time.Time
	var err error

	// 根据不同类型处理时间戳
	switch v := timestamp.(type) {
	case int64:
		t = time.Unix(v, 0)
	case int:
		t = time.Unix(int64(v), 0)
	case string:
		// 尝试解析字符串为时间
		t, err = time.Parse("2006-01-02 15:04:05", v)
		if err != nil {
			t, err = time.Parse("2006-01-02", v)
			if err != nil {
				// 如果无法解析，尝试转换为时间戳
				if ts, err := strconv.ParseInt(v, 10, 64); err == nil {
					t = time.Unix(ts, 0)
				} else {
					return v // 无法解析，返回原始值
				}
			}
		}
	default:
		return ""
	}

	return t.Format(format)
}

// SetIntToABC 将数字转换为字母
func SetIntToABC(num int) string {
	if num < 0 || num > 25 {
		return "" // 只支持0-25的数字
	}

	// A的ASCII码是65，Z是90
	return string(rune(65 + num))
}

// ExtractNotifyParams 解析通知参数
func ExtractNotifyParams(params string) map[string]string {
	result := make(map[string]string)
	paramsArr := strings.Split(params, "|")

	for _, vo := range paramsArr {
		tmp := strings.Split(vo, "_")
		if len(tmp) == 2 {
			result[tmp[0]] = tmp[1]
		}
	}

	return result
}

// PercentToNum 百分比转换为数字
func PercentToNum(n float64) float64 {
	return n / 100
}

// NumToPercent 数字转换为百分比字符串
func NumToPercent(n float64) string {
	// 保留1位小数的百分比
	return fmt.Sprintf("%.1f%%", n*100)
}

// ToPercent 计算百分比
func ToPercent(m, n float64) string {
	if n <= 0 {
		return NumToPercent(0)
	}

	return NumToPercent(m / n)
}

// GenerateSelectTree 无限极分类转为带有children的树形select结构
type SelectNode struct {
	Key      string       `json:"key"`
	Val      interface{}  `json:"val"`
	Children []SelectNode `json:"children"`
}

func GenerateSelectTree(data []map[string]interface{}, pid interface{}) []SelectNode {
	tree := []SelectNode{}

	for _, v := range data {
		// 检查v中是否存在pid字段并且等于传入的pid
		if v["pid"] == pid {
			// 创建一个新节点
			node := SelectNode{
				Key: v["key"].(string),
				Val: v["val"],
			}
			// 递归获取子节点
			node.Children = GenerateSelectTree(data, v["val"])
			tree = append(tree, node)
		}
	}

	return tree
}

// RemoveEmoji 删除字符串中的Emoji表情
func RemoveEmoji(str string) string {
	if str == "" {
		return ""
	}

	result := ""
	for _, r := range str {
		// 过滤掉4字节以上的字符（通常是Emoji表情）
		if len(string(r)) < 4 {
			result += string(r)
		}
	}

	return result
}

// GenerateListTree 无限极分类转为带有children的树形list表格结构
type TreeNode struct {
	Children []TreeNode `json:"children"`
	// 其他字段根据实际需求动态添加
}

func GenerateListTree(data []map[string]interface{}, pid interface{}, config []string) []map[string]interface{} {
	tree := []map[string]interface{}{}

	for _, v := range data {
		// 检查v中是否存在config[1]字段并且等于传入的pid
		if config[1] != "" && v[config[1]] == pid {
			// 创建一个新节点并复制原始数据
			newNode := make(map[string]interface{})
			for k, val := range v {
				newNode[k] = val
			}
			// 递归获取子节点
			if config[0] != "" {
				newNode["children"] = GenerateListTree(data, v[config[0]], config)
			}
			tree = append(tree, newNode)
		}
	}

	return tree
}

// DoOrderSn 生成订单流水号
func DoOrderSn(orderType string) string {
	now := time.Now()
	microsecond := now.UnixNano() / 1000000 % 1000 // 获取微秒部分的后三位
	randNum := rand.Intn(100)                      // 生成0-99的随机数

	return fmt.Sprintf("%s%s%03d%02d", now.Format("20060102150405"), orderType, microsecond, randNum)
}

// DeterminePaymentSource 判断支付来源
func DeterminePaymentSource(authCode string) string {
	// 微信支付以10-15开头
	if matched, _ := regexp.MatchString("^(10|11|12|13|14|15)", authCode); matched {
		return "wx_pay"
	}
	// 支付宝以25-30开头
	if matched, _ := regexp.MatchString("^(25|26|27|28|29|30)", authCode); matched {
		return "ali_pay"
	}
	return "Unknown"
}

// GetCollectType 获取收款类型
func GetCollectType(id string) []map[string]interface{} {
	collectTypes := []map[string]interface{}{
		{"val": "bank", "key": "银行卡"},
		{"val": "wechat", "key": "微信支付"},
		{"val": "alipay", "key": "支付宝"},
	}

	if id == "" {
		return collectTypes
	}

	// 如果指定了id，返回对应的名称
	for _, item := range collectTypes {
		if item["val"] == id {
			return []map[string]interface{}{item}
		}
	}

	return []map[string]interface{}{}
}

// GetCommissionType 获取分佣类型
func GetCommissionType(key string) []map[string]interface{} {
	commissionTypes := []map[string]interface{}{
		{"roletype": "agent", "title": "分销达人"},
		{"roletype": "province", "title": "省代理"},
		{"roletype": "city", "title": "市代理"},
		{"roletype": "district", "title": "区县代理"},
		{"roletype": "tuanzhang", "title": "社区代理"},
		{"roletype": "store", "title": "商家"},
	}

	if key == "" {
		return commissionTypes
	}

	// 如果指定了key，返回对应的类型
	for _, item := range commissionTypes {
		if item["roletype"] == key {
			return []map[string]interface{}{item}
		}
	}

	return []map[string]interface{}{}
}

// GetOtype 获取订单类型
func GetOtype(id int) []map[string]interface{} {
	orderTypes := []map[string]interface{}{
		{"val": 0, "key": "普通订单"},
		{"val": 1, "key": "需求订单"},
		{"val": 2, "key": "跑腿订单"},
	}

	if id < 0 {
		return orderTypes
	}

	// 如果指定了id，返回对应的名称
	for _, item := range orderTypes {
		if item["val"] == id {
			return []map[string]interface{}{item}
		}
	}

	return []map[string]interface{}{}
}

// GetCouponPtype 获取优惠券获取类型
func GetCouponPtype(id int) []map[string]interface{} {
	couponTypes := []map[string]interface{}{
		{"val": 1, "key": "领取"},
		{"val": 2, "key": "下单赠送"},
		{"val": 3, "key": "新人券"},
		{"val": 4, "key": "发放"},
	}

	if id == 0 {
		return couponTypes
	}

	// 如果指定了id，返回对应的名称
	for _, item := range couponTypes {
		if item["val"] == id {
			return []map[string]interface{}{item}
		}
	}

	return []map[string]interface{}{}
}

// GetCouponType 获取优惠券类型
func GetCouponType(id int) []map[string]interface{} {
	couponTypes := []map[string]interface{}{
		{"val": 10, "key": "代金券"},
		{"val": 20, "key": "折扣券"},
	}

	if id == 0 {
		return couponTypes
	}

	// 如果指定了id，返回对应的名称
	for _, item := range couponTypes {
		if item["val"] == id {
			return []map[string]interface{}{item}
		}
	}

	return []map[string]interface{}{}
}

// GetExpireType 获取有效期类型
func GetExpireType(id int) []map[string]interface{} {
	expireTypes := []map[string]interface{}{
		{"val": 10, "key": "即时生效"},
		{"val": 20, "key": "固定时间"},
	}

	if id == 0 {
		return expireTypes
	}

	// 如果指定了id，返回对应的名称
	for _, item := range expireTypes {
		if item["val"] == id {
			return []map[string]interface{}{item}
		}
	}

	return []map[string]interface{}{}
}

// GetColor 获取颜色类型
func GetColor(id string) []map[string]interface{} {
	colors := []map[string]interface{}{
		{"val": "blue", "key": "蓝色"},
		{"val": "red", "key": "红色"},
		{"val": "violet", "key": "紫色"},
		{"val": "yellow", "key": "黄色"},
	}

	if id == "" {
		return colors
	}

	// 如果指定了id，返回对应的名称
	for _, item := range colors {
		if item["val"] == id {
			return []map[string]interface{}{item}
		}
	}

	return []map[string]interface{}{}
}

// GetPrizeRtype 获取奖品类型
func GetPrizeRtype(id int) []map[string]interface{} {
	prizes := []map[string]interface{}{
		{"val": 1, "key": "谢谢参与"},
		{"val": 2, "key": "余额红包"},
		{"val": 3, "key": "优惠券"},
		{"val": 4, "key": "积分"},
	}

	if id == 0 {
		return prizes
	}

	// 如果指定了id，返回对应的名称
	for _, item := range prizes {
		if item["val"] == id {
			return []map[string]interface{}{item}
		}
	}

	return []map[string]interface{}{}
}

// GetOrderTplOption 获取订单模板选项
func GetOrderTplOption(id int) []map[string]interface{} {
	options := []map[string]interface{}{
		{"val": "shipping_name", "key": "下单用户"},
		{"val": "order_num_alias", "key": "订单号"},
		{"val": "total", "key": "订单金额"},
		{"val": "pay_subject", "key": "商品信息"},
		{"val": "shipping_tel", "key": "联系电话"},
		{"val": "pay_time", "key": "购买时间"},
	}

	if id == 0 {
		return options
	}

	// 如果指定了id，返回对应的名称
	for i, item := range options {
		if i+1 == id {
			return []map[string]interface{}{item}
		}
	}

	return []map[string]interface{}{}
}

// GetNotWinningPtype 获取未中奖奖品类型
func GetNotWinningPtype(id int) []map[string]interface{} {
	prizes := []map[string]interface{}{
		{"val": 1, "key": "无"},
		{"val": 2, "key": "余额红包"},
		{"val": 3, "key": "优惠券"},
		{"val": 4, "key": "积分"},
	}

	if id == 0 {
		return prizes
	}

	// 如果指定了id，返回对应的名称
	for _, item := range prizes {
		if item["val"] == id {
			return []map[string]interface{}{item}
		}
	}

	return []map[string]interface{}{}
}

// GetObtainPtype 获取获取方式类型
func GetObtainPtype(id int) []map[string]interface{} {
	types := []map[string]interface{}{
		{"val": 1, "key": "商品下单"},
		{"val": 2, "key": "拼团"},
		{"val": 3, "key": "拼团没中"},
		{"val": 4, "key": "上传商品"},
	}

	if id == 0 {
		return types
	}

	// 如果指定了id，返回对应的名称
	for _, item := range types {
		if item["val"] == id {
			return []map[string]interface{}{item}
		}
	}

	return []map[string]interface{}{}
}

// YesNo 将布尔值转换为"是"或"否"
func YesNo(b bool) string {
	if b {
		return "是"
	}
	return "否"
}

// PaymentCode 将支付代码转换为支付方式名称
func PaymentCode(code string) string {
	paymentMap := map[string]string{
		"balance_pay":       "余额支付",
		"wx_pay":            "微信支付",
		"points_pay":        "积分兑换",
		"goodsgiftcard_pay": "购物卡抵扣",
		"delivery_pay":      "货到付款",
		"alipay":            "支付宝",
		"offline_pay":       "线下交易",
		"":                  "未支付",
	}

	if name, exists := paymentMap[code]; exists {
		return name
	}
	return code
}

// TuanFoundStatus 拼团状态转换
func TuanFoundStatus(status int) string {
	switch status {
	case 0:
		return "待成团"
	case 1:
		return "拼团成功"
	case 2:
		return "拼团失败"
	default:
		return "未知状态"
	}
}

// RefundType 退款类型转换
func RefundType(id int) string {
	switch id {
	case 1:
		return "未发货退款"
	case 2:
		return "退货退款"
	case 3:
		return "换货"
	default:
		return "未知类型"
	}
}

// RefundTypeYuyue 预约退款类型转换
func RefundTypeYuyue(id int) string {
	switch id {
	case 1:
		return "未服务退款"
	case 2:
		return "不满意退款"
	case 3:
		return "返工"
	default:
		return "未知类型"
	}
}

// RefundStatus 退款状态转换
func RefundStatus(id int) string {
	switch id {
	case 0:
		return "待处理"
	case 1:
		return "已退款"
	case 2:
		return "已同意退换货"
	case 3:
		return "已拒绝"
	default:
		return "未知状态"
	}
}

// RefundStatusYuyue 预约退款状态转换
func RefundStatusYuyue(id int) string {
	switch id {
	case 0:
		return "待处理"
	case 1:
		return "已退款"
	case 2:
		return "已同意售后"
	case 3:
		return "已拒绝"
	default:
		return "未知状态"
	}
}

// ShareLevel 分享层级转换
func ShareLevel(level int) string {
	levelMap := map[int]string{
		1: "一层佣金",
		2: "二层佣金",
		3: "三层佣金",
	}

	if name, exists := levelMap[level]; exists {
		return name
	}
	return "未知层级"
}

// TimingUnitName 时间单位名称转换
func TimingUnitName(key string) string {
	unitMap := map[string]string{
		"day":   "天",
		"week":  "周",
		"month": "月",
		"year":  "年",
	}

	if name, exists := unitMap[key]; exists {
		return name
	}
	return key
}

// IsHTTP 检查是否为HTTP路径
func IsHTTP(path string) bool {
	if path == "" {
		return false
	}

	parts := strings.Split(path, "/")
	if len(parts) > 0 {
		prefix := parts[0]
		return prefix == "https:" || prefix == "http:"
	}
	return false
}

// GetCollectTypeName 获取收款类型名称
func GetCollectTypeName(vals string) string {
	if vals == "" {
		return ""
	}

	collectTypes := GetCollectType("")
	var result string

	for _, val := range strings.Split(vals, ",") {
		for _, item := range collectTypes {
			if item["val"] == val {
				if result == "" {
					result = item["key"].(string)
				} else {
					result += "," + item["key"].(string)
				}
			}
		}
	}

	return result
}

// GetGoodsDeliveryMode 获取商品配送方式
func GetGoodsDeliveryMode(id int) []map[string]interface{} {
	deliveryModes := []map[string]interface{}{
		{"val": 1, "key": "同城配送"},
		{"val": 2, "key": "到店自提"},
		{"val": 3, "key": "快递"},
		{"val": 5, "key": "社区点自提"},
	}

	if id == 0 {
		return deliveryModes
	}

	// 如果指定了id，返回对应的名称
	for _, item := range deliveryModes {
		if item["val"] == id {
			return []map[string]interface{}{item}
		}
	}

	return []map[string]interface{}{}
}

// GetGoodsDeliveryModeName 获取商品配送方式名称
func GetGoodsDeliveryModeName(ids string) string {
	if ids == "" {
		return ""
	}

	deliveryModes := GetGoodsDeliveryMode(0)
	var result string

	for _, id := range strings.Split(ids, ",") {
		idInt, _ := strconv.Atoi(id)
		for _, item := range deliveryModes {
			if item["val"] == idInt {
				if result == "" {
					result = item["key"].(string)
				} else {
					result += "," + item["key"].(string)
				}
			}
		}
	}

	return result
}

// GetCourseDeliveryMode 获取课程配送方式
func GetCourseDeliveryMode(id int) []map[string]interface{} {
	deliveryModes := []map[string]interface{}{
		{"val": 1, "key": "线上课"},
		{"val": 2, "key": "线下课"},
	}

	if id == 0 {
		return deliveryModes
	}

	// 如果指定了id，返回对应的名称
	for _, item := range deliveryModes {
		if item["val"] == id {
			return []map[string]interface{}{item}
		}
	}

	return []map[string]interface{}{}
}

// GetSourceType 获取资源类型
func GetSourceType(id int) []map[string]interface{} {
	sourceTypes := []map[string]interface{}{
		{"val": 1, "key": "图文"},
		{"val": 2, "key": "音频"},
		{"val": 3, "key": "视频"},
	}

	if id == 0 {
		return sourceTypes
	}

	// 如果指定了id，返回对应的名称
	for _, item := range sourceTypes {
		if item["val"] == id {
			return []map[string]interface{}{item}
		}
	}

	return []map[string]interface{}{}
}

// GetGoodsDeliveryModeArray 获取商品配送方式数组
func GetGoodsDeliveryModeArray(ids string) []map[string]interface{} {
	if ids == "" {
		return []map[string]interface{}{}
	}

	//modes := GetGoodsDeliveryMode()
	idList := strings.Split(ids, ",")
	result := make([]map[string]interface{}, 0, len(idList))

	// for _, id := range idList {
	// 	idInt, _ := strconv.Atoi(id)
	// 	if mode, exists := modes[idInt]; exists {
	// 		result = append(result, mode)
	// 	}
	// }

	return result
}

// GetServiceDeliveryMode 获取服务配送方式
func GetServiceDeliveryMode(id int) interface{} {
	deliveryModes := []map[string]interface{}{
		{"val": 1, "key": "上门服务"},
		{"val": 2, "key": "到店服务"},
		{"val": 4, "key": "在线服务"},
	}

	if id <= 0 {
		return deliveryModes
	}

	for _, mode := range deliveryModes {
		if mode["val"] == id {
			return mode["key"]
		}
	}

	return nil
}

// GetServiceDeliveryModeName 获取服务配送方式名称
func GetServiceDeliveryModeName(ids string) string {
	if ids == "" {
		return ""
	}

	idList := strings.Split(ids, ",")
	var result string

	for _, id := range idList {
		idInt, _ := strconv.Atoi(id)
		modeName := GetServiceDeliveryMode(idInt)
		if modeName != nil {
			if result == "" {
				result = modeName.(string)
			} else {
				result = result + "," + modeName.(string)
			}
		}
	}

	return result
}

// GetServiceDeliveryModeArray 获取服务配送方式数组
func GetServiceDeliveryModeArray(ids string) []map[string]interface{} {
	if ids == "" {
		return []map[string]interface{}{}
	}

	deliveryModes := []map[string]interface{}{
		{"val": 1, "key": "上门服务"},
		{"val": 2, "key": "到店服务"},
		{"val": 4, "key": "在线服务"},
	}

	idList := strings.Split(ids, ",")
	result := make([]map[string]interface{}, 0, len(idList))

	for _, id := range idList {
		idInt, _ := strconv.Atoi(id)
		for _, mode := range deliveryModes {
			if mode["val"] == idInt {
				result = append(result, mode)
				break
			}
		}
	}

	return result
}

// GetPtype 获取类型（商品/服务）
func GetPtype(id int) interface{} {
	types := []map[string]interface{}{
		{"val": 1, "key": "商品"},
		{"val": 2, "key": "服务"},
	}

	if id <= 0 {
		return types
	}

	for _, t := range types {
		if t["val"] == id {
			return t["key"]
		}
	}

	return nil
}

// Status 状态显示（HTML图标）
func Status(str bool) string {
	if str {
		return "<i class='fa fa-check-square'></i>"
	} else {
		return "<i class='fa fa-ban'></i>"
	}
}

// BuildOrderNo 生成唯一订单号
func BuildOrderNo(prefix string) string {
	now := time.Now()
	rand.Seed(time.Now().UnixNano())
	randNum := rand.Intn(100000000) // 8位随机数
	return fmt.Sprintf("%s%s%08d", prefix, now.Format("20060102"), randNum)
}

// APIException 抛出API异常
func APIException(statusKey string, customMsg string, data interface{}) error {
	// 在实际应用中，这里应该返回一个自定义的错误类型
	// 这里简化实现，返回一个普通错误
	return fmt.Errorf("API Error [%s]: %s", statusKey, customMsg)
}

// APISuccess 返回成功响应
func APISuccess(data interface{}, msg string) map[string]interface{} {
	if msg == "" {
		msg = "操作成功"
	}
	return map[string]interface{}{
		"status": 200,
		"msg":    msg,
		"data":   data,
	}
}

// APIError 返回错误响应
func APIError(statusKey string, customMsg string, data interface{}) map[string]interface{} {
	// 在实际应用中，应该从配置中获取状态码配置
	// 这里简化实现，使用默认值
	statusConfig := map[string]interface{}{
		"code": 500,
		"msg":  "未知错误",
	}

	// 实际项目中应该从配置文件中获取状态码
	// 这里只是示例

	result := map[string]interface{}{
		"status": statusConfig["code"],
		"msg":    statusConfig["msg"],
	}

	if customMsg != "" {
		result["msg"] = customMsg
	}

	if data != nil {
		result["data"] = data
	}

	return result
}

// StrongHTTP 将HTTP转换为HTTPS
func StrongHTTP(path string) string {
	return strings.Replace(path, "http://", "https://", 1)
}

// SetPicsView 设置图片视图
func SetPicsView(pics string) []map[string]interface{} {
	if pics == "" {
		return []map[string]interface{}{}
	}

	picList := strings.Split(pics, ",")
	result := make([]map[string]interface{}, 0, len(picList))

	for _, pic := range picList {
		if pic != "" {
			result = append(result, map[string]interface{}{
				"url": ToImg(pic),
			})
		}
	}

	return result
}

// 先不加逻辑
func ToImg(url string) string {
	if url == "" {
		return ""
	}
	return url
}

// ScriptPath 获取脚本路径
func ScriptPath(c *gin.Context) string {
	if c == nil {
		return ""
	}

	scriptName := c.Request.URL.Path
	tmpPath := strings.Split(scriptName, "/")
	var pathOne string
	var ret string

	if len(tmpPath) > 1 {
		pathOne = tmpPath[1]
		if len(tmpPath) > 2 {
			ret = "/" + tmpPath[1] + "/" + tmpPath[2]
		}
	}

	if pathOne == "addons" {
		return ret
	}
	return ""
}

// RemoveXSS 过滤XSS
func RemoveXSS(str string) string {
	// 移除控制字符
	str = regexp.MustCompile("[\\x00-\\x08\\x0B\\x0C\\x0E-\\x1F\\x7F]+").ReplaceAllString(str, "")

	// 定义需要过滤的标签和属性
	params1 := []string{"javascript", "vbscript", "expression", "applet", "meta", "xml", "blink", "script", "embed", "object", "iframe", "frame", "frameset", "ilayer", "layer", "bgsound"}
	params2 := []string{"onabort", "onactivate", "onafterprint", "onafterupdate", "onbeforeactivate", "onbeforecopy", "onbeforecut", "onbeforedeactivate", "onbeforeeditfocus", "onbeforepaste", "onbeforeprint", "onbeforeunload", "onbeforeupdate", "onblur", "onbounce", "oncellchange", "onchange", "onclick", "oncontextmenu", "oncontrolselect", "oncopy", "oncut", "ondataavailable", "ondatasetchanged", "ondatasetcomplete", "ondblclick", "ondeactivate", "ondrag", "ondragend", "ondragenter", "ondragleave", "ondragover", "ondragstart", "ondrop", "onerror", "onerrorupdate", "onfilterchange", "onfinish", "onfocus", "onfocusin", "onfocusout", "onhelp", "onkeydown", "onkeypress", "onkeyup", "onlayoutcomplete", "onload", "onlosecapture", "onmousedown", "onmouseenter", "onmouseleave", "onmousemove", "onmouseout", "onmouseover", "onmouseup", "onmousewheel", "onmove", "onmoveend", "onmovestart", "onpaste", "onpropertychange", "onreadystatechange", "onreset", "onresize", "onresizeend", "onresizestart", "onrowenter", "onrowexit", "onrowsdelete", "onrowsinserted", "onscroll", "onselect", "onselectionchange", "onselectstart", "onstart", "onstop", "onsubmit", "onunload"}

	params := append(params1, params2...)

	// 构建正则表达式并替换
	for _, param := range params {
		// 构建可以匹配各种编码形式的正则表达式
		pattern := "(" + param[0:1] + ")(?:&#[x|X]0([9][a][b]);?|&#0([9][10][13]);?)?" + param[1:]
		str = regexp.MustCompile("(?i)"+pattern).ReplaceAllString(str, "$1")
	}

	return str
}

// StrExists 检查字符串是否存在
func StrExists(str string, find string) bool {
	return strings.Contains(str, find)
}

// PassHash 对密码进行加密，类似于PHP中的pass_hash函数
func PassHash(passwordInput, salt string) string {
	// 从环境变量中读取AUTH_KEY，如果不存在则使用默认值
	authKey := "" // 默认值，确保向后兼容
	if value, exists := os.LookupEnv("AUTH_KEY"); exists {
		authKey = value
	}

	// 拼接密码、盐值和authkey
	combined := passwordInput + "-" + salt + "-" + authKey

	// 计算SHA1哈希
	sha1Hash := sha1.New()
	sha1Hash.Write([]byte(combined))

	// 将哈希结果转换为十六进制字符串
	return hex.EncodeToString(sha1Hash.Sum(nil))
}
func ToFloat64(value interface{}) (float64, error) {
	switch v := value.(type) {
	case float64:
		return v, nil
	case float32:
		return float64(v), nil
	case int:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case uint:
		return float64(v), nil
	case uint32:
		return float64(v), nil
	case uint64:
		return float64(v), nil
		// 处理字符串类型
	case string:
		return strconv.ParseFloat(v, 64)
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("无法转换类型 %T 到 float64", v)
	}
}

// ToStr 将任意类型转换为string
func ToStr(value interface{}) (string, error) {
	switch v := value.(type) {
	case string:
		return v, nil
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64), nil
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32), nil
	case int:
		return strconv.Itoa(v), nil
	case int32:
		return strconv.Itoa(int(v)), nil
	case int64:
		return strconv.FormatInt(v, 10), nil
	case uint:
		return strconv.FormatUint(uint64(v), 10), nil
	case uint32:
		return strconv.FormatUint(uint64(v), 10), nil
	case uint64:
		return strconv.FormatUint(v, 10), nil
	case bool:
		return strconv.FormatBool(v), nil
	case nil:
		return "", nil
	default:
		return fmt.Sprintf("%v", v), nil
	}
}

// ToInt 将任意类型转换为int
func ToInt(value interface{}) (int, error) {
	switch v := value.(type) {
	case int:
		return v, nil
	case int32:
		return int(v), nil
	case int64:
		return int(v), nil
	case uint:
		return int(v), nil
	case uint32:
		return int(v), nil
	case uint64:
		return int(v), nil
	case float64:
		return int(v), nil
	case float32:
		return int(v), nil
	case string:
		return strconv.Atoi(v)
	case bool:
		if v {
			return 1, nil
		} else {
			return 0, nil
		}
	case nil:
		return 0, nil
	default:
		return 0, fmt.Errorf("无法转换类型 %T 到 int", v)
	}
}

