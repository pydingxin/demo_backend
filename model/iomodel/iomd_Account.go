package iomodel

// CMD_IOMODEL_PARSE_FROM_REQUEST #request转此
// CMD_IOMODEL_CONVERT_TO_DBMODEL_Account #此转Account
type RegisterInput struct {
	User string `v:"required|length:4,40#请输入账号|账号长度为{min}到{max}位"`
	Pass string `v:"required|password3#请输入密码|密码长度6~18位，包含数字、大小写字母和特殊字符"`
}

// CMD_IOMODEL_PARSE_FROM_REQUEST #request转此
// CMD_IOMODEL_CONVERT_TO_DBMODEL_Account #此转Account
type LoginInput struct {
	User string `v:"required|length:4,40#请输入账号|账号长度为{min}到{max}位"`
	Pass string `v:"required|password3#请输入密码|密码长度6~18位，包含数字、大小写字母和特殊字符"`
}

// CMD_IOMODEL_PARSE_FROM_REQUEST #request转此
// CMD_IOMODEL_CONVERT_TO_DBMODEL_Account #此转Account
type ChangePassInput struct {
	User     string `v:"required|length:4,40#请输入账号|账号长度为{min}到{max}位"`
	Pass     string `v:"required|password3#请输入密码|密码长度6~18位，包含数字、大小写字母和特殊字符"`
	NewPass  string `v:"required|password3#请输入新密码|新密码长度6~18位，包含数字、大小写字母和特殊字符"`
	NewPass2 string `v:"required|same:NewPass#请确认新密码|两次输入的新密码不相同"`
}

// CMD_IOMODEL_PARSE_FROM_REQUEST #request转此
// CMD_IOMODEL_CONVERT_TO_DBMODEL_Account #此转Account
type DeleteInput struct {
	User string `v:"required|length:4,40#请输入账号|账号长度为{min}到{max}位"`
}

// CMD_IOMODEL_QUERY_FROM_DBMODEL_Account Account查询生成子结构体
// CMD_IOMODEL_STRUCT_SAVE_TO_EXCEL_FILE
type SmallAccount struct {
	User string
	Pass string
}
