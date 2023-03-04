package middleware

import (
	"demo_backend/pkg/response"

	"github.com/gogf/gf/v2/net/ghttp"
)

/*
鉴权中间件，鉴定是否登录
登录接口在 api.unauth.handler_api_unauth_acount_login 添加session的session_alive字段
登出接口在 api.account.handler_api_account_logout 清空session的所有字段
鉴权检验这个字段是否还在
*/
func MiddlewareAuth(r *ghttp.Request) {
	session_alive, _ := r.Session.Get("session_alive", false) //session续活
	if session_alive.Bool() {
		r.Middleware.Next()
	} else {
		response.FailMsg(r, "请先登录")
	}

}
