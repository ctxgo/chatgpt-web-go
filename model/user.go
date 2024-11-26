package model

import (
	"chatgpt-web-new-go/common/auth/password"
	"errors"

	"gorm.io/gorm"
)

// ComparePassword 检查密码是否匹配
func (user *User) ComparePassword(_password string) bool {
	return password.CheckHash(_password, user.Password)
}

// BeforeSave 保存前
func (user *User) BeforeSave(tx *gorm.DB) (err error) {
	if !tx.Statement.Changed("password") {
		return nil
	}
	var newPass string
	switch u := tx.Statement.Dest.(type) {
	case map[string]interface{}:
		newPass = u["password"].(string)
	case *User:
		newPass = u.Password
	case []*User:
		newPass = u[tx.Statement.CurDestIndex].Password
	default:
		return errors.New("User.BeforeSave hook 错误，类型断言失败")
	}
	if !password.IsHashed(newPass) {
		password := password.Hash(newPass)
		tx.Statement.SetColumn("password", password)
	}
	return
}
