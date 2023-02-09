package api

import (
	"demo_backend/model"
	"demo_backend/tool"
	"fmt"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

func handler_api_account_logout(r *ghttp.Request) {
	// /api/account/logout
	r.Session.RemoveAll()
	r.Response.WriteJsonExit(g.Map{"status": true})
}

func handler_api_account_change_password(r *ghttp.Request) {
	// /api/account/change_password
	// 获取输入并验证
	var in *model.Input_ChangePassword
	if err := r.Parse(&in); err != nil {
		r.Response.WriteJsonExit(g.Map{"status": false, "msg": err.Error()})
	}

	accountid := r.Session.MustGet("accountId", -1).Int()
	if accountid == -1 {
		r.Response.WriteJsonExit(g.Map{"status": false, "msg": "handler_api_account_change_password 会话数据丢失"})
	}
	tool.GetGormConnection().Model(&model.Account{}).Where("id = ?", accountid).Update("pass", in.Pass)
	fmt.Print("更新密码 id=", accountid, " pass=", in.Pass)
	r.Response.WriteJsonExit(g.Map{"status": true})

}

func RouterGroup_ApiAccount(group *ghttp.RouterGroup) {
	group.POST("/logout", handler_api_account_logout)
	group.POST("/change_password", handler_api_account_change_password)
}
