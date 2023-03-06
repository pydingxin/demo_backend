package iomodel

// 输入-注册账号
// CMD_IOMODEL_PARSE_FROM_REQUEST #request转此
// CMD_IOMODEL_CONVERT_TO_DBMODEL_Account #此转Account
type AccountRegisterInput struct {
	User string `v:"required|length:4,40#请输入账号|账号长度为{min}到{max}位"`
	Pass string `v:"required|password3#请输入密码|密码长度6~18位，包含数字、大小写字母和特殊字符"`
}

// 输入-登录账号
// CMD_IOMODEL_PARSE_FROM_REQUEST #request转此
// CMD_IOMODEL_CONVERT_TO_DBMODEL_Account #此转Account
type AccountLoginInput struct {
	User string `v:"required|length:4,40#请输入账号|账号长度为{min}到{max}位"`
	Pass string `v:"required|password3#请输入密码|密码长度6~18位，包含数字、大小写字母和特殊字符"`
}

// 输入-修改自己账号的密码
// CMD_IOMODEL_PARSE_FROM_REQUEST #request转此
// CMD_IOMODEL_CONVERT_TO_DBMODEL_Account #此转Account
type AccountChangePassInput struct {
	User     string `v:"required|length:4,40#请输入账号|账号长度为{min}到{max}位"`
	Pass     string `v:"required|password3#请输入密码|密码长度6~18位，包含数字、大小写字母和特殊字符"`
	NewPass  string `v:"required|password3#请输入新密码|新密码长度6~18位，包含数字、大小写字母和特殊字符"`
	NewPass2 string `v:"required|same:NewPass#请确认新密码|两次输入的新密码不相同"`
}

// 输入-删除账号
// CMD_IOMODEL_PARSE_FROM_REQUEST #request转此
// CMD_IOMODEL_CONVERT_TO_DBMODEL_Account #此转Account
type AccountDeleteInput struct {
	User string `v:"required|length:4,40#请输入账号|账号长度为{min}到{max}位"`
}

// 输出-管理员获取账号列表 输入-管理员编辑账号
// CMD_IOMODEL_PARSE_FROM_REQUEST 			#request转此
// CMD_IOMODEL_CONVERT_TO_DBMODEL_Account 	#此转Account
// CMD_IOMODEL_QUERY_FROM_DBMODEL_Account 	#Account库表生成
type SmallAccount struct {
	ID   uint
	User string `v:"required|length:4,40#请输入账号|账号长度为{min}到{max}位"`
	Pass string `v:"required|password3#请输入密码|密码长度6~18位，包含数字、大小写字母和特殊字符"`
}
