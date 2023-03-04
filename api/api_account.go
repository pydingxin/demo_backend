package api

import (
	"demo_backend/middleware"
	"demo_backend/model/do"
	"demo_backend/pkg/response"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

func handler_api_account_logout(r *ghttp.Request) {
	// /api/account/logout
	r.Session.RemoveAll()
	r.Response.WriteJsonExit(g.Map{"status": true})
}

func handler_api_account_changepass(r *ghttp.Request) {
	// 	// /api/account/change_password
	// 	// 获取输入并验证
	// 	var in *model.Input_ChangePassword
	// 	if err := r.Parse(&in); err != nil {
	// 		r.Response.WriteJsonExit(g.Map{"status": false, "msg": err.Error()})
	// 	}

	// 	accountid := r.Session.MustGet("accountId", -1).Int()
	// 	if accountid == -1 {
	// 		r.Response.WriteJsonExit(g.Map{"status": false, "msg": "handler_api_account_changepass 会话数据丢失"})
	// 	}
	// 	tool.GetGormConnection().Model(&model.Account{}).Where("id = ?", accountid).Update("pass", in.Pass)
	// 	fmt.Print("更新密码 id=", accountid, " pass=", in.Pass)
	// 	r.Response.WriteJsonExit(g.Map{"status": true})

}

func handler_api_account_register(r *ghttp.Request) {
	// 注册账号 /api/account/register
	in := do.Request_To_RegisterInput(r)
	ac := do.RegisterInput_To_Account(in)
	ac = do.Account_CreateOne(ac)
	response.DoneData(r, ac)
}

func handler_api_account_login(r *ghttp.Request) {
	// 登录账号 /api/account/change_login
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
