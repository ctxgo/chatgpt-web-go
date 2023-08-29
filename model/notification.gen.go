// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameNotification = "notification"

// Notification mapped from table <notification>
type Notification struct {
	ID         int64     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Title      string    `gorm:"column:title;not null;comment:标题" json:"title"`     // 标题
	Content    string    `gorm:"column:content;not null;comment:内容" json:"content"` // 内容
	Sort       int32     `gorm:"column:sort;not null;default:1" json:"sort"`
	Status     int32     `gorm:"column:status;not null;comment:状态" json:"status"` // 状态
	CreateTime time.Time `gorm:"column:create_time;default:CURRENT_TIMESTAMP" json:"create_time"`
	UpdateTime time.Time `gorm:"column:update_time;default:CURRENT_TIMESTAMP" json:"update_time"`
}

// TableName Notification's table name
func (*Notification) TableName() string {
	return TableNameNotification
}