package api

import (
	"demo_backend/middleware"
	"demo_backend/model/dbmodel"
	"demo_backend/model/do"
	"demo_backend/pkg/response"

	"github.com/gogf/gf/v2/net/ghttp"
)

// 获取当前账号的信息
func get_current_user_account(r *ghttp.Request) *dbmodel.Account {
	id := r.Session.MustGet("accountId").Uint()
	print("accountId=", id)
	return do.Account_QueryOneById(id)
}

// 注销登录 /api/account/logout
func handler_api_account_logout(r *ghttp.Request) {
	r.Session.RemoveAll()
	response.DoneData(r, "已登出")
}

// 修改密码 /api/account/changepass
func handler_api_account_changepass(r *ghttp.Request) {
	input := do.Request_To_ChangePassInput(r)

	curAC := get_current_user_account(r)
	if curAC.User != input.User || curAC.Pass != input.Pass {
		response.FailMsg(r, "用户名或密码错误") //当前用户是否正确输入了账号
	}
	// 用ID和更新非零字段
	do.Account_UpdateOneById(&dbmodel.Account{ID: curAC.ID, Pass: input.NewPass})
	response.DoneData(r, "修改成功")
}

// 注册账号 /api/account/register
func handler_api_account_register(r *ghttp.Request) {
	ac := do.RegisterInput_To_Account(do.Request_To_RegisterInput(r))

	//用User查询账号是否存在
	if do.Account_QueryExistsByFields(&dbmodel.Account{User: ac.User}) {
		response.FailMsg(r, "该账号已存在")
	}
	response.DoneData(r, do.Account_CreateOne(ac))
}

// 登录账号 /api/account/login
func handler_api_account_login(r *ghttp.Request) {
	ac := do.LoginInput_To_Account(do.Request_To_LoginInput(r))

	// 用User+Pass查询账号是否存在
	if do.Account_QueryExistsByFields(ac) {
		r.Session.Set("session_alive", true)
		r.Session.Set("accountId", ac.ID)
		response.DoneData(r, "登录成功")
	} else {
		response.FailMsg(r, "登录失败")
	}
}

func RouterGroup_Account(group *ghttp.RouterGroup) {
	// 不需要登录
	group.POST("/register", handler_api_account_register)
	group.POST("/login", handler_api_account_login)

	// 需要登录
	group.Group("/", func(group *ghttp.RouterGroup) {
		group.Middleware(middleware.MiddlewareAuth)
		group.POST("/logout", handler_api_account_logout)
		group.POST("/changepass", handler_api_account_changepass)

	})
}
