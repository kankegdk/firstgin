package helper

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// 接收gin上下文和要验证的结构体指针或map
// 返回验证是否成功（true表示验证通过）
func ValidateRequest(c *gin.Context, data interface{}) bool {
	// 1. 首先进行基础的ShouldBind验证（处理binding标签规则）
	if err := c.ShouldBind(data); err != nil {
		msg := ParseValidationError(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": msg})
		return false
	}

	// 2. 使用反射检查结构体中的验证标签，进行额外的自定义验证
	validator := NewValidator()
	dataValue := reflect.ValueOf(data)

	// 确保data是指针类型
	if dataValue.Kind() == reflect.Ptr {
		// 获取指针指向的元素
		dataValue = dataValue.Elem()
	}

	// 处理结构体类型的数据
	if dataValue.Kind() == reflect.Struct {
		// 动态加载验证规则
		if !loadValidationRulesFromStruct(c, validator, dataValue) {
			return false
		}

		// 将结构体转换为map用于验证
		dataMap := structToMap(dataValue)
		valid, errMsg := validator.Check(dataMap, "")
		if !valid {
			c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
			return false
		}
	}

	// 所有验证通过
	return true
}

// loadValidationRulesFromStruct 从结构体标签中加载验证规则
func loadValidationRulesFromStruct(c *gin.Context, validator *Validator, dataValue reflect.Value) bool {
	dataType := dataValue.Type()
	// 遍历结构体的所有字段
	for i := 0; i < dataValue.NumField(); i++ {
		// 获取字段信息
		field := dataType.Field(i)
		fieldName := field.Name
		fieldValue := dataValue.Field(i).Interface()
		jsonTag := field.Tag.Get("json")
		validateTag := field.Tag.Get("customvalidate")

		// 如果没有验证标签，跳过
		if validateTag == "" {
			continue
		}

		// 解析字段名（使用json标签或字段名）
		fieldKey := fieldName
		if jsonTag != "" && jsonTag != "-" {
			fieldKey = strings.Split(jsonTag, ",")[0]
		}

		// 解析验证规则（支持|或,分隔）
		rules := strings.FieldsFunc(validateTag, func(r rune) bool {
			return r == '|' || r == ','
		})
		for _, rule := range rules {
			// 动态添加验证规则
			if !addValidationRule(validator, fieldKey, rule, fieldValue) {
				return false
			}
		}
	}
	return true
}

// addValidationRule 根据规则字符串添加验证规则
func addValidationRule(validator *Validator, field, rule string, value interface{}) bool {
	ruleParts := strings.SplitN(rule, ":", 2)
	ruleName := strings.TrimSpace(ruleParts[0])
	log.Println("自定义校验方法:", ruleName)
	log.Println("设置规则:", ruleParts)
	var ruleParams string
	if len(ruleParts) > 1 {
		ruleParams = strings.TrimSpace(ruleParts[1])
	}

	// 根据规则名称设置验证函数和错误信息
	var callback func(interface{}) bool
	var message string

	// 构建消息
	getRuleMessage := func(baseMsg string) string {
		if ruleParams != "" {
			return fmt.Sprintf("%s %s", field, fmt.Sprintf(baseMsg, ruleParams))
		}
		return fmt.Sprintf("%s %s", field, baseMsg)
	}

	// 根据规则名称选择验证函数
	switch strings.ToLower(ruleName) {
	case "email":
		callback = func(v interface{}) bool {
			if val, ok := v.(string); ok {
				return IsEmail(val)
			}
			return false
		}
		message = getRuleMessage("不是有效的邮箱地址")
	case "mobile", "phone", "telephone":
		callback = func(v interface{}) bool {
			if val, ok := v.(string); ok {
				result := IsMobile(val)
				return result
			}
			return false
		}
		message = getRuleMessage("不是有效的手机号码")
	case "url":
		callback = func(v interface{}) bool {
			if val, ok := v.(string); ok {
				return IsURL(val)
			}
			return false
		}
		message = getRuleMessage("不是有效的URL地址")
	case "number":
		callback = func(v interface{}) bool {
			if val, ok := v.(string); ok {
				return IsNumber(val)
			} else if IsNumeric(v) {
				return true
			}
			return false
		}
		message = getRuleMessage("必须是数字")
	case "alpha":
		callback = func(v interface{}) bool {
			if val, ok := v.(string); ok {
				return IsAlpha(val)
			}
			return false
		}
		message = getRuleMessage("只能包含字母")
	case "alphanum":
		callback = func(v interface{}) bool {
			if val, ok := v.(string); ok {
				return IsAlphaNum(val)
			}
			return false
		}
		message = getRuleMessage("只能包含字母和数字")
	case "chinese":
		callback = func(v interface{}) bool {
			if val, ok := v.(string); ok {
				return IsChinese(val)
			}
			return false
		}
		message = getRuleMessage("只能包含中文字符")
	case "idcard":
		callback = func(v interface{}) bool {
			if val, ok := v.(string); ok {
				return IsIdCard(val)
			}
			return false
		}
		message = getRuleMessage("不是有效的身份证号码")
	case "min":
		if ruleParams != "" {
			min, err := strconv.Atoi(ruleParams)
			if err == nil {
				callback = func(v interface{}) bool {
					return Min(v, min)
				}
				message = getRuleMessage("不能小于%v")
			}
		}
	case "max":
		if ruleParams != "" {
			max, err := strconv.Atoi(ruleParams)
			if err == nil {
				// 保存max值为interface{}类型，确保Max函数参数类型匹配
				maxValue := interface{}(max)
				callback = func(v interface{}) bool {
					result := Max(v, maxValue)
					return result
				}
				message = getRuleMessage("不能大于%v")
			}
		}
	case "length":
		if ruleParams != "" {
			length, err := strconv.Atoi(ruleParams)
			if err == nil {
				callback = func(v interface{}) bool {
					return Length(v, length)
				}
				message = getRuleMessage("长度必须等于%d")
			}
		}
		// 可以根据需要添加更多的验证规则
	}

	// 如果有回调函数，则添加验证规则
	if callback != nil {
		validator.AddRule(field, ruleName, message, callback)
	}
	return true
}

// structToMap 将结构体转换为map[string]interface{}
func structToMap(dataValue reflect.Value) map[string]interface{} {
	dataMap := make(map[string]interface{})
	dataType := dataValue.Type()

	for i := 0; i < dataValue.NumField(); i++ {
		field := dataType.Field(i)
		fieldName := field.Name
		jsonTag := field.Tag.Get("json")
		fieldValue := dataValue.Field(i).Interface()

		// 使用json标签作为键（如果存在且不为-）
		fieldKey := fieldName
		if jsonTag != "" && jsonTag != "-" {
			jsonParts := strings.Split(jsonTag, ",")
			fieldKey = jsonParts[0]
		}

		dataMap[fieldKey] = fieldValue
	}
	return dataMap
}

// ParseValidationError 解析验证错误信息，返回友好的错误提示
func ParseValidationError(errMsg string) string {
	// 处理必填字段错误
	if strings.Contains(errMsg, "required") {
		// 提取字段名
		field := extractFieldName(errMsg)
		if field != "" {
			return field + "不能为空"
		}
		return "缺少必填参数"
	}

	// 处理长度验证错误
	if strings.Contains(errMsg, "len") {
		field := extractFieldName(errMsg)
		if field != "" {
			// 特殊处理手机号
			if field == "Telephone" || field == "tel" || field == "phone" || field == "mobile" {
				return "手机号格式错误"
			}
			return field + "格式错误"
		}
	}

	// 处理邮箱格式错误
	if strings.Contains(errMsg, "email") {
		return "邮箱格式错误"
	}

	// 处理数字范围错误
	if (strings.Contains(errMsg, "min") || strings.Contains(errMsg, "max")) && strings.Contains(errMsg, "number") {
		field := extractFieldName(errMsg)
		if field != "" {
			return field + "数值不在有效范围内"
		}
	}

	// 默认错误信息
	return "请求参数错误"
}

// 以下是从PHP Validate类迁移的验证规则函数

// IsEmail 验证是否为有效的邮箱地址
func IsEmail(email string) bool {
	// 使用正则表达式验证邮箱
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(pattern, email)
	return matched
}

// IsMobile 验证是否为有效的手机号码
func IsMobile(phone string) bool {
	// 使用正则表达式验证手机号（中国大陆手机号）
	pattern := `^1[3-9]\d{9}$`
	matched, _ := regexp.MatchString(pattern, phone)
	return matched
}

// IsURL 验证是否为有效的URL地址
func IsURL(url string) bool {
	// 使用正则表达式验证URL
	pattern := `^(https?://)?([\da-z.-]+)\.([a-z.]{2,6})([/\w .-]*)*/?$`
	matched, _ := regexp.MatchString(pattern, url)
	return matched
}

// IsIP 验证是否为有效的IP地址
func IsIP(ip string, version string) bool {
	var pattern string
	if version == "ipv6" {
		pattern = `^([0-9a-fA-F]{0,4}:){1,7}([0-9a-fA-F]{1,4}|localhost)$`
	} else {
		// 默认ipv4
		pattern = `^((25[0-5]|2[0-4]\d|[01]?\d\d?)\.){3}(25[0-5]|2[0-4]\d|[01]?\d\d?)$`
	}
	matched, _ := regexp.MatchString(pattern, ip)
	return matched
}

// IsDate 验证是否为有效的日期格式
func IsDate(date string) bool {
	// 支持多种日期格式
	formats := []string{
		"2006-01-02",
		"2006/01/02",
		"2006.01.02",
		"2006-01-02 15:04:05",
		"2006/01/02 15:04:05",
		"2006.01.02 15:04:05",
	}

	for _, format := range formats {
		if _, err := time.Parse(format, date); err == nil {
			return true
		}
	}
	return false
}

// IsNumber 验证是否为数字
func IsNumber(value string) bool {
	pattern := `^\d+$`
	matched, _ := regexp.MatchString(pattern, value)
	return matched
}

// IsAlpha 验证是否为字母
func IsAlpha(value string) bool {
	pattern := `^[a-zA-Z]+$`
	matched, _ := regexp.MatchString(pattern, value)
	return matched
}

// IsAlphaNum 验证是否为字母和数字组合
func IsAlphaNum(value string) bool {
	pattern := `^[a-zA-Z0-9]+$`
	matched, _ := regexp.MatchString(pattern, value)
	return matched
}

// IsAlphaDash 验证是否为字母、数字、下划线、破折号组合
func IsAlphaDash(value string) bool {
	pattern := `^[a-zA-Z0-9_-]+$`
	matched, _ := regexp.MatchString(pattern, value)
	return matched
}

// IsChinese 验证是否为中文字符
func IsChinese(value string) bool {
	pattern := `^[\p{Han}]+$`
	matched, _ := regexp.MatchString(pattern, value)
	return matched
}

// IsBetween 验证值是否在指定范围内
func IsBetween(value interface{}, min, max interface{}) bool {
	switch v := value.(type) {
	case int:
		return v >= min.(int) && v <= max.(int)
	case int64:
		return v >= min.(int64) && v <= max.(int64)
	case float64:
		return v >= min.(float64) && v <= max.(float64)
	case string:
		return len(v) >= min.(int) && len(v) <= max.(int)
	default:
		return false
	}
}

// IsIdCard 验证是否为有效的身份证号码
func IsIdCard(idCard string) bool {
	// 18位身份证验证规则
	pattern := `^(^[1-9]\d{5}(18|19|([23]\d))\d{2}((0[1-9])|(10|11|12))(([0-2][1-9])|10|20|30|31)\d{3}[0-9Xx]$)|(^[1-9]\d{5}\d{2}((0[1-9])|(10|11|12))(([0-2][1-9])|10|20|30|31)\d{3}$)$`
	matched, _ := regexp.MatchString(pattern, idCard)
	return matched
}

// IsZipCode 验证是否为有效的邮政编码
func IsZipCode(zipCode string) bool {
	// 中国邮政编码验证规则（6位数字）
	pattern := `^\d{6}$`
	matched, _ := regexp.MatchString(pattern, zipCode)
	return matched
}

// IsBoolean 验证是否为布尔值
func IsBoolean(value interface{}) bool {
	switch v := value.(type) {
	case bool:
		return true
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return v == 0 || v == 1
	case float32, float64:
		return v == 0 || v == 1
	case string:
		return v == "true" || v == "false" || v == "0" || v == "1"
	default:
		return false
	}
}

// IsAccepted 验证是否为可接受的值（如'yes', 'on', '1'）
func IsAccepted(value interface{}) bool {
	switch v := value.(type) {
	case bool:
		return v
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return v == 1
	case float32, float64:
		return v == 1
	case string:
		return v == "yes" || v == "on" || v == "1" || v == "true"
	default:
		return false
	}
}

// IsAfter 验证日期是否在指定日期之后
func IsAfter(date, afterDate string) bool {
	// 尝试解析日期
	d, err := time.Parse("2006-01-02", date)
	if err != nil {
		// 尝试其他格式
		otherFormats := []string{
			"2006/01/02",
			"2006.01.02",
			"2006-01-02 15:04:05",
		}
		var parseErr error
		for _, format := range otherFormats {
			d, parseErr = time.Parse(format, date)
			if parseErr == nil {
				break
			}
		}
		if parseErr != nil {
			return false
		}
	}

	// 尝试解析比较日期
	after, err := time.Parse("2006-01-02", afterDate)
	if err != nil {
		return false
	}

	return d.After(after)
}

// IsBefore 验证日期是否在指定日期之前
func IsBefore(date, beforeDate string) bool {
	// 尝试解析日期
	d, err := time.Parse("2006-01-02", date)
	if err != nil {
		// 尝试其他格式
		otherFormats := []string{
			"2006/01/02",
			"2006.01.02",
			"2006-01-02 15:04:05",
		}
		var parseErr error
		for _, format := range otherFormats {
			d, parseErr = time.Parse(format, date)
			if parseErr == nil {
				break
			}
		}
		if parseErr != nil {
			return false
		}
	}

	// 尝试解析比较日期
	before, err := time.Parse("2006-01-02", beforeDate)
	if err != nil {
		return false
	}

	return d.Before(before)
}

// InArray 验证值是否在指定数组中
func InArray(value interface{}, array interface{}) bool {
	switch arr := array.(type) {
	case []string:
		for _, v := range arr {
			if v == fmt.Sprintf("%v", value) {
				return true
			}
		}
	case []int:
		val, ok := value.(int)
		if ok {
			for _, v := range arr {
				if v == val {
					return true
				}
			}
		}
	case []int64:
		val, ok := value.(int64)
		if ok {
			for _, v := range arr {
				if v == val {
					return true
				}
			}
		}
	case []float64:
		val, ok := value.(float64)
		if ok {
			for _, v := range arr {
				if v == val {
					return true
				}
			}
		}
	}
	return false
}

// Min 验证值是否大于等于最小值
func Min(value interface{}, min interface{}) bool {
	switch v := value.(type) {
	case int:
		return v >= min.(int)
	case int64:
		return v >= min.(int64)
	case float64:
		return v >= min.(float64)
	case string:
		// 尝试将字符串转换为数值进行比较
		if num, err := strconv.ParseFloat(v, 64); err == nil {
			minNum, ok := min.(int)
			if ok {
				return num >= float64(minNum)
			}
		}
		// 无法转换为数值则返回验证失败
		return false
	default:
		return false
	}
}

// Max 验证值是否小于等于最大值
func Max(value interface{}, max interface{}) bool {
	switch v := value.(type) {
	case int:
		return v <= max.(int)
	case int64:
		return v <= max.(int64)
	case float64:
		return v <= max.(float64)
	case string:
		// 先尝试将字符串转换为数值进行比较
		if num, err := strconv.ParseFloat(v, 64); err == nil {
			maxNum, ok := max.(int)
			if ok {
				return num <= float64(maxNum)
			}
		}
		// 无法转换为数值则返回验证失败
		return false
	default:
		return false
	}
}

// Length 验证值的长度是否等于指定值
func Length(value interface{}, length int) bool {
	switch v := value.(type) {
	case string:
		return len(v) == length
	case []interface{}:
		return len(v) == length
	default:
		// 对于其他类型，转换为字符串后检查长度
		return len(fmt.Sprintf("%v", v)) == length
	}
}

// Confirm 验证字段值是否与另一字段值一致
func Confirm(value, confirmValue interface{}) bool {
	return fmt.Sprintf("%v", value) == fmt.Sprintf("%v", confirmValue)
}

// Different 验证字段值是否与另一字段值不同
func Different(value, differentValue interface{}) bool {
	return fmt.Sprintf("%v", value) != fmt.Sprintf("%v", differentValue)
}

// DateFormat 验证日期是否符合指定格式
func DateFormat(date, format string) bool {
	_, err := time.Parse(format, date)
	return err == nil
}

// IsExpire 验证时间是否在有效期内
func IsExpire(startTime, endTime string) bool {
	// 解析开始时间
	start, err := time.Parse("2006-01-02 15:04:05", startTime)
	if err != nil {
		// 尝试其他格式
		otherFormats := []string{
			"2006-01-02",
			"2006/01/02",
			"2006.01.02",
		}
		var parseErr error
		for _, f := range otherFormats {
			start, parseErr = time.Parse(f, startTime)
			if parseErr == nil {
				break
			}
		}
		if parseErr != nil {
			return false
		}
	}

	// 解析结束时间
	end, err := time.Parse("2006-01-02 15:04:05", endTime)
	if err != nil {
		return false
	}

	// 获取当前时间
	now := time.Now()

	// 检查当前时间是否在有效期内
	return now.After(start) && now.Before(end)
}

// AllowIP 验证IP是否在允许的IP列表中
func AllowIP(ip string, allowIPs []string) bool {
	for _, allowIP := range allowIPs {
		// 支持IP段（简单实现，仅支持CIDR格式）
		if strings.Contains(allowIP, "/") {
			if isIPInCIDR(ip, allowIP) {
				return true
			}
		} else {
			if ip == allowIP {
				return true
			}
		}
	}
	return false
}

// DenyIP 验证IP是否在禁止的IP列表中
func DenyIP(ip string, denyIPs []string) bool {
	for _, denyIP := range denyIPs {
		// 支持IP段（简单实现，仅支持CIDR格式）
		if strings.Contains(denyIP, "/") {
			if isIPInCIDR(ip, denyIP) {
				return true
			}
		} else {
			if ip == denyIP {
				return true
			}
		}
	}
	return false
}

// isIPInCIDR 检查IP是否在CIDR网段内（简单实现）
func isIPInCIDR(ip, cidr string) bool {
	// 这里简化实现，实际项目中可以使用net包的功能
	// 解析CIDR
	ipAddr, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return false
	}
	// 解析IP
	checkIP := net.ParseIP(ip)
	if checkIP == nil {
		return false
	}
	// 检查IP是否在网段内
	return ipNet.Contains(checkIP) && !ipAddr.Equal(checkIP) // 排除网络地址
}

