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
	return do.Account_QueryOneById(id)
}

// 当前账号是否管理员，以便执行特权操作，比如创建删除账号
// 如果不是管理员，则中断后续处理流程
func assure_current_user_is_admin(r *ghttp.Request) {
	//真实的管理员验证比这里复杂
	if r.Session.MustGet("accountId").Uint() != 1 {
		response.FailMsg(r, "您不是管理员")
	}
}

// 注销登录 /api/account/logout
func handler_api_account_logout(r *ghttp.Request) {
	r.Session.RemoveAll()
	response.DoneData(r, "已登出")
}

// 修改自己账号的密码 /api/account/changepass
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

// 管理员-注册账号 /api/account/register
func handler_api_account_register(r *ghttp.Request) {
	assure_current_user_is_admin(r)

	ac := do.RegisterInput_To_Account(do.Request_To_RegisterInput(r))

	//用User查询账号是否存在
	if do.Account_QueryExistsByFields(&dbmodel.Account{User: ac.User}) {
		response.FailMsg(r, "该账号已存在")
	}
	response.DoneData(r, do.Account_CreateOne(ac))
}

// 管理员-删除账号 /api/account/delete
func handler_api_account_delete(r *ghttp.Request) {
	assure_current_user_is_admin(r) //必须是管理员
	delac := do.DeleteInput_To_Account(do.Request_To_DeleteInput(r))

	//禁止管理员删除自己
	myac := get_current_user_account(r)
	if myac.User == delac.User {
		response.FailMsg(r, "不能删除自己的账号")
	}

	// 删除
	cnt := do.Account_DeleteMultiByFields(delac)
	if cnt == 1 {
		response.DoneData(r, "删除成功")
	} else {
		response.FailMsg(r, "该账号不存在")
	}
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
	group.POST("/login", handler_api_account_login)

	// 需要登录
	group.Group("/", func(group *ghttp.RouterGroup) {
		group.Middleware(middleware.MiddlewareAuth)
		group.POST("/register", handler_api_account_register) //一般管理员才能注册账号
		group.POST("/logout", handler_api_account_logout)
		group.POST("/delete", handler_api_account_delete)
		group.POST("/changepass", handler_api_account_changepass)

	})
}
