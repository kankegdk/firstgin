package structs

import (
	"time"
)

// Ad 广告模型，定义了广告数据的结构
// 使用json标签来指定JSON序列化时的字段名
type Ad struct {
	ID               int       `json:"id"`
	Weid             int       `json:"weid"`
	Ocid             int       `json:"ocid"`
	Sid              int       `json:"sid"`
	Ptype            int       `json:"ptype"`
	Url              string    `json:"url"`
	Title            string    `json:"title"`
	Pic              string    `json:"pic"`
	Sort             int       `json:"sort"`
	Status           int       `json:"status"`
	ValidPeriodStart time.Time `json:"valid_period_start"`
	ValidPeriodEnd   time.Time `json:"valid_period_end"`
	PageUrl          string    `json:"page_url"`
}