// IsImage 验证文件是否为图片类型
func IsImage(filePath string) bool {
	// 获取文件后缀
	ext := strings.ToLower(filePath[strings.LastIndex(filePath, ".")+1:])
	// 常见图片类型
	imageTypes := map[string]bool{"jpg": true, "jpeg": true, "png": true, "gif": true, "bmp": true, "webp": true}
	return imageTypes[ext]
}

// IsArray 验证值是否为数组类型
func IsArray(value interface{}) bool {
	_, ok := value.([]interface{})
	return ok
}

// IsObject 验证值是否为对象类型
func IsObject(value interface{}) bool {
	_, ok := value.(map[string]interface{})
	return ok
}

// IsNumeric 验证值是否为数值类型
func IsNumeric(value interface{}) bool {
	switch value.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return true
	case string:
		// 简单检查字符串是否为纯数字
		str := value.(string)
		if len(str) == 0 {
			return false
		}
		// 允许负号和小数点
		numericRegex := regexp.MustCompile(`^-?\d+(\.\d+)?$`)
		return numericRegex.MatchString(str)
	default:
		return false
	}
}

// ExtractErrorMsg 从验证错误中提取友好的错误提示
func ExtractErrorMsg(err error) string {
	if err == nil {
		return ""
	}

	// 提取错误信息
	msg := err.Error()
	// 可以根据需要添加更多的错误信息处理逻辑
	return msg
}

