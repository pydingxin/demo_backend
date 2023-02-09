package api

import (
	"demo_backend/model"
	"demo_backend/tool"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

//所有不登录时能访问的接口都在这里

func handler_api_unauth_acount_login(r *ghttp.Request) {
	// 账号登录 /api/unauth/account/login

	// 获取输入并验证
	var in *model.Account
	if err := r.Parse(&in); err != nil {
		r.Response.WriteJsonExit(g.Map{"status": false, "msg": err.Error()})
	}

	// 查询
	queriedAccountPtr := new(model.Account)
	result := tool.GetGormConnection().Where(" name = ? And pass = ?", in.Name, in.Pass).Find(queriedAccountPtr) //表字段名已转小写
	if result.Error != nil {
		r.Response.WriteJsonExit(g.Map{"status": false, "msg": result.Error.Error()})
	}

	// 回复
	if result.RowsAffected == 1 {
		r.Session.Set("session_alive", true) // 根据alive鉴权
		r.Session.Set("accountId", queriedAccountPtr.ID)
		// g.Dump("登录成功", "input:", in, "queriedAccount", queriedAccountPtr)
		r.Response.WriteJsonExit(g.Map{"status": true, "data": queriedAccountPtr})
	} else {
		r.Response.WriteJsonExit(g.Map{"status": false, "msg": "账号或密码错误"})
	}
}

func RouterGroup_ApiUnauth(group *ghttp.RouterGroup) {
	group.POST("/account/login", handler_api_unauth_acount_login)
}
