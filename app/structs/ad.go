package structs

import (
	"time"
)

// Ad 广告模型，定义了广告数据的结构
// 使用json标签来指定JSON序列化时的字段名
type Ad struct {
	ID               int       `gorm:"column:id" json:"id"`
	Weid             int       `gorm:"column:weid" json:"weid"`
	Ocid             int       `gorm:"column:ocid" json:"ocid"`
	Sid              int       `gorm:"column:sid" json:"sid"`
	Ptype            int       `gorm:"column:ptype" json:"ptype"`
	Url              string    `gorm:"column:url" json:"url"`
	Title            string    `gorm:"column:title" json:"title"`
	Pic              string    `gorm:"column:pic" json:"pic"`
	Sort             int       `gorm:"column:sort" json:"sort"`
	Status           int       `gorm:"column:status" json:"status"`
	ValidPeriodStart time.Time `gorm:"column:valid_period_start" json:"valid_period_start"`
	ValidPeriodEnd   time.Time `gorm:"column:valid_period_end" json:"valid_period_end"`
	PageUrl          string    `gorm:"column:page_url" json:"page_url"`
}