// ValidatorRule 定义验证规则
type ValidatorRule struct {
	Field    string
	Rule     string
	Message  string
	Callback func(interface{}) bool
}

// Validator 验证器结构体
type Validator struct {
	Rules  []ValidatorRule
	Scenes map[string][]string
}

// NewValidator 创建新的验证器实例
func NewValidator() *Validator {
	return &Validator{
		Rules:  make([]ValidatorRule, 0),
		Scenes: make(map[string][]string),
	}
}

// AddRule 添加验证规则
func (v *Validator) AddRule(field, rule, message string, callback func(interface{}) bool) {
	v.Rules = append(v.Rules, ValidatorRule{
		Field:    field,
		Rule:     rule,
		Message:  message,
		Callback: callback,
	})
}

// AddScene 添加验证场景
func (v *Validator) AddScene(scene string, fields []string) {
	v.Scenes[scene] = fields
}

// Check 执行验证
func (v *Validator) Check(data map[string]interface{}, scene string) (bool, string) {
	// 确定需要验证的字段
	var fields []string
	if scene != "" && v.Scenes[scene] != nil {
		fields = v.Scenes[scene]
	} else {
		// 默认验证所有规则中的字段
		fieldMap := make(map[string]bool)
		for _, rule := range v.Rules {
			fieldMap[rule.Field] = true
		}
		for field := range fieldMap {
			fields = append(fields, field)
		}
	}

	// 执行验证
	for _, rule := range v.Rules {
		// 检查是否在场景中需要验证该字段
		needCheck := false
		for _, field := range fields {
			if field == rule.Field {
				needCheck = true
				break
			}
		}
		if !needCheck {
			continue
		}

		// 获取字段值
		value, exists := data[rule.Field]
		if !exists {
			value = nil
		}

		// 执行验证回调
		if !rule.Callback(value) {
			return false, rule.Message
		}
	}

	return true, ""
}

