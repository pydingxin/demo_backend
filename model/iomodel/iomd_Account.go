package iomodel

// CMD_IOMODEL_PARSE_FROM_REQUEST 生成代码-request中获取iomodel
// CMD_IOMODEL_CONVERT_TO_DBMODEL_Account 生成代码-iomodel转为dbmodel
type RegisterInput struct {
	User string `v:"required|length:4,40#请输入账号|账号长度为{min}到{max}位"`
	Pass string `v:"required|length:4,40#请输入密码|密码长度为{min}到{max}位"`
}

// CMD_IOMODEL_PARSE_FROM_REQUEST #request转LoginInput
// CMD_IOMODEL_CONVERT_TO_DBMODEL_Account #LoginInput转Account
type LoginInput struct {
	User string `v:"required|length:4,40#请输入账号|账号长度为{min}到{max}位"`
	Pass string `v:"required|length:4,40#请输入密码|密码长度为{min}到{max}位"`
}

// # CMD_IOMODEL_QUERY_FROM_DBMODEL_Account Account查询生成子结构体
