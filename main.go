package main

import (
	"demo_backend/api"
	"demo_backend/middleware"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

func main() {
	s := g.Server()
	// s.EnableHTTPS("./httpsCertification/server.crt", "./httpsCertification/server.key")

	s.Use(middleware.MiddlewareLog)          //业务日志
	s.Use(middleware.MiddlewareErrorHandler) //异常日志
	s.Use(middleware.MiddlewareCORS)         //跨域

	s.Group("/api", func(group *ghttp.RouterGroup) {
		// 哪些接口需要登录，由各组自己确定
		group.Group("/account", api.RouterGroup_Account) // 账号类

	})

	s.SetPort(80)
	s.SetIndexFolder(true) //静态文件
	s.SetServerRoot("./static")
	s.Run()
}