// ValidateGin 验证Gin请求参数
func ValidateGin(c *gin.Context, v *Validator, scene string) (bool, string) {
	// 解析请求参数到map
	data := make(map[string]interface{})

	// 绑定查询参数
	queryParams := c.Request.URL.Query()
	for key, values := range queryParams {
		if len(values) > 1 {
			data[key] = values
		} else if len(values) == 1 {
			data[key] = values[0]
		}
	}

	// 绑定表单参数
	if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "PATCH" {
		if err := c.Request.ParseForm(); err == nil {
			for key, values := range c.Request.Form {
				if len(values) > 1 {
					data[key] = values
				} else if len(values) == 1 {
					data[key] = values[0]
				}
			}
		}

		// 尝试绑定JSON
		jsonData := make(map[string]interface{})
		if err := c.ShouldBindJSON(&jsonData); err == nil {
			for key, value := range jsonData {
				data[key] = value
			}
		}
	}

	// 执行验证
	return v.Check(data, scene)
}

// extractFieldName 从错误信息中提取字段名
// 简单解析常见的validator错误格式
func extractFieldName(errMsg string) string {
	// 处理形如: Key: 'StructName.Field' Error:Field validation for 'Field' failed on the 'required' tag
	if strings.Contains(errMsg, "Key: ") && strings.Contains(errMsg, "Error:") {
		parts := strings.Split(errMsg, "Key: '")
		if len(parts) > 1 {
			fieldPart := strings.Split(parts[1], "'")
			if len(fieldPart) > 0 {
				// 提取字段名，如 "StructName.Field" -> "Field"
				dotParts := strings.Split(fieldPart[0], ".")
				if len(dotParts) > 1 {
					return dotParts[len(dotParts)-1]
				}
				return fieldPart[0]
			}
		}
	}

	// 处理形如: Field validation for 'Field' failed on the 'required' tag
	if strings.Contains(errMsg, "Field validation for '") && strings.Contains(errMsg, "' failed on the") {
		parts := strings.Split(errMsg, "Field validation for '")
		if len(parts) > 1 {
			fieldPart := strings.Split(parts[1], "' failed on the")
			if len(fieldPart) > 0 {
				return fieldPart[0]
			}
		}
	}

	// 直接查找字段名（简单场景）
	if strings.Contains(errMsg, "' failed on the") {
		parts := strings.Split(errMsg, "' failed on the")
		if len(parts) > 0 {
			fieldPart := strings.TrimSuffix(parts[0], "'")
			return fieldPart
		}
	}

	return ""
}
