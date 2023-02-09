package model

import (
	"demo_backend/tool"

	"gorm.io/gorm"
)

type Account struct {
	gorm.Model
	Name string `v:"name@required|length:4,30#请输入账号|账号长度为{min}到{max}位"`
	Pass string `v:"pass@required|length:4,30#请输入密码|密码长度为{min}到{max}位"`
}

type Input_ChangePassword struct {
	Pass  string `v:"pass@required|length:4,30#请输入密码|密码长度为{min}到{max}位"`
	Pass2 string `v:"pass2@required|same:Pass#请确认密码|两次输入密码不同"`
}

func init() {
	tool.GetGormConnection().AutoMigrate(&Account{})
}
