// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameInviteRecord = "invite_record"

// InviteRecord mapped from table <invite_record>
type InviteRecord struct {
	ID         int64     `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	UserID     int64     `gorm:"column:user_id;primaryKey" json:"user_id"`
	InviteCode string    `gorm:"column:invite_code;not null;comment:邀请码" json:"invite_code"`               // 邀请码
	SuperiorID int64     `gorm:"column:superior_id;not null;comment:上级ID（一旦确定将不可修改）" json:"superior_id"`   // 上级ID（一旦确定将不可修改）
	Reward     string    `gorm:"column:reward;not null;comment:奖励" json:"reward"`                          // 奖励
	RewardType string    `gorm:"column:reward_type;not null;comment:奖励类型" json:"reward_type"`              // 奖励类型
	Status     int32     `gorm:"column:status;not null;default:3;comment:0-异常｜1-正常发放｜3-审核中" json:"status"` // 0-异常｜1-正常发放｜3-审核中
	Remarks    string    `gorm:"column:remarks;not null;comment:评论" json:"remarks"`                        // 评论
	IP         string    `gorm:"column:ip;not null" json:"ip"`
	UserAgent  string    `gorm:"column:user_agent;comment:ua" json:"user_agent"` // ua
	CreateTime time.Time `gorm:"column:create_time;not null;default:CURRENT_TIMESTAMP" json:"create_time"`
	UpdateTime time.Time `gorm:"column:update_time;not null;default:CURRENT_TIMESTAMP" json:"update_time"`
	IsDelete   int32     `gorm:"column:is_delete;not null" json:"is_delete"`
}

// TableName InviteRecord's table name
func (*InviteRecord) TableName() string {
	return TableNameInviteRecord
}
