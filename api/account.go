package api

import (
	"demo_backend/middleware"
	"demo_backend/model/dbmodel"
	"demo_backend/model/do"
	"demo_backend/pkg/response"
	"demo_backend/pkg/tool"

	"github.com/gogf/gf/v2/net/ghttp"
)

func LoginWorks(ac *dbmodel.Account, r *ghttp.Request) {
	//登录触发事件
	if id, err := r.Session.Id(); err == nil {
		tool.SessionIdSet(ac.User, id) //在tool里记录sessionid，给websocket用
	}
}
func LogoutWorks(r *ghttp.Request) {
	//登出触发事件
}

func init() {
	// 初始化account
	if do.Account_QueryAllCount() == 0 {
		do.Account_CreateOne(&dbmodel.Account{User: "admin", Pass: "Abc123."})
	}
}

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
	LogoutWorks(r)
	r.Session.RemoveAll()
	response.DoneData(r, "已登出")
}

// 修改自己账号的密码 /api/account/changepass
func handler_api_account_changepass(r *ghttp.Request) {
	input := do.Request_To_AccountChangePassInput(r)

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

	ac := do.AccountRegisterInput_To_Account(do.Request_To_AccountRegisterInput(r))

	//用User查询账号是否存在
	if do.Account_QueryExistsByFields(&dbmodel.Account{User: ac.User}) {
		response.FailMsg(r, "该账号已存在")
	}
	response.DoneData(r, do.Account_CreateOne(ac))
}

// 管理员-删除账号 /api/account/delete
func handler_api_account_delete(r *ghttp.Request) {
	assure_current_user_is_admin(r) //必须是管理员
	delac := do.AccountDeleteInput_To_Account(do.Request_To_AccountDeleteInput(r))

	//禁止管理员删除自己
	myac := get_current_user_account(r)
	if myac.User == delac.User {
		response.FailMsg(r, "管理员不能删除自己的账号")
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
	ac := do.AccountLoginInput_To_Account(do.Request_To_AccountLoginInput(r))

	// 用User+Pass查询账号是否存在
	if do.Account_QueryExistsByFields(ac) {
		r.Session.Set("session_alive", true)
		r.Session.Set("accountId", ac.ID)
		LoginWorks(ac, r)
		response.DoneData(r, "登录成功")
	} else {
		response.FailMsg(r, "登录失败")
	}
}

// 管理员获取账号列表 /api/account/list
func handler_api_account_list(r *ghttp.Request) {
	assure_current_user_is_admin(r) //必须是管理员
	msg := do.Account_QueryAll_SmallAccount()
	response.DoneData(r, msg)
}

// 管理员编辑账号 /api/account/edit
func handler_api_account_edit(r *ghttp.Request) {
	assure_current_user_is_admin(r) //必须是管理员
	ac := do.SmallAccount_To_Account(do.Request_To_SmallAccount(r))
	if do.Account_UpdateOneById(ac) == 1 {
		response.DoneData(r, "修改成功")
	} else {
		response.FailMsg(r, "该ID不存在")
	}
}

func RouterGroup_Account(group *ghttp.RouterGroup) {
	// 不需要登录
	group.POST("/login", handler_api_account_login)

	// 需要登录
	group.Group("/", func(group *ghttp.RouterGroup) {
		group.Middleware(middleware.MiddlewareAuth)
		group.POST("/logout", handler_api_account_logout)
		group.POST("/delete", handler_api_account_delete)
		group.POST("/changepass", handler_api_account_changepass)

		//管理员操作
		group.POST("/register", handler_api_account_register)
		group.POST("/list", handler_api_account_list)
		group.POST("/edit", handler_api_account_edit)

	})
}
